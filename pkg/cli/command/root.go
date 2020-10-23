// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"github.com/spf13/cobra"
	"os"
)

func GetRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                    "atomix",
		Short:                  "Atomix command line client",
		BashCompletionFunction: bashCompletion,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := getContext()
			if ctx == nil {
				ctx = newContext("atomix")
				setContext(ctx)
			}
			return ctx.run(os.Args[1:]...)
		},
	}

	addClientFlags(cmd)

	cmd.AddCommand(&cobra.Command{
		Use: "help",
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	})
	cmd.AddCommand(newCompletionCommand())
	cmd.AddCommand(newConfigCommand())
	cmd.AddCommand(newGetCommand())
	cmd.AddCommand(newCreateCommand())
	cmd.AddCommand(newDeleteCommand())
	cmd.AddCommand(newCounterCommand())
	cmd.AddCommand(newElectionCommand())
	cmd.AddCommand(newListCommand())
	cmd.AddCommand(newLockCommand())
	cmd.AddCommand(newLogCommand())
	cmd.AddCommand(newMapCommand())
	cmd.AddCommand(newSetCommand())
	cmd.AddCommand(newValueCommand())
	return cmd
}
