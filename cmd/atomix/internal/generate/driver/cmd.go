// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"github.com/atomix/cli/internal/exec"
	"github.com/atomix/sdk/pkg/version"
	"github.com/spf13/cobra"
	"os"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "driver",
		Args: cobra.MaximumNArgs(1),
		RunE: runCommand,
	}
	cmd.Flags().StringP("name", "n", "", "the driver name")
	cmd.Flags().StringP("api-version", "v", "", "the driver API version")
	cmd.Flags().StringP("module-path", "p", "", "the driver module path")
	cmd.Flags().String("github-owner", "", "the GitHub user to which to publish release artifacts")
	cmd.Flags().String("github-repo", "", "the GitHub repo to which to publish release artifacts")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("module")
	_ = cmd.MarkFlagRequired("api-version")
	return cmd
}

func runCommand(cmd *cobra.Command, args []string) error {
	var context TemplateContext

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}
	context.Driver.Name = name

	apiVersion, err := cmd.Flags().GetString("api-version")
	if err != nil {
		return err
	}
	context.Driver.APIVersion = apiVersion

	path, err := cmd.Flags().GetString("module-path")
	if err != nil {
		return err
	}
	context.Module.Path = path

	context.Runtime.Version = version.Version()

	repoOwner, err := cmd.Flags().GetString("github-owner")
	if err != nil {
		return err
	}
	context.Repo.Owner = repoOwner

	repoName, err := cmd.Flags().GetString("github-repo")
	if err != nil {
		return err
	}
	context.Repo.Name = repoName

	var dir string
	if len(args) == 1 {
		dir = args[0]
	} else {
		dir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	err = generate(dir, context)
	if err != nil {
		return err
	}

	err = exec.Run("go", exec.WithEnv(os.Environ()...), exec.WithDir(dir), exec.WithArgs("mod", "tidy"))
	if err != nil {
		return err
	}
	return nil
}
