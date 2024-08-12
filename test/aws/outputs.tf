output "logs_bucket_id" {
  value = aws_s3_bucket.logs.id
}

output "test_bucket_id" {
  value = aws_s3_bucket.test.id
}

output "instance_id" {
  value = aws_instance.example_instance.id
}

output "db_instance_id" {
  value = aws_db_instance.example_db.id
}

