# Installation

## Homebrew installation:

```sh
brew install terramaid
```

## Golang
If you have a functional go environment, you can install with:

```sh
go install github.com/RoseSecurity/terramaid@latest
```

Build from source:

```sh
git clone git@github.com:RoseSecurity/terramaid.git
cd terramaid
make build
```

## Docker Image

Run the following command to utilize the Terramaid Docker image:

```sh
docker run -it -v $(pwd):/usr/src/terramaid rosesecurity/terramaid:latest
```