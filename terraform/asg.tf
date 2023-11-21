
data "aws_ami" "amazon-linux" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn-ami-hvm-*-x86_64-ebs"]
  }
}


resource "aws_launch_configuration" "api_lc" {
  name_prefix     = "api-lc"
  image_id        = data.aws_ami.amazon-linux.id
  instance_type   = "t3a.micro"
  security_groups = [aws_security_group.ec2_sg.id]

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "api_asg" {
  launch_configuration = aws_launch_configuration.api_lc.id
  max_size             = 1
  min_size             = 1
  desired_capacity     = 1

  health_check_grace_period = 300
  health_check_type         = "ELB"
  availability_zones        = module.vpc.azs

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_attachment" "asg_attachment" {
  autoscaling_group_name = aws_autoscaling_group.api_asg.id
  elb = aws_elb.api_elb.id
}