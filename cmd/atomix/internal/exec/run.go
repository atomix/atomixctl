// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"os"
	"os/exec"
	"strings"
)

func Run(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	println(strings.Join(cmd.Args, " "))
	return cmd.Run()
}
