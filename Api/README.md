# API de Autenticación

API simplificada en Go que proporciona funcionalidades básicas de registro de usuario y login.

## Funcionalidades

- **Registro de usuario**: Permite crear nuevos usuarios con email y contraseña
- **Login**: Autenticación de usuarios existentes con JWT
- **Logout**: Invalidación de tokens JWT

## Endpoints

### Salud del servicio
- `GET /health` - Verificar estado del servicio

### Autenticación
- `POST /api/auth/signup` - Registro de nuevo usuario
- `POST /api/auth/login` - Login de usuario
- `POST /api/auth/logout` - Logout de usuario (requiere token)

## Estructura del proyecto

```
internal/
├── application/
│   └── useCase/
│       └── auth.go          # Lógica de negocio de autenticación
├── domain/
│   ├── entities/
│   │   └── user.go          # Entidad Usuario
│   ├── interfaces/
│   │   └── user_repository.go # Interface del repositorio
│   └── errors.go            # Errores del dominio
├── infrastructure/
│   ├── migrations/          # Migraciones de base de datos
│   └── repository/
│       └── user.go          # Implementación del repositorio
└── presentation/
    ├── auth.go              # Handlers de autenticación
    └── router.go            # Configuración de rutas
```

## Configuración

Variables de entorno:
- `DATABASE_URL`: URL de conexión a PostgreSQL
- `JWT_SECRET`: Secreto para firmar tokens JWT
- `PORT`: Puerto del servidor (default: 8080)

## Ejecución

```bash
go run cmd/api/main.go
```