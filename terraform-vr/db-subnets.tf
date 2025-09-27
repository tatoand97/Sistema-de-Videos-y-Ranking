###############################################################################
# Subredes privadas para RDS en 2 AZ diferentes
# - Ajusta los CIDR si ya están ocupados (aquí usamos 10.0.6.0/24 y 10.0.7.0/24)
# - Requiere que ya existan:
#     - aws_vpc.main
#     - aws_route_table.private_rt (tabla de ruteo de subredes privadas)
###############################################################################

# AZ disponibles (evita hardcodear la letra de la AZ)
# Si ya tienes otro data aws_availability_zones en tu código, puedes reutilizarlo
data "aws_availability_zones" "db_azs" {
  state = "available"
}

# Subred privada DB - AZ 1
resource "aws_subnet" "db_a" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.6.0/24"
  availability_zone = data.aws_availability_zones.db_azs.names[0]

  # Subred privada: sin IP pública automática
  map_public_ip_on_launch = false

  tags = {
    Name = "vr-db-a"
    Role = "db"
  }
}

# Subred privada DB - AZ 2
resource "aws_subnet" "db_b" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.7.0/24"
  availability_zone = data.aws_availability_zones.db_azs.names[1]

  map_public_ip_on_launch = false

  tags = {
    Name = "vr-db-b"
    Role = "db"
  }
}

# Asociaciones con la tabla de ruteo privada
resource "aws_route_table_association" "db_a_assoc" {
  subnet_id      = aws_subnet.db_a.id
  route_table_id = aws_route_table.private_rt.id
}

resource "aws_route_table_association" "db_b_assoc" {
  subnet_id      = aws_subnet.db_b.id
  route_table_id = aws_route_table.private_rt.id
}