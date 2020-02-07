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
	"github.com/atomix/go-client/pkg/client/map"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

func newMapCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "map {put,get,remove,size,clear}",
		Short: "Manage the state of a distributed map",
	}
	addClientFlags(cmd)
	cmd.PersistentFlags().StringP("name", "n", "", "the map name")
	cmd.PersistentFlags().Lookup("name").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_maps"},
	}
	cmd.MarkPersistentFlagRequired("name")
	cmd.AddCommand(newMapGetCommand())
	cmd.AddCommand(newMapPutCommand())
	cmd.AddCommand(newMapRemoveCommand())
	cmd.AddCommand(newMapKeysCommand())
	cmd.AddCommand(newMapSizeCommand())
	cmd.AddCommand(newMapClearCommand())
	cmd.AddCommand(newMapWatchCommand())
	return cmd
}

func getMapName(cmd *cobra.Command) string {
	name, _ := cmd.Flags().GetString("name")
	return name
}

func getMap(cmd *cobra.Command, name string) _map.Map {
	database := getDatabase(cmd)
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	m, err := database.GetMap(ctx, name)
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return m
}

func newMapGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "get <key>",
		Args: cobra.ExactArgs(1),
		Run:  runMapGetCommand,
	}
}

func runMapGetCommand(cmd *cobra.Command, args []string) {
	_map := getMap(cmd, getMapName(cmd))
	key := args[0]
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	value, err := _map.Get(ctx, key)
	if err != nil {
		ExitWithError(ExitError, err)
	} else if value != nil {
		ExitWithOutput(value.String())
	} else {
		ExitWithOutput(nil)
	}
}

func newMapPutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "put <key> <value>",
		Args: cobra.ExactArgs(2),
		Run:  runMapPutCommand,
	}
	cmd.Flags().Int64("version", 0, "the entry version")
	return cmd
}

func runMapPutCommand(cmd *cobra.Command, args []string) {
	m := getMap(cmd, getMapName(cmd))
	key := args[0]
	value := args[1]
	version, _ := cmd.Flags().GetInt64("version")
	opts := []_map.PutOption{}
	if version > 0 {
		opts = append(opts, _map.IfVersion(version))
	}

	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	kv, err := m.Put(ctx, key, []byte(value), opts...)
	if err != nil {
		ExitWithError(ExitError, err)
	} else if kv != nil {
		ExitWithOutput(kv.String())
	} else {
		ExitWithOutput(nil)
	}
}

func newMapRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "remove <key>",
		Args: cobra.ExactArgs(1),
		Run:  runMapRemoveCommand,
	}
	cmd.Flags().Int64("version", 0, "the entry version")
	return cmd
}

func runMapRemoveCommand(cmd *cobra.Command, args []string) {
	m := getMap(cmd, getMapName(cmd))
	key := args[0]
	version, _ := cmd.Flags().GetInt64("version")
	opts := []_map.RemoveOption{}
	if version > 0 {
		opts = append(opts, _map.IfVersion(version))
	}

	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	value, err := m.Remove(ctx, key, opts...)
	if err != nil {
		ExitWithError(ExitError, err)
	} else if value != nil {
		ExitWithOutput(value.String())
	} else {
		ExitWithOutput(nil)
	}
}

func newMapKeysCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "keys",
		Args: cobra.NoArgs,
		Run:  runMapKeysCommand,
	}
}

func runMapKeysCommand(cmd *cobra.Command, _ []string) {
	m := getMap(cmd, getMapName(cmd))
	ch := make(chan *_map.Entry)
	err := m.Entries(context.TODO(), ch)
	if err != nil {
		ExitWithError(ExitError, err)
	}
	for kv := range ch {
		println(fmt.Sprintf("%v", kv))
	}
}

func newMapSizeCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "size",
		Args: cobra.NoArgs,
		Run:  runMapSizeCommand,
	}
}

func runMapSizeCommand(cmd *cobra.Command, _ []string) {
	_map := getMap(cmd, getMapName(cmd))
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	size, err := _map.Len(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(size)
	}
}

func newMapClearCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "clear",
		Args: cobra.NoArgs,
		Run:  runMapClearCommand,
	}
}

func runMapClearCommand(cmd *cobra.Command, _ []string) {
	_map := getMap(cmd, getMapName(cmd))
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	err := _map.Clear(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithSuccess()
	}
}

func newMapWatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "watch",
		Args: cobra.NoArgs,
		Run:  runMapWatchCommand,
	}
	cmd.Flags().BoolP("replay", "r", false, "replay current map entries at start")
	return cmd
}

func runMapWatchCommand(cmd *cobra.Command, _ []string) {
	m := getMap(cmd, getMapName(cmd))
	watchCh := make(chan *_map.Event)
	opts := []_map.WatchOption{}
	replay, _ := cmd.Flags().GetBool("replay")
	if replay {
		opts = append(opts, _map.WithReplay())
	}
	if err := m.Watch(context.Background(), watchCh, opts...); err != nil {
		ExitWithError(ExitError, err)
	}

	signalCh := make(chan os.Signal, 2)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case event := <-watchCh:
			switch event.Type {
			case _map.EventNone:
				Output("Replayed: Key: %s, Value: %v, Version: %d", event.Entry.Key, event.Entry.Value, event.Entry.Version)
			case _map.EventInserted:
				Output("Inserted: Key: %s, Value: %v, Version: %d", event.Entry.Key, event.Entry.Value, event.Entry.Version)
			case _map.EventUpdated:
				Output("Updated: Key: %s, Value: %v, Version: %d", event.Entry.Key, event.Entry.Value, event.Entry.Version)
			case _map.EventRemoved:
				Output("Removed: Key: %s, Value: %v, Version: %d", event.Entry.Key, event.Entry.Value, event.Entry.Version)
			}
		case <-signalCh:
			ExitWithSuccess()
		}
	}
}
