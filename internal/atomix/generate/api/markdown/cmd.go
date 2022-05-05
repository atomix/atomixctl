// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package markdown

import (
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "markdown",
		Short:   "Generates Markdown documentation from Protobuf sources",
		Example: "atomix generate api markdown --input ../api --pattern '**/*.proto' --output ./docs",
		Aliases: []string{"md"},
		Args:    cobra.NoArgs,
		RunE:    run,
	}
	cmd.Flags().StringP("input", "i", ".", "the path to the root of the Protobuf sources")
	cmd.Flags().StringP("pattern", "f", "**/*.proto", "a pattern by which to filter Protobuf sources")
	cmd.Flags().StringP("output", "o", ".", "the path to which to write generated docs")
	return cmd
}
