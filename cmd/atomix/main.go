// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/atomix/cli/cmd/atomix/internal"
	"os"
)

func main() {
	cmd := internal.GetCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
