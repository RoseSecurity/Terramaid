# Terramaid

<p align="center">
<img width=100% height=100% src="./docs/img/Terramaid.png">
</p>

<p align="center">
  <em>A utility for creating Mermaid diagrams from Terraform configurations</em>
</p>

**Checkout the [Docs](https://rosesecurity.github.io/terramaid/) to learn more about `terramaid`**

## Introduction

Terramaid transforms your Terraform resources and plans into visually appealing Mermaid diagrams. By converting complex infrastructure into easy-to-understand diagrams, Terramaid enhances documentation, simplifies review processes, and fosters better collaboration among team members. Whether you're looking to enrich your project's documentation, streamline reviews, or just bring a new level of clarity to your Terraform configurations, Terramaid is the perfect utility to integrate into your development workflow.

## Demo

<p align="center">
<img width=100% height=100% src="./docs/img/terramaid_vhs_demo.gif">
</p>

### Output

```mermaid
flowchart TD
 subgraph Hashicorp
 subgraph Terraform
  aws_db_instance.example_db["aws_db_instance.example_db"]
  aws_instance.example_instance["aws_instance.example_instance"]
  aws_s3_bucket.logs["aws_s3_bucket.logs"]
  aws_s3_bucket.test["aws_s3_bucket.test"]
  aws_s3_bucket_policy.logs_policy["aws_s3_bucket_policy.logs_policy"]
  aws_s3_bucket_policy.test_policy["aws_s3_bucket_policy.test_policy"]
  aws_s3_bucket_policy.logs_policy --> aws_s3_bucket.logs
  aws_s3_bucket_policy.test_policy --> aws_s3_bucket.test
 end
 end
```