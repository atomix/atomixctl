// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package deps

import (
	"github.com/atomix/cli/cmd/atomix/internal/env"
	"github.com/spf13/cobra"
)

func GetCommand(environment env.Environment) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deps",
		Args: cobra.NoArgs,
	}

	switch environment {
	case env.Native:
		cmd.RunE = runInDocker
	case env.Docker:
		cmd.RunE = runNative
	}

	cmd.Flags().BoolP("check", "c", false, "check module compatibility only")
	cmd.Flags().StringP("version", "v", "", "the target runtime API version")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}
