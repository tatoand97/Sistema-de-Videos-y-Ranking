
resource "tls_private_key" "ssh" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "generated" {
  key_name   = "vr-keypair"
  public_key = tls_private_key.ssh.public_key_openssh
}

resource "local_file" "private_key_pem" {
  filename        = "${path.module}/key.pem"
  content         = tls_private_key.ssh.private_key_pem
  file_permission = "0600"
}
