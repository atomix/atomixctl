// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/atomix/cli/cmd/atomix/internal/runtime/add/atom"
	"github.com/atomix/cli/cmd/atomix/internal/runtime/add/driver"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "add",
	}
	cmd.AddCommand(
		atom.GetCommand(),
		driver.GetCommand())
	return cmd
}
