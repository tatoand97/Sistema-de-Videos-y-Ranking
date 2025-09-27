locals {
  docker_bootstrap = <<-EOT
#!/bin/bash
set -e

# Instalar Docker (Amazon Linux 2023)
dnf update -y
dnf install -y docker
systemctl enable docker
systemctl start docker

# Red Docker (idempotente)
docker network create vrnet || true
EOT
}

# FRONT (público)
resource "aws_instance" "front" {
  ami                    = data.aws_ami.al2023.id
  instance_type          = local.instance_type
  subnet_id              = aws_subnet.public.id
  key_name               = aws_key_pair.generated.key_name
  vpc_security_group_ids = [aws_security_group.front.id, aws_security_group.egress_all.id]

  user_data = <<-EOF
${local.docker_bootstrap}

# FRONT
docker pull ${local.img_front}
docker rm -f front || true
docker run -d --name front --restart unless-stopped \
  --network vrnet \
  -p 80:80 \
  ${local.img_front}
EOF

  tags = {
    Name = "vr-front"
  }
}

# API (privado)
resource "aws_instance" "api" {
  ami                    = data.aws_ami.al2023.id
  instance_type          = local.instance_type
  subnet_id              = aws_subnet.api.id
  key_name               = aws_key_pair.generated.key_name
  vpc_security_group_ids = [aws_security_group.internal.id, aws_security_group.egress_all.id]

  user_data = <<-EOF
${local.docker_bootstrap}

# API
docker pull ${local.img_api}
docker rm -f api || true
docker run -d --name api --restart unless-stopped \
  --network vrnet \
  -e POSTGRES_USER=${local.postgres_user} \
  -e POSTGRES_PASSWORD=${local.postgres_password} \
  -e POSTGRES_DB=${local.postgres_db} \
  -e POSTGRES_HOST=${aws_db_instance.postgres.address} \
  -e POSTGRES_PORT=${local.postgres_port} \
  -e REDIS_HOST=${aws_instance.redis.private_ip} \
  -e REDIS_PORT=6379 \
  -e RABBIT_HOST=${aws_instance.rabbit.private_ip} \
  -e RABBIT_PORT=5672 \
  -e MINIO_ENDPOINT=http://${aws_instance.minio.private_ip}:9000 \
  -e MINIO_ACCESS_KEY=${local.minio_root_user} \
  -e MINIO_SECRET_KEY=${local.minio_root_password} \
  -p 8080:8080 \
  ${local.img_api}
EOF

  tags = {
    Name = "vr-api"
  }
}

# WORKERS (privado) - un contenedor que agrupa los workers
resource "aws_instance" "workers" {
  ami                    = data.aws_ami.al2023.id
  instance_type          = local.instance_type
  subnet_id              = aws_subnet.workers.id
  key_name               = aws_key_pair.generated.key_name
  vpc_security_group_ids = [aws_security_group.internal.id, aws_security_group.egress_all.id]

  user_data = <<-EOF
${local.docker_bootstrap}

# WORKERS
docker pull ${local.img_workers}
docker rm -f workers || true
docker run -d --name workers --restart unless-stopped \
  --network vrnet \
  -e POSTGRES_USER=${local.postgres_user} \
  -e POSTGRES_PASSWORD=${local.postgres_password} \
  -e POSTGRES_DB=${local.postgres_db} \
  -e POSTGRES_HOST=${aws_db_instance.postgres.address} \
  -e POSTGRES_PORT=${local.postgres_port} \
  -e REDIS_HOST=${aws_instance.redis.private_ip} \
  -e REDIS_PORT=6379 \
  -e RABBIT_HOST=${aws_instance.rabbit.private_ip} \
  -e RABBIT_PORT=5672 \
  -e MINIO_ENDPOINT=http://${aws_instance.minio.private_ip}:9000 \
  -e MINIO_ACCESS_KEY=${local.minio_root_user} \
  -e MINIO_SECRET_KEY=${local.minio_root_password} \
  ${local.img_workers}
EOF

  tags = {
    Name = "vr-workers"
  }
}

# RABBITMQ (privado)
resource "aws_instance" "rabbit" {
  ami                    = data.aws_ami.al2023.id
  instance_type          = local.instance_type
  subnet_id              = aws_subnet.rabbit.id
  key_name               = aws_key_pair.generated.key_name
  vpc_security_group_ids = [aws_security_group.internal.id, aws_security_group.egress_all.id]

  user_data = <<-EOF
${local.docker_bootstrap}

# RABBITMQ
docker pull ${local.img_rabbit}
docker rm -f rabbit || true
docker run -d --name rabbit --restart unless-stopped \
  --network vrnet \
  -p 5672:5672 -p 15672:15672 \
  ${local.img_rabbit}
EOF

  tags = {
    Name = "vr-rabbit"
  }
}

# MINIO (privado)
resource "aws_instance" "minio" {
  ami                    = data.aws_ami.al2023.id
  instance_type          = local.instance_type
  subnet_id              = aws_subnet.minio.id
  key_name               = aws_key_pair.generated.key_name
  vpc_security_group_ids = [aws_security_group.internal.id, aws_security_group.egress_all.id]

  user_data = <<-EOF
${local.docker_bootstrap}

# MINIO
docker pull ${local.img_minio}
docker rm -f minio || true
mkdir -p /data/minio
docker run -d --name minio --restart unless-stopped \
  --network vrnet \
  -e MINIO_ROOT_USER=${local.minio_root_user} \
  -e MINIO_ROOT_PASSWORD=${local.minio_root_password} \
  -v /data/minio:/data \
  -p 9000:9000 -p 9001:9001 \
  ${local.img_minio} server /data --console-address ":9001"
EOF

  tags = {
    Name = "vr-minio"
  }
}

# REDIS (privado)
resource "aws_instance" "redis" {
  ami                    = data.aws_ami.al2023.id
  instance_type          = local.instance_type
  subnet_id              = aws_subnet.redis.id
  key_name               = aws_key_pair.generated.key_name
  vpc_security_group_ids = [aws_security_group.internal.id, aws_security_group.egress_all.id]

  user_data = <<-EOF
${local.docker_bootstrap}

# REDIS
docker pull ${local.img_redis}
docker rm -f redis || true
docker run -d --name redis --restart unless-stopped \
  --network vrnet \
  -p 6379:6379 \
  ${local.img_redis}
EOF

  tags = {
    Name = "vr-redis"
  }
}
