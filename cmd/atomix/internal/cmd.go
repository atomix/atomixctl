// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"github.com/atomix/cli/cmd/atomix/internal/config"
	"github.com/atomix/cli/cmd/atomix/internal/docs"
	"github.com/atomix/cli/cmd/atomix/internal/generate"
	"github.com/atomix/cli/cmd/atomix/internal/version"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "atomix",
	}
	cmd.AddCommand(
		config.GetCommand(),
		docs.GetCommand(),
		generate.GetCommand(),
		version.GetCommand())
	return cmd
}
