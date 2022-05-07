// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	_ "embed"
	"github.com/atomix/cli/internal/template"
	"github.com/iancoleman/strcase"
	"os"
	"path/filepath"
)

var (
	//go:embed templates/.gitignore.tpl
	gitIgnoreTemplate string

	//go:embed templates/.goreleaser.yaml.tpl
	goReleaserTemplate string

	//go:embed templates/Makefile.tpl
	makefileTemplate string

	//go:embed templates/go.mod.tpl
	goModTemplate string

	//go:embed templates/driver.go.tpl
	driverTemplate string

	//go:embed templates/config.proto.tpl
	configProtoTemplate string
)

type TemplateContext struct {
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

func generate(dir string, context TemplateContext) error {
	err := apply(".gitignore", gitIgnoreTemplate, filepath.Join(dir, ".gitignore"), context)
	if err != nil {
		return err
	}

	err = apply(".goreleaser.yaml", goReleaserTemplate, filepath.Join(dir, ".goreleaser.yaml"), context)
	if err != nil {
		return err
	}

	err = apply("Makefile", makefileTemplate, filepath.Join(dir, "Makefile"), context)
	if err != nil {
		return err
	}

	err = apply("go.mod", goModTemplate, filepath.Join(dir, "go.mod"), context)
	if err != nil {
		return err
	}

	pluginDir := filepath.Join(dir, "driver", strcase.ToKebab(context.Driver.Name))
	err = os.MkdirAll(pluginDir, 0755)
	if err != nil {
		return err
	}

	err = apply("driver.go", driverTemplate, filepath.Join(pluginDir, "driver.go"), context)
	if err != nil {
		return err
	}

	protoDir := filepath.Join(dir, "api", "atomix", "driver", strcase.ToKebab(context.Driver.Name), context.Driver.APIVersion)
	err = os.MkdirAll(protoDir, 0755)
	if err != nil {
		return err
	}

	err = apply("config.proto", configProtoTemplate, filepath.Join(protoDir, "config.proto"), context)
	if err != nil {
		return err
	}
	return nil
}

func apply(name string, text string, path string, args TemplateContext) error {
	tpl := template.NewTemplate(name, text)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	return tpl.Execute(file, args)
}
