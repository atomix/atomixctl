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
)

func newMapCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "map {create,put,get,remove,size,clear,delete}",
		Short: "Manage the state of a distributed map",
	}
	addClientFlags(cmd)
	cmd.AddCommand(newMapCreateCommand())
	cmd.AddCommand(newMapGetCommand())
	cmd.AddCommand(newMapPutCommand())
	cmd.AddCommand(newMapRemoveCommand())
	cmd.AddCommand(newMapKeysCommand())
	cmd.AddCommand(newMapSizeCommand())
	cmd.AddCommand(newMapClearCommand())
	cmd.AddCommand(newMapDeleteCommand())
	return cmd
}

func addMapFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("name", "n", "", "the map name")
	cmd.Flags().Lookup("name").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_maps"},
	}
	cmd.MarkPersistentFlagRequired("name")
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

func newMapCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "create <name>",
		Args: cobra.ExactArgs(1),
		Run:  runMapCreateCommand,
	}
}

func runMapCreateCommand(cmd *cobra.Command, args []string) {
	_map := getMap(cmd, args[0])
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	_map.Close(ctx)
	ExitWithOutput(fmt.Sprintf("Created %s", _map.Name().String()))
}

func newMapDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "delete <name>",
		Args: cobra.ExactArgs(1),
		Run:  runMapDeleteCommand,
	}
}

func runMapDeleteCommand(cmd *cobra.Command, args []string) {
	_map := getMap(cmd, args[0])
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	err := _map.Delete(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(fmt.Sprintf("Deleted %s", _map.Name().String()))
	}
}

func newMapGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get <key>",
		Args: cobra.ExactArgs(1),
		Run:  runMapGetCommand,
	}
	addMapFlags(cmd)
	return cmd
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
	addMapFlags(cmd)
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
	addMapFlags(cmd)
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
	cmd := &cobra.Command{
		Use:  "keys",
		Args: cobra.NoArgs,
		Run:  runMapKeysCommand,
	}
	addMapFlags(cmd)
	return cmd
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
	cmd := &cobra.Command{
		Use:  "size",
		Args: cobra.NoArgs,
		Run:  runMapSizeCommand,
	}
	addMapFlags(cmd)
	return cmd
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
	cmd := &cobra.Command{
		Use:  "clear",
		Args: cobra.NoArgs,
		Run:  runMapClearCommand,
	}
	addMapFlags(cmd)
	return cmd
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
