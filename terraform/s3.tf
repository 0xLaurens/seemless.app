resource "aws_s3_bucket" "frontend_bucket" {
  bucket = var.s3_bucket_name


  tags = {
    Name = "Vue Frontend"
  }
}

data "aws_iam_policy_document" "frontend_iam_policy" {
  version = "2012-10-17"

  statement {
    sid    = "AllowAOI"
    effect = "Allow"
    principals {
      type        = "AWS"
      identifiers = [aws_cloudfront_origin_access_identity.cloudfront_origin_identity.iam_arn]
    }
    actions   = ["s3:GetObject"]
    resources = ["arn:aws:s3:::${var.s3_bucket_name}/*"]
  }
}

resource "aws_s3_bucket_public_access_block" "block_public" {
  bucket = aws_s3_bucket.frontend_bucket.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_policy" "frontend_bucket_policy" {
  bucket = aws_s3_bucket.frontend_bucket.id
  policy = data.aws_iam_policy_document.frontend_iam_policy.json
}