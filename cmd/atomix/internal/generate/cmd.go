// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package generate

import (
	"fmt"
	"github.com/atomix/cli/cmd/atomix/internal/config"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, dir, args := args[0], args[1], args[2:]
			return runInDocker(dir, name, args...)
		},
	}

	config, _ := config.Load()
	for _, gen := range config.Generators {
		cmd.AddCommand(getGeneratorCommand(gen))
	}
	return cmd
}

func getGeneratorCommand(generator config.GeneratorConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:  generator.Name,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, args := args[1], args[2:]
			return runInDocker(dir, generator.Image, args...)
		},
	}
	cmd.DisableFlagParsing = true
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		dir, args := args[1], args[2:]
		args = append(args, "--help")
		err := runInDocker(dir, generator.Image, args...)
		if err != nil {
			os.Exit(1)
		}
	})
	return cmd
}

func runInDocker(dir string, image string, args ...string) error {
	cmd := &exec.Cmd{
		Path: "docker",
		Args: append([]string{
			"run",
			"-i",
			"-v",
			fmt.Sprintf("%s:/build", dir),
			image,
		}, args...),
		Env:    os.Environ(),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	return cmd.Run()
}
