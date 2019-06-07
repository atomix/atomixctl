# Atomix CLI

This project provides a CLI for [Atomix 4].

![Atomix CLI](https://media.giphy.com/media/JR7It5AxK2rPbP6TI0/giphy.gif)

To install the CLI, run:

```bash
> go get -u github.com/atomix/atomix-cli/cmd/atomix
```

Once the CLI has been installed, initialize the settings:

```bash
> atomix init
Created ~/.atomix/config.yaml
```

The configuration file created in your home directory is used by the CLI
to connect to the Atomix controller, provide default namespaces and application
names, etc. It's also used to store configuration changes made by the CLI.

To configure completion for the CLI, source the output of `atomix completion` with
the desired shell argument:

```bash
source <(atomix completion bash)
```

The CLI can be used to deploy and connect to Atomix controllers. To deploy a Kubernetes
controller, pipe the output of `atomix controller k8s deploy` to `kubectl`:

```bash
> atomix controller k8s deploy -s atomix-controller -n kube-system | kubectl apply -f -
```

To connect the CLI to an existing Kubernetes controller, use `k8s connect`:

```bash
> atomix controller k8s connect
atomix-controller.kube-system.svc.cluster.local:5679
```

For containerized environments like Kubernetes, a Docker image is provided. The image
can be build by simply running:

```bash
> make
go build -o build/_output/bin/atomix ./cmd/cli
docker build . -f build/Dockerfile -t atomix/atomix-cli:latest
...
```

To use the CLI in Kubernetes, run the `atomix/atomix-cli` Docker image in
a single pod deployment:

```bash
> kubectl run atomix-cli --rm -it --image atomix/atomix-cli:latest --image-pull-policy "IfNotPresent" --restart "Never"
```

This command will run the CLI image as a `Deployment` and log into the bash shell.
Once you've joined the container, be sure to connect to the Atomix controller by running:

```bash
> atomix controller k8s connect
atomix-controller.kube-system.svc.cluster.local:5679
```

Once connected, you should be able to see the partition groups deployed in the
Kubernetes cluster:

```bash
> atomix groups
...
```

Once the shell is exited, the `Deployment` will be deleted.

[Atomix 4]: https://github.com/atomix/atomix/tree/4.0
