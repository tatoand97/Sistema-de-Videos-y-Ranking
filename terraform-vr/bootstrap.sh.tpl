
#!/bin/bash
set -e

# Instalar Docker (Amazon Linux 2023)
dnf update -y
dnf install -y docker
systemctl enable docker
systemctl start docker

# Red Docker
docker network create vrnet || true
