// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"errors"
	"github.com/atomix/cli/internal/exec"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "driver",
		Args: cobra.MaximumNArgs(1),
		RunE: runCommand,
	}
	cmd.Flags().StringP("name", "n", "", "the driver name")
	cmd.Flags().StringP("module", "m", "", "the driver module path")
	cmd.Flags().StringP("api-version", "v", "", "the runtime API version for which to generate the driver")
	cmd.Flags().StringP("repo", "r", "", "the GitHub repository to which to publish release artifacts")
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

	path, err := cmd.Flags().GetString("module")
	if err != nil {
		return err
	}
	context.Module.Path = path

	apiVersion, err := cmd.Flags().GetString("api-version")
	if err != nil {
		return err
	}
	context.Runtime.Version = apiVersion

	repo, err := cmd.Flags().GetString("repo")
	if err != nil {
		return err
	}

	if repo != "" {
		repoParts := strings.Split(repo, "/")
		if len(repoParts) != 2 {
			return errors.New("invalid repository format")
		}
		repoOwner, repoName := repoParts[0], repoParts[1]
		context.Repo.Owner = repoOwner
		context.Repo.Name = repoName
	}

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
