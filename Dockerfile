# Copyright (c) RoseSecurity
# SPDX-License-Identifier: Apache-2.0

FROM golang:alpine@sha256:68932fa6d4d4059845c8f40ad7e654e626f3ebd3706eef7846f319293ab5cb7a AS builder
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

FROM alpine:3.22.0@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715
COPY --from=builder /usr/src/terramaid/terraform /usr/local/bin/terraform
COPY --from=builder /usr/src/terramaid/terramaid /usr/local/bin/terramaid

RUN apk update && apk add --no-cache git
USER nobody

#Set the working directory for Terramaid
WORKDIR /usr/src/terramaid

# Set the entrypoint and default command
ENTRYPOINT ["/usr/local/bin/terramaid"]
