package command

import (
	"fmt"
	"github.com/atomix/atomix-api/proto/atomix/primitive"
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
	group := newGroupFromEnv(cmd)
	t, _ := cmd.Flags().GetString("type")
	var primitives []*primitive.PrimitiveInfo
	var err error
	if t == "" {
		primitives, err = group.GetPrimitives(newTimeoutContext(cmd))
	} else {
		primitives, err = group.GetPrimitives(newTimeoutContext(cmd), t)
	}

	if err != nil {
		ExitWithError(ExitError, err)
	}

	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	printPrimitives(primitives, !noHeaders)
}

func printPrimitives(primitives []*primitive.PrimitiveInfo, includeHeaders bool) {
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
