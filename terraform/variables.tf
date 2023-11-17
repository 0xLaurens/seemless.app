variable "s3_bucket_name" {
  description = "s3 static website bucket name"
  type        = string
  default     = "frontend-s3-bitdash"
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "eu-central-1"
}

variable "s3_origin_id" {
  description = "Cloudfront s3 origin id"
  type        = string
  default     = "bitdash-origin-id"
}