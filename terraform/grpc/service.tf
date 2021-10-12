resource "aws_ecs_service" "grpc" {
  name                               = "grpc"
  cluster                            = var.cluster_id
  task_definition                    = aws_ecs_task_definition.grpc.arn
  launch_type                        = "FARGATE"
  deployment_minimum_healthy_percent = 0
  deployment_maximum_percent         = 200
  desired_count                      = 1

  network_configuration {
    security_groups  = [aws_security_group.ecs_grpc_task.id]
    subnets          = var.public_subnet_ids
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_alb_target_group.grpc.arn
    container_name   = "grpc"
    container_port   = 50051
  }
}
