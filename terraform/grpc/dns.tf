resource "aws_route53_record" "grpc" {
  zone_id = var.domain_zone_id
  name    = "grpc.${var.domain_name}"
  type    = "A"

  alias {
    evaluate_target_health = false
    name                   = var.alb_dns_name
    zone_id                = var.alb_zone_id
  }
}

resource "aws_alb_target_group" "grpc" {
  name             = "grpc"
  port             = 443
  protocol         = "HTTPS"
  protocol_version = "GRPC"
  vpc_id           = var.vpc_id
  target_type      = "ip"

  health_check {
    healthy_threshold   = "3"
    interval            = "30"
    protocol            = "HTTP"
    matcher             = "0,12"
    timeout             = "3"
    path                = "/"
    unhealthy_threshold = "2"
  }
}

resource "aws_alb_listener_rule" "grpc" {
  listener_arn = var.alb_listener_arn

  action {
    type             = "forward"
    target_group_arn = aws_alb_target_group.grpc.arn
  }

  condition {
    host_header {
      values = ["grpc.${var.domain_name}"]
    }
  }
}

