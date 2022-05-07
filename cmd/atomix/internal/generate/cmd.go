// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package generate

import (
	"github.com/atomix/cli/cmd/atomix/internal/generate/deps"
	"github.com/atomix/cli/cmd/atomix/internal/generate/docs"
	"github.com/atomix/cli/cmd/atomix/internal/generate/driver"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen"},
	}
	cmd.AddCommand(deps.GetCommand())
	cmd.AddCommand(docs.GetCommand())
	cmd.AddCommand(driver.GetCommand())
	return cmd
}
