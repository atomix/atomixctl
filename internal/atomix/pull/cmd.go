// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package pull

import (
	"github.com/atomix/cli/internal/atomix/pull/driver"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "pull",
	}
	cmd.AddCommand(driver.GetCommand())
	return cmd
}
