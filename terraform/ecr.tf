resource "aws_ecr_repository" "container_repo" {
  name = "container-repo"
  image_scanning_configuration {
    scan_on_push = true
  }
}