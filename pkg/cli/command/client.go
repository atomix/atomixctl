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
	"github.com/atomix/go-client/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func addClientFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("scope", "s", viper.GetString("scope"), "the application scope")
	cmd.PersistentFlags().StringP("database", "d", viper.GetString("database"), "the database name")
	cmd.PersistentFlags().String("config", "", "config file (default: $HOME/.atomix/config.yaml)")
	cmd.PersistentFlags().Duration("timeout", 15*time.Second, "the operation timeout")

	cmd.PersistentFlags().Lookup("scope").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_scopes"},
	}
	cmd.PersistentFlags().Lookup("database").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_databases"},
	}
}

func getCancelContext(cmd *cobra.Command) (context.Context, context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-ctx.Done():
			signal.Stop(sigCh)
		case <-sigCh:
			cancel()
		}
	}()
	return ctx, cancel
}

func getTimeoutContext(cmd *cobra.Command) (context.Context, context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	timeout, _ := cmd.Flags().GetDuration("timeout")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	go func() {
		select {
		case <-ctx.Done():
			signal.Stop(sigCh)
		case <-sigCh:
			cancel()
		}
	}()
	return ctx, cancel
}

func getClientController() string {
	return viper.GetString("controller")
}

func getClientNamespace() string {
	return viper.GetString("namespace")
}

func getClientDatabase(cmd *cobra.Command) string {
	database, _ := cmd.Flags().GetString("database")
	return database
}

func getClientScope(cmd *cobra.Command) string {
	app, _ := cmd.Flags().GetString("scope")
	return app
}

func getClient(cmd *cobra.Command) (*client.Client, error) {
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	return client.NewWithContext(
		ctx,
		getClientController(),
		client.WithNamespace(getClientNamespace()),
		client.WithScope(getClientScope(cmd)))
}

func getDatabase(cmd *cobra.Command) (*client.Database, error) {
	client, err := getClient(cmd)
	if err != nil {
		return nil, err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	return client.GetDatabase(ctx, getClientDatabase(cmd))
}
