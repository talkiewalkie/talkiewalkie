resource "aws_s3_bucket" "public_landing_dev" {
  bucket        = aws_route53_zone.main.name
  acl           = "public-read"
  force_destroy = true
  policy        = jsonencode({
    Version   = "2012-10-17"
    Statement = [
      {
        Sid       = "PublicReadGetObject"
        Effect    = "Allow"
        Principal = "*"
        Action    = "s3:GetObject"
        Resource  = ["arn:aws:s3:::${aws_route53_zone.main.name}/*"]
      },
    ]
  })

  website {
    index_document = "index.html"
    error_document = "index.html"
  }
}

resource "aws_route53_record" "www" {
  zone_id = aws_route53_zone.main.zone_id
  name    = "www"

  type    = "CNAME"
  ttl     = 5
  records = [aws_route53_record.apex.name]

  depends_on = [aws_s3_bucket.public_landing_dev]
}

resource "aws_route53_record" "apex" {
  name    = aws_route53_zone.main.name
  zone_id = aws_route53_zone.main.zone_id
  type    = "A"
  alias {
    name                   = aws_s3_bucket.public_landing_dev.website_domain
    zone_id                = aws_s3_bucket.public_landing_dev.hosted_zone_id
    evaluate_target_health = false
  }
}


