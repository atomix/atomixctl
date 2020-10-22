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
	"github.com/atomix/go-client/pkg/client"
	"github.com/spf13/cobra"
	"io"
	"text/tabwriter"
	"time"
)

func newGetDatabasesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "databases",
		Short: "Get a list of databases",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd)
			if err != nil {
				return err
			}
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			databases, err := client.GetDatabases(ctx)
			if err != nil {
				return err
			}
			noHeaders, _ := cmd.Flags().GetBool("no-headers")
			printDatabases(databases, !noHeaders, cmd.OutOrStdout())
			return nil
		},
	}
	cmd.PersistentFlags().Duration("timeout", 15*time.Second, "the operation timeout")
	cmd.Flags().Bool("no-headers", false, "exclude headers from the output")
	return cmd
}

func printDatabases(databases []*client.Database, includeHeaders bool, out io.Writer) {
	writer := new(tabwriter.Writer)
	writer.Init(out, 0, 0, 3, ' ', tabwriter.FilterHTML)
	if includeHeaders {
		fmt.Fprintln(writer, "NAMESPACE\tNAME")
	}
	for _, database := range databases {
		fmt.Fprintln(writer, fmt.Sprintf("%s\t%s", database.Namespace, database.Name))
	}
	writer.Flush()
}
