# Copyright RoseSecurity 2024, 2026
# SPDX-License-Identifier: Apache-2.0

FROM golang:1.25@sha256:8a7adc288b77e9b787cd2695029eb54d10ae80571b21d44fed68d067ad0a9c96 AS build

# Install Terramaid
RUN go install github.com/RoseSecurity/terramaid@latest

FROM mcr.microsoft.com/vscode/devcontainers/base:debian@sha256:e389149057371298eaeb72e736fc2c4dcb79974de222dd471d2e6d675d2f61c4

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
