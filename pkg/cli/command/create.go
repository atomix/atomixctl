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

func newCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create {counter,election,list,lock,map,set,value}",
		Short: "Create a distributed primitive",
	}
	addClientFlags(cmd)
	cmd.AddCommand(newCreateCounterCommand())
	cmd.AddCommand(newCreateElectionCommand())
	cmd.AddCommand(newCreateListCommand())
	cmd.AddCommand(newCreateLockCommand())
	cmd.AddCommand(newCreateMapCommand())
	cmd.AddCommand(newCreateSetCommand())
	cmd.AddCommand(newCreateValueCommand())
	return cmd
}

func newCreateCounterCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "counter <name>",
		Args: cobra.ExactArgs(1),
		Run:  runCreateCounterCommand,
	}
}

func runCreateCounterCommand(cmd *cobra.Command, args []string) {
	counter := getCounter(cmd, args[0])
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	counter.Close(ctx)
	ExitWithOutput(fmt.Sprintf("Created counter %s", counter.Name().String()))
}

func newCreateElectionCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "election <name>",
		Args: cobra.ExactArgs(1),
		Run:  runCreateElectionCommand,
	}
}

func runCreateElectionCommand(cmd *cobra.Command, args []string) {
	election := getElection(cmd, args[0], "")
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	election.Close(ctx)
	ExitWithOutput(fmt.Sprintf("Created election %s", election.Name().String()))
}

func newCreateListCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "list <name>",
		Args: cobra.ExactArgs(1),
		Run:  runCreateListCommand,
	}
}

func runCreateListCommand(cmd *cobra.Command, args []string) {
	list := getList(cmd, args[0])
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	list.Close(ctx)
	ExitWithOutput(fmt.Sprintf("Created list %s", list.Name().String()))
}

func newCreateLockCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "lock <name>",
		Args: cobra.ExactArgs(1),
		Run:  runCreateLockCommand,
	}
}

func runCreateLockCommand(cmd *cobra.Command, args []string) {
	lock := getLock(cmd, args[0])
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	lock.Close(ctx)
	ExitWithOutput(fmt.Sprintf("Created lock %s", lock.Name().String()))
}

func newCreateMapCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "map <name>",
		Args: cobra.ExactArgs(1),
		Run:  runCreateMapCommand,
	}
}

func runCreateMapCommand(cmd *cobra.Command, args []string) {
	_map := getMap(cmd, args[0])
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	_map.Close(ctx)
	ExitWithOutput(fmt.Sprintf("Created map %s", _map.Name().String()))
}

func newCreateSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "set <name>",
		Args: cobra.ExactArgs(1),
		Run:  runCreateSetCommand,
	}
}

func runCreateSetCommand(cmd *cobra.Command, args []string) {
	set := getSet(cmd, args[0])
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	set.Close(ctx)
	ExitWithOutput(fmt.Sprintf("Created set %s", set.Name().String()))
}

func newCreateValueCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "value <name>",
		Args: cobra.ExactArgs(1),
		Run:  runCreateValueCommand,
	}
}

func runCreateValueCommand(cmd *cobra.Command, args []string) {
	value := getMap(cmd, args[0])
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	value.Close(ctx)
	ExitWithOutput(fmt.Sprintf("Created value %s", value.Name().String()))
}
