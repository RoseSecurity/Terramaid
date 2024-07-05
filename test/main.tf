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

resource "aws_instance" "example_instance" {
  ami           = "ami-0c94855ba95c574c8"
  instance_type = "t2.micro"

  tags = {
    Name        = "Dev Example Instance"
    Environment = "dev"
  }
}

resource "aws_db_instance" "example_db" {
  allocated_storage    = 20
  storage_type         = "gp2"
  engine               = "mysql"
  engine_version       = "5.7"
  instance_class       = "db.t2.micro"
  username             = "foo"
  password             = "foobarbaz"
  parameter_group_name = "default.mysql5.7"

  tags = {
    Name        = "Dev Example DB"
    Environment = "dev"
  }
}
