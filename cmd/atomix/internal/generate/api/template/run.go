// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"github.com/atomix/cli/cmd/atomix/internal/exec"
	"github.com/spf13/cobra"
)

func runNative(cmd *cobra.Command, _ []string) error {
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	config, err := ParseConfigFile(configPath)
	if err != nil {
		return err
	}

	args, err := cmd.Flags().GetStringToString("args")
	if err != nil {
		return err
	}
	return Generate(config, args)
}

func runInDocker(cmd *cobra.Command, args []string) error {
	return exec.InDocker()
}
