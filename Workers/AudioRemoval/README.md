# AudioRemoval Worker

Worker para eliminación de audio de videos usando arquitectura hexagonal.

## Arquitectura

- **Domain**: Lógica de negocio y entidades
- **Ports**: Interfaces para comunicación externa
- **Adapters**: Implementaciones concretas (RabbitMQ, MinIO, FFmpeg)

## Configuración

Variables de entorno en `.env`:
- `RABBITMQ_URL`: URL de conexión a RabbitMQ
- `MINIO_ENDPOINT`: Endpoint de MinIO
- `MINIO_ACCESS_KEY`: Clave de acceso MinIO
- `MINIO_SECRET_KEY`: Clave secreta MinIO
- `RAW_BUCKET`: Bucket de videos originales
- `PROCESSED_BUCKET`: Bucket de videos procesados
- `QUEUE_NAME`: Nombre de la cola RabbitMQ

## Uso

```bash
go run cmd/main.go
```

El worker escucha mensajes JSON con formato:
```json
{"filename": "video.mp4"}
```