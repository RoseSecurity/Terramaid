# Copyright (c) RoseSecurity
# SPDX-License-Identifier: Apache-2.0

FROM golang:alpine@sha256:72567335df90b4ed71c01bf91fb5f8cc09fc4d5f6f21e183a085bafc7ae1bec8 AS builder
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

FROM alpine:3.23.2@sha256:865b95f46d98cf867a156fe4a135ad3fe50d2056aa3f25ed31662dff6da4eb62
COPY --from=builder /usr/src/terramaid/terraform /usr/local/bin/terraform
COPY --from=builder /usr/src/terramaid/terramaid /usr/local/bin/terramaid

RUN apk update && apk add --no-cache git
USER nobody

#Set the working directory for Terramaid
WORKDIR /usr/src/terramaid

# Set the entrypoint and default command
ENTRYPOINT ["/usr/local/bin/terramaid"]
