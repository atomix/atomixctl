// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package completion

import (
	"errors"
	"github.com/spf13/cobra"
	"os"
)

func GetCommand() *cobra.Command {
	return &cobra.Command{
		Use:       "completion <shell>",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "bash" {
				return runCompletionBash(os.Stdout, cmd.Parent())
			} else if args[0] == "zsh" {
				return runCompletionZsh(os.Stdout, cmd.Parent())
			} else {
				return errors.New("unsupported shell type " + args[0])
			}
		},
	}
}
