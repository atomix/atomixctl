// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package api

type GeneratorConfig struct {
	Templates []TemplateConfig `yaml:"templates,omitempty"`
}

type TemplateConfig struct {
	Name   string       `yaml:"name"`
	Path   string       `yaml:"path"`
	Input  InputConfig  `yaml:"input"`
	Output OutputConfig `yaml:"output"`
}

type InputConfig struct {
	FileNamePattern string `yaml:"file_name_pattern"`
}

type OutputMode string

const (
	OutputOverwrite OutputMode = "overwrite"
	OutputAppend    OutputMode = "append"
)

type OutputConfig struct {
	FileNameTemplate string     `yaml:"file_name_template"`
	Mode             OutputMode `yaml:"mode"`
}
