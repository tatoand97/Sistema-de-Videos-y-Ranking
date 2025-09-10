video land VideoRank – Frontend
================================

Resumen
- React + Vite + TypeScript.
- Consumo de la API basada en la colección Postman incluida (auth, videos, rankings, ubicación).
- Se entrega como build estático servido por Nginx en el mismo `docker-compose` del backend (puerto 8081).

Variables de entorno
- `VITE_API_BASE_URL` (default `http://localhost:8080`)

Desarrollo local
1) Instalar dependencias: `npm install`
2) Ejecutar dev server: `npm run dev`
3) App: http://localhost:5173

Cómo levantar el Backend/API
- El frontend asume la API en `http://localhost:8080` (nginx → api).
- Desde la raíz del repo, levanta los servicios necesarios:

  1) Infra de soporte: `docker compose up -d postgres redis rabbitmq minio minio-buckets`
  2) Migraciones (one-off): `docker compose up --build migrate && docker compose rm -f migrate`
  3) API + Nginx: `docker compose up -d api nginx`

- Verificación: `curl http://localhost:8080/health` → `{ "status": "ok" }`

Build de producción
1) `npm run build` → genera `frontend/dist`
2) Con `docker-compose up -d nginx` la app queda servida en http://localhost:8081

Estructura de features principales
- Autenticación (login, registro, logout; token guardado en localStorage).
- Videos (listar mis videos, detalle, subir multipart, eliminar).
- Público (listar videos públicos, votar video).
- Rankings (paginado y filtro por ciudad).

Rutas
- `/` públicos + votar
- `/rankings` listados
- `/login`, `/register`
- `/profile` (protegida)
- `/videos` (protegida)
- `/upload` (protegida)
- `/videos/:id` (protegida)

Notas
- El backend define CORS abierto (`*`) en `docker-compose.yml`, por lo que el dev server funciona sin proxy adicional.
- Para el despliegue, sólo es necesario construir (`dist/`) y dejar que Nginx lo sirva en el puerto 8081.
