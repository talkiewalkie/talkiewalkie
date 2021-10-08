resource "aws_security_group" "ecs_grpc_task" {
  name   = "grpc_task"
  vpc_id = var.vpc_id

  ingress {
    protocol         = "tcp"
    from_port        = 443
    to_port          = 50051
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


resource "aws_iam_role_policy_attachment" "grpc_task_ecr_policy" {
  for_each   = toset([
    aws_iam_role.grpc_task_role.name,
    aws_iam_role.grpc_task_execution_role.name
  ])
  role       = each.value
  policy_arn = var.ecr_policy_arn
}

resource "aws_iam_role_policy_attachment" "grpc_task_s3_policy" {
  for_each   = toset([
    aws_iam_role.grpc_task_role.name,
    aws_iam_role.grpc_task_execution_role.name
  ])
  role       = each.value
  policy_arn = var.s3_policy_arn
}

resource "aws_iam_role_policy_attachment" "grpc_task_log_policy" {
  for_each   = toset([
    aws_iam_role.grpc_task_role.name,
    aws_iam_role.grpc_task_execution_role.name
  ])
  role       = each.value
  policy_arn = var.logs_policy_arn
}
