# Go Zero to Cloud-Native Mastery: Production & Enterprise Edition

Hoja de ruta, repositorio de práctica deliberada y manual de ingeniería avanzada para dominar Go de manera idiomática, transicionando desde paradigmas orientados a objetos tradicionales (Java/Spring Boot) y asíncronos dinámicos (JavaScript/Node.js) hacia sistemas concurrentes de alto rendimiento, microservicios serverless en GCP y herramientas de infraestructura cloud-native seguras.

---

## 🎯 Objetivos de Aprendizaje

- **Composición sobre Herencia**: Desaprender jerarquías rígidas y clases abstractas de Java; dominar composición de structs e interfaces implícitas desacopladas.
- **Concurrencia Pragmática**: Superar el Event Loop mono-hilo de Node.js y los threads pesados del OS en la JVM dominando el **scheduler M:N de Go, Goroutines, Canales y sincronización atómica**.
- **Errores como Valores**: Erradicar el flujo de control oculto de excepciones `try/catch/finally`; gestionar errores deterministas con `if err != nil`, wrapping idiomático (`%w`), `errors.Is` y `errors.As`.
- **Mecánica de Memoria & Runtime**: Entender el layout de memoria, padding, alineación de structs, semántica de valor vs. puntero, escape analysis (stack vs. heap) y sintonización del GC con `GOMEMLIMIT`.
- **Ecosistema Cloud-Native Real**: Desarrollar servicios de streaming en memoria, consumers escalables de Cloud Pub/Sub, repositorios transaccionales en Cloud Datastore/Firestore y pipelines de BigQuery optimizados para Cloud Run.
- **Seguridad Defensiva (OWASP Top 10)**: Diseñar software resiliente a vulnerabilidades de inyección, SSRF, criptografía débil, Zip Bombs y ataques de denegación de servicio en memoria.

---

## 📂 Estructura del Repositorio

```text
.
├── 01-basics/            # Fundamentos, tipos, layout de memoria y errores explícitos
├── 02-intermediate/      # Concurrencia básica, interfaces, composición y testing
├── 03-advanced/          # Sync profunda, context, I/O streaming y profiling con pprof
├── 04-expert/            # Sistemas distribuidos, zero-allocation y serialización binaria
├── 05-gcp-cloud-native/  # Microservicios Cloud Run, SDK de GCP, BigQuery y Datastore
├── 06-cybersecurity/     # Criptografía moderna, auth, mTLS y estándares OWASP Top 10
├── 07-prax-challenges/   # 20 Retos reales basados en la arquitectura de producción de PRAX
├── 08-interview-prep/    # 30 Problemas y preguntas técnicas de nivel Senior / Staff
├── 09-ai-bigdata-security/ # 10 Proyectos Weekend: IA + Go + Big Data + Ciberseguridad
├── docs/                 # Documentación técnica profunda, incidentes y troubleshooting
│   └── PRAX_TROUBLESHOOTING.md # 10 Errores críticos de PRAX y 20 retos de depuración en Go
└── Makefile              # Automatización de builds, tests con -race, linters y auditorías
```

---

## 🛠️ Herramientas y Convenciones de Calidad de Código

Para mantener estándares de ingeniería de software de primer nivel durante todos los ejercicios:

```bash
# 1. Detección estricta de condiciones de carrera en tests concurrentes
go test -race -v ./...

# 2. Análisis estático estricto y linters de producción
golangci-lint run

# 3. Auditoría de seguridad estática (SAST) y vulnerabilidades en dependencias (SCA)
gosec ./...
govulncheck ./...

# 4. Formato oficial y optimización de imports
go fmt ./...
goimports -w .

# 5. Inspección de Escape Analysis (ver qué variables escapan al heap)
go build -gcflags="-m -m" ./...

# 6. Profiling de CPU y memoria interactivo
go test -bench=. -benchmem -cpuprofile=cpu.pprof -memprofile=mem.pprof
go tool pprof -http=:8081 mem.pprof
```

---

# 📚 Checklist de Proyectos, Casos de Estudio y Ejercicios

---

## Nivel 1: Fundamentos & Go Idiomático (Básico)
> **Foco:** Control de flujo, structs, slices vs. arrays, punteros y gestión de errores explícita.

### 💡 Caso de Estudio: La Trampa de las Excepciones de Java vs. Errores como Valores en Go
> **Contexto de la Industria:** En sistemas bancarios y de pagos basados en Java/Spring Boot (como los primeros microservicios de Monzo), una excepción no capturada de tipo `NullPointerException` o `DataAccessException` en un hilo de ejecución secundario provocaba respuestas HTTP 500 silenciosas o inconsistencias de estado financiero.
> 
> **La Solución en Go:** Go no posee excepciones ni `try/catch`. Tratar los errores como valores devueltos en la tupla `(resultado, error)` fuerza al desarrollador a tomar una decisión explícita en cada punto de falla. El uso de `fmt.Errorf("operación fallida: %w", err)` preserva la cadena de causalidad, permitiendo a los controladores HTTP inspeccionar la causa raíz con `errors.Is` sin romper el encapsulamiento ni asumir suposiciones en tiempo de ejecución.

- [ ] **01. Word & Char Frequency Counter**
  - Procesar texto crudo usando `bufio.Scanner`.
  - Diferenciar entre `byte` (ASCII/raw) y `rune` (UTF-8 code points).
  - Almacenar resultados en un `map[rune]int` y ordenar llaves alfabéticamente.
- [ ] **02. Two Sum & Subarray Sum (LeetCode 1 & 560)**
  - Resolver en tiempo $O(n)$ usando maps para búsquedas en tiempo constante.
  - Practicar asignación y resizing de slices con `make([]T, len, cap)` para evitar realocaciones innecesarias.
- [ ] **03. In-Memory Task Manager CLI**
  - Implementar operaciones CRUD con structs y punteros receptores.
  - Serializar y deserializar el estado hacia/desde JSON usando `encoding/json`.
- [ ] **04. Custom Error Wrapping Engine**
  - Diseñar tipos de error de dominio personalizados implementando el método `Error() string`.
  - Encapsular errores con `%w` y evaluarlos jerárquicamente con `errors.Is` y `errors.As`.
- [ ] **05. Matrix Transpose & Zero Matrix**
  - Manejo de slices bidimensionales (`[][]int`).
  - Analizar el comportamiento de punteros y encabezados de slices al mutar datos dentro de funciones.
- [ ] **06. Temperature & Unit Converter (Domain Types)**
  - Declarar tipos de dominio (`type Celsius float64`, `type Fahrenheit float64`).
  - Adjuntar métodos por valor y métodos por puntero a estos tipos escalares.
- [ ] **07. Binary Search & Slices Truncation**
  - Implementar búsqueda binaria iterativa sin recurrir al paquete `sort`.
  - Dominar las expresiones de corte completas `slice[low:high:max]` para restringir capacidad.
- [ ] **08. Fibonacci: Benchmarks de Memoria**
  - Implementar versiones: recursiva pura, memoizada con map y con slice preasignado.
  - Ejecutar `go test -bench=. -benchmem` para comparar asignaciones en heap (`allocs/op`).
- [ ] **09. Simple .env & Flags Config Parser**
  - Parsear variables de entorno y argumentos de línea de comandos (`os.Args` y paquete `flag`).
  - Mapear configuraciones a un struct de configuración usando struct tags personalizados.
- [ ] **10. LRU Cache (Single-Threaded)**
  - Implementar una lista doblemente enlazada con nodos interconectados y un map interno.
  - Abstraer la estructura mediante interfaces desacopladas sin recurrir a paquetes externos.

---

## Nivel 2: Composición, Interfaces & Concurrencia Inicial (Medio)
> **Foco:** Duck typing estático, goroutines, canales con/sin buffer, select y table-driven tests.

### 💡 Caso de Estudio: Por qué Discord abandonó el Pool de Threads de Java/Python por Goroutines
> **Contexto de la Industria:** Discord manejaba millones de conexiones WebSocket de voz y mensajería simultáneas. En un modelo tradicional multihilo (1 hilo del SO = 1 conexión), la memoria base por hilo (1 a 2 MiB en la JVM o Linux) requería cientos de gigabytes solo en stacks de memoria y producía un costo prohibitivo de context switching en el kernel.
> 
> **La Solución en Go:** Go utiliza Goroutines cuyo stack inicial es de tan solo **2 KiB**, creciendo y encogiéndose dinámicamente en el heap según sea necesario. Gracias al scheduler M:N de Go en el userspace, un solo proceso puede gestionar fácilmente 250,000 conexiones concurrentes consumiendo menos de 1 GiB de RAM y con cambios de contexto que ocurren en decenas de nanosegundos.

- [ ] **11. Generic Stack & Queue con Constraints**
  - Implementar estructuras de datos usando Go Generics (`[T any]`, `comparable`).
  - Escribir tests dirigidos por tablas (*table-driven tests*) testeando casos límite.
- [ ] **12. Polimorfismo con Interfaces Pequeñas**
  - Definir interfaces de un solo método (`io.Reader`, `io.Writer`, `io.Closer`).
  - Implementar un sistema de almacenamiento intercambiable (memoria vs. disco local) sin interfaces infladas estilo Java Spring.
- [ ] **13. Ping-Pong Channel Synchronization**
  - Coordinar dos goroutines que se intercambian un contador usando canales no amortiguados (`chan int`).
  - Detectar y prevenir deadlocks en tiempo de compilación y ejecución.
- [ ] **14. Worker Pool Pattern**
  - Crear un despachador de tareas con canales de trabajos (`jobs chan`) y resultados (`results chan`).
  - Procesar cargas de trabajo con un pool fijo de $N$ workers concurrentes sin fuga de goroutines.
- [ ] **15. Fan-In / Fan-Out Data Pipeline**
  - Múltiples goroutines leen de un solo canal (Fan-Out) y consolidan resultados en un canal único (Fan-In).
  - Cerrar canales de manera segura desde el productor sin generar *panics* por envío a canal cerrado.
- [ ] **16. Timeout con select y time.After**
  - Consumir de un canal con fallback de timeout cancelable.
  - Identificar y corregir fugas de memoria generadas por tickers huérfanos que no se detienen (`ticker.Stop()`).
- [ ] **17. Token Bucket Rate Limiter**
  - Implementar un limitador de tasa de peticiones usando `time.Ticker` y canales amortiguados.
  - Limitar el paso de peticiones concurrentes a una tasa fija por segundo con soporte para ráfagas.
- [ ] **18. Concurrent File Downloader/Reader**
  - Dividir un archivo o stream en chunks de bytes procesados concurrentemente.
  - Reensamblar las partes en el orden estricto original usando `sync.WaitGroup`.
- [ ] **19. Dynamic Pub/Sub Broker en Memoria**
  - Mecanismo de suscripción y publicación a tópicos basados en strings.
  - Manejo de suscriptores lentos evitando bloqueos globales mediante canales con buffer o descarte de mensajes.
- [ ] **20. Validador de Datos con Reflection (reflect)**
  - Inspeccionar campos y struct tags en tiempo de ejecución.
  - Construir un motor de validación de campos obligatorios (`validate:"required"`).

---

## Nivel 3: Concurrencia Profunda, Sincronización & I/O (Avanzado)
> **Foco:** Paquete `sync`, `atomic`, ciclo de vida con `context`, prevención de data races y networking.

### 💡 Caso de Estudio: Cómo Uber eliminó cuellos de botella de Garbage Collection con `sync.Pool`
> **Contexto de la Industria:** Uber procesa miles de millones de eventos geográficos y peticiones HTTP en tiempo real para el matching de conductores y pasajeros. En sus servicios Go de alta concurrencia, la constante asignación y liberación de buffers JSON y slices de bytes para cada request generaba una enorme presión en el recolector de basura (GC), disparando las pausas *Stop-The-World* y la latencia p99.
> 
> **La Solución en Go:** Implementaron pools de objetos reutilizables con `sync.Pool`. Al reciclar buffers (`bytes.Buffer`) entre peticiones, redujeron la tasa de asignación de memoria en el heap en más del **70%**, aplanando la latencia p99 y estabilizando el consumo de CPU.

- [ ] **21. Thread-Safe Cache con Sharded Locks**
  - Implementar una caché concurrente evitando un único `sync.RWMutex` global.
  - Minimizar la contención particionando las llaves en shards hash independientes (ej. 32 shards).
- [ ] **22. Context Cancellation & Deadline Propagation**
  - Encadenar árboles de cancelación usando `context.WithCancel`, `WithTimeout` y `WithValue`.
  - Garantizar la limpieza y terminación de goroutines subordinadas cuando el contexto padre aborta la ejecución.
- [ ] **23. Circuit Breaker Pattern**
  - Modelar la máquina de estados: `Closed`, `Open`, `Half-Open`.
  - Implementar transiciones de estado libres de locks bajo alta concurrencia usando `sync/atomic`.
- [ ] **24. Graceful Shutdown HTTP Server**
  - Capturar señales del sistema operativo (`SIGINT`, `SIGTERM`).
  - Drenar conexiones HTTP activas con `server.Shutdown(ctx)` garantizando que no se corten requests en vuelo.
- [ ] **25. Memory Reuse con sync.Pool**
  - Minimizar la presión sobre el recolector de basura (GC) reutilizando buffers de bytes.
  - Medir `allocs/op` antes y después de la optimización con `testing.B`.
- [ ] **26. TCP Reverse Proxy con Streaming Bidireccional**
  - Aceptar conexiones TCP crudas y redireccionar tráfico usando `io.Copy`.
  - Gestionar el cierre ordenado de conexiones en ambos sentidos sin leaks de descriptores de sockets.
- [ ] **27. Singleflight Group (Request Coalescing)**
  - Implementar un mecanismo para erradicar el problema del *Cache Stampede* (Thundering Herd).
  - Si 100 llamadas simultáneas solicitan la misma llave inexistente en caché, solo 1 ejecuta la consulta al backend y las 99 restantes esperan su resultado.
- [ ] **28. Data Race Hunting con -race y pprof**
  - Escribir un test unitario con condiciones de carrera deliberadas sobre una variable compartida.
  - Diagnosticar y resolver la contención usando `go test -race` y perfiles de contención de mutex con `pprof`.

---

## Nivel 4: Sistemas Distribuidos & Low-Level (Experto)
> **Foco:** Zero-allocation, serialización binaria, consenso distribuido y motores de persistencia.

### 💡 Caso de Estudio: VictoriaMetrics superando a Prometheus con Zero-Alloc Parsing
> **Contexto de la Industria:** Prometheus (desarrollado en Go) tradicionalmente sufría de alto consumo de memoria RAM al ingerir millones de métricas por segundo procedentes de clústeres de Kubernetes, debido a la creación constante de pequeños objetos en memoria por cada muestra de serie temporal.
> 
> **La Solución en Go:** VictoriaMetrics reescribió los componentes críticos de ingesta en Go aplicando estrictamente técnicas de *Zero-Allocation*: reutilización masiva de buffers de bytes, decodificación binaria directa sin conversiones intermedias a string, y empaquetado compacto de estructuras en memoria. El resultado fue un consumo de RAM hasta **7 veces menor** y una velocidad de ingesta 4 veces superior.

- [ ] **29. Motor Key-Value Persistente tipo Bitcask (LSM Simplificado)**
  - Escrituras secuenciales rápidas en disco mediante append-only log.
  - Indexación en memoria basada en tablas hash de offsets a disco y proceso de compactación (*merge*).
- [ ] **30. Implementación de Raft (Leader Election & Heartbeats)**
  - Implementar estados de nodos distribuidos (`Leader`, `Follower`, `Candidate`).
  - Enviar RPCs simulados y manejar timeouts aleatorios con canales para reelección de líderes.
- [ ] **31. Ring Buffer Circular Lock-Free**
  - Cola de mensajes de alto rendimiento basada en operaciones atómicas (`atomic.CompareAndSwap` / `atomic.Pointer`).
- [ ] **32. Parser de Protocolo Binario Personalizado**
  - Codificar y decodificar frames binarios directamente desde streams (`io.Reader`/`io.Writer`) usando `encoding/binary`.
- [ ] **33. Microservicio gRPC con Interceptores**
  - Definir esquemas en `.proto` y compilar stubs con `protoc-gen-go` y `protoc-gen-go-grpc`.
  - Crear interceptores unarios y de streaming para autenticación basada en metadatos, distributed tracing y métricas.

---

## Nivel 5: GCP, Big Data & Cloud Infra
> **Foco:** Integración idiomática con servicios de Google Cloud, procesamiento masivo y utilidades cloud-native.

### 💡 Caso de Estudio: Monzo Bank operando 2,500 Microservicios en Go sobre Google Cloud
> **Contexto de la Industria:** Monzo Bank construyó uno de los mayores bancos digitales del Reino Unido operando más de 2,500 microservicios en Go sobre Google Cloud Platform y Kubernetes. Necesitaban despliegues ultra rápidos, aislamiento de fallos, y consumo predecible de memoria sin las sorpresas de calentamiento de la JVM.
> 
> **La Solución en Go:** Adoptaron Go como lenguaje estándar único para backend: clientes HTTP y gRPC ligeros, propagación rigurosa de contextos (`context.Context`) a través de llamadas RPC, y consumo de BigQuery y Datastore sin capas ORM pesadas, manteniendo latencias de microservicios inferiores a 10 ms.

- [ ] **34. CLI de Gestión de Infraestructura GCP**
  - Construir una herramienta CLI modular usando `cobra`.
  - Autenticarse de forma transparente mediante Application Default Credentials (ADC) e interactuar con el SDK oficial de GCP.
- [ ] **35. Pipeline de Ingestión en Lote para BigQuery**
  - Leer streams de eventos masivos y agruparlos en memoria por tamaño (ej. 500 filas) o por ventana de tiempo.
  - Insertar registros utilizando la API de streaming de BigQuery con deduplicación garantizada mediante `InsertID`.
- [ ] **36. Repositorio de Persistencia en Google Cloud Datastore / Firestore**
  - Diseñar un patrón repositorio idiomático que desacople la lógica de negocio de los tipos nativos del SDK de Datastore.
  - Implementar transacciones multi-entidad atómicas y consultas con paginación basada en cursores (`datastore.Cursor`).
- [ ] **37. Worker Serverless en Google Cloud Run**
  - Microservicio optimizado para Cold Starts mínimos (< 250 ms).
  - Configurar límites estrictos de memoria (`GOMEMLIMIT`) y evaluar rendimiento frente a concurrencia variable.
- [ ] **38. Consumer Escalable de Cloud Pub/Sub**
  - Procesar mensajes de suscripciones Pub/Sub con confirmación explícita (`msg.Ack()`, `msg.Nack()`).
  - Controlar la tasa máxima de procesamiento concurrente por instancia para evitar saturar bases de datos aguas abajo.
- [ ] **39. Exportador de Métricas OpenTelemetry a Google Cloud Monitoring**
  - Instrumentar trazas distribuidas y métricas de latencia p95/p99 en endpoints HTTP.
  - Exportar telemetría hacia Cloud Trace y Cloud Monitoring de forma asíncrona sin bloquear la ruta crítica.

---

## Nivel 6: Ciberseguridad & DevSecOps (OWASP Top 10 en Go)
> **Foco:** Criptografía defensiva, control de identidades, prevención de vulnerabilidades web y hardening.

### 💡 Caso de Estudio: Cloudflare y la Mitigación de Vulnerabilidades de Memoria y TLS
> **Contexto de la Industria:** Cloudflare maneja cerca del 20% del tráfico web mundial. Históricamente, las vulnerabilidades críticas en parsers escritos en C/C++ (como buffer overflows o use-after-free tipo Heartbleed) representaban una amenaza constante a la seguridad global.
> 
> **La Solución en Go:** Reescribieron gran parte de sus servicios de borde, resolución DNS y herramientas de seguridad en Go. La seguridad de memoria nativa de Go (type safety, memory safety, ausencia de aritmética de punteros descontrolada) eliminó familias enteras de vulnerabilidades críticas de ejecución remota de código, permitiendo implementar pilas criptográficas modernas (`crypto/tls`, ChaCha20-Poly1305, Ed25519) con total auditabilidad.

- [ ] **40. Secret Injector con Google Secret Manager**
  - Cargar y refrescar credenciales de forma asíncrona desde Secret Manager en tiempo de ejecución.
  - Gestionar la rotación de llaves criptográficas sin reiniciar el contenedor ni interrumpir peticiones en curso.
- [ ] **41. Token Engine JWT Seguro (RS256/Ed25519)**
  - Generación, firma y validación estricta de tokens JWT.
  - Rechazo de algoritmos inseguros (`none`), verificación de claims (`exp`, `iss`, `aud`) y extracción hacia `context.Context`.
- [ ] **42. Mutual TLS (mTLS) Inter-Service Communication**
  - Configurar certificados TLS mutuos de cliente y servidor programáticamente usando `crypto/tls`.
  - Forzar verificación estricta de certificados y protocolos mínimos (`tls.VersionTLS13`).
- [ ] **43. Concurrent Port & Security Header Scanner**
  - Escáner de red concurrente con control estricto de file descriptors abiertos.
  - Protección activa contra SSRF validando y resolviendo rangos de IPs privadas y metadatos de nube antes de conectar.
- [ ] **44. Cifrado de Datos en Reposo con AES-GCM & HMAC**
  - Cifrar campos de PII sensibles (salarios, documentos de identidad) usando `crypto/aes` y `crypto/cipher` en modo GCM.
  - Validación de integridad y autenticidad antes de persistir en bases de datos NoSQL.
- [ ] **45. Hardened HTTP Reverse Proxy / WAF Ligero**
  - Middleware de inspección de payloads con límites de tamaño (`http.MaxBytesReader`).
  - Sanitización de headers, control estricto de CORS y mitigación de ataques Slowloris mediante timeouts estrictos.

---

# 🏢 Nivel 7: Prax Cloud-Native Challenges (20 Retos Reales del Ecosistema PRAX)
> **Foco:** Casos de ingeniería extraídos de la migración real del ecosistema PRAX (Java/Spring Boot/Python hacia Go/GCP).
> **Contexto de Arquitectura:** El ecosistema PRAX administra gestión de talento humano, evaluaciones de clima y cultura organizacional, cargas masivas de nóminas en Excel, persistencia en Google Cloud Datastore, sincronización de eventos vía Cloud Pub/Sub y analítica en BigQuery sobre Cloud Run con restricciones estrictas de recursos (512 MiB de RAM).

### 💡 Caso de Estudio PRAX: La Migración de `prax-employees-service` (Java) a `prax-go-employees-ms` (Go)
> **El Problema:** La arquitectura legada en Java utilizaba un job en Cloud Run que procesaba planillas de hasta 2,500 empleados. Para realizar validaciones cruzadas de datos, insertaba cada registro en Datastore como una entidad temporal (`TempEmployee`), acumulando más de **45,000 operaciones de escritura y lectura** por carga. Esto disparaba los costos en GCP, requería contenedores de más de 2 GiB de RAM por el overhead de la JVM, y tardaba minutos en completarse con frecuentes problemas de *cold start*.
> 
> **La Solución en Go:** Se rediseñó el servicio como `prax-go-employees-ms` en Go:
> 1. Se eliminó por completo el antipatrón `TempEmployee`: la planilla de Excel/TSV se transmite desde Cloud Storage y se valida enteramente **en memoria por streaming** utilizando tan solo ~2.5 MiB de RAM.
> 2. Se ajustó el contenedor de Cloud Run a **512 MiB** (`GOMEMLIMIT=460MiB`, concurrencia 10, sin CPU throttling).
> 3. Los tiempos de *cold start* cayeron de 15 segundos a menos de **250 milisegundos**.
> 4. Se implementó sincronización determinista de eventos hacia BigQuery con Pub/Sub mediante `sync.WaitGroup` y `topic.Stop()`, erradicando los `Thread.sleep` arbitrarios de Java.
> 5. Mapeo transparente en `prax-api-gateway` para mantener compatibilidad 100% con `prax-ui` sin alterar el frontend.

---

### Los 20 Retos de Producción PRAX:

- [ ] **P01. In-Memory Streaming Excel Parser (`excelize`) en Cloud Run 512 MiB**
  - Parsear archivos `.xlsx` y `.tsv` de 2,500 filas transmitidos desde un `io.Reader` de Cloud Storage.
  - Utilizar el cursor iterativo `rows.Next()` de `excelize` para mantener el uso de memoria por debajo de 5 MiB sin volcar todo el archivo en el heap.
- [ ] **P02. Generador de Reporte Excel con Marcado de Errores (Red Cell Formatting)**
  - Cuando una fila del batch contiene errores de validación (cédula inválida, email malformado), resaltar la celda en rojo en el Excel original en memoria.
  - Subir el archivo resultante con errores a Cloud Storage mediante `io.Pipe` sin guardarlo temporalmente en disco.
- [ ] **P03. Mapeador Cero-Asignaciones (Datastore Key a DTO String)**
  - Implementar funciones de conversión pura y bidireccional entre `*datastore.Key` y representaciones REST (strings/enteros primitivos) emulando `prax-go-common/mapper`.
  - Asegurar 0 asignaciones de heap en la conversión de llaves numéricas y con ancestros.
- [ ] **P04. Motor de Validación Cruzada en Memoria (Eliminación de `TempEmployee`)**
  - Implementar un motor que reciba un slice de empleados y verifique unicidad interna de identificaciones y correos en tiempo $O(n)$ usando maps.
  - Consultar en un único batch optimizado (`datastore.GetMulti`) las entidades existentes en Datastore para validar duplicados globales.
- [ ] **P05. Worker Pool con Límite de Concurrencia para `prax-workplaces-service`**
  - Para cada empleado del lote, verificar la existencia de su centro de trabajo consumiendo el microservicio REST interno `prax-workplaces-service`.
  - Controlar las peticiones salientes con un pool de 10 workers concurrentes para no saturar el servicio dependiente.
- [ ] **P06. Publicador Pub/Sub Determinista sin `Thread.sleep`**
  - Migrar la publicación de eventos hacia `bigquery-publishing-topic`.
  - Utilizar `topic.Publish()` que retorna un `*pubsub.PublishResult`, recolectando los resultados con `sync.WaitGroup` y cerrando con `topic.Stop()` para asegurar confirmación antes de retornar el HTTP 200.
- [ ] **P07. Endpoint Push Handler de Pub/Sub con Validación OIDC de Cloud Run**
  - Construir el handler HTTP para endpoints Push de Cloud Run (ej. `/bigquery`).
  - Validar el token Bearer emitido por Google Cloud en el header `Authorization` y deserializar de forma segura el mensaje envuelto en Base64 (`message.data`).
- [ ] **P08. Batch Buffer & Deduplicador para BigQuery Streaming Inserts**
  - Diseñar un agrupador en memoria que acumule registros hasta alcanzar 500 filas o 2 segundos de inactividad.
  - Calcular el hash SHA-256 de los datos de cada registro para usarlo como `InsertID` en BigQuery, garantizando deduplicación automática de eventos reintentados.
- [ ] **P09. Reconciliador de Datos Datastore vs. BigQuery en Go (Reemplazo de Python)**
  - Reescribir la lógica de `prax-data-reconciler` en Go: consultar registros de una compañía en Datastore y compararlos con la tabla materializada en BigQuery.
  - Generar un reporte en streaming indicando registros faltantes, desactualizados o duplicados.
- [ ] **P10. Exportador Datastore con Paginación Segura mediante Cursores**
  - Implementar un extractor masivo de entidades Datastore que itere colecciones de más de 50,000 registros.
  - Utilizar `query.Start(cursor)` guardando el cursor periódicamente para permitir reanudación en caso de interrupción.
- [ ] **P11. Microservicio Agregador Rápido (`cf-employee-counter` en Go)**
  - Reemplazar la Cloud Function en Python `cf-employee-counter` con un microservicio Go.
  - Utilizar `client.NewAggregationQuery().WithCount()` sobre Datastore para calcular conteos de empleados por sede con latencias inferiores a 30 ms.
- [ ] **P12. Transacciones Multi-Entidad Atómicas en Datastore**
  - Implementar una transacción con `client.RunInTransaction`: actualizar el conteo en la entidad `Company`, crear el registro `BatchTask` y guardar el log de auditoría.
  - Manejar reintentos automáticos ante conflictos de concurrencia (`concurrent transaction conflict`).
- [ ] **P13. Middleware de Compatibilidad con `prax-api-gateway` y `prax-ui`**
  - Crear un middleware HTTP en Go que capture los headers inyectados por el API Gateway (`X-User-Email`, `X-Company-Key`, `Authorization`).
  - Inyectar las credenciales parseadas en el `context.Context` de la solicitud para consumo de los handlers.
- [ ] **P14. Generador de URLs Firmadas V4 para Google Cloud Storage (GCS)**
  - Implementar una función que genere URLs de descarga firmadas (`storage.SignedURL`) con una validez de 15 minutos para que `prax-ui` descargue archivos procesados de forma segura.
- [ ] **P15. Endpoint de Sondeo de Tareas (`batch/check/{taskKey}`) con Caché en Memoria**
  - Implementar el endpoint de polling que consulta el progreso de un lote de empleados.
  - Utilizar una caché interna con expiración rápida (TTL de 2 segundos) para mitigar el bombardeo de peticiones desde el frontend sin saturar Datastore.
- [ ] **P16. Semáforo de Concurrencia de Nivel de Servicio para Cloud Run**
  - Proteger la instancia de Cloud Run limitando el procesamiento simultáneo de cargas pesadas a máximo 10 goroutines activas mediante un canal semáforo (`chan struct{}`).
  - Retornar código HTTP `429 Too Many Requests` con header `Retry-After` si el semáforo está lleno.
- [ ] **P17. Migración de DTOs Java Jackson a Structs Go con Omitempty**
  - Convertir modelos complejos de Java (con herencia y anotaciones `@JsonProperty`) a structs planos de Go.
  - Configurar tags `json:"...,omitempty"` y asegurar compatibilidad binaria 1:1 con las respuestas JSON que espera `prax-ui`.
- [ ] **P18. Pipeline de Cálculo de Clima Organizacional (`evaluation-processing`)**
  - Procesar evaluaciones masivas de empleados: leer respuestas de preguntas, agrupar por dimensiones (liderazgo, comunicación) y calcular promedios y percentiles.
  - Distribuir el cálculo de dimensiones entre múltiples goroutines coordinadas con canales.
- [ ] **P19. Instrumentación SRE: Métricas en Cloud Monitoring & Health Probes**
  - Exponer endpoints `/healthz` (liveness de Cloud Run) e `/internal/employees/health` (readiness verificando conexión a Datastore y Pub/Sub).
  - Medir y registrar la duración de cada fase del batch: parseo Excel, validación, inserción en Datastore y publicación en Pub/Sub.
- [ ] **P20. Auditoría de Descriptores de Red y Conexiones HTTP Idle**
  - Configurar un `http.Client` singleton reutilizable con `Transport` configurado (`MaxIdleConns: 100`, `IdleConnTimeout: 90s`).
  - Escribir un test con `httptest.Server` para verificar que las conexiones no queden en estado `CLOSE_WAIT` ni filtren sockets bajo ráfagas de 1,000 requests.

---

## 🔍 Lecciones de Incidentes Reales: Troubleshooting en Producción (PRAX Wiki)

Para mantener este `README.md` balanceado y ágil sin sobrecargar la vista principal con cientos de líneas adicionales de post-mortems e historiales de bugs, el análisis forense profundo de las incidencias de producción y su suite de práctica se encuentra documentado en un archivo dedicado:

👉 **[Guía Completa de Troubleshooting & 20 Retos Prácticos en Go](file:///home/yonax73/computer-science/docs/PRAX_TROUBLESHOOTING.md)**

### Resumen de los 10 Fallos Reales Analizados (Extraídos de `prax-wiki`):
1. **Incompatibilidad de Claves WebSafe Datastore entre Go y Java (`s~` vs `k~` y Protobuf v1 vs v3):** Colapso con `IllegalArgumentException: id must not be equal to zero` al deserializar claves `Eh...` en Objectify. (`ADR-20260910-01`).
2. **Bloqueo del Streaming Buffer de BigQuery y Veneno en Pub/Sub (Poison Pill):** Rechazo de mutaciones `UPDATE` sobre filas en el buffer columnar caliente y retención de mensajes por más de 96 minutos. (`ADR-20260525-streaming-buffer-pubsub-lock`).
3. **Cuello de Botella N+1 y Fuga de Memoria en Objectify Session Cache:** Caída con HTTP 503 por miles de queries individuales de red y acumulación de miles de objetos pesados en Heap de 512 MiB. (`ADR-202605-12-datastore-eager-loading.md`).
4. **Saturación de Cuotas de Concurrencia en BigQuery (100 queries) y Backpressure:** Colapso de consultas analíticas DML y rechazo de tráfico con HTTP 429 hacia Pub/Sub. (`adr-202605-11-bigquery-concurrency-sweetspot.md`).
5. **Entidades Duplicadas y Corrupción de Métricas por Migraciones sin Zero-Trust:** División por cero en el frontend al borrar sedes que aún tenían empleados o sin recalcular `occupationAtWorkplace`. (`PM-20260630-deduplicacion-entidades.md`).
6. **Desbordamiento del Límite de 1.500 Bytes en Propiedades Indexadas de Datastore:** Falla fatal `InvalidArgument` al persistir catálogos de taxonomías (O*NET/UNSPSC) con descripciones largas indexadas. (`ADR-20260916-02`).
7. **Desincronización de Tipos de Datos (String vs. Integer) entre Datastore y BigQuery:** Consultas de reconciliación retornando listas vacías por discrepancia de tipos de ID. (`PM-20260630-deduplicacion-entidades.md`).
8. **Fuga de Conexiones TCP y Descriptores de Red (CLOSE_WAIT / Sockets Exhaustion):** Clientes HTTP sin drenaje completo de `resp.Body` agotando los sockets del contenedor de Cloud Run. (`REF-202605-23`).
9. **Caída del API Gateway por Timeouts en Cascada ("El Muro" 5xx > 10/s):** Operaciones pesadas sincrónicas asfixiando las conexiones del Gateway central. (`REF-202605-23`).
10. **Agotamiento de Memoria RAM (> 85% OOM Kill) en Cloud Run 512 MiB por Parseo No-Streaming:** Caída del proceso por carga de archivos de Excel completos en el DOM en lugar de iterar por streaming con `excelize.Rows()`. (`ADR-20260911-01`).

*Consulta los 20 problemas prácticos de depuración correspondientes (T01 a T20) en [docs/PRAX_TROUBLESHOOTING.md](file:///home/yonax73/computer-science/docs/PRAX_TROUBLESHOOTING.md).*

---

## 🛡️ Retos Prácticos de Seguridad OWASP Top 10 en Go

| Vulnerabilidad OWASP | Riesgo Técnico en Go | Reto Práctico a Implementar |
| :--- | :--- | :--- |
| **A01: Broken Access Control** | Modificación de recursos de otros tenants cambiando el `companyKey` en el path. | **Middleware BOLA/IDOR Defender**: Validar que el `companyKey` del path coincida estrictamente con los claims autorizados en el token del contexto. |
| **A02: Cryptographic Failures** | Almacenar datos confidenciales (cédulas, números de cuenta) en texto plano. | **Cifrado de Campos con AES-256-GCM**: Crear un envoltorio criptográfico para structs que cifre y descifre campos PII automáticamente antes de persistir en Datastore. |
| **A03: Injection** | Inyección en consultas GQL de Datastore o queries dinámicas de BigQuery. | **Parameter-Safe Query Builder**: Refactorizar consultas dinámicas para usar binding estricto de parámetros posicionales y rechazar concatenaciones directas de strings. |
| **A04: Insecure Design** | Agotamiento de recursos por cargas masivas concurrentes en Cloud Run. | **Upload Rate Limiting & Quotas**: Implementar un rate limiter con Token Bucket por Tenant que limite el número de archivos procesables por hora. |
| **A05: Security Misconfiguration** | Headers HTTP por defecto que exponen versiones del servidor o permiten ataques XSS/Clickjacking. | **Hardened Security Headers Middleware**: Inyectar headers `Content-Security-Policy`, `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY` y deshabilitar banners de servidor. |
| **A06: Vulnerable Components** | Módulos con vulnerabilidades conocidas en `go.mod`. | **CI/CD Vulnerability Gate**: Integrar `govulncheck` y `.osv-scanner.toml` en el script de validación local y en pipelines de Cloud Build. |
| **A07: Identification & Auth** | Aceptación de tokens JWT sin firma o con algoritmo `none`. | **Strict JWT Validator**: Validador de tokens RS256 con verificación de clave pública JWKS de Google, comprobando `exp`, `nbf`, `iss` y `aud` explícitos. |
| **A08: Integrity Failures** | Ataques de **Zip Bomb** o descompresión maliciosa en archivos `.xlsx` subidos. | **Decompression Bomb Guard**: Envolver el stream de lectura de Excel con un `io.LimitReader` que aborte la descompresión si el tamaño expandido excede 50 MiB. |
| **A09: Logging Failures** | Fuga de datos sensibles en logs o ausencia de auditoría en errores críticos. | **PII-Masking Structured Logger**: Implementar un logger estructurado JSON con `log/slog` que enmascare automáticamente tokens, números de cédula y correos. |
| **A10: SSRF** | Microservicios que consumen URLs provistas por el usuario y acceden a la metadata de GCP (`169.254.169.254`). | **Safe HTTP Client con IP Filtering**: Crear un `http.Transport` personalizado con hook `Control` en el `net.Dialer` que bloquee IPs privadas (RFC 1918) y la metadata de GCP antes del handshake. |

---

# 💼 30 Problemas de Entrevistas Técnicas en Go (Nivel Senior / Cloud-Native Staff)

Esta sección contiene 30 preguntas y retos técnicos reales utilizados en entrevistas para puestos Senior y Staff Go Engineer en empresas de tecnología y cloud.

---

### Bloque A: Concurrencia, Scheduler M:N & Primitivas de Sincronización (1 a 10)

#### 1. ¿Cómo funciona el Scheduler M:N de Go y qué representan G, M y P?
- **Respuesta Esperada:** Go no mapea 1 goroutine a 1 hilo del sistema operativo. Mapea $M$ goroutines a $N$ hilos del SO a través del modelo **G-M-P**:
  - **G (Goroutine):** Representa la goroutine en ejecución (stack de 2 KiB inicial, instruction pointer, estado).
  - **M (Machine):** Hilo del sistema operativo gestionado por el runtime de Go.
  - **P (Processor):** Recurso lógico necesario para ejecutar código Go (determinado por `GOMAXPROCS`). Posee una cola local de ejecutables (*runqueue*).
  - **Work Stealing & Syscalls:** Si un procesador $P$ se queda sin trabajo en su cola local, roba la mitad de las goroutines de otro $P$. Si una goroutine realiza una syscall bloqueante, el runtime desacopla el hilo $M$ de $P$ y asigna un nuevo $M$ para continuar ejecutando otras goroutines en $P$.

#### 2. Implementar un Worker Pool con Graceful Shutdown y propagación de errores usando `golang.org/x/sync/errgroup`
- **Problema de Live Coding:** Escribir una función que procese un slice de URLs concurrentemente con un límite de 5 workers. Si una URL falla, debe cancelar inmediatamente todas las demás solicitudes y retornar el error.
- **Aspectos Evaluados:** Uso de `errgroup.WithContext(ctx)`, control del canal con semáforo o límite de workers, y cancelación temprana del contexto.

#### 3. ¿Qué sucede al intentar escribir o cerrar un canal cerrado en Go? ¿Y al leer de un canal cerrado?
- **Respuesta Esperada:**
  - **Escribir en un canal cerrado:** Genera un `panic: send on closed channel` en tiempo de ejecución.
  - **Cerrar un canal ya cerrado:** Genera un `panic: close of closed channel`.
  - **Leer de un canal cerrado:** No produce panic. Retorna inmediatamente el valor cero del tipo del canal y el booleano `ok = false` en la expresión `val, ok := <-ch` (o solo el valor cero si no se usa la tupla).
  - **Leer de un canal nil:** Bloquea la goroutine para siempre (causando deadlock si no hay más goroutines activas).

#### 4. Explicar la diferencia entre `sync.Mutex` y `sync.RWMutex`. ¿Cuándo puede `sync.RWMutex` degradar el rendimiento en comparación con `sync.Mutex`?
- **Respuesta Esperada:** `sync.RWMutex` permite múltiples lectores simultáneos (`RLock()`) o un único escritor exclusivo (`Lock()`). Sin embargo, en escenarios de escrituras muy frecuentes o en arquitecturas con decenas de cores donde los lectores saturan los registros de cache (cache line bouncing), la contención interna del contador de lectores de `RWMutex` puede ser más lenta y costosa que un `sync.Mutex` tradicional.

#### 5. Implementar un Rate Limiter con Ventana Deslizante Concurrente (Sliding Window Log)
- **Problema de Live Coding:** Diseñar un struct `RateLimiter` seguro para concurrencia que permita $N$ peticiones por segundo por cliente, calculando la tasa con precisión de milisegundos mediante timestamps.
- **Aspectos Evaluados:** Uso de `sync.Mutex` o `sync.Map`, limpieza de timestamps expirados y optimización de memoria.

#### 6. ¿Cómo funciona `sync.Once` internamente y por qué es más eficiente que proteger una inicialización con `sync.Mutex`?
- **Respuesta Esperada:** `sync.Once` utiliza primero una carga atómica rápida (`atomic.LoadUint32(&o.done)`). Si el valor ya es 1 (la función ya se ejecutó), retorna inmediatamente en una sola instrucción de CPU sin tocar ningún lock. Solo si es 0, adquiere un `sync.Mutex` interno, realiza un doble chequeo atómico y ejecuta la función.

#### 7. ¿Por qué los objetos en `sync.Pool` pueden desaparecer en cualquier momento? ¿Es adecuado usar `sync.Pool` para conexiones a bases de datos?
- **Respuesta Esperada:** `sync.Pool` está diseñado para amortiguar memoria y aliviar el Garbage Collector. Cada ciclo de GC puede limpiar total o parcialmente los objetos del pool que no estén en uso. **Nunca debe usarse para conexiones a bases de datos o sockets**, ya que no ofrece garantías de persistencia ni control de ciclo de vida ordenado; para conexiones se deben utilizar pools dedicados (como `sql.DB` o `pgxpool`).

#### 8. ¿Cómo prevenir fugas de Goroutines (Goroutine Leaks)? Mencione 3 causas comunes.
- **Respuesta Esperada:** Una goroutine leak ocurre cuando una goroutine permanece en memoria bloqueada indefinidamente. Causas comunes:
  1. Enviar datos a un canal sin buffer cuando ningún receptor está escuchando.
  2. Quedarse esperando en un `select` con canales que nunca envían datos ni tienen timeout.
  3. No consumir completamente el cuerpo de una respuesta HTTP (`resp.Body.Close()` sin haber drenado los datos o sin timeout en el cliente).

#### 9. Implementar una Cola Concurrente Lock-Free usando `atomic.Pointer`
- **Problema de Live Coding:** Implementar una cola FIFO con métodos `Enqueue` y `Dequeue` utilizando punteros atómicos sin recurrir a `sync.Mutex`.
- **Aspectos Evaluados:** Comprensión de CAS (Compare-And-Swap), punteros y prevención de referencias circulares.

#### 10. ¿Cómo propagar y cancelar contextos en cascada en un árbol de llamadas RPC/HTTP?
- **Respuesta Esperada:** Pasar `ctx context.Context` siempre como primer parámetro de cada función. Usar `context.WithTimeout(parentCtx, duration)` para poner un límite global. Al momento de que el cliente aborta la conexión, el servidor HTTP cancela `r.Context()`, y todas las operaciones subordinadas (consultas a base de datos, llamadas HTTP aguas abajo) que verifiquen `ctx.Done()` deben terminar inmediatamente y liberar recursos.

---

### Bloque B: Runtime, Gestión de Memoria & Estructuras Internas (11 a 20)

#### 11. ¿Cuál es la representación interna de un `slice` en memoria y qué sucede cuando excede su capacidad?
- **Respuesta Esperada:** Un slice en Go es un struct ligero de 24 bytes (en arquitecturas de 64 bits) compuesto por:
  ```go
  type slice struct {
      array unsafe.Pointer // Puntero al array subyacente
      len   int            // Longitud actual
      cap   int            // Capacidad máxima antes de realocar
  }
  ```
  Al hacer `append` superando la capacidad, Go asigna un nuevo array subyacente más grande (duplicando tamaño para slices pequeños, o creciendo ~1.25x para slices más grandes), copia los elementos y actualiza el puntero. El slice original sigue apuntando al array anterior si no fue reasignado.

#### 12. ¿Qué es Escape Analysis y qué patrones hacen que una variable escape del Stack al Heap?
- **Respuesta Esperada:** El compilador de Go decide en tiempo de compilación si una variable se asigna en el stack (extremadamente rápido, limpiado al retornar la función) o en el heap (gestionado por el GC). Patrones que fuerzan el escape:
  - Retornar un puntero a una variable local desde una función.
  - Asignar una variable a una interfaz (como pasarla a `fmt.Println(val)` que acepta `any`).
  - Slices cuyo tamaño no se conoce en tiempo de compilación o que superan los límites del stack.
  - Enviar punteros a través de canales.

#### 13. ¿Cómo funciona internamente un `map` en Go? ¿Por qué no es thread-safe?
- **Respuesta Esperada:** Un `map` en Go es un puntero a un struct `hmap` que contiene un array de buckets (`bmap`), donde cada bucket almacena hasta 8 pares llave/valor utilizando los 8 bits superiores del hash (*tophash*). No es thread-safe porque no tiene locks internos para maximizar la velocidad en casos de uso de un solo hilo; si el runtime detecta lecturas y escrituras simultáneas, lanza un error fatal irrecuperable: `fatal error: concurrent map writes`.

#### 14. ¿Por qué `var x *MiStruct = nil; var i MiInterface = x; i == nil` evalúa a `false`?
- **Respuesta Esperada:** En Go, una interfaz se representa internamente como una tupla de dos punteros: `(tipo, valor)`:
  - En `iface` (interfaces con métodos): puntero a la tabla de métodos (`*itab`) y puntero a los datos.
  - En `eface` (interfaz vacía `any`): puntero al descriptor de tipo (`*_type`) y puntero a los datos.
  Una variable de interfaz solo es estrictamente igual a `nil` cuando **tanto su tipo como su valor son nil**. En este caso, el tipo es `*MiStruct` y el valor es `nil`, por lo que la interfaz no es nil.

#### 15. ¿Cómo funciona el Garbage Collector de Go y qué significa que sea un recolector Trícolor Concurrente?
- **Respuesta Esperada:** Go utiliza un GC de tipo **Concurrent Mark-Sweep** con abstracción de tres colores:
  - **Blanco:** Objetos no descubiertos (candidatos a recolección).
  - **Gris:** Objetos descubiertos pero cuyos hijos aún no han sido analizados.
  - **Negro:** Objetos descubiertos cuyos hijos ya fueron analizados (no se liberan).
  Se ejecuta concurrentemente con la aplicación mediante una técnica llamada *Write Barrier*, manteniendo las pausas de *Stop-The-World (STW)* típicamente por debajo de **1 milisegundo**.

#### 16. ¿Qué impacto tienen las variables `GOMEMLIMIT` y `GOGC` en entornos con memoria restringida (como Cloud Run o Kubernetes)?
- **Respuesta Esperada:** Históricamente, `GOGC` detonaba el ciclo de recolección cuando el heap alcanzaba un porcentaje de crecimiento (por defecto 100%). En contenedores con límites estrictos de RAM (ej. 512 MiB), si el heap crecía rápidamente antes de que `GOGC` disparara la recolección, el kernel de Linux liquidaba el proceso por OOM (Out Of Memory). `GOMEMLIMIT` (introducido en Go 1.19) define un techo de memoria suave (ej. `GOMEMLIMIT=460MiB`): cuando la memoria total se acerca a ese límite, el runtime ejecuta recolecciones de basura más agresivas para evitar que el contenedor sea eliminado por OOM.

#### 17. Memory Padding & Struct Alignment: ¿Cómo afecta el orden de los campos al tamaño de un struct?
- **Problema de Live Coding:**
  ```go
  type Malo struct {
      a bool    // 1 byte (+ 7 bytes padding)
      b float64 // 8 bytes
      c bool    // 1 byte (+ 7 bytes padding)
  } // Tamaño: 24 bytes

  type Bueno struct {
      b float64 // 8 bytes
      a bool    // 1 byte
      c bool    // 1 byte (+ 6 bytes padding)
  } // Tamaño: 16 bytes
  ```
- **Aspectos Evaluados:** Comprensión de alineación de palabras en arquitecturas de 64 bits y reducción de huella de memoria en slices de millones de elementos.

#### 18. ¿Por qué `string` es inmutable en Go y cómo convertir entre `[]byte` y `string` sin asignar memoria (Zero-Copy)?
- **Respuesta Esperada:** Un `string` es inmutable para garantizar seguridad concurrente y permitir que múltiples sub-cadenas compartan el mismo array de bytes subyacente sin duplicaciones. La conversión normal `string(b)` o `[]byte(s)` asigna nueva memoria en el heap para garantizar la inmutabilidad. Para conversiones sin asignación en caminos críticos, se puede usar `unsafe.String` y `unsafe.StringData` (en Go moderno) con la precaución estricta de no mutar el slice de bytes posterior.

#### 19. Implementar una función genérica `Filter[T any](s []T, predicate func(T) bool) []T` optimizando asignaciones de memoria
- **Problema de Live Coding:** Escribir una función genérica que filtre un slice in-place reutilizando el array subyacente para lograr 0 asignaciones en heap cuando sea posible.
- **Aspectos Evaluados:** Uso de Go Generics y técnica de dos punteros sobre el slice original (`s[:0]`).

#### 20. ¿Qué ventajas y desventajas tiene el uso de `reflect` (Reflexión) en Go frente a la generación de código (`go generate`)?
- **Respuesta Esperada:** La reflexión permite inspeccionar tipos y valores en runtime (usado por `encoding/json`), pero es significativamente más lenta, impide optimizaciones del compilador y traslada errores de tipos del tiempo de compilación al tiempo de ejecución. La generación de código (ej. `easyjson`, `protobuf`) genera código tipado estático en tiempo de compilación, logrando máximo rendimiento y seguridad a cambio de un paso adicional en el build.

---

### Bloque C: Arquitectura, Redes & Cloud-Native (21 a 30)

#### 21. ¿Cómo estructurar un servidor HTTP en Go con Graceful Shutdown que termine conexiones activas correctamente?
- **Problema de Live Coding:** Escribir la función `main` que inicie un `http.Server`, escuche señales `os.Interrupt` / `syscall.SIGTERM` con `signal.NotifyContext` y ejecute `server.Shutdown(ctx)` con un timeout máximo de 10 segundos.
- **Aspectos Evaluados:** Manejo de señales del SO, orden de parada y cierre seguro de sockets sin cortar transacciones en vuelo.

#### 22. ¿Cómo evitar el leak de descriptores de sockets al consumir APIs HTTP externas con `http.Client`?
- **Respuesta Esperada:** Siempre se debe cerrar el cuerpo de la respuesta con `defer resp.Body.Close()`. Además, si no se lee todo el cuerpo hasta el final, la conexión TCP no puede ser reutilizada por el `http.Transport` para Keep-Alive. Por tanto, antes de cerrar, se debe drenar el cuerpo: `io.Copy(io.Discard, resp.Body)`.

#### 23. Implementar el patrón Singleflight desde cero para evitar Cache Stampede (Thundering Herd)
- **Problema de Live Coding:** Diseñar un struct `Group` con método `Do(key string, fn func() (any, error)) (any, error)` que garantice que múltiples llamadas concurrentes solicitando la misma llave esperen el resultado de la primera ejecución en curso.
- **Aspectos Evaluados:** Uso de `sync.Mutex`, canales para broadcast de resultados y limpieza segura de la llave del map al finalizar.

#### 24. Clean Architecture en Go: ¿Cómo organizar los paquetes sin caer en dependencias circulares?
- **Respuesta Esperada:** En Go, los ciclos de importación están prohibidos por el compilador. Para aplicar Clean Architecture:
  - Definir las interfaces en el paquete que las **consume** (consumidor define la interfaz), no en el paquete que las implementa.
  - La capa de dominio o entidades (`entity`, `model`) no debe depender de ninguna capa externa (ni base de datos, ni frameworks).
  - Usar una estructura estándar: `cmd/` (entrypoints), `internal/` (código privado de la aplicación: `service`, `repository`, `handler`), y `pkg/` solo para librerías públicas reutilizables.

#### 25. ¿Cómo implementar un Circuit Breaker básico con tres estados (`Closed`, `Open`, `Half-Open`)?
- **Problema de Live Coding:** Diseñar un struct `CircuitBreaker` que monitoree fallos en peticiones externas. Si supera un umbral de 5 fallos consecutivos, transiciona a `Open` rechazando requests durante un periodo de enfriamiento de 30 segundos, pasando luego a `Half-Open` para probar una petición piloto.
- **Aspectos Evaluados:** Máquina de estados, sincronización atómica o con mutex y control temporal.

#### 26. ¿Cuáles son las diferencias críticas entre gRPC y REST para comunicación inter-servicios en Go?
- **Respuesta Esperada:**
  - **Protocolo:** gRPC opera sobre HTTP/2 (multiplexación de streams en una única conexión TCP, compresión de cabeceras HPACK), mientras que REST tradicionalmente opera sobre HTTP/1.1.
  - **Serialización:** gRPC utiliza Protocol Buffers (binario compacto, tipado estricto, generación de stubs de alto rendimiento), mientras que REST utiliza JSON (texto plano, mayor tamaño de payload, alto costo de reflexión).
  - **Streaming:** gRPC soporta streaming unario, del servidor, del cliente y bidireccional de forma nativa.

#### 27. ¿Cómo mitigar el problema de Backpressure en un consumidor de Cloud Pub/Sub o Kafka en Go?
- **Respuesta Esperada:**
  - Configurar límites en el cliente receptor (ej. `subscription.ReceiveSettings.MaxOutstandingMessages` en Google Cloud Pub/Sub).
  - Implementar un worker pool con un canal con buffer que actúe como buffer de amortiguación.
  - Si los workers se saturan y el canal se llena, dejar de solicitar nuevos mensajes o ralentizar las lecturas para que el broker de mensajería mantenga los mensajes en cola en lugar de provocar un OOM en el consumidor.

#### 28. ¿Cómo diagnosticar cuellos de botella de CPU y contención de memoria usando `pprof` y trazas de ejecución (`go tool trace`)?
- **Respuesta Esperada:**
  - Activar endpoints `net/http/pprof` en servicios HTTP.
  - Capturar perfiles con `go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30` (CPU) o `.../heap` (memoria).
  - Analizar perfiles con vistas interactivas `top`, `list <función>` y gráficos `flamegraph` en el navegador con `-http=:8081`.
  - Usar `go tool trace trace.out` para visualizar gráficamente en la línea de tiempo el comportamiento del scheduler, bloqueos de goroutines y pausas del GC.

#### 29. ¿Cómo realizar pruebas unitarias sobre llamadas a Google Cloud Storage o Datastore sin depender de emuladores lentos?
- **Respuesta Esperada:** Aplicar inversión de dependencias: en lugar de acoplar la lógica de negocio a los clientes concretos `*storage.Client` o `*datastore.Client`, definir interfaces de servicio pequeñas en la capa de negocio (ej. `FileUploader`, `EmployeeRepository`). En los tests unitarios, suministrar implementaciones mock en memoria que simulen fallos, timeouts y respuestas exitosas en milisegundos sin I/O de red.

#### 30. Implementar un Middleware de Reintentos con Retroceso Exponencial y Jitter (Exponential Backoff with Jitter)
- **Problema de Live Coding:** Escribir una función en Go que ejecute una operación falible (ej. llamada HTTP o consulta a Datastore), reintentando hasta 3 veces con una espera que crezca exponencialmente ($2^n \times \text{base}$) más un componente aleatorio (*jitter*) para evitar la sincronización de reintentos concurrentes.
- **Aspectos Evaluados:** Respeto al `ctx.Done()`, uso de `time.Sleep` o `time.After`, y generación correcta del jitter aleatorio.

---

# 🤖 Nivel 8: AI, Big Data & Cybersecurity Weekend Projects (10 Proyectos de Alto Impacto)

Proyectos prácticos, autocontenidos y de alta complejidad arquitectónica diseñados para ser desarrollados e implementados en un **fin de semana (~10 a 16 horas)** cada uno. Cada reto combina el alto rendimiento concurrente de **Go**, el procesamiento masivo de datos (**BigQuery / Datastore / Cloud Storage**), modelos de **IA Generativa / Embeddings (Google Gemini SDK / Vertex AI)** y principios avanzados de **Ciberseguridad Defensiva y DevSecOps**.

---

### 💡 Caso de Estudio: Por qué Go es el Lenguaje Predilecto para Plataformas de IA y Seguridad Cloud-Native
> **Contexto de la Industria:** Mientras que el entrenamiento e investigación de modelos de IA ocurre mayoritariamente en Python, la **inferencia en producción a gran escala, los proxies de seguridad, los pipelines de ingesta masiva (Big Data) y los agentes de ciberseguridad** se están reescribiendo masivamente en **Go** (ej. Ollama, Kubernetes K8s, Docker, Cloudflare Workers runtime, HashiCorp Vault, Palo Alto Networks Cortex). 
> 
> **La Razón:** En arquitecturas de inferencia de IA y procesamiento de Big Data en tiempo real, Python introduce cuellos de botella severos por el GIL (Global Interpreter Lock), consumo masivo de memoria RAM y latencias impredecibles en el manejo de streams de red. Go permite crear agentes y microservicios binarios autocontenidos con tiempos de arranque inferiores a 50 milisegundos, concurrencia masiva con miles de goroutines para orquestación de embeddings y llamadas a APIs de LLMs, y huella de memoria minúscula que reduce drásticamente los costos de infraestructura en la nube.

---

### Los 10 Proyectos Weekend:

#### 1. AI-Driven Cloud Audit & Threat Detection Engine en BigQuery (SIEM Inteligente)
- **Descripción:** Motor de detección de anomalías y amenazas de seguridad en tiempo real sobre los logs de auditoría de Google Cloud.
- **Stack Tecnológico:** Go 1.22+, `cloud.google.com/go/bigquery`, Google Gemini SDK (`google.golang.org/genai`), Cloud Pub/Sub.
- **Desafío Técnico:**
  - Consultar periódicamente en BigQuery tablas particionadas de auditoría (`cloudaudit_googleapis_com_activity_*`).
  - Agrupar secuencias de eventos sospechosos (ej. creación de llaves de Service Account a deshoras, múltiples llamadas fallidas seguidas de una modificación de políticas IAM con `SetIamPolicy`, o descargas anómalas de tablas de BigQuery).
  - Utilizar Gemini con *Structured Outputs* (JSON Schema) para clasificar el incidente según la matriz **MITRE ATT&CK**, calcular un puntaje de riesgo (0–100) y generar un plan de contención automatizado enviado a un webhook o canal de Slack.
- **Entregable MVP (Fin de Semana):** Un worker en Go que se ejecute cada 5 minutos, analice los últimos 1,000 eventos de auditoría y genere un reporte HTML/JSON con los incidentes priorizados y comandos de remediación en `gcloud`.

---

#### 2. Smart PII & Data Loss Prevention (DLP) Streamer para Ingesta en Data Lake
- **Descripción:** Pipeline de sanitización en streaming que intercepta archivos masivos (CSV, JSON, TSV de empleados) antes de insertarse en BigQuery o Cloud Storage, redactando información confidencial no estructurada con IA.
- **Stack Tecnológico:** Go (`io.Pipe`, `bufio.Scanner`, `sync.Pool`), Gemini Embeddings (`text-embedding-004`) o Gemini Flash, BigQuery Streaming API, `crypto/hmac`.
- **Desafío Técnico:**
  - Los filtros tradicionales basados en Expresiones Regulares (Regex) fallan al detectar PII incrustado en texto libre (ej. comentarios de retroalimentación en encuestas de clima laboral de PRAX como *"el señor Carlos cuyo teléfono personal es..."*).
  - El pipeline en Go procesa el stream en memoria mediante buffers reciclables sin volcar el archivo completo al disco.
  - Los fragmentos de texto sospechosos son evaluados concurrentemente por un pool de workers que consulta a Gemini Flash para identificar entidades sensibles (nombres, cédulas, datos de salud) y sustituirlas por un pseudónimo o hash HMAC determinista, permitiendo análisis relacional sin violar normativas de privacidad (GDPR / Habeas Data).
- **Entregable MVP (Fin de Semana):** Herramienta CLI y microservicio HTTP en Go que reciba un archivo de 50,000 registros, lo higienice en streaming con overhead < 20 MiB de RAM y cargue las filas limpias directamente en BigQuery.

---

#### 3. Semantic WAF & Firewall Anti-Prompt Injection para APIs de LLM
- **Descripción:** Reverse Proxy de seguridad en Go que protege los endpoints de LLM corporativos contra ataques de *Prompt Injection*, *Jailbreaking* y exfiltración de instrucciones del sistema (*System Prompt Leakage*).
- **Stack Tecnológico:** Go (`net/http/httputil.ReverseProxy`), Base de vectores en memoria (Cosine Similarity concurrente), Google Gemini API, `golang.org/x/time/rate`.
- **Desafío Técnico:**
  - Interceptar peticiones entrantes dirigidas a las APIs de IA antes de que lleguen al modelo principal.
  - **Fase 1 (Fast-Path):** Calcular la similitud semántica del prompt del usuario contra un corpus vectorial precargado en memoria de técnicas de jailbreak conocidas (DAN, Token Smuggling, delimitadores falsos `--- END OF SYSTEM INSTRUCTIONS ---`). Si la similitud coseno supera 0.88, bloquear inmediatamente con HTTP 403 Forbidden (< 5 ms).
  - **Fase 2 (Deep-Scan):** Para casos ambiguos, despachar una verificación asíncrona ultrarrápida a Gemini Flash (< 120 ms) con un metaprompt defensivo.
  - Registrar métricas de ataques bloqueados y direcciones IP de origen.
- **Entregable MVP (Fin de Semana):** Un proxy inverso HTTP desplegable en Cloud Run que proteja cualquier endpoint de LLM, bloqueando los 20 vectores de ataque de inyección más comunes con latencia añadida mínima.

---

#### 4. AST Code Reviewer & SAST Heurístico Potenciado por IA para Repositorios Go
- **Descripción:** Analizador estático de seguridad que recorre el Árbol de Sintaxis Abstracta (AST) de código fuente Go y utiliza un LLM para evaluar la explotabilidad real de vulnerabilidades, eliminando falsos positivos.
- **Stack Tecnológico:** Go (`go/parser`, `go/ast`, `go/token`), Google Gemini API con modo JSON estricto, Git CLI.
- **Desafío Técnico:**
  - Herramientas tradicionales como `gosec` generan alta tasa de falsos positivos al detectar llamadas potencialmente inseguras (ej. `exec.Command` o `sql.Open`) sin comprender el flujo de sanitización previo.
  - El analizador en Go parsea el AST de los paquetes Go, detecta nodos de riesgo (llamadas al sistema, descompresión de zips, variables en queries SQL) y extrae únicamente la función afectada junto con sus tipos relacionados (reduciendo el consumo de tokens en un 95% respecto a enviar archivos enteros).
  - Consulta a Gemini enviando el AST simplificado para emitir un veredicto binario: ¿Es vulnerable? Sí/No, justificación técnica y un parche sugerido en formato `git diff`.
- **Entregable MVP (Fin de Semana):** Comando CLI ejecutable (`go-ai-sec-lint ./...`) que genere un reporte en Markdown listo para integrarse en GitHub Actions con parches de código listos para aplicar.

---

#### 5. BigQuery Exfiltration Guard & SQL Cost/Abuse Analyzer
- **Descripción:** Daemon de seguridad que vigila la ejecución de consultas en BigQuery para abortar en tiempo real intentos de exfiltración masiva de datos y consultas maliciosas o financieramente destructivas.
- **Stack Tecnológico:** Go, `cloud.google.com/go/bigquery`, Gemini 1.5 Flash, Cloud Monitoring.
- **Desafío Técnico:**
  - Monitorear el stream de jobs en `INFORMATION_SCHEMA.JOBS_BY_PROJECT` en BigQuery cada 10 segundos.
  - Detectar consultas que violen umbrales de seguridad (ej. escaneo de más de 500 GB en tablas que contienen información de nómina o empleados, o queries `SELECT *` sin cláusula `WHERE` ejecutadas por Service Accounts comprometidas).
  - Enviar el árbol sintáctico de la consulta SQL y el contexto del usuario a Gemini para clasificar la intención: "Consulta analítica legítima", "Query mal optimizada de desarrollador" o "Intento de exfiltración maliciosa".
  - Si se clasifica como exfiltración con confianza > 90%, invocar `job.Cancel()` en BigQuery de inmediato y revocar temporalmente los tokens de la cuenta.
- **Entregable MVP (Fin de Semana):** Microservicio en Go con modo `--dry-run` y `--enforce` que audita queries activas y emite alertas automáticas con cálculo de bytes escaneados e intención detectada.

---

#### 6. Agente de Principio de Menor Privilegio (Least Privilege IAM Recommender)
- **Descripción:** Agente en Go que analiza el uso histórico real de permisos en BigQuery y genera automáticamente definiciones de roles de GCP con privilegios mínimos en código Terraform.
- **Stack Tecnológico:** Go, Google Cloud Resource Manager API, BigQuery Audit Logs, Gemini API.
- **Desafío Técnico:**
  - En entornos enterprise como PRAX, muchas Service Accounts reciben roles amplios (`roles/editor` o `roles/datastore.owner`) por conveniencia operativa.
  - Este servicio en Go ejecuta consultas analíticas en BigQuery sobre los últimos 30 a 90 días de Cloud Audit Logs para mapear exactamente qué métodos de API (`storage.objects.get`, `datastore.entities.query`, `pubsub.topics.publish`) ha utilizado efectivamente cada Service Account.
  - Gemini analiza la lista de permisos consumidos y redacta un rol personalizado de GCP (`google_project_iam_custom_role`) en sintaxis **Terraform (HCL)**, verificando que no se rompan dependencias operativas y eliminando permisos peligrosos no utilizados.
- **Entregable MVP (Fin de Semana):** Herramienta CLI en Go (`iam-least-privilege --sa <email> --days 30`) que genera un archivo `iam_remediation.tf` listo para aplicar con Terraform.

---

#### 7. RAG Seguro con Búsqueda Vectorial sobre Políticas SOC 2 / ISO 27001 en BigQuery
- **Descripción:** Microservicio RAG (Retrieval-Augmented Generation) de grado empresarial que responde consultas de cumplimiento normativo y arquitectura segura para equipos de desarrollo.
- **Stack Tecnológico:** Go, BigQuery Vector Search (`VECTOR_SEARCH` SQL), Gemini Text Embeddings (`text-embedding-004`), `net/http` con mTLS y validación de Claims.
- **Desafío Técnico:**
  - Indexar políticas internas de seguridad, manuales de cumplimiento SOC 2 y guías OWASP dividiéndolos en fragmentos semánticos en Go y calculando sus embeddings vectoriales guardados en BigQuery.
  - Implementar un endpoint HTTP en Go protegido con mTLS y autenticación basada en roles.
  - Al recibir una pregunta técnica (ej. *"¿Podemos almacenar contraseñas usando SHA-1 o qué algoritmo exige nuestra política SOC 2?"*), ejecutar una búsqueda vectorial top-$K$ en BigQuery y construir un prompt con guardrails anti-alucinación donde el modelo deba citar el documento, sección y fecha de vigencia exacta de la norma.
- **Entregable MVP (Fin de Semana):** Microservicio en Go con endpoint `POST /api/v1/compliance/ask` con respuestas fundamentadas en referencias documentales exactas y latencia total < 1 segundo.

---

#### 8. HoneyToken & Canary Data Synthesizer con Monitoreo de Intrusiones en Datastore
- **Descripción:** Sistema de defensa activa que genera registros trampa (HoneyTokens) hiperrealistas mediante IA e instrumenta trampas de detección de intrusos en bases de datos NoSQL y Data Lakes.
- **Stack Tecnológico:** Go, Google Cloud Datastore SDK, Cloud Pub/Sub, Gemini API, Cloud Run.
- **Desafío Técnico:**
  - Los atacantes y empleados malintencionados suelen buscar datos sensibles (números de tarjetas, directivos VIP, credenciales).
  - El generador en Go invoca a Gemini para sintetizar entidades de prueba altamente creíbles (empleados ficticios con nombres realistas, empresas ficticias) que siguen exactamente las estructuras y tags de `prax-go-common/entity`.
  - Los HoneyTokens se insertan en Cloud Datastore marcados criptográficamente de forma imperceptible (ej. dígito de verificación especial o UUID con firma HMAC interna).
  - Un worker en Go escucha los eventos de lectura de la base de datos o queries en BigQuery: si un usuario o proceso consulta un HoneyToken, se dispara una alerta de intrusión de severidad máxima P1 y se bloquea la sesión en Cloud Armor.
- **Entregable MVP (Fin de Semana):** Comando CLI `honeytoken inject --count 5` y un webhook detector en Go que alerte instantáneamente ante cualquier intento de acceso al registro canario.

---

#### 9. Pipeline de Detección de Amenazas y Acoso en Evaluaciones de Clima (PRAX NLP)
- **Descripción:** Pipeline asíncrono para procesar decenas de miles de respuestas de texto libre en encuestas de clima y cultura organizacional, detectando banderas rojas críticas protegiendo el anonimato.
- **Stack Tecnológico:** Go (Worker Pool, streaming con canales), Cloud Pub/Sub, Cloud Datastore, Gemini 1.5 Flash, BigQuery.
- **Desafío Técnico:**
  - En encuestas de clima laboral masivas (como las gestionadas por PRAX), los colaboradores escriben comentarios abiertos donde pueden reportar situaciones graves de acoso sexual, amenazas físicas o corrupción interna que no pueden esperar semanas de revisión manual.
  - Un consumidor en Go procesa lotes de mensajes desde Cloud Pub/Sub a alta velocidad.
  - Distribuye las evaluaciones a un pool de goroutines que consultan la API de Gemini para clasificar las respuestas en categorías de riesgo ético/legal y calcular la polaridad del sentimiento.
  - **Privacidad Estricta:** Antes de persistir los hallazgos en BigQuery, el microservicio elimina o enmascara irreversiblemente cualquier identificador personal (nombre del empleado, correo, IP), preservando solo la metadata de sede y la clasificación de riesgo para auditoría.
- **Entregable MVP (Fin de Semana):** Microservicio Go capaz de procesar 1,000 respuestas de texto libre en menos de 90 segundos, clasificando incidentes críticos y visualizando métricas de riesgo ético en tablas agregadas de BigQuery.

---

#### 10. Threat Intelligence Correlation Engine & Zero-Day CVE Matcher
- **Descripción:** Hub de inteligencia de amenazas que descarga feeds continuos de vulnerabilidades públicas (CVEs, OSV, NVD), los cruza contra el inventario de dependencias de la empresa en BigQuery y evalúa la alcanzabilidad real con IA.
- **Stack Tecnológico:** Go (Fan-Out / Fan-In con canales), OSV.dev REST API, BigQuery, Gemini API.
- **Desafío Técnico:**
  - Descargar periódicamente las bases de datos de vulnerabilidades de código abierto (OSV/NVD) y persistirlas en BigQuery.
  - Cruzar las dependencias de todos los microservicios de la organización (`go.mod`, `pom.xml`, `package.json`) contra las vulnerabilidades reportadas.
  - Para cada vulnerabilidad clasificada como High o Critical, el servicio en Go extrae la firma de la función afectada del advisory y consulta a Gemini analizando el código fuente interno: ¿Nuestros microservicios realmente llaman a la función vulnerable o es una dependencia transitiva inalcanzable (*dead dependency*)?
  - Genera tickets priorizados en GitHub/Jira únicamente para vulnerabilidades con impacto real confirmado, adjuntando la línea de comando `go get package@version` para su actualización.
- **Entregable MVP (Fin de Semana):** Scanner CLI en Go (`threat-matcher --repo ./...`) que cruza las dependencias locales contra la API de OSV y genera un informe de alcanzabilidad y remediación asistido por IA.

