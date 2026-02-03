# Copyright (c) RoseSecurity
# SPDX-License-Identifier: Apache-2.0

FROM golang:1.25@sha256:4c973c7cf9e94ad40236df4a9a762b44f6680560a7fb8a4e69513b5957df7217 AS build

# Install Terramaid
RUN go install github.com/RoseSecurity/terramaid@latest

FROM mcr.microsoft.com/vscode/devcontainers/base:debian@sha256:a30da48cdf5f9144ff7f2156622e701e752fc258d77ca7bb00120624f1a95938

# Install required packages
RUN apt-get update && \
    apt-get install -y apt-utils ca-certificates curl unzip zsh-autosuggestions 

# Change shell to zsh
RUN chsh -s /bin/zsh vscode

# Install Terraform
ARG TERRAFORM_VERSION=1.9.2
RUN curl -fsSL https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip -o terraform.zip && \
    unzip terraform.zip && \
    mv terraform /usr/local/bin/ && \
    rm terraform.zip

# Install the starship prompt
RUN curl -sS https://starship.rs/install.sh | sh -s -- -y

# Copy Terramaid binary from the build stage
COPY --from=build /go/bin/terramaid /usr/local/bin/terramaid

# Install the rootfs/ configurations
COPY rootfs/ /

