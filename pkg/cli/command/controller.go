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
		Use:   "controller {set,get,deploy}",
		Short: "Deploy and manage the Atomix controller",
		Run:   runControllerGetCommand,
	}
	cmd.AddCommand(newControllerSetCommand())
	cmd.AddCommand(newControllerGetCommand())
	cmd.AddCommand(newControllerDeployCommand())
	return cmd
}

func newControllerSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <address>",
		Short: "Set the controller address",
		Args:  cobra.ExactArgs(1),
		Run:   runControllerSetCommand,
	}
	cmd.AddCommand(newControllerSetKubernetesCommand())
	return cmd
}

func runControllerSetCommand(cmd *cobra.Command, args []string) {
	if err := setClientController(args[0]); err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(getClientController())
	}
}

func newControllerSetKubernetesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "kubernetes",
		Aliases: []string{"kube", "k8s"},
		Short:   "Configure a Kubernetes controller",
		Run:     runControllerSetKubernetesCommand,
	}
	cmd.Flags().String("namespace", "kube-system", "the controller namespace")
	cmd.Flags().String("name", "atomix-controller", "the controller service name")
	return cmd
}

func runControllerSetKubernetesCommand(cmd *cobra.Command, args []string) {
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

func newControllerGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get the controller address",
		Run:   runControllerGetCommand,
	}
}

func runControllerGetCommand(cmd *cobra.Command, args []string) {
	ExitWithOutput(getClientController())
}

func newControllerDeployCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a controller",
	}
	cmd.AddCommand(newControllerDeployKubernetesCommand())
	return cmd
}

func newControllerDeployKubernetesCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "kubernetes",
		Aliases: []string{"kube", "k8s"},
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
