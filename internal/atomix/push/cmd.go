// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package push

import (
	"github.com/atomix/cli/internal/atomix/push/driver"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "push",
	}
	cmd.AddCommand(driver.GetCommand())
	return cmd
}
