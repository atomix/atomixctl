export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

.PHONY: build

all: build docker

build:
	go build -o build/_output/bin/atomix ./cmd/atomix

docker:
	docker build . -f build/Dockerfile -t atomix/atomix-cli:latest