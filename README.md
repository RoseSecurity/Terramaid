# Terramaid

<p align="center">
<img width=100% height=100% src="./docs/img/Terramaid.png">
</p>

<p align="center">
  <em>A utility for creating Mermaid diagrams from Terraform configurations</em>
</p>

## Introduction

Terramaid transforms your Terraform resources and plans into visually appealing Mermaid diagrams. By converting complex infrastructure into easy-to-understand diagrams, Terramaid enhances documentation, simplifies review processes, and fosters better collaboration among team members. Whether you're looking to enrich your project's documentation, streamline reviews, or just bring a new level of clarity to your Terraform configurations, Terramaid is the perfect utility to integrate into your development workflow.

## Demo

<p align="center">
<img width=100% height=100% src="./docs/img/terramaid_vhs_demo.gif">
</p>

### Output

```mermaid
flowchart TD;
        subgraph Terraform
                aws_db_instance.dev_example_db_instance["aws_db_instance.dev_example_db_instance"]
                aws_instance.dev_example_instance["aws_instance.dev_example_instance"]
                aws_s3_bucket.dev_logs_bucket["aws_s3_bucket.dev_logs_bucket"]
                aws_s3_bucket.dev_test_bucket["aws_s3_bucket.dev_test_bucket"]
                aws_s3_bucket_policy.dev_logs_bucket_policy["aws_s3_bucket_policy.dev_logs_bucket_policy"]
                aws_s3_bucket_policy.dev_test_bucket_policy["aws_s3_bucket_policy.dev_test_bucket_policy"]
                aws_s3_bucket_policy.dev_logs_bucket_policy --> aws_s3_bucket.dev_logs_bucket
                aws_s3_bucket_policy.dev_test_bucket_policy --> aws_s3_bucket.dev_test_bucket
        end
```

## Installation

If you have a functional go environment, you can install with:

```sh
go install github.com/RoseSecurity/terramaid@latest
```

Build from source:

```sh
git clone git@github.com:RoseSecurity/terramaid.git
cd terramaid
make build
```

### Docker Image

Run the following command to utilize the Terramaid Docker image:

```sh
docker run -it -v $(pwd):/usr/src/terramaid rosesecurity/terramaid:latest
```

## CI/CD Integrations

Terramaid is designed to easily integrate with existing pipelines and workflows. For more information on sample GitHub Actions and GitLab CI/CD Pipelines, feel free to check out [GitHub Actions Integrations](./docs/GitHub_Actions_Integration.md) and [Gitlab Pipelines Integrations](./docs/GitLab_Pipelines_Integration.md).
