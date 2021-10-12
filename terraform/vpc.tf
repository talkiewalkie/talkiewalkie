resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name        = "Main"
    Environment = local.environment
  }
}

resource "aws_security_group" "endpoint_sg" {
  name   = "vpc-endpoints-sg"
  vpc_id = aws_vpc.main.id

  ingress {
    protocol         = "tcp"
    from_port        = 0
    to_port          = 0
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  egress {
    protocol         = "-1"
    from_port        = 0
    to_port          = 0
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
}

resource "aws_vpc_endpoint" "s3" {
  vpc_id       = aws_vpc.main.id
  service_name = "com.amazonaws.${local.region}.s3"
}

resource "aws_vpc_endpoint" "ecr" {
  vpc_id              = aws_vpc.main.id
  vpc_endpoint_type   = "Interface"
  security_group_ids  = [aws_security_group.endpoint_sg.id]
  private_dns_enabled = true
  subnet_ids          = aws_subnet.public.*.id
  service_name        = "com.amazonaws.${local.region}.ecr.dkr"
}

resource "aws_vpc_endpoint" "ecr2" {
  vpc_id              = aws_vpc.main.id
  vpc_endpoint_type   = "Interface"
  security_group_ids  = [aws_security_group.endpoint_sg.id]
  private_dns_enabled = true
  subnet_ids          = aws_subnet.public.*.id
  service_name        = "com.amazonaws.${local.region}.ecr.api"
}

resource "aws_vpc_endpoint" "secrets" {
  vpc_id              = aws_vpc.main.id
  vpc_endpoint_type   = "Interface"
  security_group_ids  = [aws_security_group.endpoint_sg.id]
  private_dns_enabled = true
  subnet_ids          = aws_subnet.public.*.id
  service_name        = "com.amazonaws.${local.region}.secretsmanager"
}

resource "aws_vpc_endpoint" "ssm" {
  vpc_id              = aws_vpc.main.id
  vpc_endpoint_type   = "Interface"
  security_group_ids  = [aws_security_group.endpoint_sg.id]
  private_dns_enabled = true
  subnet_ids          = aws_subnet.public.*.id
  service_name        = "com.amazonaws.${local.region}.ssm"
}

resource "aws_vpc_endpoint" "logs" {
  vpc_id              = aws_vpc.main.id
  vpc_endpoint_type   = "Interface"
  security_group_ids  = [aws_security_group.endpoint_sg.id]
  private_dns_enabled = true
  subnet_ids          = aws_subnet.public.*.id
  service_name        = "com.amazonaws.${local.region}.logs"
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id
  tags   = {
    Environment = local.environment
  }
}

resource "aws_nat_gateway" "main" {
  count         = length(local.availability_zones)
  allocation_id = element(aws_eip.nat.*.id, count.index)
  subnet_id     = element(aws_subnet.private.*.id, count.index)
  depends_on    = [aws_internet_gateway.main]

  tags = {
    Name        = "nat-${local.environment}-${format("%03d", count.index+1)}"
    Environment = local.environment
  }
}


resource "aws_eip" "nat" {
  count = length(local.availability_zones)
  vpc   = true

  tags = { Environment = local.environment }
}


resource "aws_subnet" "public" {
  vpc_id                  = aws_vpc.main.id
  map_public_ip_on_launch = true

  cidr_block        = element(local.public_subnets, count.index)
  availability_zone = element(local.availability_zones, count.index)
  count             = length(local.availability_zones)

  tags = {
    Name        = "public-${element(local.availability_zones, count.index)}"
    Environment = local.environment
  }
}

resource "aws_subnet" "private" {
  vpc_id = aws_vpc.main.id

  cidr_block        = element(local.private_subnets, count.index)
  availability_zone = element(local.availability_zones, count.index)
  count             = length(local.availability_zones)

  tags = {
    Name        = "private-${element(local.availability_zones, count.index)}"
    Environment = local.environment
  }
}


resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id
}

resource "aws_route" "public" {
  route_table_id         = aws_route_table.public.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.main.id
}

resource "aws_route_table_association" "public" {
  count          = length(local.public_subnets)
  subnet_id      = element(aws_subnet.public.*.id, count.index)
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table" "private" {
  count  = length(local.availability_zones)
  vpc_id = aws_vpc.main.id

  tags = {
    Environment = local.environment
  }
}

resource "aws_route" "private" {
  count                  = length(compact(local.private_subnets))
  route_table_id         = element(aws_route_table.private.*.id, count.index)
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = element(aws_nat_gateway.main.*.id, count.index)
}

resource "aws_route_table_association" "private" {
  count          = length(local.private_subnets)
  subnet_id      = element(aws_subnet.private.*.id, count.index)
  route_table_id = element(aws_route_table.private.*.id, count.index)
}