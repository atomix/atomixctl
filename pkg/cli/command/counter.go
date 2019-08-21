package command

import (
	"fmt"
	"github.com/atomix/atomix-go-client/pkg/client/counter"
	"github.com/spf13/cobra"
)

func newCounterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "counter {create,get,set,increment,decrement,delete}",
		Short: "Manage the state of a distributed counter",
	}
	addClientFlags(cmd)
	cmd.PersistentFlags().StringP("name", "n", "", "the counter name")
	cmd.PersistentFlags().Lookup("name").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_counters"},
	}
	cmd.MarkPersistentFlagRequired("name")
	cmd.AddCommand(newCounterCreateCommand())
	cmd.AddCommand(newCounterGetCommand())
	cmd.AddCommand(newCounterSetCommand())
	cmd.AddCommand(newCounterIncrementCommand())
	cmd.AddCommand(newCounterDecrementCommand())
	cmd.AddCommand(newCounterDeleteCommand())
	return cmd
}

func newCounterFromName(cmd *cobra.Command) counter.Counter {
	name, _ := cmd.Flags().GetString("name")
	group := newGroupFromName(cmd, name)
	m, err := group.GetCounter(newTimeoutContext(cmd), getPrimitiveName(name))
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return m
}

func newCounterCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "create",
		Args: cobra.NoArgs,
		Run:  runCounterCreateCommand,
	}
}

func runCounterCreateCommand(cmd *cobra.Command, _ []string) {
	counter := newCounterFromName(cmd)
	counter.Close()
	ExitWithOutput(fmt.Sprintf("Created %s", counter.Name().String()))
}

func newCounterDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "delete",
		Args: cobra.NoArgs,
		Run:  runCounterDeleteCommand,
	}
}

func runCounterDeleteCommand(cmd *cobra.Command, _ []string) {
	counter := newCounterFromName(cmd)
	err := counter.Delete()
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(fmt.Sprintf("Deleted %s", counter.Name().String()))
	}
}

func newCounterGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "get",
		Args: cobra.NoArgs,
		Run:  runCounterGetCommand,
	}
}

func runCounterGetCommand(cmd *cobra.Command, _ []string) {
	counter := newCounterFromName(cmd)
	value, err := counter.Get(newTimeoutContext(cmd))
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(value)
	}
}

func newCounterSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "set",
		Args: cobra.NoArgs,
		Run:  runCounterSetCommand,
	}
	cmd.Flags().Int64P("value", "v", 0, "the value to set")
	cmd.MarkFlagRequired("value")
	return cmd
}

func runCounterSetCommand(cmd *cobra.Command, _ []string) {
	counter := newCounterFromName(cmd)
	value, _ := cmd.Flags().GetInt64("value")
	err := counter.Set(newTimeoutContext(cmd), value)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(value)
	}
}

func newCounterIncrementCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "increment",
		Args: cobra.NoArgs,
		Run:  runCounterIncrementCommand,
	}
	cmd.Flags().Int64P("delta", "d", 1, "the delta by which to increment the counter")
	return cmd
}

func runCounterIncrementCommand(cmd *cobra.Command, _ []string) {
	counter := newCounterFromName(cmd)
	delta, _ := cmd.Flags().GetInt64("delta")
	value, err := counter.Increment(newTimeoutContext(cmd), delta)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(value)
	}
}

func newCounterDecrementCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "decrement",
		Args: cobra.NoArgs,
		Run:  runCounterDecrementCommand,
	}
	cmd.Flags().Int64P("delta", "d", 1, "the delta by which to decrement the counter")
	return cmd
}

func runCounterDecrementCommand(cmd *cobra.Command, _ []string) {
	counter := newCounterFromName(cmd)
	delta, _ := cmd.Flags().GetInt64("delta")
	value, err := counter.Decrement(newTimeoutContext(cmd), delta)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(value)
	}
}
