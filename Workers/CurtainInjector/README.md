# CurtainInjector Worker

Worker service que agrega metadatos de cortinillas de entrada y salida a videos MP4 usando manipulación binaria.

## Funcionalidad

- Procesa videos MP4 agregando metadatos para cortinillas de 2 segundos
- Agrega marcadores para cortinilla de entrada al inicio del video
- Agrega marcadores para cortinilla de salida al final del video
- Procesa videos desde MinIO storage
- Utiliza RabbitMQ para recibir mensajes de procesamiento
- **Solución ligera**: No requiere FFmpeg, solo manipulación binaria de MP4

## Arquitectura

Sigue el patrón de Clean Architecture con:
- **Domain**: Entidades y interfaces de negocio
- **Application**: Casos de uso y servicios de aplicación
- **Infrastructure**: Configuración y contenedor de dependencias
- **Adapters**: Implementaciones de interfaces externas (MinIO, RabbitMQ)
- **Ports**: Interfaces para servicios externos

## Configuración

Variables de entorno requeridas:
- `RABBITMQ_URL`: URL de conexión a RabbitMQ
- `MINIO_ENDPOINT`: Endpoint de MinIO
- `MINIO_ACCESS_KEY`: Clave de acceso MinIO
- `MINIO_SECRET_KEY`: Clave secreta MinIO
- `RAW_BUCKET`: Bucket de videos originales
- `PROCESSED_BUCKET`: Bucket de videos procesados
- `QUEUE_NAME`: Nombre de la cola RabbitMQ

## Dependencias

- MinIO para almacenamiento
- RabbitMQ para mensajería
- **Sin dependencias externas**: No requiere FFmpeg ni otras herramientas

## Limitaciones

- Restringido a archivos MP4 únicamente
- Agrega metadatos de cortinillas sin modificar el contenido visual real