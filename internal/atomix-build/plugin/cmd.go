// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package plugin

import (
	"github.com/atomix/runtime-api/pkg/plugin"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "plugin",
		RunE: runCommand,
	}
	return cmd
}

func runCommand(cmd *cobra.Command, args []string) error {
	var _ plugin.Repository
	return nil
}
