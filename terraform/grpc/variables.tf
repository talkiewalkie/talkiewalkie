variable "region" { type = string }
variable "vpc_id" { type = string }
variable "cluster_id" { type = string }
variable "public_subnet_ids" { type = list(string) }

variable "domain_zone_id" { type = string }
variable "domain_name" { type = string }

variable "alb_listener_arn" { type = string }
variable "alb_dns_name" { type = string }
variable "alb_zone_id" { type = string }

variable "ecr_policy_arn" { type = string }
variable "s3_policy_arn" { type = string }
variable "logs_policy_arn" { type = string }