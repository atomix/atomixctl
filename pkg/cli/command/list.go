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
	"github.com/atomix/go-client/pkg/client/list"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"strconv"
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
			if name == "-h" || name == "--help" {
				return cmd.Help()
			}
			if len(args) == 1 {
				ctx := getContext()
				if ctx == nil {
					ctx = newContext("atomix", "list", name)
					setContext(ctx)
				} else {
					ctx = ctx.withCommand("list", name)
				}
				return ctx.run()
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
			case "help", "-h", "--help":
				if len(args) == 2 {
					helpCmd := &cobra.Command{
						Use:   fmt.Sprintf("list %s [...]", name),
						Short: "Manage the state of a distributed list",
					}
					helpCmd.AddCommand(newListGetCommand(name))
					helpCmd.AddCommand(newListAppendCommand(name))
					helpCmd.AddCommand(newListInsertCommand(name))
					helpCmd.AddCommand(newListRemoveCommand(name))
					helpCmd.AddCommand(newListItemsCommand(name))
					helpCmd.AddCommand(newListSizeCommand(name))
					helpCmd.AddCommand(newListClearCommand(name))
					helpCmd.AddCommand(newListWatchCommand(name))
					return helpCmd.Help()
				} else {
					var helpCmd *cobra.Command
					switch args[2] {
					case "get":
						helpCmd = newListGetCommand(name)
					case "append":
						helpCmd = newListAppendCommand(name)
					case "insert":
						helpCmd = newListInsertCommand(name)
					case "remove":
						helpCmd = newListRemoveCommand(name)
					case "items":
						helpCmd = newListItemsCommand(name)
					case "size":
						helpCmd = newListSizeCommand(name)
					case "clear":
						helpCmd = newListClearCommand(name)
					case "watch":
						helpCmd = newListWatchCommand(name)
					default:
						return fmt.Errorf("unknown command %s", args[2])
					}
					return helpCmd.Help()
				}
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
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			err = m.Items(ctx, ch)
			if err != nil {
				return err
			}

			lines := make([]interface{}, 0)
			for value := range ch {
				lines = append(lines, string(value))
			}

			context := getContext()
			context.Printlns(lines...)
			return nil
		},
	}
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

			ctx, cancel := getCancelContext(cmd)
			defer cancel()
			if err := l.Watch(ctx, watchCh, opts...); err != nil {
				return err
			}

			for event := range watchCh {
				bytes, err := yaml.Marshal(event)
				if err != nil {
					cmd.Println(err)
				} else {
					cmd.Println(string(bytes))
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolP("replay", "r", false, "replay current list values at start")
	return cmd
}
