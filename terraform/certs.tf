resource "aws_route53_zone" "main" {
  name = "002fa7.net"
}


resource "aws_route53_record" "nginx" {
  zone_id = aws_route53_zone.main.zone_id
  name    = "nginx.${aws_route53_zone.main.name}"
  type    = "A"

  alias {
    evaluate_target_health = false
    name                   = aws_alb.main.dns_name
    zone_id                = aws_alb.main.zone_id
  }
}

resource "aws_acm_certificate" "cert" {
  domain_name               = "*.002fa7.net"
  subject_alternative_names = [
    "*.dev.002fa7.net",
    "*.staging.002fa7.net",
  ]
  validation_method         = "DNS"

  options {
    certificate_transparency_logging_preference = "ENABLED"
  }
}
