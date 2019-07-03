package command

import (
	"context"
	"fmt"
	"github.com/atomix/atomix-go-client/pkg/client/map_"
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

func newMapFromName(name string) map_.Map {
	group := newGroupFromName(name)
	m, err := group.GetMap(newTimeoutContext(), getPrimitiveName(name))
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return m
}

func newMapCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "create <map>",
		Args: cobra.ExactArgs(1),
		Run:  runMapCreateCommand,
	}
}

func runMapCreateCommand(cmd *cobra.Command, args []string) {
	map_ := newMapFromName(args[0])
	map_.Close()
	ExitWithOutput(fmt.Sprintf("Created %s", map_.Name().String()))
}

func newMapDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "delete <map>",
		Args: cobra.ExactArgs(1),
		Run:  runMapDeleteCommand,
	}
}

func runMapDeleteCommand(cmd *cobra.Command, args []string) {
	map_ := newMapFromName(args[0])
	err := map_.Delete()
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(fmt.Sprintf("Deleted %s", map_.Name().String()))
	}
}

func newMapGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "get <map> <key>",
		Args: cobra.ExactArgs(2),
		Run:  runMapGetCommand,
	}
}

func runMapGetCommand(cmd *cobra.Command, args []string) {
	map_ := newMapFromName(args[0])
	value, err := map_.Get(newTimeoutContext(), args[1])
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
		Use:  "put <map> <key> <value>",
		Args: cobra.ExactArgs(3),
		Run:  runMapPutCommand,
	}
	cmd.Flags().Int64P("version", "v", 0, "the entry version")
	return cmd
}

func runMapPutCommand(cmd *cobra.Command, args []string) {
	m := newMapFromName(args[0])
	version, _ := cmd.Flags().GetInt64("version")
	opts := []map_.PutOption{}
	if version > 0 {
		opts = append(opts, map_.WithVersion(version))
	}

	value, err := m.Put(newTimeoutContext(), args[1], []byte(args[2]), opts...)
	if err != nil {
		ExitWithError(ExitError, err)
	} else if value != nil {
		ExitWithOutput(value.String())
	} else {
		ExitWithOutput(nil)
	}
}

func newMapRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "remove <map> <key>",
		Args: cobra.ExactArgs(2),
		Run:  runMapRemoveCommand,
	}
	cmd.Flags().Int64P("version", "v", 0, "the entry version")
	return cmd
}

func runMapRemoveCommand(cmd *cobra.Command, args []string) {
	m := newMapFromName(args[0])
	version, _ := cmd.Flags().GetInt64("version")
	opts := []map_.RemoveOption{}
	if version > 0 {
		opts = append(opts, map_.WithVersion(version))
	}

	value, err := m.Remove(newTimeoutContext(), args[1], opts...)
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
		Use:  "keys <map>",
		Args: cobra.ExactArgs(1),
		Run:  runMapKeysCommand,
	}
}

func runMapKeysCommand(cmd *cobra.Command, args []string) {
	m := newMapFromName(args[0])
	ch := make(chan *map_.KeyValue)
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
		Use:  "size <map>",
		Args: cobra.ExactArgs(1),
		Run:  runMapSizeCommand,
	}
}

func runMapSizeCommand(cmd *cobra.Command, args []string) {
	map_ := newMapFromName(args[0])
	size, err := map_.Size(newTimeoutContext())
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(size)
	}
}

func newMapClearCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "clear <map>",
		Args: cobra.ExactArgs(1),
		Run:  runMapClearCommand,
	}
}

func runMapClearCommand(cmd *cobra.Command, args []string) {
	map_ := newMapFromName(args[0])
	err := map_.Clear(newTimeoutContext())
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithSuccess()
	}
}
