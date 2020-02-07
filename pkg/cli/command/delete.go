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
	"github.com/spf13/cobra"
)

func newDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete {counter,election,list,lock,map,set,value}",
		Short: "Delete a distributed primitive",
	}
	addClientFlags(cmd)
	cmd.AddCommand(newDeleteCounterCommand())
	cmd.AddCommand(newDeleteElectionCommand())
	cmd.AddCommand(newDeleteListCommand())
	cmd.AddCommand(newDeleteLockCommand())
	cmd.AddCommand(newDeleteMapCommand())
	cmd.AddCommand(newDeleteSetCommand())
	cmd.AddCommand(newDeleteValueCommand())
	return cmd
}

func newDeleteCounterCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "counter <name>",
		Args: cobra.ExactArgs(1),
		Run:  runDeleteCounterCommand,
	}
}

func runDeleteCounterCommand(cmd *cobra.Command, args []string) {
	counter := getCounter(cmd, args[0])
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	err := counter.Delete(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(fmt.Sprintf("Deleted counter %s", counter.Name().String()))
	}
}

func newDeleteElectionCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "election <name>",
		Args: cobra.ExactArgs(1),
		Run:  runDeleteElectionCommand,
	}
}

func runDeleteElectionCommand(cmd *cobra.Command, args []string) {
	election := getElection(cmd, args[0], "")
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	err := election.Delete(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(fmt.Sprintf("Deleted election %s", election.Name().String()))
	}
}

func newDeleteListCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "list <name>",
		Args: cobra.ExactArgs(1),
		Run:  runDeleteListCommand,
	}
}

func runDeleteListCommand(cmd *cobra.Command, args []string) {
	list := getList(cmd, args[0])
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	err := list.Delete(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(fmt.Sprintf("Deleted list %s", list.Name().String()))
	}
}

func newDeleteLockCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "lock <name>",
		Args: cobra.ExactArgs(1),
		Run:  runDeleteLockCommand,
	}
}

func runDeleteLockCommand(cmd *cobra.Command, args []string) {
	lock := getLock(cmd, args[0])
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	err := lock.Delete(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(fmt.Sprintf("Deleted lock %s", lock.Name().String()))
	}
}

func newDeleteMapCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "map <name>",
		Args: cobra.ExactArgs(1),
		Run:  runDeleteMapCommand,
	}
}

func runDeleteMapCommand(cmd *cobra.Command, args []string) {
	_map := getMap(cmd, args[0])
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	err := _map.Delete(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(fmt.Sprintf("Deleted map %s", _map.Name().String()))
	}
}

func newDeleteSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "set <name>",
		Args: cobra.ExactArgs(1),
		Run:  runDeleteSetCommand,
	}
}

func runDeleteSetCommand(cmd *cobra.Command, args []string) {
	set := getSet(cmd, args[0])
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	err := set.Delete(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(fmt.Sprintf("Deleted set %s", set.Name().String()))
	}
}

func newDeleteValueCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "value <name>",
		Args: cobra.ExactArgs(1),
		Run:  runDeleteValueCommand,
	}
}

func runDeleteValueCommand(cmd *cobra.Command, args []string) {
	value := getMap(cmd, args[0])
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	err := value.Delete(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(fmt.Sprintf("Deleted value %s", value.Name().String()))
	}
}
