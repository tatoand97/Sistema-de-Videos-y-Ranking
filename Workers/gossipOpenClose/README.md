
# gossipOpenClose Worker

Worker que agrega **cortinilla de apertura y cierre** con el logotipo oficial de la Asociación Nacional de Baloncesto (ANB/NBA) al video de entrada.
La suma de ambas cortinillas no excede **5 segundos** (por defecto 2.5 s cada una).

Arquitectura por capas (idéntica al worker de Watermarking):

```
gossipOpenClose/
├─ cmd/main.go
├─ internal/
│  ├─ adapters/        # MinIO, RabbitMQ, handler de mensajes
│  ├─ application/
│  │  ├─ services/     # ffmpeg: genera intro/outro y concatena
│  │  └─ usecases/     # OpenCloseUseCase
│  ├─ domain/          # interfaces y entidades
│  └─ infrastructure/  # container, config (env)
├─ assets/
│  └─ nba-logo-removebg-preview.png
├─ Dockerfile
├─ .env.example
└─ README.md
```

## Variables de entorno

- `RABBIT_URL`  (ej: `amqp://user:pass@rabbitmq:5672/`)
- `QUEUE_NAME`  (nombre de la cola de trabajo)
- `MINIO_ENDPOINT`, `MINIO_ACCESS_KEY`, `MINIO_SECRET_KEY`, `MINIO_SSL`
- `MINIO_BUCKET_RAW`       (bucket de entrada)
- `MINIO_BUCKET_PROCESSED` (bucket de salida)
- `MAX_SECONDS`            (reserva para compatibilidad, no se usa para recortar)
- `INTRO_SECONDS`          (float, defecto `2.5`)
- `OUTRO_SECONDS`          (float, defecto `2.5`)
- `TARGET_WIDTH`           (defecto `1280`)
- `TARGET_HEIGHT`          (defecto `720`)
- `FPS`                    (defecto `30`)

## Requisitos

El contenedor debe tener `ffmpeg` instalado (el `Dockerfile` ya lo incluye).
