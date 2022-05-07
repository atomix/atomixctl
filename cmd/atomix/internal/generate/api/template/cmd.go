// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"fmt"
	"github.com/atomix/cli/cmd/atomix/internal/env"
	"github.com/spf13/cobra"
)

const exampleConfig = `
generator: go-client

files:
  - **/*.proto

templates:
  - name: atom.go
	path: templates/atom.go.tpl
	filter:
	  atoms:
		- Counter
		- *Election
	  components:
		- atom
	output: "{{ dir .File.Path }}/atom.go"
  - name: manager.go
    path: templates/manager.go.tpl
	filter:
	  atoms:
		- Counter
		- *Election
	  components:
		- manager
	output: "{{ dir .File.Path }}/manager.go"
`

func GetCommand(environment env.Environment) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "template",
		Short:   "Generates files from Protobuf sources using templates",
		Long:    fmt.Sprintf("Generates files from Protobuf sources using templates\n\nExample configuration:\n%s", exampleConfig),
		Example: "atomix generate api template --config go.yaml",
		Aliases: []string{"md"},
		Args:    cobra.NoArgs,
	}

	switch environment {
	case env.Native:
		cmd.RunE = runInDocker
	case env.Docker:
		cmd.RunE = runNative
	}

	cmd.Flags().StringP("config", "c", "", "the path to the generator configuration")
	cmd.Flags().StringToStringP("arg", "a", map[string]string{}, "additional template arguments")
	_ = cmd.MarkFlagFilename("config")
	return cmd
}
