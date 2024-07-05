# GitHub Actions Integration

```yaml
name: Terramaid

on:
  pull_request:
    paths:
      - '**/*.tf'

jobs:
  run-terraform-check:
    permissions:
      pull-requests: write

    runs-on: ubuntu-latest

    steps:
    - name: Checkout repository
      uses: actions/checkout@v4

    - name: Outputs
      id: vars
      run: |
        echo "terramaid_version=$(curl -s https://api.github.com/repos/RoseSecurity/Terramaid/releases/latest | grep tag_name | cut -d '"' -f 4)" >> $GITHUB_OUTPUT
        case "${{ runner.arch }}" in
          "X64" )
            echo "arch=x86_64" >> $GITHUB_OUTPUT
            ;;
          "ARM64" )
            echo "arch=arm64" >> $GITHUB_OUTPUT
            ;;
        esac

    - name: Setup Go
      uses: actions/setup-go@v5
      with:
        go-version: 'stable'

    - name: Setup Terramaid
      run: |
        curl -L -o /tmp/terramaid.tar.gz "https://github.com/RoseSecurity/Terramaid/releases/download/${{ steps.vars.outputs.terramaid_version }}/Terramaid_Linux_${{ steps.vars.outputs.arch }}.tar.gz"
        tar -xzvf /tmp/terramaid.tar.gz -C /tmp
        mv -v /tmp/Terramaid /usr/local/bin/terramaid
        chmod +x /usr/local/bin/terramaid

    - name: Init
      run: terraform init

    - name: Terramaid
      id: terramaid
      run: |
        /usr/local/bin/terramaid

    - name: Upload comment to PR
      uses: actions/github-script@v7
      with:
        github-token: ${{ secrets.GITHUB_TOKEN }}
        script: |
          const fs = require('fs');
          const terramaid = fs.readFileSync('Terramaid.md', 'utf8');
          github.rest.issues.createComment({
            owner: context.repo.owner,
            repo: context.repo.repo,
            issue_number: context.issue.number,
            body: `## Terraform Plan\n\n${terramaid}`
          })
```
