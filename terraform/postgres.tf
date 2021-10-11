module "postgres" {
  source = "./postgres"

  region             = local.region
  vpc_id             = aws_vpc.main.id
  cluster_id         = aws_ecs_cluster.main.id
  public_subnet_ids  = aws_subnet.public.*.id
  secrets_policy_arn = aws_iam_policy.secrets_fetcher.arn
}