resource "aws_ecs_service" "postgres" {
  name                               = "postgres"
  cluster                            = var.cluster_id
  task_definition                    = aws_ecs_task_definition.postgres.arn
  launch_type                        = "FARGATE"
  deployment_minimum_healthy_percent = 0
  deployment_maximum_percent         = 200
  desired_count                      = 1

  network_configuration {
    security_groups  = [aws_security_group.ecs_postgres_task.id]
    subnets          = var.public_subnet_ids
    assign_public_ip = false
  }
}