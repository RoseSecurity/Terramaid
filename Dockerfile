# Copyright RoseSecurity 2024, 2026
# SPDX-License-Identifier: Apache-2.0

FROM golang:alpine@sha256:f23e8b227fb4493eabe03bede4d5a32d04092da71962f1fb79b5f7d1e6c2a17f AS builder
WORKDIR /usr/src/terramaid
# Terraform version
ARG TERRAFORM_VERSION=1.10.0

# Install necessary dependencies
RUN apk update && apk add --no-cache \
    bash \
    curl \
    git \
    unzip

# Install Terraform
RUN <<EOF
    curl -fsSL https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip -o terraform.zip
    unzip terraform.zip
EOF

# Copy the source code and build
COPY . .
RUN <<EOF
    go mod download && go mod verify
    go build -v -o ./terramaid main.go
EOF

FROM alpine:3.23.4@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11
COPY --from=builder /usr/src/terramaid/terraform /usr/local/bin/terraform
COPY --from=builder /usr/src/terramaid/terramaid /usr/local/bin/terramaid

RUN apk update && apk add --no-cache git
USER nobody

#Set the working directory for Terramaid
WORKDIR /usr/src/terramaid

# Set the entrypoint and default command
ENTRYPOINT ["/usr/local/bin/terramaid"]
