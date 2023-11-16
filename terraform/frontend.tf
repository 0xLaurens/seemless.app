resource "aws_s3_bucket" "frontend-vue" {
  bucket = var.s3_bucket_name
  tags = {
    Name = "Vue Frontend"
  }
}

resource "aws_s3_bucket_website_configuration" "frontend-config" {
  bucket = var.s3_bucket_name
  index_document {
    suffix = "index.html"
  }
  error_document {
    key = "index.html"
  }
}

resource "aws_s3_bucket_policy" "static_website_policy" {
  bucket = aws_s3_bucket.frontend-vue.id
  policy = templatefile("s3-policy.json", { bucket = var.s3_bucket_name })
}
