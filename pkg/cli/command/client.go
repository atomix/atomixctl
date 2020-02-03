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
	"github.com/atomix/go-client/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"time"
)

const (
	nameSep = "."
)

func addClientFlags(cmd *cobra.Command) {
	viper.SetDefault("database", "")
	cmd.PersistentFlags().StringP("database", "d", viper.GetString("database"), fmt.Sprintf("the database name (default %s)", viper.GetString("database")))
	cmd.PersistentFlags().Duration("timeout", 15*time.Second, "the operation timeout")
	viper.BindPFlag("database", cmd.PersistentFlags().Lookup("database"))
	cmd.PersistentFlags().Lookup("database").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_databases"},
	}
}

func newTimeoutContext(cmd *cobra.Command) context.Context {
	timeout, _ := cmd.Flags().GetDuration("timeout")
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	return ctx
}

func newClientFromEnv() *client.Client {
	c, err := client.NewClient(
		getClientController(),
		client.WithNamespace(getClientNamespace()),
		client.WithApplication(getClientApp()))
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return c
}

func newDatabaseFromEnv(cmd *cobra.Command) *client.Database {
	c := newClientFromEnv()
	d, err := c.GetDatabase(newTimeoutContext(cmd), getClientDatabase())
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return d
}

func newClientFromDatabase(name string) *client.Client {
	c, err := client.NewClient(
		getClientController(),
		client.WithNamespace(getDatabaseNamespace(name)),
		client.WithApplication(getClientApp()))
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return c
}

func newClientFromName(name string) *client.Client {
	c, err := client.NewClient(getClientController(), client.WithNamespace(getClientNamespace()), client.WithApplication(getPrimitiveApp(name)))
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return c
}

func newDatabaseFromName(cmd *cobra.Command, name string) *client.Database {
	c := newClientFromName(name)
	database, err := c.GetDatabase(newTimeoutContext(cmd), getClientDatabase())
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return database
}

func splitName(name string) []string {
	return strings.Split(name, nameSep)
}

func getDatabaseNamespace(name string) string {
	nameParts := splitName(name)
	if len(nameParts) == 2 {
		return nameParts[0]
	}
	return getClientNamespace()
}

func getDatabaseName(name string) string {
	nameParts := splitName(name)
	return nameParts[len(nameParts)-1]
}

func setClientController(controller string) error {
	return setConfig("controller", controller)
}

func getClientController() string {
	return getConfig("controller")
}

func getClientNamespace() string {
	return getConfig("namespace")
}

func setClientDatabase(database string) error {
	return setConfig("database", database)
}

func getClientDatabase() string {
	return getConfig("database")
}

func getClientApp() string {
	return getConfig("app")
}

func getPrimitiveApp(name string) string {
	nameParts := splitName(name)
	if len(nameParts) == 2 {
		return nameParts[0]
	}
	return getClientApp()
}

func getPrimitiveName(name string) string {
	nameParts := splitName(name)
	return nameParts[len(nameParts)-1]
}
