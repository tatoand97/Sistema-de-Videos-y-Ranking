# Pruebas de API con Postman

## Archivos incluidos

- `VideoRank_API.postman_collection.json` - Colección unificada con endpoints y variables

## Importar en Postman

1. Abre Postman
2. Haz clic en "Import" 
3. Arrastra el archivo JSON o selecciónalo manualmente
4. Las variables están incluidas en la colección

## Endpoints disponibles

### 1. Health Check
- **GET** `/health`
- Verifica que la API esté funcionando

### 2. Register User
- **POST** `/api/auth/signup`
- Registra un nuevo usuario
- Guarda automáticamente el ID del usuario creado

### 3. Login User  
- **POST** `/api/auth/login`
- Autentica un usuario
- Guarda automáticamente el token JWT en las variables de entorno

### 4. Logout User
- **POST** `/api/auth/logout` 
- Requiere token de autenticación
- Cierra la sesión del usuario

## Flujo de pruebas recomendado

1. **Health Check** - Verificar que la API esté activa
2. **Register User** - Crear un nuevo usuario
3. **Login User** - Obtener token de autenticación  
4. **Logout User** - Cerrar sesión

## Variables de entorno

- `base_url`: URL base de la API (por defecto: http://localhost:8080)
- `auth_token`: Token JWT (se establece automáticamente al hacer login)
- `user_id`: ID del usuario (se establece automáticamente al registrarse)
- `user_email`: Email del usuario (se establece automáticamente al registrarse)

## Iniciar la API

Antes de probar, asegúrate de que la API esté ejecutándose:

```bash
docker compose up -d
```

La API estará disponible en http://localhost:8080