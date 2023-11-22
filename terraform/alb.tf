#resource "aws_elb" "api_elb" {
#  name            = "load-balancer-backend"
#  subnets         = module.vpc.public_subnets
#  security_groups = [aws_security_group.lb_sg.id]
#
#  listener {
#    instance_port     = 3000
#    instance_protocol = "http"
#    lb_port           = 80
#    lb_protocol       = "http"
#  }
#
#  health_check {
#    healthy_threshold   = 2
#    unhealthy_threshold = 2
#    timeout             = 3
#    target              = "HTTP:5000/health"
#    interval            = 30
#  }
#}

resource "aws_lb" "ecs_alb" {
  name = "ecs-alb"
  internal = false
  load_balancer_type = "application"
  security_groups = [aws_security_group.lb_sg]
  subnets = [module.vpc.public_subnets]

  tags = {
    Name = "ecs-alb"
  }
}

resource "aws_lb_listener" "ecs_alb_listener" {
  load_balancer_arn = aws_lb.ecs_alb.arn
  port  = 80
  protocol = "HTTP"

  default_action {
    type = "forward"
    target_group_arn = aws_lb_target_group.ecs_tg.arn
  }
}

resource "aws_lb_target_group" "ecs_tg" {
  name = "ecs-target-group"
  port = 80
  protocol = "HTTP"
  target_type = "ip"
  vpc_id = module.vpc.vpc_id

  health_check {
    enabled = true
    healthy_threshold = 3
    unhealthy_threshold = 3
    interval = 30
    timeout = 3
    protocol = "HTTP"
    path = "/health"
  }
}