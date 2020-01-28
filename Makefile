export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

all: build images

build: deps
	go build -o build/_output/bin/atomix ./cmd/atomix

deps:
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

images:
	docker build . -f build/Dockerfile -t atomix/cli:latest