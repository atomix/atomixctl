// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package build

import (
	"github.com/atomix/cli/internal/atomix/build/plugin"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "build",
	}
	cmd.AddCommand(plugin.GetCommand())
	return cmd
}
