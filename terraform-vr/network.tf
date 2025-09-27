resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name = "vr-vpc"
  }
}

# Subred pública (front/nat)
resource "aws_subnet" "public" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.0.0/24"
  map_public_ip_on_launch = true
  availability_zone       = "us-east-1a"

  tags = {
    Name = "vr-public"
  }
}

# Subredes privadas
resource "aws_subnet" "api" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = "us-east-1a"

  tags = {
    Name = "vr-api"
  }
}

resource "aws_subnet" "workers" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.2.0/24"
  availability_zone = "us-east-1a"

  tags = {
    Name = "vr-workers"
  }
}

resource "aws_subnet" "rabbit" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.3.0/24"
  availability_zone = "us-east-1a"

  tags = {
    Name = "vr-rabbit"
  }
}

resource "aws_subnet" "minio" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.4.0/24"
  availability_zone = "us-east-1a"

  tags = {
    Name = "vr-minio"
  }
}

resource "aws_subnet" "redis" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.5.0/24"
  availability_zone = "us-east-1a"

  tags = {
    Name = "vr-redis"
  }
}

# IGW y ruta pública
resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "vr-igw"
  }
}

resource "aws_route_table" "public_rt" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }

  tags = {
    Name = "vr-public-rt"
  }
}

resource "aws_route_table_association" "public_assoc" {
  subnet_id      = aws_subnet.public.id
  route_table_id = aws_route_table.public_rt.id
}

# NAT Gateway para privadas
resource "aws_eip" "nat_eip" {
  domain     = "vpc"
  depends_on = [aws_internet_gateway.igw]

  tags = {
    Name = "vr-nat-eip"
  }
}

resource "aws_nat_gateway" "nat" {
  allocation_id = aws_eip.nat_eip.id
  subnet_id     = aws_subnet.public.id

  tags = {
    Name = "vr-nat"
  }
}

resource "aws_route_table" "private_rt" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.nat.id
  }

  tags = {
    Name = "vr-private-rt"
  }
}

# Asociaciones privadas
resource "aws_route_table_association" "api_assoc" {
  subnet_id      = aws_subnet.api.id
  route_table_id = aws_route_table.private_rt.id
}

resource "aws_route_table_association" "workers_assoc" {
  subnet_id      = aws_subnet.workers.id
  route_table_id = aws_route_table.private_rt.id
}

resource "aws_route_table_association" "rabbit_assoc" {
  subnet_id      = aws_subnet.rabbit.id
  route_table_id = aws_route_table.private_rt.id
}

resource "aws_route_table_association" "minio_assoc" {
  subnet_id      = aws_subnet.minio.id
  route_table_id = aws_route_table.private_rt.id
}

resource "aws_route_table_association" "redis_assoc" {
  subnet_id      = aws_subnet.redis.id
  route_table_id = aws_route_table.private_rt.id
}
