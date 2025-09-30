# AdminCache Worker

Cachea y publica los rankings Top-10 desde PostgreSQL hacia Redis siguiendo las reglas del Administrador de Cache (ACR): version de esquema, locks distribuidos, stale-while-revalidate y observabilidad por ciclo.

## Que hace
- Ejecuta un ciclo programado (`REFRESH_INTERVAL_SECONDS`, por defecto 300s) que lee PostgreSQL, normaliza y valida los Top-10 globales y por ciudad.
- Cada escritura es atomica: reemplaza el conjunto completo, incluye metadatos (`as_of`, `fresh_until`, `stale_until`, `schema_version`) y aplica TTL con `stale-while-revalidate` + jitter +/-10 %.
- Usa locks con lease (`CACHE_LOCK_LEASE_SECONDS`) para que un unico worker refresque cada clave a la vez. Si el refresco falla, se mantiene el dato **stale** hasta `CACHE_MAX_STALE_SECONDS`.
- Registra metricas via logs estructurados (exitos, errores, lock contention, uso de stale) por cada ciclo.

## Claves Redis (prefijo `CACHE_PREFIX`, default `videorank:`)
- Ranking global: `rank:global:{schema_version}`
- Ranking ciudad: `rank:city:{city_slug}:{schema_version}`
- Indice de ciudades activas: `rank:index:cities:{schema_version}`
- Locks: `rank:lock:global:{schema_version}` / `rank:lock:city:{city_slug}:{schema_version}`

Cada payload almacena hasta 10 elementos con `rank`, `user_id`, `username`, `score`, y los metadatos temporales necesarios para controlar `fresh`/`stale`.

## Variables de entorno (`AdminCache/.env`)
```
# Redis
REDIS_ADDR=redis:6379
CACHE_PREFIX=videorank:
SCHEMA_VERSION=v2
CACHE_TTL_FRESH_SECONDS=900
CACHE_MAX_STALE_SECONDS=600
CACHE_JITTER_PERCENT=10
CACHE_LOCK_LEASE_SECONDS=10
CACHE_MAX_TOP_USERS=10

# Postgres (ajusta al DSN de tu DB)
POSTGRES_DSN=postgres://app_user:app_password@postgres:5432/videorank?sslmode=disable
DB_READ_TIMEOUT_SECONDS=3
DB_MAX_RETRIES=3

# Warmup
REFRESH_INTERVAL_SECONDS=300
BATCH_SIZE_CITIES=50
# Lista opcional (nombres libres) -> se normaliza a slug
WARM_CITIES=bogota,medellin,cali
```

## Ejecutar
Con todo en este folder:
```bash
docker compose up -d --build admincache
docker compose logs -f admincache
```

## Ver Redis
```bash
docker exec -it redis redis-cli -- scan 0 match 'videorank:rank:*' count 100
docker exec -it redis redis-cli get 'videorank:rank:global:v2' | jq .
```

## Notas
- El worker prioriza disponibilidad: si la extraccion desde PostgreSQL falla, mantiene la version stale hasta que expire `CACHE_MAX_STALE_SECONDS`; al vencer, emite alerta (log `stale cache expired`).
- Para invalidacion dirigida, puedes eliminar las claves versionadas (`rank:city:{slug}:{schema_version}`) y dejar que el siguiente ciclo las regenere.
- Si necesitas exponer metricas a Prometheus u otro backend, conecta los logs estructurados al pipeline correspondiente.
