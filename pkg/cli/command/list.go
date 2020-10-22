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
	"context"
	"fmt"
	"github.com/atomix/go-client/pkg/client/list"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "list <name> [...]",
		Short:              "Manage the state of a distributed list",
		Args:               cobra.MinimumNArgs(1),
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If only the name was specified, open an interactive shell
			name := args[0]
			if len(args) == 1 {
				return runShell(fmt.Sprintf("list:%s", args[0]), os.Stdin, os.Stdout, os.Stderr, append(os.Args[1:], "list", name))
			}

			// Get the command for the specified operation
			var subCmd *cobra.Command
			op := args[1]
			switch op {
			case "get":
				subCmd = newListGetCommand(name)
			case "append":
				subCmd = newListAppendCommand(name)
			case "insert":
				subCmd = newListInsertCommand(name)
			case "remove":
				subCmd = newListRemoveCommand(name)
			case "items":
				subCmd = newListItemsCommand(name)
			case "size":
				subCmd = newListSizeCommand(name)
			case "clear":
				subCmd = newListClearCommand(name)
			case "watch":
				subCmd = newListWatchCommand(name)
			}

			// Set the arguments after the name and execute the command
			subCmd.SetArgs(args[2:])
			return subCmd.Execute()
		},
	}
	addClientFlags(cmd)
	return cmd
}

func getList(cmd *cobra.Command, name string) (list.List, error) {
	database, err := getDatabase(cmd)
	if err != nil {
		return nil, err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	return database.GetList(ctx, name)
}

func newListGetCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get <index>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			list, err := getList(cmd, name)
			if err != nil {
				return err
			}
			indexStr := args[0]
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				return err
			}
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			value, err := list.Get(ctx, index)
			if value != nil {
				cmd.Println(string(value))
			}
			return err
		},
	}
	addClientFlags(cmd)
	return cmd
}

func newListAppendCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "append <value>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := getList(cmd, name)
			if err != nil {
				return err
			}
			value := args[0]
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			return l.Append(ctx, []byte(value))
		},
	}
	addClientFlags(cmd)
	return cmd
}

func newListInsertCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "insert <index> <value>",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := getList(cmd, name)
			if err != nil {
				return err
			}
			indexStr := args[0]
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				return err
			}
			value := args[1]
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			return l.Insert(ctx, int(index), []byte(value))
		},
	}
	addClientFlags(cmd)
	return cmd
}

func newListRemoveCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "remove <index>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := getList(cmd, name)
			if err != nil {
				return err
			}
			indexStr := args[0]
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				return err
			}
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			value, err := m.Remove(ctx, int(index))
			if value != nil {
				cmd.Println(string(value))
			}
			return err
		},
	}
	cmd.Flags().Int64P("version", "v", 0, "the entry version")
	return cmd
}

func newListItemsCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "items",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := getList(cmd, name)
			if err != nil {
				return err
			}
			ch := make(chan []byte)
			err = m.Items(context.TODO(), ch)
			if err != nil {
				return err
			}
			for value := range ch {
				cmd.Println(string(value))
			}
			return nil
		},
	}
	addClientFlags(cmd)
	return cmd
}

func newListSizeCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "size",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			list, err := getList(cmd, name)
			if err != nil {
				return err
			}
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			size, err := list.Len(ctx)
			if err != nil {
				return err
			}
			cmd.Println(size)
			return nil
		},
	}
	addClientFlags(cmd)
	return cmd
}

func newListClearCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "clear",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			list, err := getList(cmd, name)
			if err != nil {
				return err
			}
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			return list.Clear(ctx)
		},
	}
	addClientFlags(cmd)
	return cmd
}

func newListWatchCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "watch",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := getList(cmd, name)
			if err != nil {
				return err
			}
			watchCh := make(chan *list.Event)
			opts := []list.WatchOption{}
			replay, _ := cmd.Flags().GetBool("replay")
			if replay {
				opts = append(opts, list.WithReplay())
			}
			if err := l.Watch(context.Background(), watchCh, opts...); err != nil {
				return err
			}

			signalCh := make(chan os.Signal, 2)
			signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
			for {
				select {
				case event := <-watchCh:
					bytes, err := yaml.Marshal(event)
					if err != nil {
						cmd.Println(err)
					} else {
						cmd.Println(string(bytes))
					}
				case <-signalCh:
					return nil
				}
			}
		},
	}
	cmd.Flags().BoolP("replay", "r", false, "replay current list values at start")
	addClientFlags(cmd)
	return cmd
}
