// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/atomix/cli/internal/atomix-build"
	"os"
)

func main() {
	cmd := atomix_build.GetCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
