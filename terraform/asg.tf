data "aws_ami" "amazon-linux" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn-ami-hvm-*-x86_64-ebs"]
  }
}

resource "aws_launch_template" "ecs_lt" {
  name_prefix   = "esc-lt"
  image_id      = data.aws_ami.amazon-linux.id
  instance_type = "t3a.micro"

  key_name = "ec2ecs"
  vpc_security_group_ids = [aws_security_group.ec2_sg.id]

  lifecycle {
    create_before_destroy = true
  }

  user_data = filebase64("${path.module}/ecs.sh")

  tag_specifications {
    resource_type = "instance"
    tags = {
      Name = "ecs-instance"
    }
  }
}

resource "aws_autoscaling_group" "ecs_asg" {
  vpc_zone_identifier       = module.vpc.public_subnets
  desired_capacity = 1
  max_size         = 1
  min_size         = 1


  launch_template {
    id      = aws_launch_template.ecs_lt.id
    version = "$Latest"
  }

  lifecycle {
    create_before_destroy = true
  }
}