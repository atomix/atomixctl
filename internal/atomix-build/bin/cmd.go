// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package bin

import (
	"fmt"
	"github.com/atomix/cli/internal/exec"
	"github.com/spf13/cobra"
	"os"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "bin",
		Args: cobra.ExactArgs(1),
		RunE: runCommand,
	}
	cmd.Flags().StringP("output", "o", ".", "the output directory")
	cmd.Flags().StringP("version", "v", "", "the build version")
	cmd.Flags().StringP("commit", "c", "", "the commit hash")
	return cmd
}

func runCommand(cmd *cobra.Command, args []string) error {
	var goArgs []string
	goArgs = append(goArgs, "build")
	goArgs = append(goArgs, "-trimpath")
	goArgs = append(goArgs, "-gcflags=\"all=-N -l\"")

	goEnv := os.Environ()
	goEnv = append(goEnv, "CGO_ENABLED=1")
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}
	goArgs = append(goArgs, "-o", output)
	version, err := cmd.Flags().GetString("version")
	if err != nil {
		return err
	}
	commit, err := cmd.Flags().GetString("commit")
	if err != nil {
		return err
	}
	goArgs = append(goArgs, fmt.Sprintf("-ldflags=\"-s -w -X github.com/atomix/cli/internal/atomix/version.version=%s -X github.com/atomix/cli/internal/atomix/version.commit=%s\"", version, commit))
	goArgs = append(goArgs, args...)
	return exec.Run("go", exec.WithEnv(goEnv...), exec.WithArgs(goArgs...))
}
