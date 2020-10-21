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
	"errors"
	"github.com/atomix/go-client/pkg/client/log"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func newLogCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log {get,append,remove,entries,size,clear,watch}",
		Short: "Manage the state of a distributed log",
	}
	addClientFlags(cmd)
	cmd.PersistentFlags().StringP("name", "n", "", "the log name")
	cmd.PersistentFlags().Lookup("name").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_logs"},
	}
	cmd.MarkPersistentFlagRequired("name")
	cmd.AddCommand(newLogGetCommand())
	cmd.AddCommand(newLogAppendCommand())
	cmd.AddCommand(newLogRemoveCommand())
	cmd.AddCommand(newLogEntriesCommand())
	cmd.AddCommand(newLogSizeCommand())
	cmd.AddCommand(newLogClearCommand())
	cmd.AddCommand(newLogWatchCommand())
	return cmd
}

func getLogName(cmd *cobra.Command) string {
	name, _ := cmd.Flags().GetString("name")
	return name
}

func getLog(cmd *cobra.Command, name string) log.Log {
	database := getDatabase(cmd)
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	l, err := database.GetLog(ctx, name)
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return l
}

func newLogGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "get <index>",
		Args: cobra.ExactArgs(1),
		Run:  runLogGetCommand,
	}
}

func runLogGetCommand(cmd *cobra.Command, args []string) {
	l := getLog(cmd, getLogName(cmd))
	indexStr := args[0]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		ExitWithError(ExitBadArgs, err)
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	entry, err := l.Get(ctx, log.Index(index))
	if err != nil {
		ExitWithError(ExitError, err)
	} else if entry != nil {
		ExitWithOutput("Index: %d, Value: %v", entry.Index, entry.Value)
	} else {
		ExitWithOutput("<none>")
	}
}

func newLogAppendCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "append <value>",
		Args: cobra.ExactArgs(1),
		Run:  runLogAppendCommand,
	}
}

func runLogAppendCommand(cmd *cobra.Command, args []string) {
	l := getLog(cmd, getLogName(cmd))
	value := args[0]
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	entry, err := l.Append(ctx, []byte(value))
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput("Index: %d, Value: %v", entry.Index, entry.Value)
	}
}

func newLogRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "remove <index>",
		Args: cobra.ExactArgs(1),
		Run:  runLogRemoveCommand,
	}
	cmd.Flags().Int64P("version", "v", 0, "the entry version")
	return cmd
}

func runLogRemoveCommand(cmd *cobra.Command, args []string) {
	m := getLog(cmd, getLogName(cmd))
	indexStr := args[0]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		ExitWithError(ExitBadArgs, err)
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	entry, err := m.Remove(ctx, log.Index(index))
	if err != nil {
		ExitWithError(ExitError, err)
	} else if entry != nil {
		ExitWithOutput("Index: %d, Value: %v", entry.Index, entry.Value)
	} else {
		ExitWithOutput("<none>")
	}
}

func newLogEntriesCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "entries",
		Args: cobra.NoArgs,
		Run:  runLogEntriesCommand,
	}
}

func runLogEntriesCommand(cmd *cobra.Command, _ []string) {
	ExitWithError(ExitBadFeature, errors.New("not implemented"))
}

func newLogSizeCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "size",
		Args: cobra.NoArgs,
		Run:  runLogSizeCommand,
	}
}

func runLogSizeCommand(cmd *cobra.Command, _ []string) {
	log := getLog(cmd, getLogName(cmd))
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	size, err := log.Size(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput("%d", size)
	}
}

func newLogClearCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "clear",
		Args: cobra.NoArgs,
		Run:  runLogClearCommand,
	}
}

func runLogClearCommand(cmd *cobra.Command, _ []string) {
	log := getLog(cmd, getLogName(cmd))
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	err := log.Clear(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithSuccess()
	}
}

func newLogWatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "watch",
		Args: cobra.NoArgs,
		Run:  runLogWatchCommand,
	}
	cmd.Flags().BoolP("replay", "r", false, "replay current log entries at start")
	return cmd
}

func runLogWatchCommand(cmd *cobra.Command, _ []string) {
	l := getLog(cmd, getLogName(cmd))
	watchCh := make(chan *log.Event)
	opts := []log.WatchOption{}
	replay, _ := cmd.Flags().GetBool("replay")
	if replay {
		opts = append(opts, log.WithReplay())
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
			case log.EventNone:
				Output("Replayed: Index: %d, Value: %v", event.Entry.Index, event.Entry.Value)
			case log.EventAppended:
				Output("Appended: Index: %d, Value: %v", event.Entry.Index, event.Entry.Value)
			case log.EventRemoved:
				Output("Removed: Index: %d, Value: %v", event.Entry.Index, event.Entry.Value)
			}
		case <-signalCh:
			ExitWithSuccess()
		}
	}
}
