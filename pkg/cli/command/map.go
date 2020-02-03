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
	cmd.PersistentFlags().StringP("name", "n", "", "the map name")
	cmd.PersistentFlags().Lookup("name").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_maps"},
	}
	cmd.MarkPersistentFlagRequired("name")
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

func newMapFromName(cmd *cobra.Command) _map.Map {
	name, _ := cmd.Flags().GetString("name")
	database := newDatabaseFromName(cmd, name)
	m, err := database.GetMap(newTimeoutContext(cmd), getPrimitiveName(name))
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return m
}

func newMapCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "create",
		Args: cobra.NoArgs,
		Run:  runMapCreateCommand,
	}
}

func runMapCreateCommand(cmd *cobra.Command, _ []string) {
	_map := newMapFromName(cmd)
	_map.Close()
	ExitWithOutput(fmt.Sprintf("Created %s", _map.Name().String()))
}

func newMapDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "delete",
		Args: cobra.NoArgs,
		Run:  runMapDeleteCommand,
	}
}

func runMapDeleteCommand(cmd *cobra.Command, _ []string) {
	_map := newMapFromName(cmd)
	err := _map.Delete()
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(fmt.Sprintf("Deleted %s", _map.Name().String()))
	}
}

func newMapGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Args: cobra.NoArgs,
		Run:  runMapGetCommand,
	}
	cmd.Flags().StringP("key", "k", "", "the key to get")
	cmd.MarkFlagRequired("key")
	return cmd
}

func runMapGetCommand(cmd *cobra.Command, _ []string) {
	_map := newMapFromName(cmd)
	key, _ := cmd.Flags().GetString("key")
	value, err := _map.Get(newTimeoutContext(cmd), key)
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
		Use:  "put",
		Args: cobra.NoArgs,
		Run:  runMapPutCommand,
	}
	cmd.Flags().StringP("key", "k", "", "the key to put")
	cmd.MarkFlagRequired("key")
	cmd.Flags().StringP("value", "v", "", "the value to put into the map")
	cmd.MarkFlagRequired("value")
	cmd.Flags().Int64("version", 0, "the entry version")
	return cmd
}

func runMapPutCommand(cmd *cobra.Command, _ []string) {
	m := newMapFromName(cmd)
	key, _ := cmd.Flags().GetString("key")
	value, _ := cmd.Flags().GetString("value")
	version, _ := cmd.Flags().GetInt64("version")
	opts := []_map.PutOption{}
	if version > 0 {
		opts = append(opts, _map.IfVersion(version))
	}

	kv, err := m.Put(newTimeoutContext(cmd), key, []byte(value), opts...)
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
		Use:  "remove",
		Args: cobra.NoArgs,
		Run:  runMapRemoveCommand,
	}
	cmd.Flags().StringP("key", "k", "", "the key to remove")
	cmd.MarkFlagRequired("key")
	cmd.Flags().Int64("version", 0, "the entry version")
	return cmd
}

func runMapRemoveCommand(cmd *cobra.Command, _ []string) {
	m := newMapFromName(cmd)
	key, _ := cmd.Flags().GetString("key")
	version, _ := cmd.Flags().GetInt64("version")
	opts := []_map.RemoveOption{}
	if version > 0 {
		opts = append(opts, _map.IfVersion(version))
	}

	value, err := m.Remove(newTimeoutContext(cmd), key, opts...)
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
	m := newMapFromName(cmd)
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
	_map := newMapFromName(cmd)
	size, err := _map.Len(newTimeoutContext(cmd))
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
	_map := newMapFromName(cmd)
	err := _map.Clear(newTimeoutContext(cmd))
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithSuccess()
	}
}
