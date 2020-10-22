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
	"fmt"
	"github.com/abiosoft/ishell"
	"github.com/abiosoft/readline"
	"github.com/spf13/cobra"
	"io"
	"os"
)

func GetRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                    "atomix",
		Short:                  "Atomix command line client",
		BashCompletionFunction: bashCompletion,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runShell("atomix", os.Stdin, os.Stdout, os.Stderr, os.Args)
		},
	}

	addClientFlags(cmd)

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

func runShell(name string, stdin io.ReadCloser, stdout io.Writer, stderr io.Writer, args []string) error {
	ctx := getContext()
	historyFile, err := getConfigFile("history")
	if err != nil {
		panic(err)
	}
	shell := ishell.NewWithConfig(&readline.Config{
		Prompt:      fmt.Sprintf("%s>", name),
		HistoryFile: historyFile,
		Stdin:       stdin,
		Stdout:      stdout,
		Stderr:      stderr,
	})
	shell.NotFound(func(context *ishell.Context) {
		setContextFunc(func(ctx *commandContext) {
			ctx.shellCtx = context
		})
		cmd := GetRootCommand()
		cmd.SetArgs(append(args, context.RawArgs...))
		err := cmd.Execute()
		if err != nil {
			context.Println(err)
		}
	})
	setContextFunc(func(context *commandContext) {
		context.isShell = true
		context.shell = shell
	})
	shell.Run()
	setContext(ctx)
	return nil
}
