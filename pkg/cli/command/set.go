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
	"github.com/atomix/go-client/pkg/client/set"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"os"
)

func newSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "set <name> [...]",
		Short:              "Manage the state of a distributed set",
		Args:               cobra.MinimumNArgs(1),
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If only the name was specified, open an interactive shell
			name := args[0]
			if name == "-h" || name == "--help" {
				return cmd.Help()
			}
			if len(args) == 1 {
				return runShell(cmd, fmt.Sprintf("set:%s", args[0]), os.Stdin, os.Stdout, os.Stderr, append(os.Args[1:], "set", name))
			}

			// Get the command for the specified operation
			var subCmd *cobra.Command
			op := args[1]
			switch op {
			case "add":
				subCmd = newSetAddCommand(name)
			case "contains":
				subCmd = newSetContainsCommand(name)
			case "remove":
				subCmd = newSetRemoveCommand(name)
			case "size":
				subCmd = newSetSizeCommand(name)
			case "elements":
				subCmd = newSetElementsCommand(name)
			case "clear":
				subCmd = newSetClearCommand(name)
			case "watch":
				subCmd = newSetWatchCommand(name)
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

func getSet(cmd *cobra.Command, name string) (set.Set, error) {
	database, err := getDatabase(cmd)
	if err != nil {
		return nil, err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	return database.GetSet(ctx, name)
}

func newSetAddCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "add <value>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			set, err := getSet(cmd, name)
			if err != nil {
				return err
			}
			value := args[0]
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			added, err := set.Add(ctx, value)
			if err != nil {
				return err
			}
			cmd.Println(added)
			return nil
		},
	}
	return cmd
}

func newSetContainsCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "contains <value>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			set, err := getSet(cmd, name)
			if err != nil {
				return err
			}
			value := args[0]
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			contains, err := set.Contains(ctx, value)
			if err != nil {
				return err
			}
			cmd.Println(contains)
			return nil
		},
	}
	return cmd
}

func newSetRemoveCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "remove <value>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			set, err := getSet(cmd, name)
			if err != nil {
				return err
			}
			value := args[0]
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			removed, err := set.Remove(ctx, value)
			if err != nil {
				return err
			}
			cmd.Println(removed)
			return nil
		},
	}
	return cmd
}

func newSetSizeCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "size",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			set, err := getSet(cmd, name)
			if err != nil {
				return err
			}
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			size, err := set.Len(ctx)
			if err != nil {
				return err
			}
			cmd.Println(size)
			return nil
		},
	}
	return cmd
}

func newSetElementsCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "elements",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := getSet(cmd, name)
			if err != nil {
				return err
			}

			ch := make(chan string)
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			err = m.Elements(ctx, ch)
			if err != nil {
				return err
			}

			context := getContext()
			if context.isShell {
				var buf bytes.Buffer
				for value := range ch {
					buf.WriteString(value)
					buf.WriteByte('\n')
				}
				return context.shellCtx.ShowPaged(buf.String())
			}

			for value := range ch {
				cmd.Println(value)
			}
			return nil
		},
	}
	return cmd
}

func newSetClearCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "clear",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			set, err := getSet(cmd, name)
			if err != nil {
				return err
			}
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			return set.Clear(ctx)
		},
	}
	return cmd
}

func newSetWatchCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "watch",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := getSet(cmd, name)
			if err != nil {
				return err
			}

			watchCh := make(chan *set.Event)
			opts := []set.WatchOption{}
			replay, _ := cmd.Flags().GetBool("replay")
			if replay {
				opts = append(opts, set.WithReplay())
			}

			ctx, cancel := getCancelContext(cmd)
			defer cancel()
			if err := s.Watch(ctx, watchCh, opts...); err != nil {
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
	cmd.Flags().BoolP("replay", "r", false, "replay current set values at start")
	return cmd
}
