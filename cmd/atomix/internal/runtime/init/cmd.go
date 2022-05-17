// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package init

import (
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "init",
		RunE: run,
	}
	cmd.Flags().Bool("headless", false, "run the runtime in headless mode")
	return cmd
}
