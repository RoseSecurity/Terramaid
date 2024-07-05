FROM golang:alpine

# Install necessary dependencies
RUN apk update && apk add --no-cache \
    bash \
    curl \
    git \
    unzip

# Install Terraform
RUN curl -fsSL https://releases.hashicorp.com/terraform/1.9.1/terraform_1.9.1_linux_amd64.zip -o terraform.zip && \
    unzip terraform.zip && \
    mv terraform /usr/local/bin/ && \
    rm terraform.zip

# Set the working directory for Terramaid
WORKDIR /usr/src/terramaid

# Copy the source code and build
COPY . .
RUN go mod download && go mod verify
RUN go build -v -o /usr/local/bin/terramaid main.go

# Set the entrypoint and default command
ENTRYPOINT ["/usr/local/bin/terramaid"]

