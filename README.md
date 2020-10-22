[![Build Status](https://travis-ci.org/atomix/cli.svg?branch=master)](https://travis-ci.org/atomix/cli)
[![Integration Test Status](https://img.shields.io/travis/atomix/go-framework?label=Integration%20Tests&logo=Integration)](https://travis-ci.org/onosproject/onos-test)
[![Go Report Card](https://goreportcard.com/badge/github.com/atomix/cli)](https://goreportcard.com/report/github.com/atomix/cli)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/gojp/goreportcard/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/atomix/kubernetes-benchmarks?status.svg)](https://godoc.org/github.com/atomix/cli)


# Atomix CLI

This project provides a CLI for [Atomix 4].

![Atomix CLI](https://media.giphy.com/media/cImqHbP1Bt2u5ZTVGg/giphy.gif)

## Installation

To install the CLI, run:

```bash
$ go get -u github.com/atomix/cli/cmd/atomix
```

## Configuration

To configure completion for the CLI, source the output of `atomix completion` with
the desired shell argument:

```bash
$ source <(atomix completion bash)
```

```bash
$ source $(atomix completion zsh)
```

To run the CLI in a [Kubernetes] cluster:

```bash
$ kubectl run atomix-cli --rm -it --image atomix/cli:latest --restart Never
```

Once the CLI has been installed, you can configure the CLI using the `atomix config`
suite of sub-commands:

```bash
$ atomix config get controller
atomix-controller.default.svc.cluster.local:5679
$ atomix config set controller atomix-controller.kube-system.svc.cluster.local:5679
atomix-controller.kube-system.svc.cluster.local:5679
```

## Usage

The CLI provides commands for managing distributed primitives in an Atomix database.
To see the primitives supported by the CLI, use the `--help` flag:

```bash
$ atomix --help
```

The CLI and primitive commands can be run from two different contexts. Using the shell,
primitives can be queried and modified by name:

```bash
$ atomix map "my-map" put "foo" "Hello world!"
key: "foo"
value: "Hello world!"
version: 1
$ atomix map "my-map" get "foo"
value: "Hello world!"
version: 1
```

The CLI also provides an interactive shell for operating within a specific context,
e.g. for commands operating on a specific primitive:

```bash
$ atomix map "my-map"
map:my-map> put "foo" "Hello world!"
key: "foo"
value: "Hello world!"
version: 1
map:my-map> get "foo"
value: "Hello world!"
version: 1
map:my-map> exit
```

[Atomix 4]: https://github.com/atomix/atomix/tree/4.0
