// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"github.com/atomix/cli/internal/atomix"
	"github.com/spf13/cobra/doc"
	"os"
)

func main() {
	cmd := atomix.GetCommand()
	err := doc.GenMarkdownTree(cmd, "docs/atomix")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
