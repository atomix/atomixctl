// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"github.com/atomix/cli/cmd/atomix/internal/exec"
	"github.com/atomix/cli/internal/template"
	"github.com/atomix/cli/pkg/version"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func runNative(cmd *cobra.Command, args []string) error {
	var context Context

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

	context.Runtime.Version = version.SDKVersion()

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

	err = exec.Run("go", "mod", "tidy")
	if err != nil {
		return err
	}
	return nil
}

func generate(dir string, context Context) error {
	err := apply(".gitignore", filepath.Join(dir, ".gitignore"), context)
	if err != nil {
		return err
	}

	err = apply(".goreleaser.yaml", filepath.Join(dir, ".goreleaser.yaml"), context)
	if err != nil {
		return err
	}

	err = apply("Makefile", filepath.Join(dir, "Makefile"), context)
	if err != nil {
		return err
	}

	err = apply("go.mod", filepath.Join(dir, "go.mod"), context)
	if err != nil {
		return err
	}

	pluginDir := filepath.Join(dir, "driver", strcase.ToKebab(context.Driver.Name))
	err = os.MkdirAll(pluginDir, 0755)
	if err != nil {
		return err
	}

	err = apply("driver.go", filepath.Join(pluginDir, "driver.go"), context)
	if err != nil {
		return err
	}

	protoDir := filepath.Join(dir, "api", "atomix", "driver", strcase.ToKebab(context.Driver.Name), context.Driver.APIVersion)
	err = os.MkdirAll(protoDir, 0755)
	if err != nil {
		return err
	}

	err = apply("config.proto", filepath.Join(protoDir, "config.proto"), context)
	if err != nil {
		return err
	}
	return nil
}

type Context struct {
	Driver  DriverContext
	Module  ModuleContext
	Runtime RuntimeContext
	Repo    RepoContext
}

type DriverContext struct {
	Name       string
	APIVersion string
}

type ModuleContext struct {
	Path string
}

type RuntimeContext struct {
	Version string
}

type RepoContext struct {
	Owner string
	Name  string
}

func apply(templatePath string, outputPath string, context Context) error {
	tpl := template.New(filepath.Join("/etc/atomix/templates/", templatePath))
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	return tpl.Execute(file, context)
}

func runInDocker(cmd *cobra.Command, args []string) error {
	return exec.InDocker()
}
