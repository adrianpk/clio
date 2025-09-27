### 2025-09-27

**Tarea:** Implementar un sistema de configuración dinámica para el proyecto.

**Resumen de la tarea:**
Se solicitó un sumario de las variables de entorno y claves de configuración existentes para sentar las bases de un sistema de configuración dinámica. Esto implicó:
1.  Relevar variables de `.envrc`, `makefile` y `internal/am/key.go`.
2.  Documentar estas variables en `docs/drafts/envar.md`.
3.  Crear una nueva entidad `Param` para almacenar valores de configuración dinámicos en la base de datos.
4.  Integrar `Param` en todas las capas del proyecto (migración, modelo Go, repositorio, servicio, API Handler y Router).
5.  Crear un script de cURL para probar los endpoints de la API de `Param`.

**Detalles de la implementación:**

1.  **Documentación de Variables de Entorno:**
    *   Se creó `docs/drafts/envar.md` con un listado de variables de entorno y claves de configuración.
    *   Se tradujo el documento al inglés.
    *   Se ajustó el formato de la segunda sección a `ENVAR_NAME => property.name`.
    *   Se añadió una nota indicando que las variables `CLIO_SSG_BLOCKS_MAXITEMS`, `CLIO_SSG_INDEX_MAXITEMS`, `CLIO_SSG_SEARCH_GOOGLE_ENABLED`, `CLIO_SSG_SEARCH_GOOGLE_ID` serán configurables vía web.
    *   Se corrigió la descripción de cómo las variables de entorno se mapean a las propiedades internas en `internal/am/key.go`.

2.  **Implementación de la Entidad `Param`:**
    *   **Migración de Base de Datos:**
        *   Se creó el archivo `assets/migration/sqlite/20250927xxxxxx-create-param-table.sql` (donde `xxxxxx` es un timestamp).
        *   La tabla `param` incluye `id`, `name`, `description`, `value`, `ref_key` (para sobrescribir propiedades existentes), y campos de auditoría (`created_by`, `updated_by`, `created_at`, `updated_at`).
        *   Se incluyeron índices en `name` y `ref_key` para optimizar las búsquedas.
        *   Se corrigió el archivo de migración para incluir las directivas `-- +migrate Up` y `-- +migrate Down` que son necesarias para el migrador.
    *   **Modelo Go:**
        *   Se creó `internal/feat/ssg/param.go` con la definición de la estructura `Param`, siguiendo el patrón de otros modelos del proyecto.
        *   Se renombró de `Config` a `Param` para evitar confusiones con el sistema de configuración general.
    *   **Interfaz del Repositorio:**
        *   Se añadió la interfaz `ParamRepo` con métodos CRUD (`CreateParam`, `GetParam`, `GetParamByName`, `GetParamByRefKey`, `ListParams`, `UpdateParam`, `DeleteParam`) a `internal/feat/ssg/repo.go`.
    *   **Implementación del Repositorio SQLite:**
        *   Se implementaron los métodos CRUD para `Param` en `internal/repo/sqlite/ssg.go`, utilizando el patrón de acceso a datos existente.
        *   Se definió la constante `resParam = "param"`.
    *   **Consultas SQL:**
        *   Se creó `assets/query/sqlite/ssg/param.sql` con las consultas SQL para todas las operaciones CRUD de `Param`.
    *   **Capa de Servicio:**
        *   Se añadieron los métodos CRUD para `Param` a la interfaz `Service` y a la implementación `BaseService` en `internal/feat/ssg/service.go`.
    *   **API Handler y Router:**
        *   Se añadieron constantes `resParamName` y `resParamNameCap` a `internal/feat/ssg/apihandler.go`.
        *   Se implementaron los handlers de la API para `Param` (ej. `CreateParam`, `GetParam`, `ListParams`, `UpdateParam`, `DeleteParam`, `GetParamByName`, `GetParamByRefKey`) en `internal/feat/ssg/apihandler.go`.
        *   Se actualizó el método `wrapData` en `apihandler.go` para manejar los tipos `Param` y `[]Param`.
        *   Se registraron las rutas de la API para `Param` en `internal/feat/ssg/apirouter.go`.

3.  **Script de Prueba:**
    *   Se creó `scripts/curl/ssg/param.sh` para probar el ciclo CRUD completo de la entidad `Param` a través de la API.
    *   El script fue ejecutado con éxito, confirmando la funcionalidad de los endpoints.

**Próximos pasos:**
La entidad `Param` está lista para ser utilizada. El siguiente paso sería implementar la lógica en el sistema de configuración de la aplicación para leer estos valores dinámicos y aplicarlos, sobrescribiendo los valores por defecto de las variables de entorno cuando corresponda.
