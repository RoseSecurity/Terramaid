output "logs_bucket_id" {
  value = aws_s3_bucket.logs.id
}

output "test_bucket_id" {
  value = aws_s3_bucket.test.id
}

output "instance_id" {
  value = one(aws_instance.web_server[*].id)
}

output "db_instance_id" {
  value = aws_db_instance.main_db.id
}

