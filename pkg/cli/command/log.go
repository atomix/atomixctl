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
	"github.com/atomix/go-client/pkg/client/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"strconv"
)

func newLogCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "log <name> [...]",
		Short:              "Manage the state of a distributed log",
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
					ctx = newContext("atomix", "log", name)
					setContext(ctx)
				} else {
					ctx = ctx.withCommand("log", name)
				}
				return ctx.run()
			}

			// Get the command for the specified operation
			var subCmd *cobra.Command
			op := args[1]
			switch op {
			case "get":
				subCmd = newLogGetCommand(name)
			case "append":
				subCmd = newLogAppendCommand(name)
			case "remove":
				subCmd = newLogRemoveCommand(name)
			case "entries":
				subCmd = newLogEntriesCommand(name)
			case "size":
				subCmd = newLogSizeCommand(name)
			case "clear":
				subCmd = newLogClearCommand(name)
			case "watch":
				subCmd = newLogWatchCommand(name)
			case "help", "-h", "--help":
				if len(args) == 2 {
					helpCmd := &cobra.Command{
						Use:   fmt.Sprintf("log %s [...]", name),
						Short: "Manage the state of a distributed log",
					}
					helpCmd.AddCommand(newLogGetCommand(name))
					helpCmd.AddCommand(newLogAppendCommand(name))
					helpCmd.AddCommand(newLogRemoveCommand(name))
					helpCmd.AddCommand(newLogEntriesCommand(name))
					helpCmd.AddCommand(newLogSizeCommand(name))
					helpCmd.AddCommand(newLogClearCommand(name))
					helpCmd.AddCommand(newLogWatchCommand(name))
					return helpCmd.Help()
				} else {
					var helpCmd *cobra.Command
					switch args[2] {
					case "get":
						helpCmd = newLogGetCommand(name)
					case "append":
						helpCmd = newLogAppendCommand(name)
					case "remove":
						helpCmd = newLogRemoveCommand(name)
					case "entries":
						helpCmd = newLogEntriesCommand(name)
					case "size":
						helpCmd = newLogSizeCommand(name)
					case "clear":
						helpCmd = newLogClearCommand(name)
					case "watch":
						helpCmd = newLogWatchCommand(name)
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

func getLog(cmd *cobra.Command, name string) (log.Log, error) {
	database, err := getDatabase(cmd)
	if err != nil {
		return nil, err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	return database.GetLog(ctx, name)
}

func newLogGetCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get <index>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := getLog(cmd, name)
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
			entry, err := l.Get(ctx, log.Index(index))
			if err != nil {
				return err
			} else if entry != nil {
				bytes, err := yaml.Marshal(entry)
				if err != nil {
					return err
				}
				cmd.Println(string(bytes))
			}
			return nil
		},
	}
	return cmd
}

func newLogAppendCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "append <value>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := getLog(cmd, name)
			if err != nil {
				return err
			}

			value := args[0]
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			entry, err := l.Append(ctx, []byte(value))
			if err != nil {
				return err
			} else if entry != nil {
				bytes, err := yaml.Marshal(entry)
				if err != nil {
					return err
				}
				cmd.Println(string(bytes))
			}
			return nil
		},
	}
	return cmd
}

func newLogRemoveCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "remove <index>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := getLog(cmd, name)
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
			entry, err := l.Remove(ctx, log.Index(index))
			if err != nil {
				return err
			} else if entry != nil {
				bytes, err := yaml.Marshal(entry)
				if err != nil {
					return err
				}
				cmd.Println(string(bytes))
			}
			return nil
		},
	}
	cmd.Flags().Int64P("version", "v", 0, "the entry version")
	return cmd
}

func newLogEntriesCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "entries",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := getLog(cmd, name)
			if err != nil {
				return err
			}

			ch := make(chan *log.Entry)
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			err = l.Entries(ctx, ch)
			if err != nil {
				return err
			}

			context := getContext()
			lines := make([]interface{}, 0)
			for entry := range ch {
				bytes, err := yaml.Marshal(entry)
				if err != nil {
					context.Println(err)
				} else {
					lines = append(lines, string(bytes))
				}
			}
			context.Printlns(lines...)
			return nil
		},
	}
	return cmd
}

func newLogSizeCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "size",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			log, err := getLog(cmd, name)
			if err != nil {
				return err
			}

			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			size, err := log.Size(ctx)
			if err != nil {
				return err
			}
			cmd.Println(size)
			return nil
		},
	}
	return cmd
}

func newLogClearCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "clear",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			log, err := getLog(cmd, name)
			if err != nil {
				return err
			}
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			return log.Clear(ctx)
		},
	}
	return cmd
}

func newLogWatchCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "watch",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := getLog(cmd, name)
			if err != nil {
				return err
			}

			watchCh := make(chan *log.Event)
			opts := []log.WatchOption{}
			replay, _ := cmd.Flags().GetBool("replay")
			if replay {
				opts = append(opts, log.WithReplay())
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
	cmd.Flags().BoolP("replay", "r", false, "replay current log entries at start")
	return cmd
}
