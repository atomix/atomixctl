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
	"github.com/atomix/go-client/pkg/client/counter"
	"github.com/atomix/go-client/pkg/client/election"
	"github.com/atomix/go-client/pkg/client/indexedmap"
	"github.com/atomix/go-client/pkg/client/leader"
	"github.com/atomix/go-client/pkg/client/list"
	"github.com/atomix/go-client/pkg/client/lock"
	"github.com/atomix/go-client/pkg/client/log"
	"github.com/atomix/go-client/pkg/client/map"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/atomix/go-client/pkg/client/set"
	"github.com/atomix/go-client/pkg/client/value"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"io"
	"text/tabwriter"
)

func newGetPrimitivesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "primitives [args]",
		Short: "List primitives in a database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getPrimitives(cmd, "")
		},
	}
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func newGetCountersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "counters [args]",
		Short: "List counters in a database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getPrimitives(cmd, counter.Type)
		},
	}
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func newGetElectionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "elections [args]",
		Short: "List elections in a database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getPrimitives(cmd, election.Type)
		},
	}
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func newGetIndexedMapsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "indexed-maps [args]",
		Short: "List indexed maps in a database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getPrimitives(cmd, indexedmap.Type)
		},
	}
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func newGetLeaderLatchesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "leader-latches [args]",
		Short: "List leader latches in a database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getPrimitives(cmd, leader.Type)
		},
	}
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func newGetListsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lists [args]",
		Short: "List lists in a database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getPrimitives(cmd, list.Type)
		},
	}
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func newGetLocksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locks [args]",
		Short: "List locks in a database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getPrimitives(cmd, lock.Type)
		},
	}
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func newGetLogsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs [args]",
		Short: "List logs in a database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getPrimitives(cmd, log.Type)
		},
	}
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func newGetMapsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "maps [args]",
		Short: "List maps in a database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getPrimitives(cmd, _map.Type)
		},
	}
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func newGetSetsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sets [args]",
		Short: "List sets in a database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getPrimitives(cmd, set.Type)
		},
	}
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func newGetValuesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "values [args]",
		Short: "List values in a database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getPrimitives(cmd, value.Type)
		},
	}
	cmd.Flags().Bool("no-headers", false, "exclude headers from output")
	return cmd
}

func getPrimitives(cmd *cobra.Command, primitiveType primitive.Type) error {
	database, err := getDatabase(cmd)
	if err != nil {
		return err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	var primitives []primitive.Metadata
	if primitiveType == "" {
		primitives, err = database.GetPrimitives(ctx)
	} else {
		primitives, err = database.GetPrimitives(ctx, primitive.WithPrimitiveType(primitiveType))
	}

	if err != nil {
		return err
	}

	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	return printPrimitives(primitives, !noHeaders, cmd.OutOrStdout())
}

func printPrimitives(primitives []primitive.Metadata, includeHeaders bool, out io.Writer) error {
	writer := new(tabwriter.Writer)
	writer.Init(out, 0, 0, 3, ' ', tabwriter.FilterHTML)
	if includeHeaders {
		fmt.Fprintln(writer, "NAME\tSCOPE\tTYPE")
	}
	for _, primitive := range primitives {
		fmt.Fprintln(writer, fmt.Sprintf("%s\t%s\t%s", primitive.Name.Name, primitive.Name.Scope, strcase.ToKebab(string(primitive.Type))))
	}
	return writer.Flush()
}
