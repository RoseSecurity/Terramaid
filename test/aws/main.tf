# S3 Buckets
resource "aws_s3_bucket" "logs" {
  bucket = "dev-logs"
  tags = {
    Name        = "Dev Logs"
    Environment = "dev"
  }
}

resource "aws_s3_bucket_policy" "logs_policy" {
  bucket = aws_s3_bucket.logs.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect    = "Allow"
        Principal = "*"
        Action    = ["s3:GetObject"]
        Resource  = ["${aws_s3_bucket.logs.arn}/*"]
      },
    ]
  })
}

resource "aws_s3_bucket" "test" {
  bucket = "dev-test"
  tags = {
    Name        = "Dev Test"
    Environment = "dev"
  }
}

resource "aws_s3_bucket_policy" "test_policy" {
  bucket = aws_s3_bucket.test.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect    = "Allow"
        Principal = "*"
        Action    = ["s3:GetObject"]
        Resource  = ["${aws_s3_bucket.test.arn}/*"]
      },
    ]
  })
}

# EC2 Instances
resource "aws_instance" "web_server" {
  count         = 2
  ami           = "ami-0c94855ba95c574c8"
  instance_type = "t2.micro"
  tags = {
    Name        = "Dev Web Server ${count.index + 1}"
    Environment = "dev"
  }
}

resource "aws_instance" "app_server" {
  ami           = "ami-0c94855ba95c574c8"
  instance_type = "t2.small"
  tags = {
    Name        = "Dev App Server"
    Environment = "dev"
  }
}

# RDS Database
resource "aws_db_instance" "main_db" {
  identifier           = "dev-main-db"
  allocated_storage    = 20
  storage_type         = "gp2"
  engine               = "mysql"
  engine_version       = "5.7"
  instance_class       = "db.t2.micro"
  username             = "admin"
  password             = "password123"
  parameter_group_name = "default.mysql5.7"
  skip_final_snapshot  = true
  tags = {
    Name        = "Dev Main DB"
    Environment = "dev"
  }
}

# VPC and Networking
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
  tags = {
    Name        = "Dev VPC"
    Environment = "dev"
  }
}

resource "aws_subnet" "public" {
  count             = 2
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.${count.index + 1}.0/24"
  availability_zone = "us-west-2${["a", "b"][count.index]}"
  tags = {
    Name        = "Dev Public Subnet ${count.index + 1}"
    Environment = "dev"
  }
}

resource "aws_subnet" "private" {
  count             = 2
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.${count.index + 3}.0/24"
  availability_zone = "us-west-2${["a", "b"][count.index]}"
  tags = {
    Name        = "Dev Private Subnet ${count.index + 1}"
    Environment = "dev"
  }
}

# Security Groups
resource "aws_security_group" "web" {
  name        = "allow_web"
  description = "Allow inbound web traffic"
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "HTTP from anywhere"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "Dev Web SG"
    Environment = "dev"
  }
}

resource "aws_security_group" "db" {
  name        = "allow_db"
  description = "Allow inbound database traffic"
  vpc_id      = aws_vpc.main.id

  ingress {
    description     = "MySQL from web servers"
    from_port       = 3306
    to_port         = 3306
    protocol        = "tcp"
    security_groups = [aws_security_group.web.id]
  }

  tags = {
    Name        = "Dev DB SG"
    Environment = "dev"
  }
}

# Load Balancer
resource "aws_lb" "web" {
  name               = "dev-web-lb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.web.id]
  subnets            = aws_subnet.public[*].id

  tags = {
    Name        = "Dev Web LB"
    Environment = "dev"
  }
}

resource "aws_lb_target_group" "web" {
  name     = "dev-web-tg"
  port     = 80
  protocol = "HTTP"
  vpc_id   = aws_vpc.main.id
}

resource "aws_lb_listener" "web" {
  load_balancer_arn = aws_lb.web.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.web.arn
  }
}

# Attach instances to target group
resource "aws_lb_target_group_attachment" "web" {
  count            = 2
  target_group_arn = aws_lb_target_group.web.arn
  target_id        = aws_instance.web_server[count.index].id
  port             = 80
}
