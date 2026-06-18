# Copyright RoseSecurity 2024, 2026
# SPDX-License-Identifier: Apache-2.0

FROM golang:1.25@sha256:918cfb186a8fc7b64a91a61a2b83c11f529e9002ba5a3010986f4710085ca30d AS build

# Install Terramaid
RUN go install github.com/RoseSecurity/terramaid@latest

FROM mcr.microsoft.com/vscode/devcontainers/base:debian@sha256:b60cd6e70901c42c23dc312b8e0cf2c704cad3acf65c55739ff1ef33a30166d2

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
