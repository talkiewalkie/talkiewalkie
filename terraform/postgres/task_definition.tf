resource "aws_ecs_task_definition" "postgres" {
  family                   = "postgres"
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.postgres_task_execution_role.arn
  task_role_arn            = aws_iam_role.postgres_task_role.arn
  requires_compatibilities = ["FARGATE"]

  container_definitions = jsonencode([
    {
      name  = "postgres"
      image = "postgres:13"

      essential        = true
      portMappings     = [
        {
          containerPort = 5432
          hostPort      = 5432
        }
      ]
      environment      = [
        { name : "POSTGRES_USER", value : "talkiewalkie" },
        { name : "POSTGRES_DB", value : "talkiewalkie" },
      ],
      logConfiguration = {
        logDriver = "awslogs",
        options   = {
          awslogs-group         = aws_cloudwatch_log_group.postgres.name,
          awslogs-region        = var.region,
          awslogs-stream-prefix = "ecs"
        }
      },
      secrets          = [
        {
          name      = "POSTGRES_PASSWORD"
          valueFrom = aws_secretsmanager_secret.postgres_task_password.arn
        }
      ]
    },
  ])
}

resource "aws_cloudwatch_log_group" "postgres" {
  name = "postgres"
}

resource "aws_iam_role" "postgres_task_execution_role" {
  name = "postgres-ecsTaskExecutionRole"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role" "postgres_task_role" {
  name = "postgres_task_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "postgres_task_secrets_policy" {
  for_each   = toset([
    aws_iam_role.postgres_task_role.name,
    aws_iam_role.postgres_task_execution_role.name
  ])
  role       = each.value
  policy_arn = var.secrets_policy_arn
}
