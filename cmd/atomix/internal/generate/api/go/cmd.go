// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package _go

import (
	"github.com/atomix/cli/cmd/atomix/internal/env"
	"github.com/spf13/cobra"
)

func GetCommand(environment env.Environment) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "go",
		Short:   "Generates Go sources from Protobuf sources",
		Example: "atomix generate api go --input ../api --pattern '**/*.proto' --output ./pkg --package github.com/my-user/my-repo",
		Aliases: []string{"golang"},
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
	cmd.Flags().StringP("go-path", "d", ".", "the relative path to the documentation root")
	cmd.Flags().StringP("import-path", "i", "", "the base Go path for generated sources")
	_ = cmd.MarkFlagFilename("config")
	return cmd
}
