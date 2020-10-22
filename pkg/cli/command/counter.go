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
	"github.com/atomix/go-client/pkg/client/counter"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

func newCounterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "counter <name> [...]",
		Short:              "Manage the state of a distributed counter",
		Args:               cobra.MinimumNArgs(1),
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If only the name was specified, open an interactive shell
			name := args[0]
			if name == "-h" || name == "--help" {
				return cmd.Help()
			}
			if len(args) == 1 {
				return runShell(fmt.Sprintf("counter:%s", args[0]), os.Stdin, os.Stdout, os.Stderr, append(os.Args[1:], "counter", name))
			}

			// Get the command for the specified operation
			var subCmd *cobra.Command
			op := args[1]
			switch op {
			case "get":
				subCmd = newCounterGetCommand(name)
			case "set":
				subCmd = newCounterSetCommand(name)
			case "increment":
				subCmd = newCounterIncrementCommand(name)
			case "decrement":
				subCmd = newCounterDecrementCommand(name)
			case "-h", "--help":
				return cmd.Help()
			default:
				return fmt.Errorf("unknown command %s", op)
			}
			addClientFlags(subCmd)

			// Set the arguments after the name and execute the command
			subCmd.SetArgs(args[2:])
			return subCmd.Execute()
		},
	}
	return cmd
}

func getCounter(cmd *cobra.Command, name string) (counter.Counter, error) {
	database, err := getDatabase(cmd)
	if err != nil {
		return nil, err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	return database.GetCounter(ctx, name)
}

func newCounterGetCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			counter, err := getCounter(cmd, name)
			if err != nil {
				return err
			}
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			value, err := counter.Get(ctx)
			if err != nil {
				return err
			}
			cmd.Println(value)
			return nil
		},
	}
	return cmd
}

func newCounterSetCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "set <value>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			counter, err := getCounter(cmd, name)
			if err != nil {
				return err
			}
			value, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			err = counter.Set(ctx, int64(value))
			if err != nil {
				return err
			}
			cmd.Println(value)
			return nil
		},
	}
	return cmd
}

func newCounterIncrementCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "increment [delta]",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			counter, err := getCounter(cmd, name)
			if err != nil {
				return err
			}
			var delta int64
			if len(args) > 0 {
				value, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}
				delta = int64(value)
			}
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			value, err := counter.Increment(ctx, delta)
			if err != nil {
				return err
			}
			cmd.Println(value)
			return nil
		},
	}
	return cmd
}

func newCounterDecrementCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "decrement [delta]",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			counter, err := getCounter(cmd, name)
			if err != nil {
				return err
			}
			var delta int64
			if len(args) > 0 {
				value, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}
				delta = int64(value)
			}
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			value, err := counter.Decrement(ctx, delta)
			if err != nil {
				return err
			}
			cmd.Println(value)
			return nil
		},
	}
	return cmd
}
