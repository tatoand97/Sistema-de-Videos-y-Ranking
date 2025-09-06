# AudioRemoval Worker

## Descripción
Worker que remueve el audio de videos MP4 usando FFmpeg para garantizar compatibilidad y calidad.

## Características
- **Procesamiento**: FFmpeg con copia de stream de video (sin re-encoding)
- **Comando**: `ffmpeg -i input.mp4 -c:v copy -an -y output.mp4`
- **Rendimiento**: Procesamiento rápido sin re-encoding de video
- **Memoria**: ~50MB por instancia (incluye FFmpeg)
- **Tamaño imagen**: ~80MB (Alpine + FFmpeg)

## Construcción
```bash
docker build -t audioremoval .
```

## Variables de Entorno
- `RABBITMQ_URL`: URL de conexión a RabbitMQ
- `QUEUE_NAME`: Cola de entrada (default: audio_removal_queue)
- `RAW_BUCKET`: Bucket de videos de entrada
- `PROCESSED_BUCKET`: Bucket de videos procesados
- `STATE_MACHINE_QUEUE`: Cola para notificar al orquestador

## Limitaciones
- Solo soporta archivos MP4
- Requiere FFmpeg instalado en el contenedor

## Ventajas
- Procesamiento confiable con FFmpeg
- No re-encoding de video (solo copia)
- Escalabilidad horizontal
- Manejo robusto de formatos MP4