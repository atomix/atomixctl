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
	viper.SetDefault("group", "")
	cmd.PersistentFlags().StringP("group", "g", viper.GetString("group"), fmt.Sprintf("the partition group name (default %s)", viper.GetString("group")))
	cmd.PersistentFlags().Duration("timeout", 15*time.Second, "the operation timeout")
	viper.BindPFlag("group", cmd.PersistentFlags().Lookup("group"))
	cmd.PersistentFlags().Lookup("group").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_groups"},
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

func newGroupFromEnv(cmd *cobra.Command) *client.PartitionGroup {
	c := newClientFromEnv()
	g, err := c.GetGroup(newTimeoutContext(cmd), getClientGroup())
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return g
}

func newClientFromGroup(name string) *client.Client {
	c, err := client.NewClient(
		getClientController(),
		client.WithNamespace(getGroupNamespace(name)),
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

func newGroupFromName(cmd *cobra.Command, name string) *client.PartitionGroup {
	c := newClientFromName(name)
	group, err := c.GetGroup(newTimeoutContext(cmd), getClientGroup())
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return group
}

func splitName(name string) []string {
	return strings.Split(name, nameSep)
}

func getGroupNamespace(name string) string {
	nameParts := splitName(name)
	if len(nameParts) == 2 {
		return nameParts[0]
	}
	return getClientNamespace()
}

func getGroupName(name string) string {
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

func setClientGroup(group string) error {
	return setConfig("group", group)
}

func getClientGroup() string {
	return getConfig("group")
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
