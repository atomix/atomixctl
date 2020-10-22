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

func runShell(cmd *cobra.Command, name string, stdin io.ReadCloser, stdout io.Writer, stderr io.Writer, args []string) error {
	parentCtx := getContext()
	if parentCtx.isShell {
		name = fmt.Sprintf("%s:%s", parentCtx.shellName, name)
	}

	historyFile, err := getConfigFile("history.log")
	if err != nil {
		return err
	}

	shell := ishell.NewWithConfig(&readline.Config{
		Prompt:      fmt.Sprintf("%s> ", name),
		HistoryFile: historyFile,
		Stdin:       stdin,
		Stdout:      stdout,
		Stderr:      stderr,
	})

	shell.NotFound(func(context *ishell.Context) {
		setContextFunc(func(ctx *commandContext) {
			ctx.isRoot = false
			ctx.shellCtx = context
		})
		cmd := GetRootCommand()
		cmd.SetArgs(append(args, context.RawArgs...))
		err := cmd.Execute()
		if err != nil {
			context.Println(err)
		}
	})
	setContextFunc(func(ctx *commandContext) {
		ctx.isShell = true
		ctx.shellName = name
		ctx.shell = shell
		ctx.shellCmd = cmd
	})
	shell.Interrupt(func(context *ishell.Context, count int, input string) {
		if parentCtx.isShell {
			context.Stop()
		} else if count >= 2 {
			context.Println("Interrupted")
			os.Exit(1)
		} else {
			context.Println("Input Ctrl-c once more to exit")
		}
	})
	shell.EOF(func(context *ishell.Context) {
		context.Stop()
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "help",
		Func: func(context *ishell.Context) {
			cmd.Help()
		},
	})
	shell.Run()
	setContext(parentCtx)
	return nil
}
