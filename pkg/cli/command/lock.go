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
		Use:   "lock {lock,get}",
		Short: "Manage the state of a distributed lock",
	}
	addClientFlags(cmd)
	cmd.PersistentFlags().StringP("name", "n", "", "the lock name")
	cmd.PersistentFlags().Lookup("name").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_locks"},
	}
	cmd.MarkPersistentFlagRequired("name")
	cmd.AddCommand(newLockLockCommand())
	cmd.AddCommand(newLockGetCommand())
	return cmd
}

func getLockName(cmd *cobra.Command) string {
	name, _ := cmd.Flags().GetString("name")
	return name
}

func getLock(cmd *cobra.Command, name string) lock.Lock {
	database := getDatabase(cmd)
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	m, err := database.GetLock(ctx, name)
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return m
}

func newLockLockCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "lock",
		Args: cobra.NoArgs,
		Run:  runLockLockCommand,
	}
}

func runLockLockCommand(cmd *cobra.Command, _ []string) {
	lock := getLock(cmd, getLockName(cmd))
	ctx, cancel := getTimeoutContext(cmd)
	version, err := lock.Lock(ctx)
	cancel()
	if err != nil {
		ExitWithError(ExitError, err)
	}
	fmt.Println(fmt.Sprintf("Acquired lock %d", version))
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	ctx, cancel = getTimeoutContext(cmd)
	_, err = lock.Unlock(ctx)
	cancel()
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(fmt.Sprintf("Released lock %d", version))
	}
}

func newLockGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Args: cobra.NoArgs,
		Run:  runLockGetCommand,
	}
	cmd.Flags().Uint64P("version", "v", 0, "the lock version")
	return cmd
}

func runLockGetCommand(cmd *cobra.Command, _ []string) {
	l := getLock(cmd, getLockName(cmd))
	version, _ := cmd.Flags().GetUint64("version")
	if version == 0 {
		ctx, cancel := getTimeoutContext(cmd)
		defer cancel()
		locked, err := l.IsLocked(ctx)
		if err != nil {
			ExitWithError(ExitError, err)
		} else {
			if locked {
				ExitWithOutput("<locked>")
			} else {
				ExitWithOutput("<unlocked>")
			}
		}
	} else {
		ctx, cancel := getTimeoutContext(cmd)
		defer cancel()
		locked, err := l.IsLocked(ctx, lock.IfVersion(version))
		if err != nil {
			ExitWithError(ExitError, err)
		} else {
			if locked {
				ExitWithOutput("<locked>")
			} else {
				ExitWithOutput("<unlocked>")
			}
		}
	}
}
