provider "aws" {
  alias  = "cloudfront-cert"
  region = "us-east-1"
}

resource "aws_route53_zone" "seemless-domain" {
  name = "seemless.app"

  lifecycle {
    # prevent destroy the dns has to be manually be reconfigured
    prevent_destroy = true
  }
}

resource "aws_acm_certificate" "cert" {
  provider                  = aws.cloudfront-cert
  domain_name               = "seemless.app"
  validation_method         = "DNS"
  subject_alternative_names = ["*.seemless.app"]
}

resource "aws_route53_record" "cert_record" {
  for_each = {
    for dvo in aws_acm_certificate.cert.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = aws_route53_zone.seemless-domain.zone_id
}

resource "aws_acm_certificate_validation" "cert_validation" {
  provider                = aws.cloudfront-cert
  certificate_arn         = aws_acm_certificate.cert.arn
  validation_record_fqdns = [for record in aws_route53_record.cert_record : record.fqdn]
}

resource "aws_route53_record" "www" {
  zone_id = aws_route53_zone.seemless-domain.zone_id
  type    = "A"
  name    = "www.seemless.app"

  alias {
    evaluate_target_health = false
    name                   = aws_cloudfront_distribution.s3_distribution.domain_name
    zone_id                = aws_cloudfront_distribution.s3_distribution.hosted_zone_id
  }

  depends_on = [aws_acm_certificate_validation.cert_validation]
}

resource "aws_route53_record" "base" {
  zone_id = aws_route53_zone.seemless-domain.zone_id
  type    = "A"
  name    = "seemless.app"

  alias {
    evaluate_target_health = false
    name                   = aws_cloudfront_distribution.s3_distribution.domain_name
    zone_id                = aws_cloudfront_distribution.s3_distribution.hosted_zone_id
  }

  depends_on = [aws_acm_certificate_validation.cert_validation]
}