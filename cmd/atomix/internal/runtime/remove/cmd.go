// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package remove

import (
	"github.com/atomix/cli/cmd/atomix/internal/runtime/remove/atom"
	"github.com/atomix/cli/cmd/atomix/internal/runtime/remove/driver"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "remove",
	}
	cmd.AddCommand(
		atom.GetCommand(),
		driver.GetCommand())
	return cmd
}
