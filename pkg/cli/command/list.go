package command

import (
	"context"
	"fmt"
	"github.com/atomix/atomix-go-client/pkg/client/list"
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

func newListFromName(name string) list.List {
	group := newGroupFromName(name)
	m, err := group.GetList(newTimeoutContext(), getPrimitiveName(name))
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return m
}

func newListCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "create <list>",
		Args: cobra.ExactArgs(1),
		Run:  runListCreateCommand,
	}
}

func runListCreateCommand(cmd *cobra.Command, args []string) {
	list := newListFromName(args[0])
	list.Close()
	ExitWithOutput(fmt.Sprintf("Created %s", list.Name().String()))
}

func newListDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "delete <list>",
		Args: cobra.ExactArgs(1),
		Run:  runListDeleteCommand,
	}
}

func runListDeleteCommand(cmd *cobra.Command, args []string) {
	list := newListFromName(args[0])
	err := list.Delete()
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(fmt.Sprintf("Deleted %s", list.Name().String()))
	}
}

func newListGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "get <list> <index>",
		Args: cobra.ExactArgs(2),
		Run:  runListGetCommand,
	}
}

func runListGetCommand(cmd *cobra.Command, args []string) {
	list := newListFromName(args[0])
	index, err := strconv.ParseInt(args[1], 0, 32)
	if err != nil {
		ExitWithError(ExitError, err)
	}
	value, err := list.Get(newTimeoutContext(), int(index))
	if err != nil {
		ExitWithError(ExitError, err)
	} else if value != "" {
		ExitWithOutput(value)
	} else {
		ExitWithOutput(nil)
	}
}

func newListAppendCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "append <list> <value>",
		Args: cobra.ExactArgs(2),
		Run:  runListAppendCommand,
	}
}

func runListAppendCommand(cmd *cobra.Command, args []string) {
	l := newListFromName(args[0])
	err := l.Append(newTimeoutContext(), args[1])
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(nil)
	}
}

func newListInsertCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "insert <list> <index> <value>",
		Args: cobra.ExactArgs(2),
		Run:  runListInsertCommand,
	}
}

func runListInsertCommand(cmd *cobra.Command, args []string) {
	l := newListFromName(args[0])
	index, err := strconv.ParseInt(args[1], 0, 32)
	if err != nil {
		ExitWithError(ExitError, err)
	}
	value := args[2]
	err = l.Insert(newTimeoutContext(), int(index), value)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(nil)
	}
}

func newListRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "remove <list> <index>",
		Args: cobra.ExactArgs(2),
		Run:  runListRemoveCommand,
	}
	cmd.Flags().Int64P("version", "v", 0, "the entry version")
	return cmd
}

func runListRemoveCommand(cmd *cobra.Command, args []string) {
	m := newListFromName(args[0])
	index, err := strconv.ParseInt(args[1], 0, 32)
	if err != nil {
		ExitWithError(ExitError, err)
	}
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
		Use:  "items <list>",
		Args: cobra.ExactArgs(1),
		Run:  runListItemsCommand,
	}
}

func runListItemsCommand(cmd *cobra.Command, args []string) {
	m := newListFromName(args[0])
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
		Use:  "size <list>",
		Args: cobra.ExactArgs(1),
		Run:  runListSizeCommand,
	}
}

func runListSizeCommand(cmd *cobra.Command, args []string) {
	list := newListFromName(args[0])
	size, err := list.Size(newTimeoutContext())
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(size)
	}
}

func newListClearCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "clear <list>",
		Args: cobra.ExactArgs(1),
		Run:  runListClearCommand,
	}
}

func runListClearCommand(cmd *cobra.Command, args []string) {
	list := newListFromName(args[0])
	err := list.Clear(newTimeoutContext())
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithSuccess()
	}
}
