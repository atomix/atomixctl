// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"github.com/atomix/cli/pkg/version"
	"os"
	"os/exec"
)

func InDocker() error {
	var image string
	if version.IsRelease() {
		image = fmt.Sprintf("atomix/cli:%s", version.Version())
	} else {
		image = "atomix/cli:latest"
	}

	cmd := &exec.Cmd{
		Path: "docker",
		Args: append([]string{
			"run",
			"-i",
			"-v",
			"`pwd`:/atomix",
			image,
		}, os.Args[1:]...),
		Env:    os.Environ(),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	return cmd.Run()
}
