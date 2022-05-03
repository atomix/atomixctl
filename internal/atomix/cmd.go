// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package atomix

import (
	"github.com/atomix/cli/internal/atomix/build"
	"github.com/atomix/cli/internal/atomix/completion"
	"github.com/atomix/cli/internal/atomix/generate"
	"github.com/atomix/cli/internal/atomix/version"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                    "atomix",
		BashCompletionFunction: completion.BashCompletionFunction,
	}
	cmd.AddCommand(
		build.GetCommand(),
		completion.GetCommand(),
		generate.GetCommand(),
		version.GetCommand())
	return cmd
}
