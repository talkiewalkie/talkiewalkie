resource "aws_ecs_task_definition" "grpc" {
  family                   = "grpc"
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.grpc_task_execution_role.arn
  task_role_arn            = aws_iam_role.grpc_task_role.arn
  requires_compatibilities = ["FARGATE"]

  container_definitions = jsonencode([
    {
      name  = "grpc"
      #      image = "344467466667.dkr.ecr.eu-west-3.amazonaws.com/talkiewalkie:latest"
      #      image = "grpc/java-example-hostname"
      image = "344467466667.dkr.ecr.eu-west-3.amazonaws.com/talkiewalkie:grpc-helloworld"
      #      image = "docker.pkg.github.com/aws-samples/grpc-examples/greeter_server:v0.1.0"

      essential    = true
      portMappings = [
        {
          #          containerPort = 8080
          containerPort = 50051
          #          hostPort      = 8080
          hostPort      = 50051
        }
      ]
      environment  = [
        { name : "POSTGRES_HOST", value : "postgres" },
        { name : "POSTGRES_USER", value : "postgres" },
        { name : "POSTGRES_PASSWORD", value : "verysecret" }
      ]
      #      logConfiguration = {
      #        logDriver = "awslogs",
      #        options   = {
      #          awslogs-group         = aws_cloudwatch_log_group.grpc.name,
      #          awslogs-region        = var.region,
      #          awslogs-stream-prefix = "ecs"
      #        }
      #      }
    },
  ])
}

resource "aws_cloudwatch_log_group" "grpc" {
  name = "grpc"
}

resource "aws_iam_role" "grpc_task_execution_role" {
  name = "grpc-ecsTaskExecutionRole"

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

resource "aws_iam_role" "grpc_task_role" {
  name = "grpc_task_role"

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