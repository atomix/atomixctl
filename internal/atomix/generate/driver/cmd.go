// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"github.com/spf13/cobra"
	"os"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "driver",
		Args: cobra.MaximumNArgs(1),
		RunE: runCommand,
	}
	cmd.Flags().StringP("name", "n", "", "the driver name")
	cmd.Flags().StringP("module", "m", "", "the driver module path")
	cmd.Flags().StringP("runtime", "r", "", "the runtime version for which to generate the driver")
	return cmd
}

func runCommand(cmd *cobra.Command, args []string) error {
	var context TemplateContext

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}
	context.Driver.Name = name

	path, err := cmd.Flags().GetString("module")
	if err != nil {
		return err
	}
	context.Module.Path = path

	runtimeVersion, err := cmd.Flags().GetString("runtime")
	if err != nil {
		return err
	}
	context.Runtime.Version = runtimeVersion

	var dir string
	if len(args) == 1 {
		dir = args[0]
	} else {
		dir, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	return generate(dir, context)
}
