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

func newCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create {counter,election,list,lock,map,set,value}",
		Short: "Create a distributed primitive",
	}
	cmd.AddCommand(newCreateCounterCommand())
	cmd.AddCommand(newCreateElectionCommand())
	cmd.AddCommand(newCreateListCommand())
	cmd.AddCommand(newCreateLockCommand())
	cmd.AddCommand(newCreateLogCommand())
	cmd.AddCommand(newCreateMapCommand())
	cmd.AddCommand(newCreateSetCommand())
	cmd.AddCommand(newCreateValueCommand())
	return cmd
}

func newCreateCounterCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "counter <name>",
		Args: cobra.ExactArgs(1),
		RunE: runCreateCounterCommand,
	}
}

func runCreateCounterCommand(cmd *cobra.Command, args []string) error {
	counter, err := getCounter(cmd, args[0])
	if err != nil {
		return err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	counter.Close(ctx)
	cmd.Println(counter.Name().String())
	return nil
}

func newCreateElectionCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "election <name>",
		Args: cobra.ExactArgs(1),
		RunE: runCreateElectionCommand,
	}
}

func runCreateElectionCommand(cmd *cobra.Command, args []string) error {
	election, err := getElection(cmd, args[0], "")
	if err != nil {
		return err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	election.Close(ctx)
	cmd.Println(election.Name().String())
	return nil
}

func newCreateListCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "list <name>",
		Args: cobra.ExactArgs(1),
		RunE: runCreateListCommand,
	}
}

func runCreateListCommand(cmd *cobra.Command, args []string) error {
	list, err := getList(cmd, args[0])
	if err != nil {
		return err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	list.Close(ctx)
	cmd.Println(list.Name().String())
	return nil
}

func newCreateLockCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "lock <name>",
		Args: cobra.ExactArgs(1),
		RunE: runCreateLockCommand,
	}
}

func runCreateLockCommand(cmd *cobra.Command, args []string) error {
	lock, err := getLock(cmd, args[0])
	if err != nil {
		return err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	lock.Close(ctx)
	cmd.Println(lock.Name().String())
	return nil
}

func newCreateLogCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "log <name>",
		Args: cobra.ExactArgs(1),
		RunE: runCreateLogCommand,
	}
}

func runCreateLogCommand(cmd *cobra.Command, args []string) error {
	log, err := getLog(cmd, args[0])
	if err != nil {
		return err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	log.Close(ctx)
	cmd.Println(log.Name().String())
	return nil
}

func newCreateMapCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "map <name>",
		Args: cobra.ExactArgs(1),
		RunE: runCreateMapCommand,
	}
}

func runCreateMapCommand(cmd *cobra.Command, args []string) error {
	_map, err := getMap(cmd, args[0])
	if err != nil {
		return err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	_map.Close(ctx)
	cmd.Println(_map.Name().String())
	return nil
}

func newCreateSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "set <name>",
		Args: cobra.ExactArgs(1),
		RunE: runCreateSetCommand,
	}
}

func runCreateSetCommand(cmd *cobra.Command, args []string) error {
	set, err := getSet(cmd, args[0])
	if err != nil {
		return err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	set.Close(ctx)
	cmd.Println(set.Name().String())
	return nil
}

func newCreateValueCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "value <name>",
		Args: cobra.ExactArgs(1),
		RunE: runCreateValueCommand,
	}
}

func runCreateValueCommand(cmd *cobra.Command, args []string) error {
	value, err := getValue(cmd, args[0])
	if err != nil {
		return err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	value.Close(ctx)
	cmd.Println(value.Name().String())
	return nil
}
