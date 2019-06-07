package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
)

const (
	kubeControllerUrl = "https://raw.githubusercontent.com/atomix/atomix-k8s-controller/master/deploy/atomix-controller.yaml"
)

func newControllerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "controller <type>",
		Short: "Deploy and manage the Atomix controller",
		Run:   runControllerCommand,
	}
	cmd.AddCommand(newControllerKubernetesCommand())
	return cmd
}

func runControllerCommand(cmd *cobra.Command, args []string) {
	controller := getClientController()
	ExitWithOutput(controller)
}

func newControllerKubernetesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "kubernetes [connect,deploy]",
		Aliases: []string{"kube", "k8s"},
		Short: "Deploy and manage a Kubernetes controller",
	}
	cmd.AddCommand(newControllerKubernetesConnectCommand())
	cmd.AddCommand(newControllerKubernetesDeployCommand())
	return cmd
}

func newControllerKubernetesConnectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "connect",
		Short: "Connect to a Kubernetes controller",
		Run: runControllerConnectKubernetesCommand,
	}
	cmd.Flags().String("namespace",  "kube-system", "the namespace in which to deploy the controller")
	cmd.Flags().String("name", "atomix-controller", "the controller name")
	return cmd
}

func runControllerConnectKubernetesCommand(cmd *cobra.Command, args []string) {
	namespace, _ := cmd.Flags().GetString("namespace")
	name, _ := cmd.Flags().GetString("name")
	service := fmt.Sprintf("%s.%s.svc.cluster.local:5679", name, namespace)
	viper.Set("controller", service)
	if err := viper.WriteConfig(); err != nil {
		ExitWithError(ExitError, err)
	} else {
		value := viper.Get("controller")
		ExitWithOutput(value)
	}
}

func newControllerKubernetesDeployCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "deploy",
		Short:   "Deploy a Kubernetes controller",
		Run:     runControllerDeployKubernetesCommand,
	}
}

func runControllerDeployKubernetesCommand(cmd *cobra.Command, args []string) {
	manifest, err := downloadFile(kubeControllerUrl)
	if err != nil {
		ExitWithError(ExitError, err)
	}
	io.Copy(os.Stdout, manifest)
	ExitWithSuccess()
}

func downloadFile(url string) (io.ReadCloser, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}
