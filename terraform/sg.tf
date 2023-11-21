resource "aws_security_group" "lb_sg" {
  name        = "load balancer security group"
  description = "Allow HTTP to be accessed"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description      = "Allow HTTP requests"
    protocol         = "tcp"
    from_port        = 80
    to_port          = 80
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  #  ingress {
  #    description      = "Allow HTTPS requests"
  #    protocol         = "tcp"
  #    from_port        = 443
  #    to_port          = 443
  #    cidr_blocks      = ["0.0.0.0/0"]
  #    ipv6_cidr_blocks = ["::/0"]
  #  }

  egress {
    description      = "Allow API traffic out"
    from_port        = 5000
    to_port          = 5000
    protocol         = "tcp"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
}

resource "aws_security_group" "ec2_sg" {
  ingress {
    description = "Allow request to the API port"
    protocol    = "tcp"
    from_port   = 5000
    to_port     = 5000
  }
}