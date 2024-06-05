resource "aws_s3_bucket" "dev_logs_bucket" {
  bucket = "dev-logs-bucket"

  tags = {
    Name        = "Dev Logs Bucket"
    Environment = "dev"
  }
}

resource "aws_s3_bucket_policy" "dev_logs_bucket_policy" {
  bucket = aws_s3_bucket.dev_logs_bucket.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect    = "Allow"
        Principal = "*"
        Action    = ["s3:GetObject"]
        Resource  = ["${aws_s3_bucket.dev_logs_bucket.arn}/*"]
      },
    ]
  })
}

resource "aws_s3_bucket" "dev_test_bucket" {
  bucket = "dev-test-bucket"

  tags = {
    Name        = "Dev Test Bucket"
    Environment = "dev"
  }

}

resource "aws_s3_bucket_policy" "dev_test_bucket_policy" {
  bucket = aws_s3_bucket.dev_test_bucket.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect    = "Allow"
        Principal = "*"
        Action    = ["s3:GetObject"]
        Resource  = ["${aws_s3_bucket.dev_test_bucket.arn}/*"]
      },
    ]
  })

}

resource "aws_instance" "dev_example_instance" {
  ami           = "ami-0c94855ba95c574c8"
  instance_type = "t2.micro"

  tags = {
    Name        = "Dev Example Instance"
    Environment = "dev"
  }

}

resource "aws_db_instance" "dev_example_db_instance" {
  allocated_storage    = 20
  storage_type         = "gp2"
  engine               = "mysql"
  engine_version       = "5.7"
  instance_class       = "db.t2.micro"
  username             = "foo"
  password             = "foobarbaz"
  parameter_group_name = "default.mysql5.7"

  tags = {
    Name        = "Dev Example DB Instance"
    Environment = "dev"
  }

}


