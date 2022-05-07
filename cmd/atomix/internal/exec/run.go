// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"os"
	"os/exec"
)

func Run(command string, args ...string) error {
	cmd := &exec.Cmd{
		Path:   command,
		Args:   args,
		Env:    os.Environ(),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	return cmd.Run()
}
