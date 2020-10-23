export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

ifdef VERSION
CLI_VERSION := $(VERSION)
else
CLI_VERSION := latest
endif

build: deps
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/amd64/atomix ./cmd/atomix
	GOOS=linux GOARCH=386 go build -o bin/linux/386/atomix ./cmd/atomix
	GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64/atomix ./cmd/atomix
	GOOS=windows GOARCH=386 go build -o bin/windows/386/atomix.exe ./cmd/atomix
	GOOS=windows GOARCH=amd64 go build -o bin/windows/amd64/atomix.exe ./cmd/atomix

deps:
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

generate-docs:
	go run github.com/atomix/cli/cmd/atomix-generate-docs

image: build
	docker build . -f build/Dockerfile -t atomix/cli:${CLI_VERSION}

push: # @HELP push kubernetes-controller Docker image
	docker push atomix/cli:${CLI_VERSION}
