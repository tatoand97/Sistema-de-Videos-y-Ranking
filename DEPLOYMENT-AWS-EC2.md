# Despliegue Distribuido en AWS EC2

## Arquitectura de Despliegue

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   EC2-WORKERS   │    │    EC2-API      │    │  EC2-FRONTEND   │
│                 │    │                 │    │                 │
│ • RabbitMQ      │◄──►│ • API Service   │◄──►│ • Nginx         │
│ • All Workers   │    │ • Redis Cache   │    │ • Frontend      │
│ • States Machine│    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       │
┌─────────────────┐    ┌─────────────────┐              │
│   EC2-STORAGE   │    │  EC2-DATABASE   │              │
│                 │    │                 │              │
│ • MinIO         │    │ • PostgreSQL    │              │
│ • Object Store  │    │ • Migrations    │              │
└─────────────────┘    └─────────────────┘              │
         ▲                                               │
         └───────────────────────────────────────────────┘
```

## Requisitos Previos

- Cuenta AWS con acceso a EC2 y VPC
- 5 instancias t3.micro (Free Tier)
- VPC personalizada configurada
- Security Groups con principio de menor privilegio
- Key Pair para SSH

---

## 0. CONFIGURACIÓN DE VPC Y SEGURIDAD

### 0.1 Crear VPC Personalizada

```bash
# Crear VPC
aws ec2 create-vpc \
  --cidr-block 10.0.0.0/16 \
  --tag-specifications 'ResourceType=vpc,Tags=[{Key=Name,Value=VideoRank-VPC}]'

# Anotar VPC ID: vpc-xxxxxxxxx
VPC_ID=vpc-xxxxxxxxx
```

### 0.2 Crear Subnets

```bash
# Subnet Pública (Frontend)
aws ec2 create-subnet \
  --vpc-id $VPC_ID \
  --cidr-block 10.0.1.0/24 \
  --availability-zone us-east-1a \
  --tag-specifications 'ResourceType=subnet,Tags=[{Key=Name,Value=VideoRank-Public-Subnet}]'

# Subnet Privada - Aplicación (API, Workers)
aws ec2 create-subnet \
  --vpc-id $VPC_ID \
  --cidr-block 10.0.2.0/24 \
  --availability-zone us-east-1a \
  --tag-specifications 'ResourceType=subnet,Tags=[{Key=Name,Value=VideoRank-App-Subnet}]'

# Subnet Privada - Datos (Database, Storage)
aws ec2 create-subnet \
  --vpc-id $VPC_ID \
  --cidr-block 10.0.3.0/24 \
  --availability-zone us-east-1a \
  --tag-specifications 'ResourceType=subnet,Tags=[{Key=Name,Value=VideoRank-Data-Subnet}]'

# Anotar Subnet IDs
PUBLIC_SUBNET_ID=subnet-xxxxxxxxx
APP_SUBNET_ID=subnet-xxxxxxxxx
DATA_SUBNET_ID=subnet-xxxxxxxxx
```

### 0.3 Configurar Internet Gateway

```bash
# Crear Internet Gateway
aws ec2 create-internet-gateway \
  --tag-specifications 'ResourceType=internet-gateway,Tags=[{Key=Name,Value=VideoRank-IGW}]'

# Anotar IGW ID
IGW_ID=igw-xxxxxxxxx

# Asociar IGW a VPC
aws ec2 attach-internet-gateway \
  --internet-gateway-id $IGW_ID \
  --vpc-id $VPC_ID
```

### 0.4 Configurar NAT Gateway (para subnets privadas)

```bash
# Crear Elastic IP para NAT Gateway
aws ec2 allocate-address --domain vpc

# Anotar Allocation ID
ALLOC_ID=eipalloc-xxxxxxxxx

# Crear NAT Gateway en subnet pública
aws ec2 create-nat-gateway \
  --subnet-id $PUBLIC_SUBNET_ID \
  --allocation-id $ALLOC_ID \
  --tag-specifications 'ResourceType=nat-gateway,Tags=[{Key=Name,Value=VideoRank-NAT}]'

# Anotar NAT Gateway ID
NAT_GW_ID=nat-xxxxxxxxx
```

### 0.5 Configurar Route Tables

```bash
# Route Table para subnet pública
aws ec2 create-route-table \
  --vpc-id $VPC_ID \
  --tag-specifications 'ResourceType=route-table,Tags=[{Key=Name,Value=VideoRank-Public-RT}]'

PUBLIC_RT_ID=rtb-xxxxxxxxx

# Ruta a Internet Gateway
aws ec2 create-route \
  --route-table-id $PUBLIC_RT_ID \
  --destination-cidr-block 0.0.0.0/0 \
  --gateway-id $IGW_ID

# Asociar subnet pública
aws ec2 associate-route-table \
  --subnet-id $PUBLIC_SUBNET_ID \
  --route-table-id $PUBLIC_RT_ID

# Route Table para subnets privadas
aws ec2 create-route-table \
  --vpc-id $VPC_ID \
  --tag-specifications 'ResourceType=route-table,Tags=[{Key=Name,Value=VideoRank-Private-RT}]'

PRIVATE_RT_ID=rtb-xxxxxxxxx

# Ruta a NAT Gateway
aws ec2 create-route \
  --route-table-id $PRIVATE_RT_ID \
  --destination-cidr-block 0.0.0.0/0 \
  --nat-gateway-id $NAT_GW_ID

# Asociar subnets privadas
aws ec2 associate-route-table \
  --subnet-id $APP_SUBNET_ID \
  --route-table-id $PRIVATE_RT_ID

aws ec2 associate-route-table \
  --subnet-id $DATA_SUBNET_ID \
  --route-table-id $PRIVATE_RT_ID
```

### 0.6 Crear Security Groups con Principio de Menor Privilegio

#### SG-SSH (Acceso SSH desde tu IP)
```bash
aws ec2 create-security-group \
  --group-name VideoRank-SSH-SG \
  --description "SSH access from admin IP" \
  --vpc-id $VPC_ID

SSH_SG_ID=sg-xxxxxxxxx

# Permitir SSH solo desde tu IP
aws ec2 authorize-security-group-ingress \
  --group-id $SSH_SG_ID \
  --protocol tcp \
  --port 22 \
  --cidr YOUR_PUBLIC_IP/32
```

#### SG-Frontend (Acceso público HTTP/HTTPS)
```bash
aws ec2 create-security-group \
  --group-name VideoRank-Frontend-SG \
  --description "Public HTTP access" \
  --vpc-id $VPC_ID

FRONTEND_SG_ID=sg-xxxxxxxxx

# HTTP público
aws ec2 authorize-security-group-ingress \
  --group-id $FRONTEND_SG_ID \
  --protocol tcp \
  --port 80 \
  --cidr 0.0.0.0/0

# Frontend alternativo
aws ec2 authorize-security-group-ingress \
  --group-id $FRONTEND_SG_ID \
  --protocol tcp \
  --port 8081 \
  --cidr 0.0.0.0/0
```

#### SG-API (Acceso desde Frontend y público)
```bash
aws ec2 create-security-group \
  --group-name VideoRank-API-SG \
  --description "API access" \
  --vpc-id $VPC_ID

API_SG_ID=sg-xxxxxxxxx

# API público
aws ec2 authorize-security-group-ingress \
  --group-id $API_SG_ID \
  --protocol tcp \
  --port 8080 \
  --cidr 0.0.0.0/0

# Redis solo desde Workers
aws ec2 authorize-security-group-ingress \
  --group-id $API_SG_ID \
  --protocol tcp \
  --port 6379 \
  --source-group $WORKERS_SG_ID
```

#### SG-Workers (Acceso desde API)
```bash
aws ec2 create-security-group \
  --group-name VideoRank-Workers-SG \
  --description "Workers and RabbitMQ access" \
  --vpc-id $VPC_ID

WORKERS_SG_ID=sg-xxxxxxxxx

# RabbitMQ desde API
aws ec2 authorize-security-group-ingress \
  --group-id $WORKERS_SG_ID \
  --protocol tcp \
  --port 5672 \
  --source-group $API_SG_ID

# RabbitMQ UI desde tu IP
aws ec2 authorize-security-group-ingress \
  --group-id $WORKERS_SG_ID \
  --protocol tcp \
  --port 15672 \
  --cidr YOUR_PUBLIC_IP/32
```

#### SG-Storage (Acceso desde API y Workers)
```bash
aws ec2 create-security-group \
  --group-name VideoRank-Storage-SG \
  --description "MinIO storage access" \
  --vpc-id $VPC_ID

STORAGE_SG_ID=sg-xxxxxxxxx

# MinIO API desde API y Workers
aws ec2 authorize-security-group-ingress \
  --group-id $STORAGE_SG_ID \
  --protocol tcp \
  --port 9000 \
  --source-group $API_SG_ID

aws ec2 authorize-security-group-ingress \
  --group-id $STORAGE_SG_ID \
  --protocol tcp \
  --port 9000 \
  --source-group $WORKERS_SG_ID

# MinIO Console desde tu IP
aws ec2 authorize-security-group-ingress \
  --group-id $STORAGE_SG_ID \
  --protocol tcp \
  --port 9001 \
  --cidr YOUR_PUBLIC_IP/32
```

#### SG-Database (Acceso solo desde API y Workers)
```bash
aws ec2 create-security-group \
  --group-name VideoRank-Database-SG \
  --description "Database access" \
  --vpc-id $VPC_ID

DATABASE_SG_ID=sg-xxxxxxxxx

# PostgreSQL desde API
aws ec2 authorize-security-group-ingress \
  --group-id $DATABASE_SG_ID \
  --protocol tcp \
  --port 5432 \
  --source-group $API_SG_ID

# PostgreSQL desde Workers
aws ec2 authorize-security-group-ingress \
  --group-id $DATABASE_SG_ID \
  --protocol tcp \
  --port 5432 \
  --source-group $WORKERS_SG_ID
```

### 0.7 Matriz de Comunicación (Principio de Menor Privilegio)

| Origen | Destino | Puerto | Protocolo | Justificación |
|--------|---------|--------|-----------|---------------|
| Internet | Frontend | 80, 8081 | TCP | Acceso público web |
| Internet | API | 8080 | TCP | API REST pública |
| Frontend | API | 8080 | TCP | Llamadas API |
| API | Database | 5432 | TCP | Consultas DB |
| API | Storage | 9000 | TCP | Subida/descarga archivos |
| API | Workers | 5672 | TCP | Envío mensajes RabbitMQ |
| Workers | Database | 5432 | TCP | Actualización estados |
| Workers | Storage | 9000 | TCP | Procesamiento archivos |
| Workers | API | 6379 | TCP | Cache Redis |
| Admin IP | All | 22 | TCP | Administración SSH |
| Admin IP | Workers | 15672 | TCP | RabbitMQ UI |
| Admin IP | Storage | 9001 | TCP | MinIO Console |

**Comunicaciones NO permitidas:**
- Frontend ↔ Database (directo)
- Frontend ↔ Storage (directo)
- Frontend ↔ Workers (directo)
- Storage ↔ Database (directo)
- Cualquier acceso no listado arriba

---

## 1. EC2-DATABASE (PostgreSQL + Migrations)

### 1.1 Crear instancia EC2
```bash
# Crear instancia en subnet de datos
aws ec2 run-instances \
  --image-id ami-0abcdef1234567890 \
  --count 1 \
  --instance-type t3.micro \
  --key-name YOUR_KEY_PAIR \
  --security-group-ids $SSH_SG_ID $DATABASE_SG_ID \
  --subnet-id $DATA_SUBNET_ID \
  --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=VideoRank-Database}]'

# Anotar Instance ID y IP privada
DATABASE_INSTANCE_ID=i-xxxxxxxxx
DATABASE_PRIVATE_IP=10.0.3.x
```

### 1.2 Instalar Docker
```bash
sudo yum update -y
sudo yum install -y docker git
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -a -G docker ec2-user
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
exit
```

### 1.3 Clonar proyecto y preparar migraciones
```bash
git clone <tu-repositorio>
cd Sistema-de-Videos-y-Ranking
mkdir -p database/migrations
cp Api/internal/infrastructure/migrations/* database/migrations/
```

### 1.4 Crear docker-compose-db.yml
```yaml
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: ${DB_USER:-app_user}
      POSTGRES_PASSWORD: ${DB_PASS:-secure_password_123}
      POSTGRES_DB: ${DB_NAME:-videorank}
      LC_ALL: C.UTF-8
      LANG: C.UTF-8
    ports:
      - "5432:5432"
    volumes:
      - pg-data:/var/lib/postgresql/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL","pg_isready -U $$POSTGRES_USER -d $$POSTGRES_DB"]
      interval: 10s
      timeout: 5s
      retries: 5

  migrate:
    image: migrate/migrate:4
    depends_on:
      postgres:
        condition: service_healthy
    volumes:
      - ./database/migrations:/migrations:ro
    environment:
      DATABASE_URL: postgres://${DB_USER:-app_user}:${DB_PASS:-secure_password_123}@postgres:5432/${DB_NAME:-videorank}?sslmode=disable
    entrypoint: ["/bin/sh","-c"]
    command: |
      echo "Running migrations..." \
      && /usr/local/bin/migrate -path=/migrations -database "$$DATABASE_URL" up
    restart: "no"

volumes:
  pg-data:
```

### 1.5 Variables de entorno (.env)
```bash
DB_USER=app_user
DB_PASS=secure_password_123
DB_NAME=videorank
```

### 1.6 Desplegar
```bash
docker-compose -f docker-compose-db.yml up -d
```

---

## 2. EC2-STORAGE (MinIO)

### 2.1 Crear instancia EC2
```bash
# Crear instancia en subnet de datos
aws ec2 run-instances \
  --image-id ami-0abcdef1234567890 \
  --count 1 \
  --instance-type t3.micro \
  --key-name YOUR_KEY_PAIR \
  --security-group-ids $SSH_SG_ID $STORAGE_SG_ID \
  --subnet-id $DATA_SUBNET_ID \
  --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=VideoRank-Storage}]'

# Anotar Instance ID y IP privada
STORAGE_INSTANCE_ID=i-xxxxxxxxx
STORAGE_PRIVATE_IP=10.0.3.x
```

### 2.2 Instalar Docker (mismo proceso que 1.2)

### 2.3 Crear docker-compose-minio.yml
```yaml
services:
  minio:
    image: minio/minio:latest
    command: ["server","/data","--console-address",":9001"]
    environment:
      MINIO_ROOT_USER: ${MINIO_USER:-minio}
      MINIO_ROOT_PASSWORD: ${MINIO_PASS:-minio12345}
    ports:
      - "9000:9000"
      - "9001:9001"
    volumes:
      - minio-data:/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 10s
      retries: 3

  minio-buckets:
    image: minio/mc
    depends_on:
      minio:
        condition: service_healthy
    environment:
      MINIO_ROOT_USER: ${MINIO_USER:-minio}
      MINIO_ROOT_PASSWORD: ${MINIO_PASS:-minio12345}
    entrypoint: >
      /bin/sh -c "
      until /usr/bin/mc alias set myminio http://minio:9000 $$MINIO_ROOT_USER $$MINIO_ROOT_PASSWORD; do
        echo 'Waiting for MinIO...'
        sleep 2
      done;
      /usr/bin/mc mb myminio/raw-videos --ignore-existing;
      /usr/bin/mc mb myminio/processed-videos-trim --ignore-existing;
      /usr/bin/mc mb myminio/processed-videos-edit --ignore-existing;
      /usr/bin/mc mb myminio/processed-videos-audio-removal --ignore-existing;
      /usr/bin/mc mb myminio/processed-videos-watermarking --ignore-existing;
      /usr/bin/mc mb myminio/processed-videos --ignore-existing;
      /usr/bin/mc anonymous set public myminio/processed-videos;
      echo 'Buckets created';
      exit 0;
      "
    restart: "no"

volumes:
  minio-data:
```

### 2.4 Variables de entorno (.env)
```bash
MINIO_USER=minio
MINIO_PASS=minio12345
```

### 2.5 Desplegar
```bash
docker-compose -f docker-compose-minio.yml up -d
```

---

## 3. EC2-WORKERS (Workers + RabbitMQ)

### 3.1 Crear instancia EC2
```bash
# Crear instancia en subnet de aplicación
aws ec2 run-instances \
  --image-id ami-0abcdef1234567890 \
  --count 1 \
  --instance-type t3.micro \
  --key-name YOUR_KEY_PAIR \
  --security-group-ids $SSH_SG_ID $WORKERS_SG_ID \
  --subnet-id $APP_SUBNET_ID \
  --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=VideoRank-Workers}]'

# Anotar Instance ID y IP privada
WORKERS_INSTANCE_ID=i-xxxxxxxxx
WORKERS_PRIVATE_IP=10.0.2.x
```

### 3.2 Instalar Docker (mismo proceso que 1.2)

### 3.3 Clonar proyecto
```bash
git clone <tu-repositorio>
cd Sistema-de-Videos-y-Ranking
```

### 3.4 Crear docker-compose-workers.yml
```yaml
services:
  rabbitmq:
    image: rabbitmq:3.13-management-alpine
    environment:
      RABBITMQ_DEFAULT_USER: admin
      RABBITMQ_DEFAULT_PASS: ${RABBITMQ_PASS:-admin123}
    ports:
      - "5672:5672"
      - "15672:15672"
    volumes:
      - rabbitmq-data:/var/lib/rabbitmq
      - ./rabbitmq/enabled_plugins:/etc/rabbitmq/enabled_plugins:ro
      - ./rabbitmq/rabbitmq.conf:/etc/rabbitmq/rabbitmq.conf:ro
      - ./rabbitmq/definitions.json:/etc/rabbitmq/definitions.json:ro
    restart: unless-stopped
    healthcheck:
      test: ["CMD","rabbitmq-diagnostics","-q","ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  trim-video:
    build:
      context: ./Workers
      dockerfile: TrimVideo/Dockerfile
    depends_on:
      rabbitmq:
        condition: service_healthy
    environment:
      RABBITMQ_URL: amqp://admin:${RABBITMQ_PASS:-admin123}@rabbitmq:5672/
      MINIO_ENDPOINT: ${MINIO_HOST}:9000
      MINIO_ACCESS_KEY: ${MINIO_USER:-minio}
      MINIO_SECRET_KEY: ${MINIO_PASS:-minio12345}
      RAW_BUCKET: raw-videos
      PROCESSED_BUCKET: processed-videos-trim
      QUEUE_NAME: trim_video_queue
      STATE_MACHINE_QUEUE: states_machine_queue
      MAX_SECONDS: "30"
      MAX_RETRIES: "3"
    restart: unless-stopped

  edit-video:
    build:
      context: ./Workers
      dockerfile: EditVideo/Dockerfile
    depends_on:
      rabbitmq:
        condition: service_healthy
    environment:
      RABBITMQ_URL: amqp://admin:${RABBITMQ_PASS:-admin123}@rabbitmq:5672/
      MINIO_ENDPOINT: ${MINIO_HOST}:9000
      MINIO_ACCESS_KEY: ${MINIO_USER:-minio}
      MINIO_SECRET_KEY: ${MINIO_PASS:-minio12345}
      RAW_BUCKET: processed-videos-trim
      PROCESSED_BUCKET: processed-videos-edit
      QUEUE_NAME: edit_video_queue
      STATE_MACHINE_QUEUE: states_machine_queue
      MAX_SECONDS: "30"
      MAX_RETRIES: "3"
    restart: unless-stopped

  audio-removal:
    build:
      context: ./Workers
      dockerfile: AudioRemoval/Dockerfile
    depends_on:
      rabbitmq:
        condition: service_healthy
    environment:
      RABBITMQ_URL: amqp://admin:${RABBITMQ_PASS:-admin123}@rabbitmq:5672/
      MINIO_ENDPOINT: ${MINIO_HOST}:9000
      MINIO_ACCESS_KEY: ${MINIO_USER:-minio}
      MINIO_SECRET_KEY: ${MINIO_PASS:-minio12345}
      RAW_BUCKET: processed-videos-edit
      PROCESSED_BUCKET: processed-videos-audio-removal
      QUEUE_NAME: audio_removal_queue
      STATE_MACHINE_QUEUE: states_machine_queue
      MAX_RETRIES: "3"
    restart: unless-stopped

  watermarking:
    build:
      context: ./Workers
      dockerfile: Watermarking/Dockerfile
    depends_on:
      rabbitmq:
        condition: service_healthy
    environment:
      RABBITMQ_URL: amqp://admin:${RABBITMQ_PASS:-admin123}@rabbitmq:5672/
      MINIO_ENDPOINT: ${MINIO_HOST}:9000
      MINIO_ACCESS_KEY: ${MINIO_USER:-minio}
      MINIO_SECRET_KEY: ${MINIO_PASS:-minio12345}
      RAW_BUCKET: processed-videos-audio-removal
      PROCESSED_BUCKET: processed-videos-watermarking
      QUEUE_NAME: watermarking_queue
      STATE_MACHINE_QUEUE: states_machine_queue
      MAX_SECONDS: "30"
      MAX_RETRIES: "3"
    restart: unless-stopped

  gossip-open-close:
    build:
      context: ./Workers
      dockerfile: gossipOpenClose/Dockerfile
    depends_on:
      rabbitmq:
        condition: service_healthy
    environment:
      RABBITMQ_URL: amqp://admin:${RABBITMQ_PASS:-admin123}@rabbitmq:5672/
      MINIO_ENDPOINT: ${MINIO_HOST}:9000
      MINIO_ACCESS_KEY: ${MINIO_USER:-minio}
      MINIO_SECRET_KEY: ${MINIO_PASS:-minio12345}
      MINIO_BUCKET_RAW: processed-videos-watermarking
      MINIO_BUCKET_PROCESSED: processed-videos
      QUEUE_NAME: gossip_open_close_queue
      MAX_SECONDS: "30"
      INTRO_SECONDS: "2.5"
      OUTRO_SECONDS: "2.5"
      TARGET_WIDTH: "1280"
      TARGET_HEIGHT: "720"
      FPS: "30"
      LOGO_PATH: ./assets/nba-logo-removebg-preview.png
      MAX_RETRIES: "3"
    restart: unless-stopped

  admin-cache:
    build:
      context: ./Workers/AdminCache
    depends_on:
      rabbitmq:
        condition: service_healthy
    environment:
      DATABASE_URL: postgres://${DB_USER:-app_user}:${DB_PASS:-secure_password_123}@${DB_HOST}:5432/${DB_NAME:-videorank}?sslmode=disable
      RABBITMQ_URL: amqp://admin:${RABBITMQ_PASS:-admin123}@rabbitmq:5672/
      REDIS_URL: redis://${REDIS_HOST}:6379
    restart: unless-stopped

  states-machine:
    build:
      context: ./Workers
      dockerfile: StatesMachine/Dockerfile
    depends_on:
      rabbitmq:
        condition: service_healthy
    environment:
      RABBITMQ_URL: amqp://admin:${RABBITMQ_PASS:-admin123}@rabbitmq:5672/
      QUEUE_NAME: states_machine_queue
      EDIT_VIDEO_QUEUE: edit_video_queue
      AUDIO_REMOVAL_QUEUE: audio_removal_queue
      WATERMARKING_QUEUE: watermarking_queue
      DATABASE_URL: postgres://${DB_USER:-app_user}:${DB_PASS:-secure_password_123}@${DB_HOST}:5432/${DB_NAME:-videorank}?sslmode=disable
      MAX_RETRIES: "3"
      RETRY_DELAY_MINUTES: "5"
    restart: unless-stopped

volumes:
  rabbitmq-data:
```

### 3.5 Variables de entorno (.env)
```bash
# IPs privadas de otras instancias (usar las IPs reales asignadas)
DB_HOST=$DATABASE_PRIVATE_IP
REDIS_HOST=$API_PRIVATE_IP
MINIO_HOST=$STORAGE_PRIVATE_IP

# Credenciales
DB_USER=app_user
DB_PASS=secure_password_123
DB_NAME=videorank
RABBITMQ_PASS=admin123
MINIO_USER=minio
MINIO_PASS=minio12345
```

### 3.6 Desplegar
```bash
docker-compose -f docker-compose-workers.yml up -d
```

---

## 4. EC2-API (API + Redis)

### 4.1 Crear instancia EC2
```bash
# Crear instancia en subnet de aplicación
aws ec2 run-instances \
  --image-id ami-0abcdef1234567890 \
  --count 1 \
  --instance-type t3.micro \
  --key-name YOUR_KEY_PAIR \
  --security-group-ids $SSH_SG_ID $API_SG_ID \
  --subnet-id $APP_SUBNET_ID \
  --associate-public-ip-address \
  --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=VideoRank-API}]'

# Anotar Instance ID, IP privada y pública
API_INSTANCE_ID=i-xxxxxxxxx
API_PRIVATE_IP=10.0.2.x
API_PUBLIC_IP=x.x.x.x
```

### 4.2 Instalar Docker (mismo proceso que 1.2)

### 4.3 Clonar proyecto
```bash
git clone <tu-repositorio>
cd Sistema-de-Videos-y-Ranking
```

### 4.4 Crear docker-compose-api.yml
```yaml
services:
  redis:
    image: redis:7-alpine
    command: ["redis-server","--appendonly","yes"]
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD","redis-cli","ping"]
      interval: 10s
      timeout: 3s
      retries: 5

  api:
    build:
      context: ./Api
    depends_on:
      redis:
        condition: service_healthy
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://${DB_USER:-app_user}:${DB_PASS:-secure_password_123}@${DB_HOST}:5432/${DB_NAME:-videorank}?sslmode=disable
      JWT_SECRET: ${JWT_SECRET:-6EPOc2/6veZ71FAJSF68iv21ho83NLQSaycHEGTTjGO9TBmRsphMp5JqgieFTcGn}
      PORT: "8080"
      CORS_ORIGIN: "*"
      REDIS_ADDR: redis:6379
      CACHE_PREFIX: "videorank:"
      CACHE_TTL_SECONDS: "120"
      MINIO_ENDPOINT: ${MINIO_HOST}:9000
      MINIO_ACCESS_KEY: ${MINIO_USER:-minio}
      MINIO_SECRET_KEY: ${MINIO_PASS:-minio12345}
      MINIO_BUCKET: raw-videos
      MINIO_USE_SSL: "false"
      RABBITMQ_URL: amqp://admin:${RABBITMQ_PASS:-admin123}@${RABBITMQ_HOST}:5672/
      STATES_MACHINE_QUEUE: states_machine_queue
      RABBITMQ_QUEUE_MAXLEN: "1000"
    restart: unless-stopped
    healthcheck:
      test: ["CMD","wget","--no-verbose","--tries=1","--spider","http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  redis-data:
```

### 4.5 Variables de entorno (.env)
```bash
# IPs privadas (usar las IPs reales asignadas)
DB_HOST=$DATABASE_PRIVATE_IP
RABBITMQ_HOST=$WORKERS_PRIVATE_IP
MINIO_HOST=$STORAGE_PRIVATE_IP

# Credenciales
DB_USER=app_user
DB_PASS=secure_password_123
DB_NAME=videorank
RABBITMQ_PASS=admin123
MINIO_USER=minio
MINIO_PASS=minio12345
JWT_SECRET=your_super_secure_jwt_secret_key_here
```

### 4.6 Desplegar
```bash
docker-compose -f docker-compose-api.yml up -d
```

---

## 5. EC2-FRONTEND (Frontend + Nginx)

### 5.1 Crear instancia EC2
```bash
# Crear instancia en subnet pública
aws ec2 run-instances \
  --image-id ami-0abcdef1234567890 \
  --count 1 \
  --instance-type t3.micro \
  --key-name YOUR_KEY_PAIR \
  --security-group-ids $SSH_SG_ID $FRONTEND_SG_ID \
  --subnet-id $PUBLIC_SUBNET_ID \
  --associate-public-ip-address \
  --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=VideoRank-Frontend}]'

# Anotar Instance ID, IP privada y pública
FRONTEND_INSTANCE_ID=i-xxxxxxxxx
FRONTEND_PRIVATE_IP=10.0.1.x
FRONTEND_PUBLIC_IP=x.x.x.x
```

### 5.2 Instalar Docker y Node.js
```bash
sudo yum update -y
sudo yum install -y docker git nodejs npm
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -a -G docker ec2-user
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
exit
```

### 5.3 Construir frontend
```bash
git clone <tu-repositorio>
cd Sistema-de-Videos-y-Ranking/frontend
npm install
VITE_API_BASE_URL=http://API_PUBLIC_IP:8080 npm run build
cd ..
```

### 5.4 Crear nginx-frontend.conf
```nginx
events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    server {
        listen 80;
        server_name _;
        root /usr/share/nginx/html;
        index index.html;

        location / {
            try_files $uri $uri/ /index.html;
        }

        location /api/ {
            proxy_pass http://API_PUBLIC_IP:8080/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        }
    }

    server {
        listen 8081;
        server_name _;
        root /usr/share/nginx/html;
        index index.html;

        location / {
            try_files $uri $uri/ /index.html;
        }
    }
}
```

### 5.5 Crear docker-compose-frontend.yml
```yaml
services:
  nginx:
    image: nginx:1.27-alpine
    ports:
      - "80:80"
      - "8081:8081"
    volumes:
      - ./nginx-frontend.conf:/etc/nginx/nginx.conf:ro
      - ./frontend/dist:/usr/share/nginx/html:ro
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:80"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### 5.6 Desplegar
```bash
# Reemplazar API_PUBLIC_IP en nginx-frontend.conf
sed -i "s/API_PUBLIC_IP/$API_PUBLIC_IP/g" nginx-frontend.conf
docker-compose -f docker-compose-frontend.yml up -d
```

---

## Orden de Despliegue

1. **EC2-DATABASE**: Iniciar primero
2. **EC2-STORAGE**: Iniciar segundo
3. **EC2-WORKERS**: Iniciar tercero (necesita DB y Storage)
4. **EC2-API**: Iniciar cuarto (necesita DB, Storage, Workers)
5. **EC2-FRONTEND**: Iniciar último (necesita API)

## Asignación de Security Groups por EC2

### EC2-DATABASE (Subnet: Data - 10.0.3.0/24)
- Security Groups: SSH-SG + Database-SG
- Acceso permitido:
  - SSH desde Admin IP
  - PostgreSQL (5432) desde API-SG y Workers-SG únicamente

### EC2-STORAGE (Subnet: Data - 10.0.3.0/24)
- Security Groups: SSH-SG + Storage-SG
- Acceso permitido:
  - SSH desde Admin IP
  - MinIO API (9000) desde API-SG y Workers-SG únicamente
  - MinIO Console (9001) desde Admin IP

### EC2-WORKERS (Subnet: App - 10.0.2.0/24)
- Security Groups: SSH-SG + Workers-SG
- Acceso permitido:
  - SSH desde Admin IP
  - RabbitMQ (5672) desde API-SG únicamente
  - RabbitMQ UI (15672) desde Admin IP

### EC2-API (Subnet: App - 10.0.2.0/24)
- Security Groups: SSH-SG + API-SG
- Acceso permitido:
  - SSH desde Admin IP
  - API (8080) desde Internet (0.0.0.0/0)
  - Redis (6379) desde Workers-SG únicamente

### EC2-FRONTEND (Subnet: Public - 10.0.1.0/24)
- Security Groups: SSH-SG + Frontend-SG
- Acceso permitido:
  - SSH desde Admin IP
  - HTTP (80, 8081) desde Internet (0.0.0.0/0)

## Comandos de Verificación

```bash
# Estado de servicios
docker-compose ps

# Logs en tiempo real
docker-compose logs -f [service]

# Recursos del sistema
htop
df -h

# Conectividad entre servicios
telnet <ip-privada> <puerto>

# Health checks
curl http://localhost:8080/health  # API
curl http://localhost:9000/minio/health/live  # MinIO
```

## URLs de Acceso

- **Frontend**: http://FRONTEND_PUBLIC_IP/
- **API**: http://API_PUBLIC_IP:8080/
- **RabbitMQ UI**: http://WORKERS_PUBLIC_IP:15672/
- **MinIO Console**: http://STORAGE_PUBLIC_IP:9001/

## Costos Estimados (Free Tier)

- 5 x t3.micro: $0 (12 meses)
- 5 x 30GB EBS: $0 (30GB gratuitos)
- Transferencia de datos: Mínima dentro de VPC
- **Total**: $0 por 12 meses

---

## BUENAS PRÁCTICAS ADICIONALES (FREE TIER)

### A. API Gateway (RECOMENDADO)

#### A.1 Crear API Gateway REST
```bash
# Crear API Gateway
aws apigateway create-rest-api \
  --name VideoRank-API \
  --description "VideoRank API Gateway" \
  --endpoint-configuration types=REGIONAL

# Anotar API ID
API_GW_ID=xxxxxxxxxx

# Obtener root resource ID
aws apigateway get-resources --rest-api-id $API_GW_ID
ROOT_RESOURCE_ID=xxxxxxxxxx
```

#### A.2 Configurar Proxy Resource
```bash
# Crear resource proxy
aws apigateway create-resource \
  --rest-api-id $API_GW_ID \
  --parent-id $ROOT_RESOURCE_ID \
  --path-part "{proxy+}"

PROXY_RESOURCE_ID=xxxxxxxxxx

# Crear método ANY
aws apigateway put-method \
  --rest-api-id $API_GW_ID \
  --resource-id $PROXY_RESOURCE_ID \
  --http-method ANY \
  --authorization-type NONE

# Configurar integración con EC2
aws apigateway put-integration \
  --rest-api-id $API_GW_ID \
  --resource-id $PROXY_RESOURCE_ID \
  --http-method ANY \
  --type HTTP_PROXY \
  --integration-http-method ANY \
  --uri "http://$API_PRIVATE_IP:8080/{proxy}"
```

#### A.3 Configurar Rate Limiting y Throttling
```bash
# Crear usage plan (1000 requests/day, 10 req/sec)
aws apigateway create-usage-plan \
  --name VideoRank-Basic \
  --description "Basic usage plan" \
  --throttle burstLimit=20,rateLimit=10 \
  --quota limit=1000,period=DAY

USAGE_PLAN_ID=xxxxxxxxxx

# Crear API Key
aws apigateway create-api-key \
  --name VideoRank-Key \
  --description "API Key for VideoRank"

API_KEY_ID=xxxxxxxxxx

# Asociar API Key con Usage Plan
aws apigateway create-usage-plan-key \
  --usage-plan-id $USAGE_PLAN_ID \
  --key-id $API_KEY_ID \
  --key-type API_KEY
```

#### A.4 Desplegar API Gateway
```bash
# Crear deployment
aws apigateway create-deployment \
  --rest-api-id $API_GW_ID \
  --stage-name prod

# URL final: https://$API_GW_ID.execute-api.us-east-1.amazonaws.com/prod
```

### B. CloudFront CDN (RECOMENDADO)

#### B.1 Crear distribución para Frontend
```bash
# Crear CloudFront distribution
aws cloudfront create-distribution \
  --distribution-config '{
    "CallerReference": "VideoRank-Frontend-'$(date +%s)'",
    "Comment": "VideoRank Frontend CDN",
    "DefaultCacheBehavior": {
      "TargetOriginId": "VideoRank-Frontend",
      "ViewerProtocolPolicy": "redirect-to-https",
      "TrustedSigners": {
        "Enabled": false,
        "Quantity": 0
      },
      "ForwardedValues": {
        "QueryString": false,
        "Cookies": {
          "Forward": "none"
        }
      },
      "MinTTL": 0,
      "DefaultTTL": 86400,
      "MaxTTL": 31536000
    },
    "Origins": {
      "Quantity": 1,
      "Items": [
        {
          "Id": "VideoRank-Frontend",
          "DomainName": "'$FRONTEND_PUBLIC_IP'",
          "CustomOriginConfig": {
            "HTTPPort": 80,
            "HTTPSPort": 443,
            "OriginProtocolPolicy": "http-only"
          }
        }
      ]
    },
    "Enabled": true,
    "PriceClass": "PriceClass_100"
  }'

# Anotar CloudFront Domain Name
CLOUDFRONT_DOMAIN=xxxxxxxxxx.cloudfront.net
```

### C. CloudWatch Monitoring (FREE TIER)

#### C.1 Configurar CloudWatch Agent en EC2
```bash
# En cada EC2, instalar CloudWatch Agent
wget https://s3.amazonaws.com/amazoncloudwatch-agent/amazon_linux/amd64/latest/amazon-cloudwatch-agent.rpm
sudo rpm -U ./amazon-cloudwatch-agent.rpm

# Crear configuración básica
sudo tee /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json > /dev/null <<EOF
{
  "metrics": {
    "namespace": "VideoRank/EC2",
    "metrics_collected": {
      "cpu": {
        "measurement": [
          "cpu_usage_idle",
          "cpu_usage_iowait",
          "cpu_usage_user",
          "cpu_usage_system"
        ],
        "metrics_collection_interval": 300
      },
      "disk": {
        "measurement": [
          "used_percent"
        ],
        "metrics_collection_interval": 300,
        "resources": [
          "*"
        ]
      },
      "mem": {
        "measurement": [
          "mem_used_percent"
        ],
        "metrics_collection_interval": 300
      }
    }
  }
}
EOF

# Iniciar agent
sudo /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl \
  -a fetch-config -m ec2 -c file:/opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json -s
```

#### C.2 Crear CloudWatch Alarms
```bash
# Alarm para CPU alta en API
aws cloudwatch put-metric-alarm \
  --alarm-name "VideoRank-API-HighCPU" \
  --alarm-description "API EC2 High CPU" \
  --metric-name CPUUtilization \
  --namespace AWS/EC2 \
  --statistic Average \
  --period 300 \
  --threshold 80 \
  --comparison-operator GreaterThanThreshold \
  --dimensions Name=InstanceId,Value=$API_INSTANCE_ID \
  --evaluation-periods 2

# Alarm para memoria alta en Workers
aws cloudwatch put-metric-alarm \
  --alarm-name "VideoRank-Workers-HighMemory" \
  --alarm-description "Workers EC2 High Memory" \
  --metric-name mem_used_percent \
  --namespace VideoRank/EC2 \
  --statistic Average \
  --period 300 \
  --threshold 85 \
  --comparison-operator GreaterThanThreshold \
  --dimensions Name=InstanceId,Value=$WORKERS_INSTANCE_ID \
  --evaluation-periods 2
```

### D. AWS Systems Manager (Gestión Remota)

#### D.1 Configurar SSM Agent
```bash
# Crear IAM Role para EC2
aws iam create-role \
  --role-name VideoRank-EC2-SSM-Role \
  --assume-role-policy-document '{
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Principal": {
          "Service": "ec2.amazonaws.com"
        },
        "Action": "sts:AssumeRole"
      }
    ]
  }'

# Asociar política SSM
aws iam attach-role-policy \
  --role-name VideoRank-EC2-SSM-Role \
  --policy-arn arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore

# Crear instance profile
aws iam create-instance-profile \
  --instance-profile-name VideoRank-EC2-SSM-Profile

aws iam add-role-to-instance-profile \
  --instance-profile-name VideoRank-EC2-SSM-Profile \
  --role-name VideoRank-EC2-SSM-Role

# Asociar a instancias existentes
aws ec2 associate-iam-instance-profile \
  --instance-id $API_INSTANCE_ID \
  --iam-instance-profile Name=VideoRank-EC2-SSM-Profile
```

### E. Backup Automático

#### E.1 Configurar snapshots automáticos
```bash
# Crear lifecycle policy para snapshots
aws dlm create-lifecycle-policy \
  --execution-role-arn arn:aws:iam::ACCOUNT:role/AWSDataLifecycleManagerDefaultRole \
  --description "VideoRank Daily Snapshots" \
  --state ENABLED \
  --policy-details '{
    "ResourceTypes": ["VOLUME"],
    "TargetTags": [
      {
        "Key": "Project",
        "Value": "VideoRank"
      }
    ],
    "Schedules": [
      {
        "Name": "DailySnapshots",
        "CreateRule": {
          "Interval": 24,
          "IntervalUnit": "HOURS",
          "Times": ["03:00"]
        },
        "RetainRule": {
          "Count": 7
        },
        "CopyTags": true
      }
    ]
  }'

# Etiquetar volúmenes para backup
aws ec2 create-tags \
  --resources $(aws ec2 describe-instances --instance-ids $DATABASE_INSTANCE_ID --query 'Reservations[0].Instances[0].BlockDeviceMappings[0].Ebs.VolumeId' --output text) \
  --tags Key=Project,Value=VideoRank
```

### F. Configuraciones de Seguridad Adicionales

#### F.1 Habilitar VPC Flow Logs
```bash
# Crear IAM role para Flow Logs
aws iam create-role \
  --role-name flowlogsRole \
  --assume-role-policy-document '{
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Principal": {
          "Service": "vpc-flow-logs.amazonaws.com"
        },
        "Action": "sts:AssumeRole"
      }
    ]
  }'

# Crear log group
aws logs create-log-group --log-group-name VPCFlowLogs

# Habilitar Flow Logs
aws ec2 create-flow-logs \
  --resource-type VPC \
  --resource-ids $VPC_ID \
  --traffic-type ALL \
  --log-destination-type cloud-watch-logs \
  --log-group-name VPCFlowLogs
```

#### F.2 Configurar Instance Metadata Service v2
```bash
# Forzar IMDSv2 en todas las instancias
for instance in $DATABASE_INSTANCE_ID $STORAGE_INSTANCE_ID $WORKERS_INSTANCE_ID $API_INSTANCE_ID $FRONTEND_INSTANCE_ID; do
  aws ec2 modify-instance-metadata-options \
    --instance-id $instance \
    --http-tokens required \
    --http-put-response-hop-limit 1
done
```

### G. Optimizaciones de Costos

#### G.1 Configurar Instance Scheduler
```bash
# Crear tags para auto-scheduling (Lun-Vie 8AM-6PM)
aws ec2 create-tags \
  --resources $API_INSTANCE_ID $WORKERS_INSTANCE_ID \
  --tags Key=Schedule,Value=office-hours

# Instalar Instance Scheduler (CloudFormation template gratuito)
# Ahorra ~65% en costos para instancias no críticas
```

## Arquitectura Final con Buenas Prácticas

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CloudFront    │    │  API Gateway    │    │      ALB        │
│      CDN        │    │   + Throttling  │    │  (Alternativa)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  EC2-FRONTEND   │    │    EC2-API      │    │  EC2-WORKERS    │
│   (Public)      │    │   (Private)     │    │   (Private)     │
│ + CloudWatch    │    │ + CloudWatch    │    │ + CloudWatch    │
│ + SSM           │    │ + SSM           │    │ + SSM           │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                       │
                                ▼                       ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │  EC2-STORAGE    │    │  EC2-DATABASE   │
                       │   (Private)     │    │   (Private)     │
                       │ + Snapshots     │    │ + Snapshots     │
                       │ + Encryption    │    │ + Encryption    │
                       └─────────────────┘    └─────────────────┘
```

## Costos Estimados con Buenas Prácticas (Free Tier)

### Servicios Gratuitos (12 meses):
- 5 x t3.micro EC2: $0
- 5 x 30GB EBS: $0
- API Gateway: 1M requests/mes gratis
- CloudFront: 50GB transferencia/mes gratis
- CloudWatch: 10 métricas personalizadas gratis
- Systems Manager: Gratis
- VPC Flow Logs: Primeros 10GB gratis

### Costos mínimos después de Free Tier:
- CloudWatch Logs: ~$0.50/mes
- **Total estimado**: ~$0.50/mes después del Free Tier

## Beneficios de las Buenas Prácticas:

1. **API Gateway**: Rate limiting, throttling, monitoreo, HTTPS automático
2. **CloudFront**: CDN global, cache, reducción de latencia
3. **CloudWatch**: Monitoreo proactivo, alertas automáticas
4. **SSM**: Gestión remota sin SSH, mayor seguridad
5. **Snapshots**: Backup automático, recuperación ante desastres
6. **VPC Flow Logs**: Auditoría de tráfico, detección de anomalías
7. **Instance Scheduler**: Ahorro automático de costos
8. **IMDSv2**: Mayor seguridad en metadatos de instancia

## Troubleshooting

### Problemas comunes:
1. **Conectividad**: Verificar Security Groups
2. **DNS**: Usar IPs privadas para comunicación interna
3. **Memoria**: Monitorear uso con `htop`
4. **Logs**: Revisar con `docker-compose logs -f`

### Comandos útiles:
```bash
# Reiniciar servicio específico
docker-compose restart [service]

# Reconstruir imagen
docker-compose build [service]

# Ver uso de recursos
docker stats

# Limpiar sistema
docker system prune -f

# Conectar via SSM (sin SSH)
aws ssm start-session --target $INSTANCE_ID

# Ver métricas CloudWatch
aws cloudwatch get-metric-statistics \
  --namespace AWS/EC2 \
  --metric-name CPUUtilization \
  --dimensions Name=InstanceId,Value=$API_INSTANCE_ID \
  --statistics Average \
  --start-time 2024-01-01T00:00:00Z \
  --end-time 2024-01-01T23:59:59Z \
  --period 3600
```