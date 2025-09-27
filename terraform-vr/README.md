
# Terraform - Video Ranking (AWS, EC2 + Docker + RDS)

Este proyecto crea **toda** la infraestructura sin variables externas:

- **Región:** `us-east-1` (puedes cambiar en `provider` si lo deseas).
- **Red:** 1 VPC `10.0.0.0/16` con:
  - 1 Subred pública `10.0.0.0/24` (Front/NAT/IGW)
  - 5 Subredes privadas:
    - `10.0.1.0/24` API
    - `10.0.2.0/24` Workers
    - `10.0.3.0/24` RabbitMQ
    - `10.0.4.0/24` MinIO
    - `10.0.5.0/24` Redis
  - **IGW** y **NAT Gateway** para que las instancias privadas puedan hacer `docker pull`.
- **RDS PostgreSQL** privado (usuario `app_user`, pass `app_password`, BD `videorank`).
- **EC2** (Amazon Linux 2023) que instalan Docker y corren contenedores:
  - Front: `ricardoandres97/cloud-front:latest` (HTTP 80 público)
  - API: `ricardoandres97/cloud-api:latest`
  - Workers: `ricardoandres97/cloud-workers:latest`
  - RabbitMQ: `ricardoandres97/cloud-rabbit:latest` (5672/15672 privados)
  - MinIO: `eddauni/minio:latest` (9000/9001 privados)
  - Redis: `eddauni/redis:latest` (6379 privado)

## Requisitos
- Terraform >= 1.5
- Una cuenta AWS con permisos para VPC/EC2/RDS/Elastic IP/etc.

## Uso
```bash
terraform init
terraform apply -auto-approve
```

Al finalizar, revisa los **outputs**:
- `front_public_url` → URL HTTP del Front
- `rds_endpoint` → Endpoint de PostgreSQL (privado)
- `minio_console_url_private` → URL privada del console de MinIO
- `rabbit_console_url_private` → URL privada de RabbitMQ
- `ssh_key_path` → Ruta del archivo `key.pem` generado

La clave privada se guarda localmente en `key.pem` (no se sube a AWS). Si deseas SSH:
```bash
chmod 600 key.pem
ssh -i key.pem ec2-user@<IP_PUBLICA_FRONT>
```

> **Costos**: Este stack genera costos (NAT, RDS, EC2). El NAT Gateway tiene costo por hora y tráfico. Para demo, puedes apagarlo con `terraform destroy` al terminar.

## Notas
- Si deseas exponer MinIO/Rabbit/Redis, mueve sus instancias a la subred pública o agrega un Load Balancer/port-forwarding. Por defecto permanecen privados.
- Puedes ajustar tipos de instancia en `locals`.
