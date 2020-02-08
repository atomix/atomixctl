# Installation

The Atomix CLI is a command line tool for managing distributed primitives.

To install the CLI, use the `go` CLI:

```bash
$ go get -u github.com/atomix/cli/cmd/atomix
```

Alternatively, you can deploy the Atomix CLI image directly to a Kubernetes cluster:

```bash
$ kubectl run atomix-cli --rm -it --image atomix/cli:latest --restart Never
```
