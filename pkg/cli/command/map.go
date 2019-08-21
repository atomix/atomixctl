package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/atomix/atomix-go-client/pkg/client/map"
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
	cmd.MarkFlagRequired("name")
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
	name, _ := cmd.PersistentFlags().GetString("name")
	if name == "" {
		ExitWithError(ExitBadArgs, errors.New("--name is a required flag"))
	}

	group := newGroupFromName(cmd, name)
	m, err := group.GetMap(newTimeoutContext(cmd), getPrimitiveName(name))
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
	cmd.Flags().StringP("value", "v", "", "the value to put into the map")
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
		opts = append(opts, _map.WithVersion(version))
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
	cmd.Flags().Int64("version", 0, "the entry version")
	return cmd
}

func runMapRemoveCommand(cmd *cobra.Command, _ []string) {
	m := newMapFromName(cmd)
	key, _ := cmd.Flags().GetString("key")
	version, _ := cmd.Flags().GetInt64("version")
	opts := []_map.RemoveOption{}
	if version > 0 {
		opts = append(opts, _map.WithVersion(version))
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
	ch := make(chan *_map.KeyValue)
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
