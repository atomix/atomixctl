// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Path() (string, error) {
	dir, err := homedir.Dir()
	if err != nil {
		dir, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}
	path := filepath.Join(dir, ".atomix", "config.yaml")
	return path, nil
}

func Load() Config {
	var config Config
	path, err := Path()
	if err != nil {
		return config
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return config
	}
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return config
	}
	return config
}

func Store(config Config) error {
	bytes, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	path, err := Path()
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bytes, 0755)
}

type Config struct {
	Generators []GeneratorConfig `yaml:"generators"`
}

type GeneratorConfig struct {
	Name  string `yaml:"name"`
	Image string `yaml:"image"`
}
