
locals {
  # AMI y tamaño por defecto
  instance_type = "t3.medium"

  # Imágenes Docker
  img_front   = "ricardoandres97/cloud-front:latest"
  img_api     = "ricardoandres97/cloud-api:latest"
  img_workers = "ricardoandres97/cloud-workers:latest"
  img_rabbit  = "ricardoandres97/cloud-rabbit:latest"
  img_minio   = "eddauni/minio:latest"
  img_redis   = "eddauni/redis:latest"

  # Credenciales y config de app
  postgres_user     = "app_user"
  postgres_password = "app_password"
  postgres_db       = "videorank"
  postgres_port     = 5432

  minio_root_user     = "minioadmin"
  minio_root_password = "minioadmin"
}
