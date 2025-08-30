# Datos Seed - Sistema de Videos y Ranking

## Información General

Los datos seed se cargan automáticamente durante el `docker-compose up` a través del servicio `migrate`.

## Usuarios Creados

### Credenciales por Defecto
**Contraseña para todos los usuarios:** `password123`

### Usuarios del Sistema

| Email | Nombre | Rol | Descripción |
|-------|--------|-----|-------------|
| admin@videorank.com | Admin Sistema | admin | Administrador con acceso completo |
| moderador@videorank.com | Carlos Moderador | moderator | Moderador de contenido |
| viewer@videorank.com | Sofia Viewer | viewer | Usuario solo lectura |

### Jugadores (Players)

| Email | Nombre | Ciudad | Videos |
|-------|--------|--------|--------|
| juan.perez@email.com | Juan Pérez | Bogotá | 2 videos |
| maria.garcia@email.com | María García | Medellín | 1 video |
| pedro.lopez@email.com | Pedro López | Buenos Aires | 1 video |
| ana.martinez@email.com | Ana Martínez | Ciudad de México | 1 video |
| luis.rodriguez@email.com | Luis Rodríguez | Madrid | 1 video |
| carmen.fernandez@email.com | Carmen Fernández | São Paulo | 1 video |
| diego.gonzalez@email.com | Diego González | Santiago | 1 video |

## Estructura de Datos Cargados

### Países y Ciudades
- **8 países:** Colombia, Argentina, México, España, Estados Unidos, Brasil, Chile, Perú
- **15 ciudades** distribuidas en estos países

### Sistema de Roles y Permisos
- **4 roles:** admin, moderator, player, viewer
- **12 privilegios** diferentes asignados según el rol

### Videos de Ejemplo
- **8 videos** de baloncesto con diferentes estados de procesamiento ANB
- Videos distribuidos entre los jugadores
- Archivos procesados con marca de agua ANB y cortinillas
- Estados incluyen: trimming, removing_audio, adding_watermark, adding_intro_outro, processed

### Votos
- **32 votos** distribuidos entre los videos
- Sistema de votación funcional para testing

### Estados de Video
- **trimming:** Recortando duración a máximo 30 segundos
- **adjusting_resolution:** Ajustando resolución y formato de aspecto
- **adding_watermark:** Agregando marca de agua ANB
- **removing_audio:** Eliminando audio del video
- **adding_intro_outro:** Agregando cortinillas de apertura y cierre ANB
- **processed:** Video procesado exitosamente, listo para evaluación
- **failed:** Error en el procesamiento del video

## Hash de Contraseñas

Las contraseñas están hasheadas usando **bcrypt** con cost 10, compatible con el sistema de autenticación Go:

```go
bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
```

Hash utilizado: `$2a$10$dbG1mEbaLN9iCH9EmYflK.ddkyJ8aQhw52k5tkbwCzKoEktKOkiQ.`

## Ejecución

Los datos se cargan automáticamente cuando ejecutas:

```bash
docker-compose up
```

El servicio `migrate` ejecutará:
1. `0001_create_tables.up.sql` - Creación de tablas
2. `0002_seed_data.up.sql` - Carga de datos iniciales

## Rollback

Para hacer rollback de los datos seed:

```bash
docker-compose exec migrate /usr/local/bin/migrate -path=/migrations -database "postgres://app_user:app_password@postgres:5432/videorank?sslmode=disable" down 1
```

## Testing

Puedes usar cualquiera de las credenciales listadas para probar:
- Login con diferentes roles
- Funcionalidad de votación
- Gestión de videos
- Sistema de permisos

## Notas Importantes

- Todos los `INSERT` usan `ON CONFLICT DO NOTHING` para evitar duplicados
- Los datos son seguros para re-ejecutar
- Las relaciones están correctamente establecidas
- Los IDs se resuelven dinámicamente usando subconsultas