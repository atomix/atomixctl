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
	"errors"
	"fmt"
	"github.com/atomix/go-client/pkg/client/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func newLogCommand() *cobra.Command {
	return &cobra.Command{
		Use:                "log <name> [...]",
		Short:              "Manage the state of a distributed log",
		Args:               cobra.MinimumNArgs(1),
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If only the name was specified, open an interactive shell
			name := args[0]
			if len(args) == 1 {
				return runShell(fmt.Sprintf("log:%s", args[0]), os.Stdin, os.Stdout, os.Stderr, append(os.Args[1:], "log", name))
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
			}

			// Set the arguments after the name and execute the command
			subCmd.SetArgs(args[2:])
			return subCmd.Execute()
		},
	}
}

func getLog(cmd *cobra.Command, name string) (log.Log, error) {
	database := getDatabase(cmd)
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	return database.GetLog(ctx, name)
}

func newLogGetCommand(name string) *cobra.Command {
	return &cobra.Command{
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
}

func newLogAppendCommand(name string) *cobra.Command {
	return &cobra.Command{
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
	return &cobra.Command{
		Use:  "entries",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("not implemented")
		},
	}
}

func newLogSizeCommand(name string) *cobra.Command {
	return &cobra.Command{
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
}

func newLogClearCommand(name string) *cobra.Command {
	return &cobra.Command{
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
			if err := l.Watch(context.Background(), watchCh, opts...); err != nil {
				ExitWithError(ExitError, err)
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
	cmd.Flags().BoolP("replay", "r", false, "replay current log entries at start")
	return cmd
}
