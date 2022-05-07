// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package docs

import (
	"github.com/atomix/cli/cmd/atomix/internal/env"
	"github.com/spf13/cobra"
)

func GetCommand(_ env.Environment) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "docs",
		Args: cobra.NoArgs,
		RunE: run,
	}
	cmd.Flags().StringP("output", "o", ".", "the path to which to output the docs")
	cmd.Flags().Bool("markdown", false, "generate docs in markdown format")
	cmd.Flags().Bool("man", false, "generate docs in man format")
	cmd.Flags().Bool("yaml", false, "generate docs in YAML format")
	_ = cmd.MarkFlagDirname("output")
	return cmd
}
