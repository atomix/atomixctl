// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package bin

import "github.com/spf13/cobra"

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "bin",
		RunE: runCommand,
	}
	return cmd
}

func runCommand(cmd *cobra.Command, args []string) error {
	return nil
}
