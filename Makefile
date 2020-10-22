export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

build: deps
	GOOS=linux GOARCH=amd64 go build -o build/_output/bin/atomix ./cmd/atomix

deps:
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

generate-docs:
	go run github.com/atomix/cli/cmd/atomix-generate-docs

images: build
	docker build . -f build/Dockerfile -t atomix/cli:latest
