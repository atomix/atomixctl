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
	"github.com/atomix/go-client/pkg/client/election"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

func newElectionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "election {create,enter,get,delete}",
		Short: "Managed the state of a distributed leader election",
	}
	addClientFlags(cmd)
	cmd.AddCommand(newElectionCreateCommand())
	cmd.AddCommand(newElectionGetCommand())
	cmd.AddCommand(newElectionEnterCommand())
	cmd.AddCommand(newElectionDeleteCommand())
	return cmd
}

func addElectionFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("name", "n", "", "the election name")
	cmd.Flags().Lookup("name").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_elections"},
	}
	cmd.MarkPersistentFlagRequired("name")
}

func getElectionName(cmd *cobra.Command) string {
	name, _ := cmd.Flags().GetString("name")
	return name
}

func getElection(cmd *cobra.Command, name string, id string) election.Election {
	database := getDatabase(cmd)
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	m, err := database.GetElection(ctx, name, election.WithID(id))
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return m
}

func newElectionCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "create <name>",
		Args: cobra.ExactArgs(1),
		Run:  runElectionCreateCommand,
	}
}

func runElectionCreateCommand(cmd *cobra.Command, args []string) {
	election := getElection(cmd, args[0], "")
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	election.Close(ctx)
	ExitWithOutput(fmt.Sprintf("Created %s", election.Name().String()))
}

func newElectionDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "delete <name>",
		Args: cobra.ExactArgs(1),
		Run:  runElectionDeleteCommand,
	}
}

func runElectionDeleteCommand(cmd *cobra.Command, args []string) {
	election := getElection(cmd, args[0], "")
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	err := election.Delete(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(fmt.Sprintf("Deleted %s", election.Name().String()))
	}
}

func newElectionGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Args: cobra.NoArgs,
		Run:  runElectionGetCommand,
	}
	addElectionFlags(cmd)
	return cmd
}

func runElectionGetCommand(cmd *cobra.Command, _ []string) {
	election := getElection(cmd, getElectionName(cmd), "")
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	term, err := election.GetTerm(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(term)
	}
}

func newElectionEnterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "enter <id>",
		Args: cobra.ExactArgs(1),
		Run:  runElectionEnterCommand,
	}
	addElectionFlags(cmd)
	return cmd
}

func runElectionEnterCommand(cmd *cobra.Command, args []string) {
	watchCh := make(chan *election.Event)
	election := getElection(cmd, getElectionName(cmd), args[0])
	ctx, cancel := getTimeoutContext(cmd)

	// Create a watch on the election
	err := election.Watch(context.Background(), watchCh)
	if err != nil {
		ExitWithError(ExitError, err)
	}

	// Enter the election
	_, err = election.Enter(ctx)
	cancel()
	if err != nil {
		ExitWithError(ExitError, err)
	}

	// Once we've successfully entered the election, wait for watch events
	Output("Entered election")
	signalCh := make(chan os.Signal, 2)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case event := <-watchCh:
			Output("Election state changed; Term: %d, Leader: %s, Candidates %v", event.Term.ID, event.Term.Leader, event.Term.Candidates)
		case <-signalCh:
			Output("Leaving election")
			ctx, cancel := getTimeoutContext(cmd)
			_, err = election.Leave(ctx)
			cancel()
			if err != nil {
				ExitWithError(ExitError, err)
			} else {
				ExitWithSuccess()
			}
		}
	}
}
