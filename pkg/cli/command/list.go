package command

import (
	"context"
	"errors"
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
	name, _ := cmd.PersistentFlags().GetString("name")
	if name == "" {
		ExitWithError(ExitBadArgs, errors.New("--name is a required flag"))
	}

	group := newGroupFromName(name)
	m, err := group.GetList(newTimeoutContext(), getPrimitiveName(name))
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
	return cmd
}

func runListGetCommand(cmd *cobra.Command, _ []string) {
	list := newListFromName(cmd)
	if !cmd.Flags().Changed("index") {
		ExitWithError(ExitBadArgs, errors.New("--index is a required flag"))
	}
	index, _ := cmd.Flags().GetInt("index")
	value, err := list.Get(newTimeoutContext(), index)
	if err != nil {
		ExitWithError(ExitError, err)
	} else if value != "" {
		ExitWithOutput(value)
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
	return cmd
}

func runListAppendCommand(cmd *cobra.Command, _ []string) {
	l := newListFromName(cmd)
	if !cmd.Flags().Changed("value") {
		ExitWithError(ExitBadArgs, errors.New("--value is a required flag"))
	}
	value, _ := cmd.Flags().GetString("value")
	err := l.Append(newTimeoutContext(), value)
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
	cmd.Flags().StringP("value", "v", "", "the value to insert")
	return cmd
}

func runListInsertCommand(cmd *cobra.Command, _ []string) {
	l := newListFromName(cmd)
	if !cmd.Flags().Changed("index") {
		ExitWithError(ExitBadArgs, errors.New("--index is a required flag"))
	}
	if !cmd.Flags().Changed("value") {
		ExitWithError(ExitBadArgs, errors.New("--value is a required flag"))
	}
	index, _ := cmd.Flags().GetInt("index")
	value, _ := cmd.Flags().GetString("value")
	err := l.Insert(newTimeoutContext(), int(index), value)
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
	cmd.Flags().Int64P("version", "v", 0, "the entry version")
	return cmd
}

func runListRemoveCommand(cmd *cobra.Command, _ []string) {
	m := newListFromName(cmd)
	if !cmd.Flags().Changed("index") {
		ExitWithError(ExitBadArgs, errors.New("--index is a required flag"))
	}
	index, _ := cmd.Flags().GetInt("index")
	value, err := m.Remove(newTimeoutContext(), int(index))
	if err != nil {
		ExitWithError(ExitError, err)
	} else if value != "" {
		ExitWithOutput(value)
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
	ch := make(chan string)
	err := m.Items(context.TODO(), ch)
	if err != nil {
		ExitWithError(ExitError, err)
	}
	for value := range ch {
		println(value)
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
	size, err := list.Size(newTimeoutContext())
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
	err := list.Clear(newTimeoutContext())
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithSuccess()
	}
}
