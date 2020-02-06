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
	"fmt"
	"github.com/atomix/go-client/pkg/client/election"
	"github.com/spf13/cobra"
)

func newElectionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "election {create,enter,get,leave,delete}",
		Short: "Managed the state of a distributed leader election",
	}
	addClientFlags(cmd)
	cmd.PersistentFlags().StringP("name", "n", "", "the election name")
	cmd.PersistentFlags().Lookup("name").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_elections"},
	}
	cmd.MarkPersistentFlagRequired("name")
	cmd.AddCommand(newElectionCreateCommand())
	cmd.AddCommand(newElectionGetCommand())
	cmd.AddCommand(newElectionEnterCommand())
	cmd.AddCommand(newElectionLeaveCommand())
	cmd.AddCommand(newElectionDeleteCommand())
	return cmd
}

func getElection(cmd *cobra.Command) election.Election {
	name, _ := cmd.Flags().GetString("name")
	database := getDatabase(cmd)
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	m, err := database.GetElection(ctx, name)
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
	election := getElection(cmd)
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	election.Close(ctx)
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
	election := getElection(cmd)
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
	return &cobra.Command{
		Use:  "get",
		Args: cobra.NoArgs,
		Run:  runElectionGetCommand,
	}
}

func runElectionGetCommand(cmd *cobra.Command, _ []string) {
	election := getElection(cmd)
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
	return &cobra.Command{
		Use:  "enter",
		Args: cobra.NoArgs,
		Run:  runElectionEnterCommand,
	}
}

func runElectionEnterCommand(cmd *cobra.Command, _ []string) {
	election := getElection(cmd)
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	term, err := election.Enter(ctx)
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
	election := getElection(cmd)
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	_, err := election.Leave(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithSuccess()
	}
}
