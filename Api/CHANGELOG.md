# Changelog - Depuración del API

## Cambios Realizados

### ✅ Funcionalidades Mantenidas
- **Registro de usuario** (`POST /api/auth/signup`)
- **Login de usuario** (`POST /api/auth/login`) 
- **Logout de usuario** (`POST /api/auth/logout`)
- **Health check** (`GET /health`)
- Middleware JWT para autenticación
- Validación de contraseñas con bcrypt
- Manejo de tokens JWT con expiración

### ❌ Funcionalidades Eliminadas
- Sistema de videos y ranking
- Categorías y tareas
- Gestión de archivos estáticos
- Tablas de países, ciudades, roles, privilegios
- Sistema de votación
- Procesamiento de videos
- Workers y colas de trabajo

### 📁 Archivos Eliminados
- `internal1/` - Implementación alternativa
- `static/` - Archivos estáticos
- `internal/infrastructure/storage/` - Almacenamiento de archivos
- Migraciones complejas con múltiples tablas

### 📝 Archivos Modificados
- `cmd/api/main.go` - Simplificado para solo autenticación
- `internal/presentation/router.go` - Solo rutas de auth
- `internal/presentation/auth.go` - Handlers simplificados
- `internal/application/useCase/auth.go` - Lógica básica de auth
- `go.mod` - Dependencias mínimas necesarias
- `Dockerfile` - Imagen optimizada

### 📋 Archivos Nuevos
- `README.md` - Documentación de la API
- `.env.example` - Configuración de ejemplo
- `CHANGELOG.md` - Este archivo
- Migraciones simplificadas solo para usuarios

### 🗄️ Base de Datos
La base de datos ahora solo contiene:
- Tabla `users` con campos básicos (id, first_name, last_name, email, password_hash, created_at)
- Usuarios de prueba: admin@site.com y user@site.com

### 🔧 Configuración
Variables de entorno necesarias:
- `DATABASE_URL` - Conexión a PostgreSQL
- `JWT_SECRET` - Secreto para tokens JWT  
- `PORT` - Puerto del servidor (opcional, default: 8080)

El proyecto ahora es una API minimalista enfocada únicamente en autenticación de usuarios.