// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func ParseConfigFile(path string) (Config, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	return ParseConfig(bytes)
}

func ParseConfig(bytes []byte) (Config, error) {
	var config Config
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return config, err
	}
	return config, nil
}

type Config struct {
	Generator string           `yaml:"generator,omitempty"`
	Proto     ProtoConfig      `yaml:"proto,omitempty"`
	Templates []TemplateConfig `yaml:"templates,omitempty"`
}

type ProtoConfig struct {
	Path     string   `yaml:"path,omitempty"`
	Patterns []string `yaml:"patterns,omitempty"`
}

type ModuleConfig struct {
	Path    string `yaml:"path,omitempty"`
	Version string `yaml:"version,omitempty"`
}

type TemplateConfig struct {
	Name   string       `yaml:"name,omitempty"`
	Path   string       `yaml:"path,omitempty"`
	Data   string       `yaml:"data,omitempty"`
	Filter FilterConfig `yaml:"filter,omitempty"`
	Output string       `yaml:"output,omitempty"`
}

type FilterConfig struct {
	Atoms      []string `yaml:"atoms,omitempty"`
	Components []string `yaml:"components,omitempty"`
}
