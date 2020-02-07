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
	"context"
	"fmt"
	"github.com/atomix/go-client/pkg/client/list"
	"github.com/spf13/cobra"
	"strconv"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list {create,put,get,append,insert,remove,items,size,clear,delete}",
		Short: "Manage the state of a distributed list",
	}
	addClientFlags(cmd)
	cmd.AddCommand(newListCreateCommand())
	cmd.AddCommand(newListGetCommand())
	cmd.AddCommand(newListAppendCommand())
	cmd.AddCommand(newListInsertCommand())
	cmd.AddCommand(newListRemoveCommand())
	cmd.AddCommand(newListItemsCommand())
	cmd.AddCommand(newListSizeCommand())
	cmd.AddCommand(newListClearCommand())
	cmd.AddCommand(newListDeleteCommand())
	return cmd
}

func addListFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("name", "n", "", "the list name")
	cmd.Flags().Lookup("name").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_lists"},
	}
	cmd.MarkPersistentFlagRequired("name")
}

func getListName(cmd *cobra.Command) string {
	name, _ := cmd.Flags().GetString("name")
	return name
}

func getList(cmd *cobra.Command, name string) list.List {
	database := getDatabase(cmd)
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	m, err := database.GetList(ctx, name)
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return m
}

func newListCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "create <name>",
		Args: cobra.ExactArgs(1),
		Run:  runListCreateCommand,
	}
}

func runListCreateCommand(cmd *cobra.Command, args []string) {
	list := getList(cmd, args[0])
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	list.Close(ctx)
	ExitWithOutput(fmt.Sprintf("Created %s", list.Name().String()))
}

func newListDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "delete <name>",
		Args: cobra.ExactArgs(1),
		Run:  runListDeleteCommand,
	}
}

func runListDeleteCommand(cmd *cobra.Command, args []string) {
	list := getList(cmd, args[0])
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	err := list.Delete(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(fmt.Sprintf("Deleted %s", list.Name().String()))
	}
}

func newListGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get <index>",
		Args: cobra.ExactArgs(1),
		Run:  runListGetCommand,
	}
	addListFlags(cmd)
	return cmd
}

func runListGetCommand(cmd *cobra.Command, args []string) {
	list := getList(cmd, getListName(cmd))
	indexStr := args[0]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		ExitWithError(ExitBadArgs, err)
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	value, err := list.Get(ctx, index)
	if err != nil {
		ExitWithError(ExitError, err)
	} else if value != nil {
		ExitWithOutput(string(value))
	} else {
		ExitWithOutput(nil)
	}
}

func newListAppendCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "append <value>",
		Args: cobra.ExactArgs(1),
		Run:  runListAppendCommand,
	}
	addListFlags(cmd)
	return cmd
}

func runListAppendCommand(cmd *cobra.Command, args []string) {
	l := getList(cmd, getListName(cmd))
	value := args[0]
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	err := l.Append(ctx, []byte(value))
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(nil)
	}
}

func newListInsertCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "insert <index> <value>",
		Args: cobra.ExactArgs(2),
		Run:  runListInsertCommand,
	}
	addListFlags(cmd)
	return cmd
}

func runListInsertCommand(cmd *cobra.Command, args []string) {
	l := getList(cmd, getListName(cmd))
	indexStr := args[0]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		ExitWithError(ExitBadArgs, err)
	}
	value := args[1]
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	err = l.Insert(ctx, int(index), []byte(value))
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(nil)
	}
}

func newListRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "remove <index>",
		Args: cobra.ExactArgs(1),
		Run:  runListRemoveCommand,
	}
	addListFlags(cmd)
	cmd.Flags().Int64P("version", "v", 0, "the entry version")
	return cmd
}

func runListRemoveCommand(cmd *cobra.Command, args []string) {
	m := getList(cmd, getListName(cmd))
	indexStr := args[0]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		ExitWithError(ExitBadArgs, err)
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	value, err := m.Remove(ctx, int(index))
	if err != nil {
		ExitWithError(ExitError, err)
	} else if value != nil {
		ExitWithOutput(string(value))
	} else {
		ExitWithOutput(nil)
	}
}

func newListItemsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "items",
		Args: cobra.NoArgs,
		Run:  runListItemsCommand,
	}
	addListFlags(cmd)
	return cmd
}

func runListItemsCommand(cmd *cobra.Command, _ []string) {
	m := getList(cmd, getListName(cmd))
	ch := make(chan []byte)
	err := m.Items(context.TODO(), ch)
	if err != nil {
		ExitWithError(ExitError, err)
	}
	for value := range ch {
		println(string(value))
	}
}

func newListSizeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "size",
		Args: cobra.NoArgs,
		Run:  runListSizeCommand,
	}
	addListFlags(cmd)
	return cmd
}

func runListSizeCommand(cmd *cobra.Command, _ []string) {
	list := getList(cmd, getListName(cmd))
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	size, err := list.Len(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(size)
	}
}

func newListClearCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "clear",
		Args: cobra.NoArgs,
		Run:  runListClearCommand,
	}
	addListFlags(cmd)
	return cmd
}

func runListClearCommand(cmd *cobra.Command, _ []string) {
	list := getList(cmd, getListName(cmd))
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	err := list.Clear(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithSuccess()
	}
}
