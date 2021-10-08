locals {
  region             = "eu-west-3"
  availability_zones = ["${local.region}a", "${local.region}b"]
  private_subnets    = ["10.0.0.0/20", "10.0.32.0/20"]
  public_subnets     = ["10.0.16.0/20", "10.0.48.0/20"]
  environment        = "test"
}