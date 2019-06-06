export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

.PHONY: build

all: build

build:
	go build -o build/_output/bin/atomixctl ./cmd/atomixctl
	docker build . -f build/Dockerfile -t atomix/atomixctl:latest