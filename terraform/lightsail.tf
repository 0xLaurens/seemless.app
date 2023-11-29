resource "aws_lightsail_container_service" "lightsail_container_service" {
  name  = "lightsail-container-service"
  power = "micro"
  scale = 1

  private_registry_access {
    ecr_image_puller_role {
      is_active = true
    }
  }
}

data "aws_iam_policy_document" "lightsail_policy" {
  statement {
    effect = "Allow"

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    actions = [
      "ecr:BatchGetImage",
      "ecr:GetDownloadUrlForLayer",
    ]
  }
}

resource "aws_ecr_repository_policy" "default" {
  repository = aws_ecr_repository.container_repo.name
  policy     = data.aws_iam_policy_document.lightsail_policy.json
}

resource "aws_lightsail_container_service_deployment_version" "lightsail_container_deployment" {
  service_name = aws_lightsail_container_service.lightsail_container_service.name
  container {
    container_name = "seemless-backend"
    image          = "${aws_ecr_repository.container_repo.repository_url}:latest"
    ports          = { "3000" = "HTTP" }
  }

  public_endpoint {
    container_name = "seemless-backend"
    container_port = 3000

    health_check {
      healthy_threshold   = 2
      unhealthy_threshold = 2
      timeout_seconds     = 5
      interval_seconds    = 30
      path                = "/health"
      success_codes       = "200"
    }
  }
}