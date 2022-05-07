// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package generate

import (
	"github.com/atomix/cli/cmd/atomix/internal/env"
	"github.com/atomix/cli/cmd/atomix/internal/generate/api"
	"github.com/atomix/cli/cmd/atomix/internal/generate/deps"
	"github.com/atomix/cli/cmd/atomix/internal/generate/docs"
	"github.com/atomix/cli/cmd/atomix/internal/generate/driver"
	"github.com/spf13/cobra"
)

func GetCommand(env env.Environment) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen"},
	}
	cmd.AddCommand(api.GetCommand(env))
	cmd.AddCommand(deps.GetCommand(env))
	cmd.AddCommand(docs.GetCommand(env))
	cmd.AddCommand(driver.GetCommand(env))
	return cmd
}
