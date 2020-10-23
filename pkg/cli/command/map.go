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
	"bytes"
	"fmt"
	"github.com/atomix/go-client/pkg/client/map"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"os"
)

func newMapCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "map <name> [...]",
		Short:              "Manage a distributed map",
		Args:               cobra.MinimumNArgs(1),
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If only the name was specified, open an interactive shell
			name := args[0]
			if name == "-h" || name == "--help" {
				return cmd.Help()
			}
			if len(args) == 1 {
				return runShell(cmd, fmt.Sprintf("map:%s", args[0]), os.Stdin, os.Stdout, os.Stderr, []string{"map", name})
			}

			// Get the command for the specified operation
			var subCmd *cobra.Command
			op := args[1]
			switch op {
			case "get":
				subCmd = newMapGetCommand(name)
			case "put":
				subCmd = newMapPutCommand(name)
			case "remove":
				subCmd = newMapRemoveCommand(name)
			case "keys":
				subCmd = newMapKeysCommand(name)
			case "entries":
				subCmd = newMapEntriesCommand(name)
			case "size":
				subCmd = newMapSizeCommand(name)
			case "clear":
				subCmd = newMapClearCommand(name)
			case "watch":
				subCmd = newMapWatchCommand(name)
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

func getMap(cmd *cobra.Command, name string) (_map.Map, error) {
	database, err := getDatabase(cmd)
	if err != nil {
		return nil, err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	return database.GetMap(ctx, name)
}

func newMapGetCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get <key>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_map, err := getMap(cmd, name)
			if err != nil {
				return err
			}
			key := args[0]
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			value, err := _map.Get(ctx, key)
			if err != nil {
				return err
			} else if value != nil {
				cmd.Println(value.String())
			}
			return nil
		},
	}
	return cmd
}

func newMapPutCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "put <key> <value>",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := getMap(cmd, name)
			if err != nil {
				return err
			}
			key := args[0]
			value := args[1]
			version, _ := cmd.Flags().GetInt64("version")
			opts := []_map.PutOption{}
			if version > 0 {
				opts = append(opts, _map.IfVersion(_map.Version(version)))
			}

			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			kv, err := m.Put(ctx, key, []byte(value), opts...)
			if err != nil {
				return err
			} else if kv != nil {
				bytes, err := yaml.Marshal(kv)
				if err != nil {
					return err
				}
				cmd.Println(string(bytes))
			}
			return nil
		},
	}
	cmd.Flags().Int64("version", 0, "the entry version")
	return cmd
}

func newMapRemoveCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "remove <key>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := getMap(cmd, name)
			if err != nil {
				return err
			}
			key := args[0]
			version, _ := cmd.Flags().GetInt64("version")
			opts := []_map.RemoveOption{}
			if version > 0 {
				opts = append(opts, _map.IfVersion(_map.Version(version)))
			}

			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			value, err := m.Remove(ctx, key, opts...)
			if err != nil {
				return err
			} else if value != nil {
				bytes, err := yaml.Marshal(value)
				if err != nil {
					return err
				}
				cmd.Println(string(bytes))
			}
			return nil
		},
	}
	cmd.Flags().Int64("version", 0, "the entry version")
	return cmd
}

func newMapKeysCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "keys",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := getMap(cmd, name)
			if err != nil {
				return err
			}

			ch := make(chan *_map.Entry)
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			err = m.Entries(ctx, ch)
			if err != nil {
				return err
			}

			for kv := range ch {
				bytes, err := yaml.Marshal(kv)
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

func newMapSizeCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "size",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_map, err := getMap(cmd, name)
			if err != nil {
				return err
			}

			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			size, err := _map.Len(ctx)
			if err != nil {
				return err
			}
			cmd.Println(size)
			return nil
		},
	}
	return cmd
}

func newMapClearCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "clear",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_map, err := getMap(cmd, name)
			if err != nil {
				return err
			}
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			return _map.Clear(ctx)
		},
	}
	return cmd
}

func newMapEntriesCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "entries",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := getMap(cmd, name)
			if err != nil {
				return err
			}

			ch := make(chan *_map.Entry)
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			err = m.Entries(ctx, ch)
			if err != nil {
				return err
			}

			context := getContext()
			if context.isShell {
				var buf bytes.Buffer
				for entry := range ch {
					bytes, err := yaml.Marshal(entry)
					if err != nil {
						cmd.Println(err)
					} else {
						buf.Write(bytes)
						buf.WriteByte('\n')
					}
				}
				return context.shellCtx.ShowPaged(buf.String())
			}

			for entry := range ch {
				bytes, err := yaml.Marshal(entry)
				if err != nil {
					cmd.Println(err)
				} else {
					cmd.Println(string(bytes))
				}
			}
			return nil
		},
	}
	return cmd
}

func newMapWatchCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "watch",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := getMap(cmd, name)
			if err != nil {
				return err
			}

			watchCh := make(chan *_map.Event)
			opts := []_map.WatchOption{}
			replay, _ := cmd.Flags().GetBool("replay")
			if replay {
				opts = append(opts, _map.WithReplay())
			}

			ctx, cancel := getCancelContext(cmd)
			defer cancel()
			if err := m.Watch(ctx, watchCh, opts...); err != nil {
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
	cmd.Flags().BoolP("replay", "r", false, "replay current map entries at start")
	return cmd
}
