// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"github.com/atomix/cli/cmd/atomix/internal/env"
	"github.com/atomix/cli/cmd/atomix/internal/generate/api/docs"
	_go "github.com/atomix/cli/cmd/atomix/internal/generate/api/go"
	"github.com/atomix/cli/cmd/atomix/internal/generate/api/template"
	"github.com/spf13/cobra"
)

func GetCommand(env env.Environment) *cobra.Command {
	cmd := &cobra.Command{
		Use: "api",
	}
	cmd.AddCommand(_go.GetCommand(env))
	cmd.AddCommand(docs.GetCommand(env))
	cmd.AddCommand(template.GetCommand(env))
	return cmd
}
