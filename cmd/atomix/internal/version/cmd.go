// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package version

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/atomix/cli/pkg/version"
	"github.com/spf13/cobra"
	"strings"
)

type Info struct {
	Version string `json:"version" yaml:"version"`
	Commit  string `json:"commit" yaml:"commit"`
}

func (i Info) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("Version: %s", i.Version))
	lines = append(lines, fmt.Sprintf("Commit: %s", i.Commit))
	return strings.Join(lines, "\n")
}

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "version",
		RunE: runCommand,
	}
	cmd.Flags().StringP("output", "o", "", "the output format")
	return cmd
}

func runCommand(cmd *cobra.Command, args []string) error {
	format, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	info := Info{
		Version: version.Version(),
		Commit:  version.Commit(),
	}

	var bytes []byte
	switch format {
	case "json":
		bytes, err = json.Marshal(&info)
		if err != nil {
			return err
		}
	case "yaml":
		bytes, err = json.Marshal(&info)
		if err != nil {
			return err
		}
	case "":
		bytes = []byte(info.String())
	default:
		return errors.New(fmt.Sprintf("unknown output format '%s'", format))
	}

	_, err = fmt.Fprintln(cmd.OutOrStdout(), bytes)
	return err
}
