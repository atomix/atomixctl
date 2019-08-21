package command

import (
	"errors"
	"fmt"
	"github.com/atomix/atomix-go-client/pkg/client/election"
	"github.com/spf13/cobra"
)

func newElectionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "election {create,enter,get,leave,delete}",
		Short: "Managed the state of a distributed leader election",
	}
	addClientFlags(cmd)
	cmd.PersistentFlags().StringP("name", "n", "", "the election name")
	cmd.AddCommand(newElectionCreateCommand())
	cmd.AddCommand(newElectionGetCommand())
	cmd.AddCommand(newElectionEnterCommand())
	cmd.AddCommand(newElectionLeaveCommand())
	cmd.AddCommand(newElectionDeleteCommand())
	return cmd
}

func newElectionFromName(cmd *cobra.Command) election.Election {
	name, _ := cmd.PersistentFlags().GetString("name")
	if name == "" {
		ExitWithError(ExitBadArgs, errors.New("--name is a required flag"))
	}

	group := newGroupFromName(name)
	m, err := group.GetElection(newTimeoutContext(), getPrimitiveName(name))
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return m
}

func newElectionCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "create",
		Args: cobra.NoArgs,
		Run:  runElectionCreateCommand,
	}
}

func runElectionCreateCommand(cmd *cobra.Command, _ []string) {
	election := newElectionFromName(cmd)
	election.Close()
	ExitWithOutput(fmt.Sprintf("Created %s", election.Name().String()))
}

func newElectionDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "delete",
		Args: cobra.NoArgs,
		Run:  runElectionDeleteCommand,
	}
}

func runElectionDeleteCommand(cmd *cobra.Command, _ []string) {
	election := newElectionFromName(cmd)
	err := election.Delete()
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(fmt.Sprintf("Deleted %s", election.Name().String()))
	}
}

func newElectionGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "get",
		Args: cobra.NoArgs,
		Run:  runElectionGetCommand,
	}
}

func runElectionGetCommand(cmd *cobra.Command, _ []string) {
	election := newElectionFromName(cmd)
	term, err := election.GetTerm(newTimeoutContext())
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(term)
	}
}

func newElectionEnterCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "enter",
		Args: cobra.NoArgs,
		Run:  runElectionEnterCommand,
	}
}

func runElectionEnterCommand(cmd *cobra.Command, _ []string) {
	election := newElectionFromName(cmd)
	term, err := election.Enter(newTimeoutContext())
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(term)
	}
}

func newElectionLeaveCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "leave",
		Args: cobra.NoArgs,
		Run:  runElectionLeaveCommand,
	}
}

func runElectionLeaveCommand(cmd *cobra.Command, _ []string) {
	election := newElectionFromName(cmd)
	err := election.Leave(newTimeoutContext())
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithSuccess()
	}
}
