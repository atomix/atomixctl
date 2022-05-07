// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package docs

import (
	"github.com/atomix/cli/cmd/atomix/internal/env"
	"github.com/spf13/cobra"
)

func GetCommand(environment env.Environment) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "docs",
		Short:   "Generates documentation from Protobuf sources",
		Example: "atomix generate api docs --format markdown --input ../api --pattern '**/*.proto' --output ./docs",
		Aliases: []string{"doc"},
		Args:    cobra.NoArgs,
	}

	switch environment {
	case env.Native:
		cmd.RunE = runInDocker
	case env.Docker:
		cmd.RunE = runNative
	}

	cmd.Flags().StringP("config", "c", "", "the path to the generator configuration")
	cmd.Flags().StringP("proto-path", "p", ".", "the relative path to the Protobuf API root")
	cmd.Flags().StringSliceP("proto-pattern", "f", []string{"**/*.proto"}, "a pattern by which to filter Protobuf sources")
	cmd.Flags().StringP("docs-path", "d", ".", "the relative path to the documentation root")
	cmd.Flags().String("docs-format", "markdown", "the documentation format")
	_ = cmd.MarkFlagFilename("config")
	return cmd
}
