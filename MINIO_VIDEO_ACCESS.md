# Acceso Directo a Videos Procesados

## URLs de Acceso

Los videos procesados en el bucket `processed-videos` están disponibles mediante URLs directas:

```
http://localhost:8081/processed-videos/nombre-del-archivo.mp4
```

## Ejemplo de Uso

Si tienes un video procesado llamado `video-final.mp4` en el bucket `processed-videos`, puedes accederlo directamente en:

```
http://localhost:8081/processed-videos/video-final.mp4
```

## Características

- ✅ Streaming de video con soporte para `Accept-Ranges`
- ✅ Cache HTTP configurado (1 hora)
- ✅ Acceso directo sin autenticación
- ✅ Compatible con reproductores HTML5

## Puertos Configurados

- `8080`: API Principal
- `8081`: **Archivos MinIO (NUEVO)**
- `8082`: Consola MinIO
- `8083`: RabbitMQ Management

## Reiniciar Servicios

Después de los cambios, reinicia los contenedores:

```bash
docker-compose down
docker-compose up -d
```