# Copyright (c) RoseSecurity
# SPDX-License-Identifier: Apache-2.0

FROM golang:alpine@sha256:660f0b83cf50091e3777e4730ccc0e63e83fea2c420c872af5c60cb357dcafb2 AS builder
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

FROM alpine:3.23.3@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659
COPY --from=builder /usr/src/terramaid/terraform /usr/local/bin/terraform
COPY --from=builder /usr/src/terramaid/terramaid /usr/local/bin/terramaid

RUN apk update && apk add --no-cache git
USER nobody

#Set the working directory for Terramaid
WORKDIR /usr/src/terramaid

# Set the entrypoint and default command
ENTRYPOINT ["/usr/local/bin/terramaid"]
