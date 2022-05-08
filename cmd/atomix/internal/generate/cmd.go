// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package generate

import (
	"fmt"
	"github.com/atomix/cli/cmd/atomix/internal/config"
	"github.com/atomix/cli/cmd/atomix/internal/exec"
	"github.com/spf13/cobra"
	"os"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen"},
	}

	config := config.Load()
	for _, gen := range config.Generators {
		cmd.AddCommand(getGeneratorCommand(gen))
	}
	return cmd
}

func getGeneratorCommand(generator config.GeneratorConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use: generator.Name,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return exec.Run("docker", append([]string{
				"run",
				"--rm",
				"-i",
				"-v",
				fmt.Sprintf("%s:/build", dir),
				generator.Image,
			}, args...)...)
		},
	}
	cmd.DisableFlagParsing = true
	return cmd
}
