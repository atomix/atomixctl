// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package atomix

import (
	"github.com/atomix/cli/internal/atomix/generate"
	"github.com/atomix/cli/internal/atomix/pull"
	"github.com/atomix/cli/internal/atomix/push"
	"github.com/atomix/cli/internal/atomix/version"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "atomix",
	}
	cmd.AddCommand(
		generate.GetCommand(),
		pull.GetCommand(),
		push.GetCommand(),
		version.GetCommand())
	return cmd
}
