module "grpc" {
  source = "./grpc"

  region            = local.region
  vpc_id            = aws_vpc.main.id
  cluster_id        = aws_ecs_cluster.main.id
  public_subnet_ids = aws_subnet.public.*.id

  domain_zone_id = aws_route53_zone.main.zone_id
  domain_name    = aws_route53_zone.main.name

  alb_listener_arn = aws_alb_listener.https.arn
  alb_dns_name     = aws_lb.main.dns_name
  alb_zone_id      = aws_lb.main.zone_id

  ecr_policy_arn  = aws_iam_policy.ecr.arn
  s3_policy_arn   = aws_iam_policy.s3.arn
  logs_policy_arn = aws_iam_policy.logging_writer.arn
}