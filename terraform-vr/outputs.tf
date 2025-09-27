
output "front_public_ip" {
  value = aws_instance.front.public_ip
}

output "front_public_url" {
  value = "http://${aws_instance.front.public_ip}"
}

output "rds_endpoint" {
  value = aws_db_instance.postgres.address
}

output "minio_console_url_private" {
  value = "http://${aws_instance.minio.private_ip}:9001"
}

output "rabbit_console_url_private" {
  value = "http://${aws_instance.rabbit.private_ip}:15672"
}

output "ssh_key_path" {
  value = local_file.private_key_pem.filename
}
