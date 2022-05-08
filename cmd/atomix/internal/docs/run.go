// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package docs

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
)

func run(cmd *cobra.Command, args []string) error {
	outputDir, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	if enabled, err := cmd.Flags().GetBool("markdown"); err != nil {
		return err
	} else if enabled {
		if err := doc.GenMarkdownTree(cmd.Root(), outputDir); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if enabled, err := cmd.Flags().GetBool("man"); err != nil {
		return err
	} else if enabled {
		header := doc.GenManHeader{
			Title: "atomix",
		}
		if err := doc.GenManTree(cmd.Root(), &header, outputDir); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if enabled, err := cmd.Flags().GetBool("yaml"); err != nil {
		return err
	} else if enabled {
		if err := doc.GenYamlTree(cmd.Root(), outputDir); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	return nil
}
