variable "region" { type = string }
variable "vpc_id" { type = string }
variable "cluster_id" { type = string }
variable "public_subnet_ids" { type = list(string) }

variable "secrets_policy_arn" { type = string }
