// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package generate

import (
	"github.com/atomix/cli/internal/atomix/generate/api"
	"github.com/atomix/cli/internal/atomix/generate/client"
	"github.com/atomix/cli/internal/atomix/generate/driver"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen"},
	}
	cmd.AddCommand(api.GetCommand())
	cmd.AddCommand(client.GetCommand())
	cmd.AddCommand(driver.GetCommand())
	return cmd
}
