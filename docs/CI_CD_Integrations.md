# CI/CD Integration

## GitHub Actions

```yaml
name: Terramaid

on:
  pull_request:
    paths:
      - '**/*.tf'

jobs:
  run-terraform-check:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout repository
      uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.22'

    - name: Download Go binary
      run: |
        curl -L -o /usr/local/bin/terramaid https://github.com/RoseSecurity/Terramaid/releases/download/v0.1.0/Terramaid_0.1.0_linux_amd64
        chmod +x /usr/local/bin/terramaid

    - name: Init
      run: terraform init

    - name: Terramaid
      id: terramaid
      run: |
        ./usr/local/bin/terramaid

    - name: Upload comment to PR
      uses: actions/github-script@v6
      with:
        script: |
          const fs = require('fs');
          const terramaid = fs.readFileSync('Terramaid.md', 'utf8');
          github.rest.issues.createComment({
            owner: context.repo.owner,
            repo: context.repo.repo,
            issue_number: context.issue.number,
            body: `## Terraform Plan\n\n${terramaid}`
          });
```

## GitLab CI/CD Pipelines

**CI/CD Example:**

<p align="left">
<img width=60% height=60% src="./docs/img/Terramaid-Demo.png">
</p>