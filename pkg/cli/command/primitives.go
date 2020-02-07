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
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

func newGetPrimitivesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "primitives [args]",
		Short: "List primitives in a database",
		Run:   runGetPrimitivesCommand,
	}
	addClientFlags(cmd)
	cmd.Flags().StringP("type", "t", "", "the type of primitives to list")
	cmd.Flags().Lookup("type").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_primitive_types"},
	}
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func runGetPrimitivesCommand(cmd *cobra.Command, _ []string) {
	typeName, _ := cmd.Flags().GetString("type")
	getPrimitives(cmd, typeName)
}

func newGetCountersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "counters [args]",
		Short: "List counters in a database",
		Run:   runGetCountersCommand,
	}
	addClientFlags(cmd)
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func runGetCountersCommand(cmd *cobra.Command, _ []string) {
	getPrimitives(cmd, "counter")
}

func newGetElectionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "elections [args]",
		Short: "List elections in a database",
		Run:   runGetElectionsCommand,
	}
	addClientFlags(cmd)
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func runGetElectionsCommand(cmd *cobra.Command, _ []string) {
	getPrimitives(cmd, "election")
}

func newGetIndexedMapsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "indexed-maps [args]",
		Short: "List indexed maps in a database",
		Run:   runGetIndexedMapsCommand,
	}
	addClientFlags(cmd)
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func runGetIndexedMapsCommand(cmd *cobra.Command, _ []string) {
	getPrimitives(cmd, "indexed-map")
}

func newGetLeaderLatchesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "leader-latches [args]",
		Short: "List leader latches in a database",
		Run:   runGetLeaderLatchesCommand,
	}
	addClientFlags(cmd)
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func runGetLeaderLatchesCommand(cmd *cobra.Command, _ []string) {
	getPrimitives(cmd, "leader-latch")
}

func newGetListsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lists [args]",
		Short: "List lists in a database",
		Run:   runGetListsCommand,
	}
	addClientFlags(cmd)
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func runGetListsCommand(cmd *cobra.Command, _ []string) {
	getPrimitives(cmd, "list")
}

func newGetLocksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locks [args]",
		Short: "List locks in a database",
		Run:   runGetLocksCommand,
	}
	addClientFlags(cmd)
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func runGetLocksCommand(cmd *cobra.Command, _ []string) {
	getPrimitives(cmd, "lock")
}

func newGetLogsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs [args]",
		Short: "List logs in a database",
		Run:   runGetLogsCommand,
	}
	addClientFlags(cmd)
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func runGetLogsCommand(cmd *cobra.Command, _ []string) {
	getPrimitives(cmd, "log")
}

func newGetMapsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "maps [args]",
		Short: "List maps in a database",
		Run:   runGetMapsCommand,
	}
	addClientFlags(cmd)
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func runGetMapsCommand(cmd *cobra.Command, _ []string) {
	getPrimitives(cmd, "map")
}

func newGetSetsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sets [args]",
		Short: "List sets in a database",
		Run:   runGetSetsCommand,
	}
	addClientFlags(cmd)
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func runGetSetsCommand(cmd *cobra.Command, _ []string) {
	getPrimitives(cmd, "set")
}

func newGetValuesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "values [args]",
		Short: "List values in a database",
		Run:   runGetValuesCommand,
	}
	addClientFlags(cmd)
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func runGetValuesCommand(cmd *cobra.Command, _ []string) {
	getPrimitives(cmd, "value")
}

func getPrimitives(cmd *cobra.Command, typeName string) {
	database := getDatabase(cmd)
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	var primitives []primitive.Metadata
	var err error
	if typeName == "" {
		primitives, err = database.GetPrimitives(ctx)
	} else {
		primitives, err = database.GetPrimitives(ctx, primitive.WithPrimitiveType(primitive.Type(strcase.ToCamel(typeName))))
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
		fmt.Fprintln(writer, "NAME\tSCOPE\tTYPE")
	}
	for _, primitive := range primitives {
		fmt.Fprintln(writer, fmt.Sprintf("%s\t%s\t%s", primitive.Name.Name, primitive.Name.Scope, strcase.ToKebab(string(primitive.Type))))
	}
	writer.Flush()
}
