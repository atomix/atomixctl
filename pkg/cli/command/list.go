package command

import (
	"context"
	"fmt"
	"github.com/atomix/atomix-go-client/pkg/client/list"
	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list {create,put,get,append,insert,remove,items,size,clear,delete}",
		Short: "Manage the state of a distributed list",
	}
	addClientFlags(cmd)
	cmd.PersistentFlags().StringP("name", "n", "", "the list name")
	cmd.PersistentFlags().Lookup("name").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_lists"},
	}
	cmd.MarkPersistentFlagRequired("name")
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

func newListFromName(cmd *cobra.Command) list.List {
	name, _ := cmd.Flags().GetString("name")
	group := newGroupFromName(cmd, name)
	m, err := group.GetList(newTimeoutContext(cmd), getPrimitiveName(name))
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return m
}

func newListCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "create",
		Args: cobra.NoArgs,
		Run:  runListCreateCommand,
	}
}

func runListCreateCommand(cmd *cobra.Command, _ []string) {
	list := newListFromName(cmd)
	list.Close()
	ExitWithOutput(fmt.Sprintf("Created %s", list.Name().String()))
}

func newListDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "delete",
		Args: cobra.NoArgs,
		Run:  runListDeleteCommand,
	}
}

func runListDeleteCommand(cmd *cobra.Command, _ []string) {
	list := newListFromName(cmd)
	err := list.Delete()
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(fmt.Sprintf("Deleted %s", list.Name().String()))
	}
}

func newListGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Args: cobra.NoArgs,
		Run:  runListGetCommand,
	}
	cmd.Flags().IntP("index", "i", -1, "the index to get")
	cmd.MarkFlagRequired("index")
	return cmd
}

func runListGetCommand(cmd *cobra.Command, _ []string) {
	list := newListFromName(cmd)
	index, _ := cmd.Flags().GetInt("index")
	value, err := list.Get(newTimeoutContext(cmd), index)
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
		Use:  "append",
		Args: cobra.NoArgs,
		Run:  runListAppendCommand,
	}
	cmd.Flags().StringP("value", "v", "", "the value to append")
	cmd.MarkFlagRequired("value")
	return cmd
}

func runListAppendCommand(cmd *cobra.Command, _ []string) {
	l := newListFromName(cmd)
	value, _ := cmd.Flags().GetString("value")
	err := l.Append(newTimeoutContext(cmd), []byte(value))
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(nil)
	}
}

func newListInsertCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "insert",
		Args: cobra.NoArgs,
		Run:  runListInsertCommand,
	}
	cmd.Flags().IntP("index", "i", -1, "the index to which to insert the value")
	cmd.MarkFlagRequired("index")
	cmd.Flags().StringP("value", "v", "", "the value to insert")
	cmd.MarkFlagRequired("value")
	return cmd
}

func runListInsertCommand(cmd *cobra.Command, _ []string) {
	l := newListFromName(cmd)
	index, _ := cmd.Flags().GetInt("index")
	value, _ := cmd.Flags().GetString("value")
	err := l.Insert(newTimeoutContext(cmd), int(index), []byte(value))
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(nil)
	}
}

func newListRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "remove",
		Args: cobra.NoArgs,
		Run:  runListRemoveCommand,
	}
	cmd.Flags().IntP("index", "i", -1, "the index to remove")
	cmd.MarkFlagRequired("index")
	cmd.Flags().Int64P("version", "v", 0, "the entry version")
	return cmd
}

func runListRemoveCommand(cmd *cobra.Command, _ []string) {
	m := newListFromName(cmd)
	index, _ := cmd.Flags().GetInt("index")
	value, err := m.Remove(newTimeoutContext(cmd), int(index))
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
	m := newListFromName(cmd)
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
	list := newListFromName(cmd)
	size, err := list.Len(newTimeoutContext(cmd))
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
	list := newListFromName(cmd)
	err := list.Clear(newTimeoutContext(cmd))
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithSuccess()
	}
}
