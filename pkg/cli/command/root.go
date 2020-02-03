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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                    "atomix",
		Short:                  "Atomix command line client",
		BashCompletionFunction: bashCompletion,
	}

	viper.SetDefault("controller", ":5679")
	viper.SetDefault("namespace", "default")
	viper.SetDefault("app", "default")

	cmd.PersistentFlags().String("controller", viper.GetString("controller"), "the controller address")
	cmd.PersistentFlags().String("namespace", viper.GetString("namespace"), "the partition group namespace")
	cmd.PersistentFlags().StringP("app", "a", viper.GetString("app"), "the application name")
	cmd.PersistentFlags().String("config", "", "config file (default: $HOME/.atomix/config.yaml)")

	viper.BindPFlag("controller", cmd.PersistentFlags().Lookup("controller"))
	viper.BindPFlag("namespace", cmd.PersistentFlags().Lookup("namespace"))
	viper.BindPFlag("app", cmd.PersistentFlags().Lookup("app"))

	cmd.AddCommand(newCompletionCommand())
	cmd.AddCommand(newConfigCommand())
	cmd.AddCommand(newDatabasesCommand())
	cmd.AddCommand(newPrimitivesCommand())
	cmd.AddCommand(newCounterCommand())
	cmd.AddCommand(newListCommand())
	cmd.AddCommand(newElectionCommand())
	cmd.AddCommand(newLockCommand())
	cmd.AddCommand(newMapCommand())
	cmd.AddCommand(newSetCommand())
	return cmd
}
