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
	"github.com/atomix/go-client/pkg/client/set"
	"github.com/spf13/cobra"
)

func newSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set {add,contains,remove,size,clear}",
		Short: "Manage the state of a distributed set",
	}
	addClientFlags(cmd)
	cmd.PersistentFlags().StringP("name", "n", "", "the set name")
	cmd.PersistentFlags().Lookup("name").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_sets"},
	}
	cmd.MarkPersistentFlagRequired("name")
	cmd.AddCommand(newSetAddCommand())
	cmd.AddCommand(newSetContainsCommand())
	cmd.AddCommand(newSetRemoveCommand())
	cmd.AddCommand(newSetSizeCommand())
	cmd.AddCommand(newSetClearCommand())
	return cmd
}

func getSetName(cmd *cobra.Command) string {
	name, _ := cmd.Flags().GetString("name")
	return name
}

func getSet(cmd *cobra.Command, name string) set.Set {
	database := getDatabase(cmd)
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	m, err := database.GetSet(ctx, name)
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return m
}

func newSetAddCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "add <value>",
		Args: cobra.ExactArgs(1),
		Run:  runSetAddCommand,
	}
}

func runSetAddCommand(cmd *cobra.Command, args []string) {
	set := getSet(cmd, getSetName(cmd))
	value := args[0]
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	added, err := set.Add(ctx, value)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(added)
	}
}

func newSetContainsCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "contains <value>",
		Args: cobra.ExactArgs(1),
		Run:  runSetContainsCommand,
	}
}

func runSetContainsCommand(cmd *cobra.Command, args []string) {
	set := getSet(cmd, getSetName(cmd))
	value := args[0]
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	contains, err := set.Contains(ctx, value)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(contains)
	}
}

func newSetRemoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "remove <value>",
		Args: cobra.ExactArgs(1),
		Run:  runSetRemoveCommand,
	}
}

func runSetRemoveCommand(cmd *cobra.Command, args []string) {
	set := getSet(cmd, getSetName(cmd))
	value := args[0]
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	removed, err := set.Remove(ctx, value)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(removed)
	}
}

func newSetSizeCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "size",
		Args: cobra.NoArgs,
		Run:  runSetSizeCommand,
	}
}

func runSetSizeCommand(cmd *cobra.Command, _ []string) {
	set := getSet(cmd, getSetName(cmd))
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	size, err := set.Len(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(size)
	}
}

func newSetClearCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "clear",
		Args: cobra.NoArgs,
		Run:  runSetClearCommand,
	}
}

func runSetClearCommand(cmd *cobra.Command, _ []string) {
	set := getSet(cmd, getSetName(cmd))
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	err := set.Clear(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithSuccess()
	}
}
