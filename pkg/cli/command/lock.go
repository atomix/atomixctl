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
	"github.com/atomix/go-client/pkg/client/lock"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

func newLockCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "lock <name> [...]",
		Short:              "Manage the state of a distributed lock",
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
					ctx = newContext("atomix", "lock", name)
					setContext(ctx)
				} else {
					ctx = ctx.withCommand("lock", name)
				}
				return ctx.run()
			}

			// Get the command for the specified operation
			var subCmd *cobra.Command
			op := args[1]
			switch op {
			case "get":
				subCmd = newLockGetCommand(name)
			case "lock":
				subCmd = newLockLockCommand(name)
			case "help", "-h", "--help":
				if len(args) == 2 {
					helpCmd := &cobra.Command{
						Use:   fmt.Sprintf("lock %s [...]", name),
						Short: "Manage the state of a distributed lock",
					}
					helpCmd.AddCommand(newLockGetCommand(name))
					helpCmd.AddCommand(newLockLockCommand(name))
					return helpCmd.Help()
				} else {
					var helpCmd *cobra.Command
					switch args[2] {
					case "get":
						helpCmd = newLockGetCommand(name)
					case "lock":
						helpCmd = newLockLockCommand(name)
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

func getLock(cmd *cobra.Command, name string) (lock.Lock, error) {
	database, err := getDatabase(cmd)
	if err != nil {
		return nil, err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	return database.GetLock(ctx, name)
}

func newLockLockCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lock",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			lock, err := getLock(cmd, name)
			if err != nil {
				return err
			}

			ctx, cancel := getTimeoutContext(cmd)
			version, err := lock.Lock(ctx)
			cancel()
			if err != nil {
				return err
			}
			cmd.Println(version)

			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
			<-ch

			ctx, cancel = getTimeoutContext(cmd)
			_, err = lock.Unlock(ctx)
			cancel()
			return err
		},
	}
	return cmd
}

func newLockGetCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := getLock(cmd, name)
			if err != nil {
				return err
			}

			version, _ := cmd.Flags().GetUint64("version")
			var opts []lock.IsLockedOption
			if version != 0 {
				opts = append(opts, lock.IfVersion(version))
			}

			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			locked, err := l.IsLocked(ctx, opts...)
			if err != nil {
				return err
			}

			if locked {
				cmd.Println("locked")
			} else {
				cmd.Println("unlocked")
			}
			return nil
		},
	}
	cmd.Flags().Uint64P("version", "v", 0, "the lock version")
	return cmd
}
