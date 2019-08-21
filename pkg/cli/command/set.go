package command

import (
	"errors"
	"fmt"
	"github.com/atomix/atomix-go-client/pkg/client/set"
	"github.com/spf13/cobra"
)

func newSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set {create,add,contains,remove,size,clear,delete}",
		Short: "Manage the state of a distributed set",
	}
	addClientFlags(cmd)
	cmd.PersistentFlags().StringP("name", "n", "", "the set name")
	cmd.PersistentFlags().Lookup("name").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_sets"},
	}
	cmd.MarkFlagRequired("name")
	cmd.AddCommand(newSetCreateCommand())
	cmd.AddCommand(newSetAddCommand())
	cmd.AddCommand(newSetContainsCommand())
	cmd.AddCommand(newSetRemoveCommand())
	cmd.AddCommand(newSetSizeCommand())
	cmd.AddCommand(newSetClearCommand())
	cmd.AddCommand(newSetDeleteCommand())
	return cmd
}

func newSetFromName(cmd *cobra.Command) set.Set {
	name, _ := cmd.Flags().GetString("name")
	if name == "" {
		ExitWithError(ExitBadArgs, errors.New("--name is a required flag"))
	}

	group := newGroupFromName(cmd, name)
	m, err := group.GetSet(newTimeoutContext(cmd), getPrimitiveName(name))
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return m
}

func newSetCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "create",
		Args: cobra.NoArgs,
		Run:  runSetCreateCommand,
	}
}

func runSetCreateCommand(cmd *cobra.Command, _ []string) {
	set := newSetFromName(cmd)
	set.Close()
	ExitWithOutput(fmt.Sprintf("Created %s", set.Name().String()))
}

func newSetDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "delete",
		Args: cobra.NoArgs,
		Run:  runSetDeleteCommand,
	}
}

func runSetDeleteCommand(cmd *cobra.Command, _ []string) {
	set := newSetFromName(cmd)
	err := set.Delete()
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(fmt.Sprintf("Deleted %s", set.Name().String()))
	}
}

func newSetAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "add",
		Args: cobra.NoArgs,
		Run:  runSetAddCommand,
	}
	cmd.Flags().StringP("value", "v", "", "the value to add")
	return cmd
}

func runSetAddCommand(cmd *cobra.Command, _ []string) {
	set := newSetFromName(cmd)
	if !cmd.Flags().Changed("value") {
		ExitWithError(ExitBadArgs, errors.New("--value is a required flag"))
	}
	value, _ := cmd.Flags().GetString("value")
	added, err := set.Add(newTimeoutContext(cmd), value)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(added)
	}
}

func newSetContainsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "contains",
		Args: cobra.NoArgs,
		Run:  runSetContainsCommand,
	}
	cmd.Flags().StringP("value", "v", "", "the value to check")
	return cmd
}

func runSetContainsCommand(cmd *cobra.Command, _ []string) {
	set := newSetFromName(cmd)
	if !cmd.Flags().Changed("value") {
		ExitWithError(ExitBadArgs, errors.New("--value is a required flag"))
	}
	value, _ := cmd.Flags().GetString("value")
	contains, err := set.Contains(newTimeoutContext(cmd), value)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(contains)
	}
}

func newSetRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "remove",
		Args: cobra.NoArgs,
		Run:  runSetRemoveCommand,
	}
	cmd.Flags().StringP("value", "v", "", "the value to remove")
	return cmd
}

func runSetRemoveCommand(cmd *cobra.Command, _ []string) {
	set := newSetFromName(cmd)
	if !cmd.Flags().Changed("value") {
		ExitWithError(ExitBadArgs, errors.New("--value is a required flag"))
	}
	value, _ := cmd.Flags().GetString("value")
	removed, err := set.Remove(newTimeoutContext(cmd), value)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(removed)
	}
}

func newSetSizeCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "size",
		Args: cobra.NoArgs,
		Run:  runSetSizeCommand,
	}
}

func runSetSizeCommand(cmd *cobra.Command, _ []string) {
	set := newSetFromName(cmd)
	size, err := set.Len(newTimeoutContext(cmd))
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(size)
	}
}

func newSetClearCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "clear",
		Args: cobra.NoArgs,
		Run:  runSetClearCommand,
	}
}

func runSetClearCommand(cmd *cobra.Command, _ []string) {
	set := newSetFromName(cmd)
	err := set.Clear(newTimeoutContext(cmd))
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithSuccess()
	}
}
