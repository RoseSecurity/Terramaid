FROM golang:1.22 AS terramaid

ARG TERRAFORM_VERSION=1.9.2

# Install Terramaid
RUN go install github.com/RoseSecurity/terramaid@latest

# Install required packages
RUN apt-get update && \
    apt-get install -y zsh-autosuggestions 

# Change shell to zsh
RUN chsh -s /bin/zsh vscode

# Install Terraform
RUN curl -fsSL https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip -o terraform.zip && \
    unzip terraform.zip && \
    mv terraform /usr/local/bin/ && \
    rm terraform.zip

# Install the starship prompt
RUN curl -sS https://starship.rs/install.sh | sh -s -- -y

# Install the rootfs/ configurations
COPY rootfs/ /

