# Changelog

Todas las modificaciones notables de este proyecto se documentarán en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.1.0/),
y este proyecto se adhiere a [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Planned
- Implementación del ejercicio `01-basics/01-word-frequency/` (Word & Char Frequency Counter con `bufio.Scanner` y UTF-8 runes).
- Suite de pruebas unitarias basadas en tablas (*Table-Driven Tests*) para el módulo de conteo.

---

## [0.3.0] - 2026-09-20

### Added
- **Skill Pedagógico Antigravity (`.agents/skills/go-tutor/SKILL.md`):** Tutoría activa basada en Mayéutica Socrática, Técnica Feynman y evaluación estricta con rúbrica del 90%+ para desbloqueo de secciones.
- **Reglas del Entorno (`GEMINI.md`):** Directrices persistentes de tutoría para garantizar que el estudiante investigue e implemente su propio código sin soluciones automáticas prefabricadas.
- **Sincronización con Repositorio Remoto:** Vinculación y publicación del repositorio oficial en GitHub: [YonatanQuintero/go-zero-to-mastery](https://github.com/YonatanQuintero/go-zero-to-mastery) a través de la identidad SSH personal `github-personal`.

### Changed
- Normalización del directorio del workspace a `computer-science` (eliminación de espacios para compatibilidad con estándares POSIX).
- Actualización de autoría de Git en commits a `Yonatan Quintero <yonatan.a.quintero.r@gmail.com>`.

---

## [0.2.0] - 2026-09-20

### Added
- **Compendio de Troubleshooting Forense (`docs/PRAX_TROUBLESHOOTING.md`):**
  - Análisis de causa raíz de los 10 errores y cuellos de botella más comunes en producción documentados en `prax-wiki` (WebSafe keys `s~` vs `k~`, BigQuery streaming buffer lock & poison pill en Pub/Sub, N+1 memory leak en Objectify Session Cache, Gateway 502/504 "El Muro", OOM de 512 MiB en Cloud Run, límite de 1.500 bytes de Datastore, sockets en `CLOSE_WAIT`, etc.).
  - 20 nuevos problemas prácticos de depuración en Go (T01 a T20) divididos en 4 bloques: Almacenamiento, Resiliencia de Pipelines, Zero-Trust e Infraestructura Cloud Run / SRE.
- Enlaces y resumen ejecutivo de incidentes integrados en el `README.md`.

---

## [0.1.0] - 2026-09-20

### Added
- **Hoja de Ruta Maestra (`README.md`):**
  - Currículum progresivo de Niveles 1 al 6 (45 ejercicios prácticos desde fundamentos de memoria y punteros hasta sistemas distribuidos y microservicios GCP).
  - Casos de estudio de la industria previos a cada nivel (Monzo Bank, Discord Goroutines, Uber GC, VictoriaMetrics Zero-Alloc, Cloudflare TLS).
  - Nivel 7: 20 retos de producción modelados sobre el ecosistema empresarial de PRAX.
  - Matriz de Ciberseguridad Defensiva con los 10 riesgos OWASP adaptados al runtime de Go.
  - 30 problemas y preguntas de entrevista técnica Senior / Staff (Live Coding en concurrencia, memoria, escape analysis y arquitectura distribuida).
  - Nivel 8: 10 proyectos de fin de semana integrando Go + Big Data + Gemini AI + Ciberseguridad.
- **Automatización & Calidad:**
  - `Makefile` con targets para tests con detector de carreras de datos (`test-race`), linters (`golangci-lint`), escaneo de vulnerabilidades (`govulncheck`) y auditoría estática (`gosec`).
  - `.gitignore` especializado para herramientas, compilados y dumps de profiling de Go.
- **Estructura del Proyecto:** Scaffolding inicial de carpetas modulares (`01-basics/` a `09-ai-bigdata-security/`).
