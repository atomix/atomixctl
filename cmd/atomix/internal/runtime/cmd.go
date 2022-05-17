// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package runtime

import (
	"github.com/atomix/cli/cmd/atomix/internal/runtime/add"
	"github.com/atomix/cli/cmd/atomix/internal/runtime/init"
	"github.com/atomix/cli/cmd/atomix/internal/runtime/remove"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "runtime",
	}
	cmd.AddCommand(
		init.GetCommand(),
		add.GetCommand(),
		remove.GetCommand())
	return cmd
}
