package cli

import (
	"fmt"
	"github.com/atomix/atomix-cli/pkg/cli/command"
	"os"
)

func Execute() {
	rootCmd := command.GetRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
