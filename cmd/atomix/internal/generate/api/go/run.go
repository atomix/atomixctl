// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package _go

import (
	"github.com/atomix/cli/cmd/atomix/internal/exec"
	"github.com/spf13/cobra"
	"os"
)

func runNative(cmd *cobra.Command, args []string) error {
	var config Config
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	if configPath != "" {
		config, err = ParseConfigFile(configPath)
		if err != nil {
			return err
		}
	}

	protoPath, err := cmd.Flags().GetString("proto-path")
	if err != nil {
		return err
	}
	config.Proto.Path = protoPath

	protoPatterns, err := cmd.Flags().GetStringSlice("proto-pattern")
	if err != nil {
		return err
	}
	config.Proto.Patterns = protoPatterns

	goPath, err := cmd.Flags().GetString("go-path")
	if err != nil {
		return err
	}
	err = os.MkdirAll(goPath, 0755)
	if err != nil {
		return err
	}
	config.Go.Path = goPath

	importPath, err := cmd.Flags().GetString("import-path")
	if err != nil {
		return err
	}
	config.Go.ImportPath = importPath
	return Generate(config)
}

func runInDocker(cmd *cobra.Command, args []string) error {
	return exec.InDocker()
}
