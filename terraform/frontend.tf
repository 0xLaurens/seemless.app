resource "aws_s3_bucket" "frontend_bucket" {
  bucket = var.s3_bucket_name
  tags = {
    Name = "Vue Frontend"
  }
}

resource "aws_s3_bucket_ownership_controls" "frontend_bucket_ownership_controls" {
  bucket = aws_s3_bucket.frontend_bucket.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_public_access_block" "aws_s3_bucket_public_access" {
  bucket = aws_s3_bucket.frontend_bucket.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_acl" "frontend_aws_s3_bucket_acl" {
  depends_on = [
    aws_s3_bucket_ownership_controls.frontend_bucket_ownership_controls,
    aws_s3_bucket_public_access_block.aws_s3_bucket_public_access,
  ]

  bucket = aws_s3_bucket.frontend_bucket.id
  acl    = "public-read"
}

resource "aws_s3_bucket_website_configuration" "frontend_bucket_config" {
  bucket = aws_s3_bucket.frontend_bucket.id
  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "index.html"
  }
}