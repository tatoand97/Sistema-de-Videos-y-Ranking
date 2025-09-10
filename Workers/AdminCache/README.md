# AdminCache Worker

Cachea el ranking en Redis con TTL corto y refresco periódico. No consume colas ni usa RabbitMQ.

## Qué hace
- Ejecuta un warmup cada `REFRESH_INTERVAL_SECONDS` para calcular el ranking (global y por ciudad) desde Postgres y lo guarda paginado en Redis con TTL `CACHE_TTL_SECONDS`.
- No hay invalidación por eventos; el refresco es únicamente periódico.

## Claves Redis
- Global: `rank:global:page:{N}:size:{S}`
- Por ciudad: `rank:city:{slug}:page:{N}:size:{S}`
- Se almacena con prefijo `CACHE_PREFIX` (por defecto `videorank:`).

## Variables (.env)
Revisa `AdminCache/.env`:
```
# Redis
REDIS_ADDR=redis:6379
CACHE_PREFIX=videorank:
CACHE_TTL_SECONDS=120

# Postgres (ajusta al DSN de tu DB)
POSTGRES_DSN=postgres://app_user:app_password@postgres:5432/videorank?sslmode=disable

# Warmup
REFRESH_INTERVAL_SECONDS=60
WARM_PAGES=3
WARM_CITIES=bogota,medellin,cali
PAGE_SIZE_DEFAULT=20
```

> Si quieres usar el mismo Postgres/Redis de vote-system-docker, ajusta `POSTGRES_DSN`.

## Ejecutar
Con todo en este folder:
```bash
docker compose up -d --build admincache
docker compose logs -f admincache
```

## Ver Redis
```bash
docker exec -it redis redis-cli -- scan 0 match 'videorank:rank:*' count 100
docker exec -it redis redis-cli get 'videorank:rank:global:page:1:size:20' | jq .
```

## Notas
- Si quieres que el SQL sea parametrizable por env (como en vote-system), puedo agregar `RANK_SQL` y fallback a un default alineado a tu esquema.
- Para reflejo inmediato, se podría agregar recomputo eager o modo incremental con ZSET.

