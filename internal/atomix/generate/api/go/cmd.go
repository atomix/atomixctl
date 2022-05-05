// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package _go

import (
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "go",
		Short:   "Generates Go sources from Protobuf sources",
		Example: "atomix generate api go --input ../api --pattern '**/*.proto' --output ./pkg --package github.com/my-user/my-repo",
		Aliases: []string{"golang"},
		Args:    cobra.NoArgs,
		RunE:    run,
	}
	cmd.Flags().StringP("input", "i", ".", "the path to the root of the Protobuf sources")
	cmd.Flags().StringP("pattern", "f", "**/*.proto", "a pattern by which to filter Protobuf sources")
	cmd.Flags().StringP("output", "o", ".", "the path to which to write generated Go sources")
	cmd.Flags().StringP("path", "p", "", "the base Go path for generated sources")
	return cmd
}
