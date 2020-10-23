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
	"gopkg.in/yaml.v2"
)

func newElectionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "election <name> [...]",
		Short:              "Manage the state of a distributed leader election",
		Args:               cobra.MinimumNArgs(1),
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If only the name was specified, open an interactive shell
			name := args[0]
			if name == "-h" || name == "--help" {
				return cmd.Help()
			}
			if len(args) == 1 {
				ctx := getContext()
				if ctx == nil {
					ctx = newContext("atomix", "election", name)
					setContext(ctx)
				} else {
					ctx = ctx.withCommand("election", name)
				}
				return ctx.run()
			}

			// Get the command for the specified operation
			var subCmd *cobra.Command
			op := args[1]
			switch op {
			case "enter":
				subCmd = newElectionEnterCommand(name)
			case "get":
				subCmd = newElectionGetCommand(name)
			case "watch":
				subCmd = newElectionWatchCommand(name)
			case "help", "-h", "--help":
				if len(args) == 2 {
					helpCmd := &cobra.Command{
						Use:   fmt.Sprintf("election %s [...]", name),
						Short: "Manage the state of a distributed leader election",
					}
					helpCmd.AddCommand(newElectionEnterCommand(name))
					helpCmd.AddCommand(newElectionGetCommand(name))
					helpCmd.AddCommand(newElectionWatchCommand(name))
					return helpCmd.Help()
				} else {
					var helpCmd *cobra.Command
					switch args[2] {
					case "enter":
						helpCmd = newElectionEnterCommand(name)
					case "get":
						helpCmd = newElectionGetCommand(name)
					case "watch":
						helpCmd = newElectionWatchCommand(name)
					default:
						return fmt.Errorf("unknown command %s", args[2])
					}
					return helpCmd.Help()
				}
			default:
				return fmt.Errorf("unknown command %s", op)
			}
			addClientFlags(subCmd)

			// Set the arguments after the name and execute the command
			subCmd.SetArgs(args[2:])
			return subCmd.Execute()
		},
	}
	return cmd
}

func getElection(cmd *cobra.Command, name string, id string) (election.Election, error) {
	database, err := getDatabase(cmd)
	if err != nil {
		return nil, err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	return database.GetElection(ctx, name, election.WithID(id))
}

func newElectionGetCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use: "get {leader,term}",
	}
	cmd.AddCommand(newElectionGetLeaderCommand(name))
	cmd.AddCommand(newElectionGetTermCommand(name))
	return cmd
}

func newElectionGetLeaderCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use: "leader [options]",
		RunE: func(cmd *cobra.Command, args []string) error {
			election, err := getElection(cmd, name, "")
			if err != nil {
				return err
			}
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			term, err := election.GetTerm(ctx)
			if err != nil {
				return err
			} else if term != nil {
				cmd.Println(term.Leader)
			}
			return nil
		},
	}
	return cmd
}

func newElectionGetTermCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use: "term [options]",
		RunE: func(cmd *cobra.Command, args []string) error {
			election, err := getElection(cmd, name, "")
			if err != nil {
				return err
			}
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			term, err := election.GetTerm(ctx)
			if err != nil {
				return err
			} else if term != nil {
				bytes, err := yaml.Marshal(term)
				if err != nil {
					return err
				}
				cmd.Println(string(bytes))
			}
			return nil
		},
	}
	return cmd
}

func newElectionEnterCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "enter <id>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			watchCh := make(chan *election.Event)
			election, err := getElection(cmd, name, args[0])
			if err != nil {
				return err
			}

			// Create a watch on the election
			watchCtx, watchCancel := getCancelContext(cmd)
			defer watchCancel()
			err = election.Watch(watchCtx, watchCh)
			if err != nil {
				return err
			}

			// Enter the election
			ctx, cancel := getTimeoutContext(cmd)
			_, err = election.Enter(ctx)
			cancel()
			if err != nil {
				return err
			}

			// Once we've successfully entered the election, wait for watch events
			for event := range watchCh {
				bytes, err := yaml.Marshal(event)
				if err != nil {
					cmd.Println(err)
				} else {
					cmd.Println(string(bytes))
				}
			}
			return nil
		},
	}
	return cmd
}

func newElectionWatchCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "watch",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			e, err := getElection(cmd, name, "")
			if err != nil {
				return err
			}

			watchCh := make(chan *election.Event)
			ctx, cancel := getCancelContext(cmd)
			defer cancel()
			if err := e.Watch(ctx, watchCh); err != nil {
				return err
			}

			for event := range watchCh {
				bytes, err := yaml.Marshal(event)
				if err != nil {
					cmd.Println(err)
				} else {
					cmd.Println(string(bytes))
				}
			}
			return nil
		},
	}
	return cmd
}
