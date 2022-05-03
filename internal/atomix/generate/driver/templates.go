// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	_ "embed"
	"github.com/atomix/cli/internal/template"
	"os"
	"path/filepath"
)

var (
	//go:embed templates/go.mod.tpl
	goModTemplate string

	//go:embed templates/main.go.tpl
	mainTemplate string
)

type TemplateContext struct {
	Driver  DriverContext
	Module  ModuleContext
	Runtime RuntimeContext
}

type DriverContext struct {
	Name string
}

type ModuleContext struct {
	Path string
}

type RuntimeContext struct {
	Version string
}

func generate(dir string, context TemplateContext) error {
	err := apply("go.mod", goModTemplate, filepath.Join(dir, "go.mod"), context)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(dir, "cmd", context.Driver.Name), 0755)
	if err != nil {
		return err
	}

	err = apply("main.go", mainTemplate, filepath.Join(dir, "cmd", context.Driver.Name, "main.go"), context)
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
