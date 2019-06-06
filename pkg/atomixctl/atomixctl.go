package atomixctl

import (
	"fmt"
	"github.com/atomix/atomixctl/pkg/atomixctl/command"
	"os"
)

func Execute() {
	rootCmd := command.GetRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
