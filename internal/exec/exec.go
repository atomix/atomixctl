// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"os"
	"os/exec"
	"strings"
)

type Option func(options *Options)

type Options struct {
	Dir  string
	Env  []string
	Name string
	Args []string
}

func WithDir(dir string) Option {
	return func(options *Options) {
		options.Dir = dir
	}
}

func WithEnv(env ...string) Option {
	return func(options *Options) {
		options.Env = append(options.Env, env...)
	}
}

func WithArgs(args ...string) Option {
	return func(options *Options) {
		options.Args = args
	}
}

func Run(name string, opts ...Option) error {
	var options Options
	for _, opt := range opts {
		opt(&options)
	}

	println(strings.Join(append([]string{name}, options.Args...), " "))
	cmd := exec.Command(name, options.Args...)
	cmd.Dir = options.Dir
	cmd.Env = append(cmd.Env, options.Env...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
