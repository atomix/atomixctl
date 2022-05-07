// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"fmt"
	"github.com/atomix/cli/cmd/atomix/internal/exec"
	"github.com/bmatcuk/doublestar/v4"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func Generate(config Config, args map[string]string) error {
	return NewGenerator(config).Generate(args)
}

func NewGenerator(config Config) *Generator {
	return &Generator{
		Config: config,
	}
}

type Generator struct {
	Config Config
}

func (g *Generator) Generate(args map[string]string) error {
	for _, pattern := range g.Config.Proto.Patterns {
		if err := NewGlob(g, pattern).Generate(args); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) gen(file string, spec Spec) error {
	var path []string
	path = append(path, ".")
	path = append(path, g.Config.Proto.Path)
	path = append(path, filepath.Join(os.Getenv("GOPATH"), "src/github.com/gogo/protobuf"))

	var args []string
	args = append(args, "-I", strings.Join(path, ":"))
	args = append(args, "--template_out=%s", spec.String())
	args = append(args, file)

	return exec.Run("protoc", args...)
}

func NewGlob(generator *Generator, pattern string) *GlobGenerator {
	return &GlobGenerator{
		Generator: generator,
		Pattern:   pattern,
	}
}

type GlobGenerator struct {
	*Generator
	Pattern string
}

func (g *GlobGenerator) Generate(args map[string]string) error {
	return doublestar.GlobWalk(os.DirFS(g.Config.Proto.Path), g.Pattern, func(path string, info fs.DirEntry) error {
		if info.IsDir() {
			return nil
		}
		return NewFile(g, path).Generate(args)
	})
}

func NewFile(parent *GlobGenerator, file string) *FileGenerator {
	return &FileGenerator{
		GlobGenerator: parent,
		File:          file,
	}
}

type FileGenerator struct {
	*GlobGenerator
	File string
}

func (g *FileGenerator) Generate(args map[string]string) error {
	for _, template := range g.Config.Templates {
		if err := NewTemplate(g, template).Generate(args); err != nil {
			return err
		}
	}
	return nil
}

func NewTemplate(parent *FileGenerator, template TemplateConfig) *TemplateGenerator {
	return &TemplateGenerator{
		FileGenerator: parent,
		Template:      template,
	}
}

type TemplateGenerator struct {
	*FileGenerator
	Template TemplateConfig
}

func (g *TemplateGenerator) Generate(args map[string]string) error {
	if g.Template.Filter.Atoms != nil && g.Template.Filter.Components != nil {
		for _, atom := range g.Template.Filter.Atoms {
			for _, component := range g.Template.Filter.Components {
				err := g.gen(g.File, Spec{
					Template:  g.Template.Path,
					Output:    g.Template.Output,
					Atom:      atom,
					Component: component,
					Args:      args,
				})
				if err != nil {
					return err
				}
			}
		}
	} else if g.Template.Filter.Atoms != nil {
		for _, atom := range g.Template.Filter.Atoms {
			err := g.gen(g.File, Spec{
				Template: g.Template.Path,
				Output:   g.Template.Output,
				Atom:     atom,
				Args:     args,
			})
			if err != nil {
				return err
			}
		}
	} else if g.Template.Filter.Components != nil {
		for _, component := range g.Template.Filter.Components {
			err := g.gen(g.File, Spec{
				Template:  g.Template.Path,
				Output:    g.Template.Output,
				Component: component,
				Args:      args,
			})
			if err != nil {
				return err
			}
		}
	}
	return g.gen(g.File, Spec{
		Template: g.Template.Path,
		Output:   g.Template.Output,
		Args:     args,
	})
}

type Spec struct {
	Template  string
	Output    string
	Atom      string
	Component string
	Args      map[string]string
}

func (s Spec) String() string {
	var elems []string
	elems = append(elems, formatArg("template", s.Template))
	elems = append(elems, formatArg("output", s.Output))
	if s.Args != nil {
		for key, value := range s.Args {
			elems = append(elems, formatArg(key, value))
		}
	}
	if s.Atom != "" {
		elems = append(elems, formatArg("atom", s.Atom))
	}
	if s.Component != "" {
		elems = append(elems, formatArg("component", s.Component))
	}
	return strings.Join(elems, ",")
}

func formatArg(key, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
}
