// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	_go "github.com/atomix/cli/internal/atomix/generate/api/go"
	"github.com/atomix/cli/internal/atomix/generate/api/markdown"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "api",
		RunE: runCommand,
	}
	cmd.AddCommand(_go.GetCommand())
	cmd.AddCommand(markdown.GetCommand())
	return cmd
}

func runCommand(cmd *cobra.Command, args []string) error {
	return nil
}
