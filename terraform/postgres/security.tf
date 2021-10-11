resource "aws_security_group" "ecs_postgres_task" {
  name   = "postgres_task"
  vpc_id = var.vpc_id

  ingress {
    protocol         = "tcp"
    from_port        = 0
    to_port          = 0
    # TODO: once things are working, restrict to correct cidr block
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  egress {
    protocol         = "-1"
    from_port        = 0
    to_port          = 0
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
}

resource "aws_secretsmanager_secret" "postgres_task_password" {
  name = "postgres-task-password"
}

resource "random_password" "database_password" {
  length  = 32
  special = false
}

resource "aws_secretsmanager_secret_version" "database_password_secret_version" {
  secret_id     = aws_secretsmanager_secret.postgres_task_password.id
  secret_string = random_password.database_password.result
}
