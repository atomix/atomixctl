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
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

func newPrimitivesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "primitives [args]",
		Short: "List primitives in a partition group",
		Run:   runPrimitivesCommand,
	}
	addClientFlags(cmd)
	cmd.Flags().StringP("type", "t", "", "the type of primitives to list")
	cmd.Flags().Lookup("type").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_primitive_types"},
	}
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func runPrimitivesCommand(cmd *cobra.Command, _ []string) {
	database := getDatabase(cmd)
	t, _ := cmd.Flags().GetString("type")
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	var primitives []primitive.Metadata
	var err error
	if t == "" {
		primitives, err = database.GetPrimitives(ctx)
	} else {
		primitives, err = database.GetPrimitives(ctx, primitive.WithPrimitiveType(primitive.Type(t)))
	}

	if err != nil {
		ExitWithError(ExitError, err)
	}

	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	printPrimitives(primitives, !noHeaders)
}

func printPrimitives(primitives []primitive.Metadata, includeHeaders bool) {
	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 0, 3, ' ', tabwriter.FilterHTML)
	if includeHeaders {
		fmt.Fprintln(writer, "NAME\tAPP\tTYPE")
	}
	for _, primitive := range primitives {
		fmt.Fprintln(writer, fmt.Sprintf("%s\t%s\t%s", primitive.Name.Name, primitive.Name.Namespace, primitive.Type))
	}
	writer.Flush()
}
