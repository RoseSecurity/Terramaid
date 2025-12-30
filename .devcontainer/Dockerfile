# Copyright (c) RoseSecurity
# SPDX-License-Identifier: Apache-2.0

FROM golang:1.25@sha256:31c1e53dfc1cc2d269deec9c83f58729fa3c53dc9a576f6426109d1e319e9e9a AS build

# Install Terramaid
RUN go install github.com/RoseSecurity/terramaid@latest

FROM mcr.microsoft.com/vscode/devcontainers/base:debian@sha256:30b0a0c004ca94d36c323ee993361a7e0ae25ea255ea125201e8a9587501c324

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

