
resource "aws_db_subnet_group" "pg_subnets" {
  name       = "vr-pg-subnets"
  subnet_ids = [aws_subnet.db_a.id, aws_subnet.db_b.id]
  tags       = { Name = "vr-pg-subnets" }
}

resource "aws_db_instance" "postgres" {
  identifier              = "vr-postgres"
  engine                  = "postgres"
  engine_version          = "16"
  instance_class          = "db.t3.micro"
  allocated_storage       = 20
  username                = local.postgres_user
  password                = local.postgres_password
  db_name                 = local.postgres_db
  port                    = local.postgres_port
  db_subnet_group_name    = aws_db_subnet_group.pg_subnets.name
  vpc_security_group_ids  = [aws_security_group.rds.id]
  skip_final_snapshot     = true
  publicly_accessible     = false
  deletion_protection     = false
  storage_encrypted       = true
  backup_retention_period = 0

  tags = { Name = "vr-postgres" }
}
