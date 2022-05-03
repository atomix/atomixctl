package main

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
)

var (
    version string
    commit string
)

func main() {
    cmd := &cobra.Command{
        Use: {{ .Driver.Name | toKebab | quote }},
        RunE: func(cmd *cobra.Command, args []string) error {
            return nil
        },
    }

    if err := cmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
