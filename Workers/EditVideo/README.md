# EditVideo Worker

Worker que recorta videos MP4 a un máximo de `MAX_SECONDS` (por defecto 30 s),
usando la **misma arquitectura por capas** que `AudioRemoval`:

```
EditVideo/
├─ cmd/main.go
├─ internal/
│  ├─ adapters/ (MinIO, RabbitMQ, repos, handler)
│  ├─ application/
│  │  ├─ services/ (procesamiento mp4 con ffmpeg)
│  │  └─ usecases/ (ProcessVideoUseCase)
│  ├─ domain/ (entidades, puertos)
│  ├─ infrastructure/ (config y container)
│  └─ ports/ (interfaces Messaging/Storage)
└─ Dockerfile
```

## Variables
Ver `.env` (incluido). Reutiliza los nombres de `AudioRemoval` + `MAX_SECONDS`.

## Probar
1. Construir el servicio en docker-compose (añade el servicio `trim-video`).
2. Sube `sample.mp4` a `raw-videos-trim` en MinIO.
3. Publica en RabbitMQ (cola `trim_video_queue`) el JSON:
```json
{"filename":"sample.mp4"}
```
4. Verifica `processed-videos-trim/processed_sample.mp4`.
