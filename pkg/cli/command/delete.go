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
	"github.com/spf13/cobra"
)

func newDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete {counter,election,list,lock,map,set,value}",
		Short: "Delete a distributed primitive",
	}
	cmd.AddCommand(newDeleteCounterCommand())
	cmd.AddCommand(newDeleteElectionCommand())
	cmd.AddCommand(newDeleteListCommand())
	cmd.AddCommand(newDeleteLockCommand())
	cmd.AddCommand(newDeleteLogCommand())
	cmd.AddCommand(newDeleteMapCommand())
	cmd.AddCommand(newDeleteSetCommand())
	cmd.AddCommand(newDeleteValueCommand())
	return cmd
}

func newDeleteCounterCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "counter <name>",
		Args: cobra.ExactArgs(1),
		RunE: runDeleteCounterCommand,
	}
}

func runDeleteCounterCommand(cmd *cobra.Command, args []string) error {
	counter, err := getCounter(cmd, args[0])
	if err != nil {
		return err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	return counter.Delete(ctx)
}

func newDeleteElectionCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "election <name>",
		Args: cobra.ExactArgs(1),
		RunE: runDeleteElectionCommand,
	}
}

func runDeleteElectionCommand(cmd *cobra.Command, args []string) error {
	election, err := getElection(cmd, args[0], "")
	if err != nil {
		return err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	return election.Delete(ctx)
}

func newDeleteListCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "list <name>",
		Args: cobra.ExactArgs(1),
		RunE: runDeleteListCommand,
	}
}

func runDeleteListCommand(cmd *cobra.Command, args []string) error {
	list, err := getList(cmd, args[0])
	if err != nil {
		return err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	return list.Delete(ctx)
}

func newDeleteLockCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "lock <name>",
		Args: cobra.ExactArgs(1),
		RunE: runDeleteLockCommand,
	}
}

func runDeleteLockCommand(cmd *cobra.Command, args []string) error {
	lock, err := getLock(cmd, args[0])
	if err != nil {
		return err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	return lock.Delete(ctx)
}

func newDeleteLogCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "log <name>",
		Args: cobra.ExactArgs(1),
		RunE: runDeleteLogCommand,
	}
}

func runDeleteLogCommand(cmd *cobra.Command, args []string) error {
	log, err := getLog(cmd, args[0])
	if err != nil {
		return err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	return log.Delete(ctx)
}

func newDeleteMapCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "map <name>",
		Args: cobra.ExactArgs(1),
		RunE: runDeleteMapCommand,
	}
}

func runDeleteMapCommand(cmd *cobra.Command, args []string) error {
	_map, err := getMap(cmd, args[0])
	if err != nil {
		return err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	return _map.Delete(ctx)
}

func newDeleteSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "set <name>",
		Args: cobra.ExactArgs(1),
		RunE: runDeleteSetCommand,
	}
}

func runDeleteSetCommand(cmd *cobra.Command, args []string) error {
	set, err := getSet(cmd, args[0])
	if err != nil {
		return err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	return set.Delete(ctx)
}

func newDeleteValueCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "value <name>",
		Args: cobra.ExactArgs(1),
		RunE: runDeleteValueCommand,
	}
}

func runDeleteValueCommand(cmd *cobra.Command, args []string) error {
	value, err := getMap(cmd, args[0])
	if err != nil {
		return err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	return value.Delete(ctx)
}
