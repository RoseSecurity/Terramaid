output "ubuntu_ami" {
  value       = data.aws_ami.ubuntu.id
  description = "Ubuntu AWS AMI ID"
}
