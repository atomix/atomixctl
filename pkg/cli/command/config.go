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
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
)

var (
	configFile = ""
)

func init() {
	initConfig()
}

func newConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config {set,get,delete} [args]",
		Short: "Read and update CLI configuration options",
	}
	cmd.AddCommand(newConfigGetCommand())
	cmd.AddCommand(newConfigSetCommand())
	cmd.AddCommand(newConfigDeleteCommand())
	return cmd
}

func newConfigGetCommand() *cobra.Command {
	validArgs := []string{
		"controller",
		"namespace",
		"database",
		"scope",
	}
	return &cobra.Command{
		Use:       "get <key>",
		Args:      cobra.ExactArgs(1),
		ValidArgs: validArgs,
		Run:       runConfigGetCommand,
	}
}

func runConfigGetCommand(cmd *cobra.Command, args []string) {
	value := viper.Get(args[0])
	ExitWithOutput("%v", value)
}

func newConfigSetCommand() *cobra.Command {
	validArgs := []string{
		"controller",
		"namespace",
		"database",
		"scope",
	}
	return &cobra.Command{
		Use:       "set <key> <value>",
		Args:      cobra.ExactArgs(2),
		ValidArgs: validArgs,
		Run:       runConfigSetCommand,
	}
}

func runConfigSetCommand(cmd *cobra.Command, args []string) {
	viper.Set(args[0], args[1])
	if err := viper.WriteConfig(); err != nil {
		ExitWithError(ExitError, err)
	} else {
		value := viper.Get(args[0])
		ExitWithOutput("%v", value)
	}
}

func newConfigDeleteCommand() *cobra.Command {
	validArgs := []string{
		"controller",
		"namespace",
		"database",
		"scope",
	}
	return &cobra.Command{
		Use:       "delete <key>",
		Args:      cobra.ExactArgs(1),
		ValidArgs: validArgs,
		Run:       runConfigDeleteCommand,
	}
}

func runConfigDeleteCommand(cmd *cobra.Command, args []string) {
	viper.Set(args[0], nil)
	if err := viper.WriteConfig(); err != nil {
		ExitWithError(ExitError, err)
	} else {
		value := viper.Get(args[0])
		ExitWithOutput("%v", value)
	}
}

func setConfig(key string, value string) error {
	viper.Set(key, value)
	return viper.WriteConfig()
}

func getConfig(key string) string {
	return viper.GetString(key)
}

func getConfigFile(name string) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".atomix", name), nil
}

func initConfig() {
	viper.SetDefault("controller", "atomix-controller.kube-system.svc.cluster.local:5679")
	viper.SetDefault("namespace", "default")
	viper.SetDefault("scope", "default")
	viper.SetDefault("database", "")

	home, err := homedir.Dir()
	if err != nil {
		ExitWithError(ExitError, err)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(home + "/.atomix")
	viper.AddConfigPath("/etc/atomix")
	viper.AddConfigPath(".")

	viper.ReadInConfig()
}
