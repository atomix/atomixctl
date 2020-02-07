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
	"github.com/atomix/go-client/pkg/client/counter"
	"github.com/spf13/cobra"
	"strconv"
)

func newCounterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "counter {create,get,set,increment,decrement,delete}",
		Short: "Manage the state of a distributed counter",
	}
	addClientFlags(cmd)
	cmd.AddCommand(newCounterCreateCommand())
	cmd.AddCommand(newCounterGetCommand())
	cmd.AddCommand(newCounterSetCommand())
	cmd.AddCommand(newCounterIncrementCommand())
	cmd.AddCommand(newCounterDecrementCommand())
	cmd.AddCommand(newCounterDeleteCommand())
	return cmd
}

func addCounterFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("name", "n", "", "the list name")
	cmd.Flags().Lookup("name").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_counters"},
	}
	cmd.MarkPersistentFlagRequired("name")
}

func getCounterName(cmd *cobra.Command) string {
	name, _ := cmd.Flags().GetString("name")
	return name
}

func getCounter(cmd *cobra.Command, name string) counter.Counter {
	database := getDatabase(cmd)
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	m, err := database.GetCounter(ctx, name)
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return m
}

func newCounterCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "create <name>",
		Args: cobra.NoArgs,
		Run:  runCounterCreateCommand,
	}
}

func runCounterCreateCommand(cmd *cobra.Command, args []string) {
	counter := getCounter(cmd, args[0])
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	counter.Close(ctx)
	ExitWithOutput(fmt.Sprintf("Created %s", counter.Name().String()))
}

func newCounterDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "delete <name>",
		Args: cobra.NoArgs,
		Run:  runCounterDeleteCommand,
	}
}

func runCounterDeleteCommand(cmd *cobra.Command, args []string) {
	counter := getCounter(cmd, args[0])
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	err := counter.Delete(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(fmt.Sprintf("Deleted %s", counter.Name().String()))
	}
}

func newCounterGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Args: cobra.NoArgs,
		Run:  runCounterGetCommand,
	}
	addCounterFlags(cmd)
	return cmd
}

func runCounterGetCommand(cmd *cobra.Command, _ []string) {
	counter := getCounter(cmd, getCounterName(cmd))
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	value, err := counter.Get(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(value)
	}
}

func newCounterSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "set <value>",
		Args: cobra.ExactArgs(1),
		Run:  runCounterSetCommand,
	}
	addCounterFlags(cmd)
	return cmd
}

func runCounterSetCommand(cmd *cobra.Command, args []string) {
	counter := getCounter(cmd, getCounterName(cmd))
	value, err := strconv.Atoi(args[0])
	if err != nil {
		ExitWithError(ExitBadArgs, err)
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	err = counter.Set(ctx, int64(value))
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(value)
	}
}

func newCounterIncrementCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "increment [delta]",
		Args: cobra.MaximumNArgs(1),
		Run:  runCounterIncrementCommand,
	}
	addCounterFlags(cmd)
	return cmd
}

func runCounterIncrementCommand(cmd *cobra.Command, args []string) {
	counter := getCounter(cmd, getCounterName(cmd))
	var delta int64
	if len(args) > 0 {
		value, err := strconv.Atoi(args[0])
		if err != nil {
			ExitWithError(ExitBadArgs, err)
		}
		delta = int64(value)
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	value, err := counter.Increment(ctx, delta)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(value)
	}
}

func newCounterDecrementCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "decrement [delta]",
		Args: cobra.MaximumNArgs(1),
		Run:  runCounterDecrementCommand,
	}
	addCounterFlags(cmd)
	return cmd
}

func runCounterDecrementCommand(cmd *cobra.Command, args []string) {
	counter := getCounter(cmd, getCounterName(cmd))
	var delta int64
	if len(args) > 0 {
		value, err := strconv.Atoi(args[0])
		if err != nil {
			ExitWithError(ExitBadArgs, err)
		}
		delta = int64(value)
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	value, err := counter.Decrement(ctx, delta)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(value)
	}
}
