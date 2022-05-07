// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"github.com/atomix/cli/cmd/atomix/internal/env"
	"github.com/spf13/cobra"
)

func GetCommand(environment env.Environment) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "driver",
		Args: cobra.NoArgs,
	}

	switch environment {
	case env.Native:
		cmd.RunE = runInDocker
	case env.Docker:
		cmd.RunE = runNative
	}

	cmd.Flags().StringP("name", "n", "", "the driver name")
	cmd.Flags().StringP("module-path", "p", "", "the driver module path")
	cmd.Flags().String("github-owner", "", "the GitHub user to which to publish release artifacts")
	cmd.Flags().String("github-repo", "", "the GitHub repo to which to publish release artifacts")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("module")
	return cmd
}
