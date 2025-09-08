# Diagnóstico: Upload API sin publicación en cola

## Problema Identificado
- API responde exitoso ✅
- Registro en BD se crea ✅  
- No hay publicación en cola ❌
- No hay logs asociados ❌

## Puntos de Verificación

### 1. Variables de Entorno
Verificar que estas variables estén configuradas:
```bash
RABBITMQ_URL=amqp://admin:admin@rabbitmq:5672/
STATES_MACHINE_QUEUE=states_machine_queue
```

### 2. Conexión RabbitMQ
El publisher se inicializa en main.go pero puede fallar silenciosamente:
```go
if rabbitURL != "" {
    p, err := infraMessaging.NewRabbitMQPublisher(rabbitURL)
    if err != nil {
        log.Printf("warning: rabbitmq publisher init failed: %v", err)
        // ⚠️ Continúa sin publisher
    }
}
```

### 3. Publicación Asíncrona
La publicación se hace en goroutine y errores solo se imprimen:
```go
go func() {
    if publishErr := uc.publisher.Publish(uc.queue, b); publishErr != nil {
        fmt.Printf("Warning: Failed to publish message to queue %s: %v\n", uc.queue, publishErr)
    }
}()
```

## Comandos de Verificación

### 1. Verificar logs del contenedor API
```bash
docker logs app_api -f
```

### 2. Verificar estado RabbitMQ
```bash
# Acceder a RabbitMQ Management
http://localhost:8083
# Usuario: admin, Password: admin
```

### 3. Verificar cola states_machine_queue
```bash
docker exec -it sistema-de-videos-y-ranking-rabbitmq-1 rabbitmqctl list_queues
```

### 4. Verificar conexión desde API a RabbitMQ
```bash
docker exec -it app_api ping rabbitmq
```

## Soluciones Propuestas

### 1. Agregar Logs Detallados
Modificar uploads_use_case.go para agregar más logging:

```go
// Antes de publicar
log.Printf("Attempting to publish message for video ID: %d to queue: %s", video.VideoID, uc.queue)

// En la goroutine
go func() {
    log.Printf("Publishing message to queue %s: %s", uc.queue, string(b))
    if publishErr := uc.publisher.Publish(uc.queue, b); publishErr != nil {
        log.Printf("ERROR: Failed to publish message to queue %s: %v", uc.queue, publishErr)
    } else {
        log.Printf("SUCCESS: Message published to queue %s", uc.queue)
    }
}()
```

### 2. Verificar Publisher no es nil
```go
if uc.publisher == nil {
    log.Printf("WARNING: Publisher is nil, cannot publish message for video ID: %d", video.VideoID)
    return // o manejar según necesidad
}
```

### 3. Hacer Publicación Síncrona (temporal para debug)
```go
// Remover go func() para hacer síncrono temporalmente
if publishErr := uc.publisher.Publish(uc.queue, b); publishErr != nil {
    log.Printf("ERROR: Failed to publish message to queue %s: %v", uc.queue, publishErr)
    // Opcional: retornar error si es crítico
} else {
    log.Printf("SUCCESS: Message published to queue %s", uc.queue)
}
```

## Verificación Paso a Paso

1. **Verificar variables de entorno en contenedor**
2. **Verificar logs de inicialización del publisher**  
3. **Verificar logs durante upload**
4. **Verificar estado de colas en RabbitMQ**
5. **Verificar conectividad de red entre contenedores**