resource "aws_elb" "elastic-lb" {
  name            = "load-balancer for the backend"
  subnets         = module.vpc.public_subnets
  security_groups = [aws_security_group.lb_sg.id]

  listener {
    instance_port     = 80
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "HTTP:80/"
    interval            = 30
  }
}