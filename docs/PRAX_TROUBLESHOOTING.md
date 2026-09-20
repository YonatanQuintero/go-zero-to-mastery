# 🔍 PRAX Production Troubleshooting & Real-World Practice Cases in Go

Este documento contiene el análisis arquitectónico de los **10 errores y problemas más comunes documentados en `prax-wiki`** (basado en incidentes reales, ADRs y playbooks de producción de PRAX), seguido de **20 problemas prácticos en Go (T01 a T20)** diseñados para dominar la depuración, resiliencia y optimización de sistemas concurrentes y cloud-native.

---

## 📑 Tabla de Contenido
1. [Análisis de los 10 Errores Más Comunes en PRAX](#1-los-10-errores-y-problemas-más-comunes-en-prax)
2. [20 Problemas Prácticos de Troubleshooting en Go (T01 a T20)](#2-20-problemas-prácticos-de-troubleshooting-en-go-t01-a-t20)
   - [Bloque 1: Interoperabilidad & Gobernanza de Almacenamiento (T01 a T05)](#bloque-1-interoperabilidad--gobernanza-de-almacenamiento-t01-a-t05)
   - [Bloque 2: Resiliencia de Pipelines de Datos & Concurrencia (T06 a T10)](#bloque-2-resiliencia-de-pipelines-de-datos--concurrencia-t06-a-t10)
   - [Bloque 3: Integridad de Datos, Deduplicación & Zero-Trust (T11 a T15)](#bloque-3-integridad-de-datos-deduplicación--zero-trust-t11-a-t15)
   - [Bloque 4: Redes, SRE & Optimización de Recursos en Cloud Run (T16 a T20)](#bloque-4-redes-sre--optimización-de-recursos-en-cloud-run-t16-a-t20)

---

# 1. Los 10 Errores y Problemas Más Comunes en PRAX

### 1. Incompatibilidad de Claves WebSafe Datastore entre Go y Java (`s~` vs `k~` y Protobuf v1 vs v3)
- **Referencia Documental:** `ADR-20260910-01`, QA Caso TC-B02 / Hallazgo 23 (`prax-ops-tools`).
- **Manifestación del Error:** Al crear o actualizar empleados en Go (`prax-go-employees-ms`) y luego intentar abrir la ficha sociodemográfica en Java (`prax-employees-service`), Java detonaba una excepción irrecuperable: `java.lang.IllegalArgumentException: id must not be equal to zero (HTTP 500)`.
- **Causa Raíz:**
  - El SDK nativo de Go (`cloud.google.com/go/datastore`) codifica claves con `Key.Encode()`, produciendo Base64 URL de Google Cloud Datastore v1 (prefijo `Eh...`).
  - El ecosistema Java downstream utiliza Objectify 6 y Google App Engine SDK v1.x, el cual espera claves WebSafe de Protobuf `Reference` v3 (prefijo obligatorio `ag...`).
  - Adicionalmente, las particiones de App ID diferían: en producción (`prax-prod`) se opera bajo `s~prax-prod` (`agtz...`), mientras que en desarrollo (`prax-dev`) se opera bajo `k~prax-dev` (`agpr...`). Al recibir `Eh...`, el parser de Objectify fallaba silenciosamente retornando `0L`, intentando buscar entidades por ID 0.
- **Solución en Go:** Prohibir `Key.Encode()` para referencias expuestas y utilizar obligatoriamente `gcputils.EncodeLegacyWebSafeKey(appID, namespace, key)` de `prax-go-common`, resolviendo el `appID` dinámicamente según la partición de la empresa.

---

### 2. Bloqueo del Streaming Buffer de BigQuery y Veneno en Pub/Sub (Poison Pill por UPDATE prematuro)
- **Referencia Documental:** `ADR-20260525-streaming-buffer-pubsub-lock`, Alerta de Producción `Oldest unacked message age`.
- **Manifestación del Error:** Mensajes en la suscripción de Pub/Sub quedaban retenidos hasta por 96 minutos en un bucle infinito de reintentos fallidos, bloqueando la sincronización de evaluaciones de trabajadores hacia BigQuery.
- **Causa Raíz:**
  - Un evento insertaba un empleado en BigQuery mediante la API de Streaming (`insertAll`).
  - Segundos después, otro evento intentaba ejecutar un comando SQL `UPDATE` o `MERGE` sobre ese mismo empleado.
  - BigQuery rechazaba la mutación con el error: `UPDATE or DELETE statement over table... would affect rows in the streaming buffer`. En BigQuery, los datos en el Streaming Buffer son temporalmente inmutables mientras se transfieren al almacenamiento columnar (tarda de 30 a 90 minutos).
  - Al fallar con HTTP 500, Pub/Sub realizaba un NACK inmediato, convirtiendo el mensaje en una "píldora venenosa" (Poison Pill).
- **Solución en Go:** Adoptar el patrón analítico **Append-Only** (hacer siempre `INSERT` con timestamp `updated_at` y consultar mediante vistas deduplicadas con `ROW_NUMBER() OVER(PARTITION BY id ORDER BY updated_at DESC) = 1`), o desviar fallos transitorios de streaming buffer a una Dead Letter Queue (DLQ) con retraso de reintento.

---

### 3. Cuello de Botella N+1 y Fuga de Memoria en Objectify Session Cache (OOM en Cloud Run)
- **Referencia Documental:** `ADR-202605-12-datastore-eager-loading.md`, Incidente masivo cliente Sodexo (~5,000 empleados).
- **Manifestación del Error:** Al generar reportes de evaluación masivos en Excel, el contenedor de Cloud Run moría por falta de memoria (`HTTP 503` / `Exit Code 137 OOM`).
- **Causa Raíz:**
  - **N+1 en Red:** Para verificar si una pregunta de consentimiento aplicaba, el servicio lanzaba una consulta individual a Datastore por cada uno de los 5,000 empleados, saturando el pool de conexiones SSL.
  - **Fuga en Session Cache:** El ORM de Objectify cargaba páginas de 500 registros manteniendo referencias fuertes en su memoria interna sin invocar `ofy().clear()`, acumulando decenas de miles de objetos en el Heap de la JVM (512 MiB).
- **Solución en Go:** Implementar **Eager Loading** cargando los catálogos en un mapa en memoria $O(1)$ previo a la iteración, procesar las filas en streaming mediante iteradores desacoplados y utilizar buffers reciclables con `sync.Pool`.

---

### 4. Saturación de Cuotas de Concurrencia en BigQuery (100 consultas simultáneas) y Rate Limits DML
- **Referencia Documental:** `adr-202605-11-bigquery-concurrency-sweetspot.md`.
- **Manifestación del Error:** En eventos de alta carga, BigQuery colapsaba con `Exceeded rate limits: too many concurrent queries`. Cuando un desarrollador restringió Cloud Run a `--concurrency=1`, Cloud Run empezó a rechazar peticiones legítimas de Pub/Sub con `HTTP 429 Too Many Requests`.
- **Causa Raíz:** BigQuery impone un límite estricto de 100 consultas concurrentes por proyecto y máximo 20 operaciones DML simultáneas por tabla. El autoescalado libre de Cloud Run (80 solicitudes simultáneas por contenedor $\times$ $N$ instancias) superaba por mucho este umbral analítico.
- **Solución en Go:** Aplicar la arquitectura del **"Punto Dulce"**: fijar `--max-instances=1` en el publisher, sintonizar una concurrencia interna estricta de 15 workers gobernados por un semáforo (`chan struct{}` de tamaño 15), y ajustar el Backoff mínimo de Pub/Sub a 10 segundos.

---

### 5. Entidades Duplicadas y Corrupción de Métricas por Migraciones sin Zero-Trust
- **Referencia Documental:** `PM-20260630-deduplicacion-entidades.md`.
- **Manifestación del Error:** Selects del frontend mostrando sedes y cargos duplicados; tras una limpieza con scripts externos, se produjeron errores de `División por Cero` en el cálculo de reportes de ocupación.
- **Causa Raíz:**
  - Fallos en validaciones de unicidad en cargas masivas CSV históricas crearon miles de registros duplicados en `Employees`, `Positions` y `Workplaces`.
  - Al ejecutar scripts de borrado masivo en Python sin validación de integridad referencial sincrónica, se eliminaron sedes "perdedoras" que aún tenían empleados asociados.
  - Se omitió recalcular manualmente métricas derivadas como `occupationAtWorkplace`, que normalmente gestionaba el código de aplicación.
- **Solución en Go:** Aplicar el principio **Zero-Trust**: antes de cualquier borrado (`client.DeleteMulti`), verificar sincrónicamente en Datastore que la entidad no tenga hijos huérfanos, encapsular en transacciones atómicas y recalcular triggers de negocio de forma determinista.

---

### 6. Desbordamiento del Límite de 1.500 Bytes en Propiedades Indexadas de Datastore (`InvalidArgument`)
- **Referencia Documental:** `ADR-20260916-02-motor-ingesta-taxonomias-go-cloud-run-jobs.md` (Catálogos O*NET y UNSPSC).
- **Manifestación del Error:** Al intentar insertar entidades con descripciones largas (más de 1,500 caracteres UTF-8), Datastore abortaba la transacción con error fatal `InvalidArgument: value is too long for indexed property`.
- **Causa Raíz:** Google Cloud Datastore impone un límite rígido de **1,500 bytes** para cualquier campo que tenga habilitada la indexación automática. Descripciones de cargos, tareas ocupacionales y textos de normativas superan fácilmente este tamaño.
- **Solución en Go:** Configurar de forma estricta tags `datastore:",noindex"` en los structs de Go, o usar `datastore.Property{NoIndex: true}` dinámico, separando campos descriptivos de los índices de búsqueda (`searchTokens []string`).

---

### 7. Desincronización de Tipos de Datos (String vs. Integer) entre Datastore y BigQuery
- **Referencia Documental:** `PM-20260630-deduplicacion-entidades.md`.
- **Manifestación del Error:** Scripts de reconciliación y queries analíticas reportaban falsos faltantes o retornaban conjuntos vacíos de datos.
- **Causa Raíz:** Datastore permite tipado flexible a nivel de entidad. Históricamente, identificadores como `companyId` o `identificationNumber` se guardaron como `String` en algunas colecciones y como `Integer` (int64) en otras. Al consultar en BigQuery, el motor infería `INTEGER`, provocando que filtros en Go o Python (`PropertyFilter("companyId", "=", companyID)`) fallaran silenciosamente por discrepancia de tipos.
- **Solución en Go:** Mappers bidireccionales con tipado explícito estricto y funciones de coerción determinista en `prax-go-common/mapper` que normalizan antes de interactuar con el SDK.

---

### 8. Fuga de Conexiones TCP y Descriptores de Red (CLOSE_WAIT / Sockets Exhaustion)
- **Referencia Documental:** `REF-202605-23-inventario-alertas-produccion.md`, Colapso Interno 5xx.
- **Manifestación del Error:** Microservicios en Cloud Run dejando de responder repentinamente tras atender varios miles de peticiones externas, arrojando errores `dial tcp: i/o timeout` o `too many open files`.
- **Causa Raíz:** Clientes HTTP que realizaban llamadas a microservicios downstream ejecutando `defer resp.Body.Close()` pero sin leer completamente los datos residuales del socket (`io.Copy(io.Discard, resp.Body)`), impidiendo que el transporte HTTP reutilizara la conexión TCP en el pool de Keep-Alive y dejando conexiones huérfanas en estado `CLOSE_WAIT`.
- **Solución en Go:** Configuración de Singletons de `*http.Client` con `http.Transport` sintonizado (`MaxIdleConns: 100`, `IdleConnTimeout: 90s`) y utilidades que garanticen el drenaje completo del body.

---

### 9. Caída del API Gateway por Timeouts en Cascada y Tráfico Bloqueante ("El Muro" 5xx > 10/s)
- **Referencia Documental:** `REF-202605-23-inventario-alertas-produccion.md` (Alerta P1: "El Muro").
- **Manifestación del Error:** Caída generalizada de la plataforma con respuestas `502 Bad Gateway` y `504 Gateway Timeout` emitidas por `prax-api-gateway`.
- **Causa Raíz:** Peticiones HTTP sincrónicas para operaciones batch que tardaban más de 60 segundos en responder. Al recibirse múltiples solicitudes simultáneas de diferentes empresas, los hilos de red del Gateway se agotaban esperando a los microservicios de backend, asfixiando el enrutador central.
- **Solución en Go:** Desacoplamiento asíncrono obligatorio: responder `HTTP 202 Accepted` en menos de 200 ms retornando un `taskId`, procesar el lote en segundo plano con goroutines/PubSub y exponer endpoints livianos de polling con caché.

---

### 10. Agotamiento de Memoria RAM (> 85% OOM Kill) en Cloud Run 512 MiB por Parseo No-Streaming
- **Referencia Documental:** `ADR-20260911-01-migracion-batch-employees-in-memory-go.md` y `REF-202605-23`.
- **Manifestación del Error:** Al subir planillas Excel de empleados superiores a 5,000 registros, el contenedor de Cloud Run era liquidado de inmediato por el kernel de Linux (`Exit Code 137 OOM`).
- **Causa Raíz:** Carga completa del árbol DOM del archivo `.xlsx` en memoria RAM utilizando bibliotecas convencionales de parseo, expandiendo un archivo comprimido de 10 MB a más de 300 MB en memoria y superando el límite de 512 MiB junto con el runtime.
- **Solución en Go:** Parseo por streaming fila por fila utilizando `excelize.Rows()`, configuración obligatoria de `GOMEMLIMIT=460MiB` en Cloud Run para detonar GC agresivo previo al límite de cgroup, y reciclaje de structs con `sync.Pool`.

---

# 2. 20 Problemas Prácticos de Troubleshooting en Go (T01 a T20)

---

### Bloque 1: Interoperabilidad & Gobernanza de Almacenamiento (T01 a T05)

#### T01. Decodificador y Codificador Bidireccional de Claves WebSafe Multiversión
- **Escenario:** Un microservicio en Go debe recibir claves primarias generadas por un sistema legacy Java (`agtz...`) y por sistemas modernos (`Eh...`), resolver si pertenecen a la partición `s~prax-prod` o `k~prax-dev` e interactuar con Datastore sin provocar `ParseException`.
- **Reto de Implementación:**
  - Escribir un paquete `keysafe` con funciones `DecodeAnyKey(raw string) (*datastore.Key, string, error)` y `EncodeLegacy(appID, namespace string, key *datastore.Key) (string, error)`.
  - Desempaquetar el binario Protobuf v3 inspeccionando el Tag 13 para extraer el `appID`.
  - Crear pruebas unitarias con tablas comparando claves reales de Java y Go.

#### T02. Guardián de Tamaño de Propiedades en Datastore (1.500 Bytes NoIndex Guard)
- **Escenario:** En la carga masiva de catálogos O*NET y descripciones de puestos de trabajo, textos de más de 1,500 bytes detonaban caídas `InvalidArgument: value is too long for indexed property`.
- **Reto de Implementación:**
  - Implementar la interfaz `datastore.PropertyLoadSaver` en un struct genérico `CatalogEntity`.
  - En el método `Save()`, inspeccionar el tamaño en bytes de cada campo string; si supera 1,500 bytes UTF-8, forzar dinámicamente `Property{NoIndex: true}`.
  - Asegurar que los campos cortos (`code`, `name`) permanezcan indexados para permitir búsquedas.

#### T03. Normalizador de Tipos Débiles NoSQL a Tipos Fuertes Go (String vs. Int64)
- **Escenario:** Al leer entidades históricas de Datastore, el campo `companyId` aparece en algunos registros como entero `int64` y en otros como cadena `string`, provocando errores de deserialización al mapear directamente a un struct.
- **Reto de Implementación:**
  - Crear un tipo personalizado `FlexID int64` que implemente las interfaces `json.Unmarshaler` y `datastore.PropertyLoadSaver`.
  - Soportar la lectura transparente tanto si el valor en la base de datos es `int64` como si es `string` numérico (`"12345"`).

#### T04. Paginador Resiliente de Datastore con Resguardo de Cursor y Backoff
- **Escenario:** Extraer 100,000 registros de Datastore mediante un worker en Go. Si ocurre un fallo de red a mitad del proceso, el worker debe reanudar desde el último cursor guardado sin reiniciar desde cero.
- **Reto de Implementación:**
  - Diseñar una función `StreamEntities(ctx context.Context, kind string, batchSize int, fn func([]Entity) error) error`.
  - Extraer el `datastore.Cursor`, persistirlo tras cada batch exitoso y aplicar retroceso exponencial si Datastore retorna un error transitorio `Unavailable`.

#### T05. Transacción Atómica con Aislamiento y Detección de Conflictos Concurrente
- **Escenario:** Dos peticiones simultáneas intentan actualizar el inventario de empleados activos de una sede en Datastore.
- **Reto de Implementación:**
  - Implementar una función `UpdateWorkplaceCount(ctx context.Context, client *datastore.Client, wpKey *datastore.Key, delta int) error` utilizando `client.RunInTransaction`.
  - Simular colisión concurrente con múltiples goroutines y verificar que el runtime reintente automáticamente ante errores de contención transaccional (`concurrent transaction conflict`).

---

### Bloque 2: Resiliencia de Pipelines de Datos & Concurrencia (T06 a T10)

#### T06. Controlador de "Punto Dulce" para BigQuery Publisher (Concurrencia Estricta 15)
- **Escenario:** Proteger a BigQuery de saturación de consultas DML evitando errores de cuota (`too many concurrent queries`) y previniendo errores HTTP 429 hacia Pub/Sub.
- **Reto de Implementación:**
  - Construir un servicio HTTP en Go que reciba eventos Push de Pub/Sub.
  - Utilizar un canal semáforo (`chan struct{}`) con capacidad exacta de 15 tokens para gobernar las ejecuciones simultáneas hacia el cliente de BigQuery.
  - Si el semáforo está lleno, encolar con timeout de 5 segundos; si expira, retornar HTTP 429 con cabecera `Retry-After: 10` para que Pub/Sub asuma el backpressure.

#### T07. Pipeline Anti-Poison Pill con Detección de Streaming Buffer Lock en BigQuery
- **Escenario:** Capturar el error 500 cuando se intenta hacer `UPDATE` a una fila recién insertada en el Streaming Buffer de BigQuery y evitar el bucle infinito de reintentos en Pub/Sub.
- **Reto de Implementación:**
  - Interceptar el error del SDK `googleapi.Error` verificando si el mensaje contiene `would affect rows in the streaming buffer`.
  - Si coincide, hacer ACK inmediato ante Pub/Sub para destapar la tubería principal y publicar el evento en un tópico secundario de retraso (Delay Topic / DLQ) con atributo de tiempo programado de 45 minutos.

#### T08. Deduplicador en Memoria con Ventana Deslizante y Hashing FNV/SHA-256
- **Escenario:** Procesar un flujo masivo de eventos Pub/Sub donde los reintentos de red generan eventos duplicados dentro de una ventana de 60 segundos.
- **Reto de Implementación:**
  - Construir un struct `Deduplicator` seguro para hilos que almacene hashes de eventos recientes.
  - Limpiar automáticamente las claves con antigüedad superior a 60 segundos utilizando una goroutine de barrido y un `sync.RWMutex`.
  - Descartar eventos duplicados con costo computacional $O(1)$ sin consultar bases de datos externas.

#### T09. Dispatcher Asíncrono con HTTP 202 y Polling en Caché (Anti-Muro 502/504)
- **Escenario:** Eliminar los bloqueos del API Gateway provocados por peticiones pesadas de empleados.
- **Reto de Implementación:**
  - Endpoint `POST /api/v1/batch/employees` que valida el header, genera un `taskID`, lanza el procesamiento pesado en una goroutine desacoplada y retorna inmediatamente `HTTP 202 Accepted` con payload `{"taskId": "..."}` en menos de 100 ms.
  - Endpoint `GET /api/v1/batch/status/{taskId}` que consulta el estado desde una caché en memoria rápida (`sync.Map`) con TTL de 5 segundos.

#### T10. Worker Pool con Límite de Memoria Adaptativo y Pausa por Presión de GC
- **Escenario:** Un pipeline de procesamiento de archivos que pausa la ingestión de nuevos registros si la memoria del contenedor se acerca al límite de Cloud Run.
- **Reto de Implementación:**
  - Monitorear `runtime.ReadMemStats` periódicamente.
  - Si `Alloc` supera los 400 MiB (en un contenedor de 512 MiB), el despachador de workers debe frenar el consumo del canal de entrada y forzar `runtime.GC()` hasta que la memoria retorne a niveles seguros.

---

### Bloque 3: Integridad de Datos, Deduplicación & Zero-Trust (T11 a T15)

#### T11. Limpiador Zero-Trust de Entidades Huérfanas en Datastore
- **Escenario:** Replicar con seguridad la limpieza de sedes y cargos sin dejar empleados desasociados.
- **Reto de Implementación:**
  - Escribir una función en Go que reciba un slice de IDs de sedes a eliminar.
  - Antes de ejecutar `DeleteMulti`, realizar una consulta de validación: `SELECT __key__ FROM Employee WHERE workplaceKey = @targetKey LIMIT 1`.
  - Si se detecta al menos un empleado asignado, abortar la eliminación de esa sede específica y reportar la violación de integridad referencial sin interrumpir las demás.

#### T12. Recalculador Determinista de Triggers de Negocio (`occupationAtWorkplace`)
- **Escenario:** Al importar masivamente empleados saltándose los métodos tradicionales de la aplicación, el porcentaje de ocupación de las sedes queda desactualizado.
- **Reto de Implementación:**
  - Construir una rutina en Go que consulte todas las sedes de una empresa y ejecute una consulta de agregación `NewAggregationQuery().WithCount()` sobre los empleados activos de cada sede.
  - Actualizar el campo `occupationAtWorkplace` en la entidad `Workplace` dentro de una transacción atómica, asegurando manejo de casos de capacidad cero (previniendo divisiones por cero).

#### T13. Sincronizador de BigQuery Append-Only con Vistas de Deduplicación
- **Escenario:** Reemplazar las operaciones `UPDATE` destructivas en BigQuery por inserciones inmutables compatibles con Streaming API.
- **Reto de Implementación:**
  - Escribir un publicador en Go que inserte registros con columnas `entity_id`, `payload_json`, `operation_type` (`INSERT`/`UPDATE`/`DELETE`) y `timestamp_utc`.
  - Generar el script DDL en Go para crear una vista en BigQuery que materialice el estado actual de cada entidad utilizando `ROW_NUMBER() OVER(PARTITION BY entity_id ORDER BY timestamp_utc DESC)`.

#### T14. Detector y Fusionador de Empleados Duplicados por Huella Digital
- **Escenario:** Identificar empleados duplicados creados con variaciones en su tipo o número de documento (ej. espacios en blanco o ceros a la izquierda).
- **Reto de Implementación:**
  - Normalizar números de documento (eliminación de caracteres no numéricos) y correos electrónicos (minúsculas y trim).
  - Generar un hash determinista SHA-256 de la identidad normalizada.
  - Construir un script en Go que agrupe registros por hash y genere un plan de fusión (*Merge Plan*) preservando la fecha de inicio más antigua (`workplaceStartDate`) sin sobreescribir datos históricos.

#### T15. Validador Estricto de Esquemas YAML para Cargas de Catálogos Masivos
- **Escenario:** Validar que los archivos de taxonomías (O*NET / UNSPSC) cumplan los estándares de nomenclatura de PRAX (PascalCase para Kinds, camelCase para propiedades) antes de tocar la base de datos.
- **Reto de Implementación:**
  - Construir un validador en Go que lea el archivo de mapeo y compruebe por reflexión que ninguna propiedad indexada exceda 1,500 bytes y que los nombres de campos sigan estrictamente las convenciones de la organización.

---

### Bloque 4: Redes, SRE & Optimización de Recursos en Cloud Run (T16 a T20)

#### T16. Cliente HTTP Singleton con Drenaje de Sockets Anti-CLOSE_WAIT
- **Escenario:** Eliminar la fuga de descriptores de sockets al invocar microservicios internos.
- **Reto de Implementación:**
  - Implementar un wrapper `SafeHTTPClient` en Go con `http.Transport` configurado (`MaxIdleConnsPerHost: 20`, `IdleConnTimeout: 90s`).
  - Crear una función helper `DoAndDrain(req *http.Request) ([]byte, int, error)` que garantice la lectura completa con `io.Copy(io.Discard, resp.Body)` y el cierre en bloque `defer resp.Body.Close()`.
  - Crear una prueba con 1,000 solicitudes concurrentes verificando que no queden sockets residuales.

#### T17. Sintonizador Dinámico de Memoria en Cloud Run (`GOMEMLIMIT` y Cgroups v2)
- **Escenario:** Evitar que un contenedor Go en Cloud Run de 512 MiB muera por OOM Kill ante ráfagas repentinas de asignación de memoria.
- **Reto de Implementación:**
  - Escribir una función de inicio en `main.go` que lea el límite de memoria del contenedor directamente desde `/sys/fs/cgroup/memory.max` (cgroups v2).
  - Si el archivo existe, configurar automáticamente la variable de depuración del runtime `debug.SetMemoryLimit` al 85% de dicho valor (ej. 435 MiB para un contenedor de 512 MiB), permitiendo que el GC se active de forma agresiva antes de que el kernel envíe `SIGKILL`.

#### T18. Health Probes Inteligentes con Diagnóstico de Descriptores de Red (`/healthz`)
- **Escenario:** Informar al balanceador de Cloud Run que una instancia no debe recibir más tráfico si sus conexiones internas a Datastore o Pub/Sub están degradadas o si está sufriendo fuga de descriptores.
- **Reto de Implementación:**
  - Endpoint `/healthz` (Liveness) que verifique que el proceso esté vivo.
  - Endpoint `/readyz` (Readiness) que inspeccione la cantidad de descriptores de archivo abiertos en `/proc/self/fd` (en Linux); si supera el 80% del límite `ulimit`, retornar `HTTP 503 Service Unavailable` para que Cloud Run redirija el tráfico a otra instancia sana.

#### T19. Parseador en Streaming de Excel con Huella de Memoria < 30 MiB
- **Escenario:** Parsear archivos Excel de 50,000 filas de empleados sin superar el límite de RAM de Cloud Run.
- **Reto de Implementación:**
  - Consumir el archivo desde un `io.Reader` conectado a Google Cloud Storage.
  - Utilizar el iterador `excelize.Rows()` procesando fila por fila sin volcar el libro completo a memoria.
  - Validar los datos y despacharlos a un canal de lotes (`chan []Employee`, tamaño 400), liberando la memoria de cada fila inmediatamente tras su lectura.

#### T20. Reconciliador Continuo Datastore vs. BigQuery con Detección de Desincronización
- **Escenario:** Detectar en segundo plano registros que existen en Datastore pero no se reflejaron en BigQuery debido a fallos en Pub/Sub.
- **Reto de Implementación:**
  - Un Cloud Run Job en Go que tome el `companyKey` como argumento.
  - Consulte las claves de empleados en Datastore mediante una proyección de sólo claves (`query.KeysOnly()`).
  - Ejecute una consulta en BigQuery `SELECT id FROM EVALUATIONDONE.employees WHERE company_id = @companyKey`.
  - Calcule la diferencia simétrica entre ambos conjuntos y emita eventos de reparación automática hacia el tópico de Pub/Sub para los registros desincronizados.
