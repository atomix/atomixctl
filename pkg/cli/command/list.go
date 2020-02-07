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
	"github.com/atomix/go-client/pkg/client/list"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list {put,get,append,insert,remove,items,size,clear,watch}",
		Short: "Manage the state of a distributed list",
	}
	addClientFlags(cmd)
	cmd.PersistentFlags().StringP("name", "n", "", "the list name")
	cmd.PersistentFlags().Lookup("name").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_lists"},
	}
	cmd.MarkPersistentFlagRequired("name")
	cmd.AddCommand(newListGetCommand())
	cmd.AddCommand(newListAppendCommand())
	cmd.AddCommand(newListInsertCommand())
	cmd.AddCommand(newListRemoveCommand())
	cmd.AddCommand(newListItemsCommand())
	cmd.AddCommand(newListSizeCommand())
	cmd.AddCommand(newListClearCommand())
	cmd.AddCommand(newListWatchCommand())
	return cmd
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

func newListGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "get <index>",
		Args: cobra.ExactArgs(1),
		Run:  runListGetCommand,
	}
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
	return &cobra.Command{
		Use:  "append <value>",
		Args: cobra.ExactArgs(1),
		Run:  runListAppendCommand,
	}
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
	return &cobra.Command{
		Use:  "insert <index> <value>",
		Args: cobra.ExactArgs(2),
		Run:  runListInsertCommand,
	}
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
	return &cobra.Command{
		Use:  "items",
		Args: cobra.NoArgs,
		Run:  runListItemsCommand,
	}
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
	return &cobra.Command{
		Use:  "size",
		Args: cobra.NoArgs,
		Run:  runListSizeCommand,
	}
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
	return &cobra.Command{
		Use:  "clear",
		Args: cobra.NoArgs,
		Run:  runListClearCommand,
	}
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

func newListWatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "watch",
		Args: cobra.NoArgs,
		Run:  runListWatchCommand,
	}
	cmd.Flags().BoolP("replay", "r", false, "replay current list values at start")
	return cmd
}

func runListWatchCommand(cmd *cobra.Command, _ []string) {
	l := getList(cmd, getListName(cmd))
	watchCh := make(chan *list.Event)
	opts := []list.WatchOption{}
	replay, _ := cmd.Flags().GetBool("replay")
	if replay {
		opts = append(opts, list.WithReplay())
	}
	if err := l.Watch(context.Background(), watchCh, opts...); err != nil {
		ExitWithError(ExitError, err)
	}

	signalCh := make(chan os.Signal, 2)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case event := <-watchCh:
			switch event.Type {
			case list.EventNone:
				Output("Replayed: %v", event.Value)
			case list.EventInserted:
				Output("Inserted: %v", event.Value)
			case list.EventRemoved:
				Output("Removed: %v", event.Value)
			}
		case <-signalCh:
			ExitWithSuccess()
		}
	}
}
