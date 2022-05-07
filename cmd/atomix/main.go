// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/atomix/cli/cmd/atomix/internal"
	"github.com/atomix/cli/cmd/atomix/internal/env"
	"os"
)

var environment env.Environment

func main() {
	cmd := internal.GetCommand(environment)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
