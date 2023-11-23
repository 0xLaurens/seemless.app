resource "aws_ecr_repository" "container_repo" {
  name = "seemless-backend"
  image_scanning_configuration {
    scan_on_push = true
  }
}