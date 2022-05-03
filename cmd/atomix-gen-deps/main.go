// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"github.com/atomix/cli/internal/exec"
	"github.com/rogpeppe/go-internal/modfile"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	cmd := getCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func getCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "plugin",
		Args: cobra.MaximumNArgs(1),
		RunE: runCommand,
	}
	cmd.Flags().BoolP("check", "c", false, "check module compatibility only")
	cmd.Flags().StringP("version", "v", "", "the target runtime API version")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}

func runCommand(cmd *cobra.Command, args []string) error {
	var path string
	if len(args) == 1 {
		path = args[0]
	} else {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		path = dir
	}

	target, err := cmd.Flags().GetString("version")
	if err != nil {
		return err
	}

	checkOnly, err := cmd.Flags().GetBool("read-only")
	if err != nil {
		return err
	}

	if checkOnly {
		fmt.Fprintf(cmd.OutOrStdout(), "Checking plugin module constraints against target API version %s\n", target)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "Updating plugin module constraints for target API version %s\n", target)
	}

	err = exec.Run("go", exec.WithArgs("mod", "tidy"))
	if err != nil {
		return err
	}

	rootPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	srcModFile := filepath.Join(rootPath, "go.mod")
	srcModBytes, err := ioutil.ReadFile(srcModFile)
	if err != nil {
		return err
	}

	srcMod, err := modfile.Parse("go.mod", srcModBytes, nil)
	if err != nil {
		return err
	}

	tmpDir, err := ioutil.TempDir(rootPath, "atomix")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	tgtModFile := filepath.Join(tmpDir, "go.mod")
	tgtModURL := fmt.Sprintf("https://raw.githubusercontent.com/atomix/runtime-api/%s/go.mod", target)
	err = exec.Run("wget", exec.WithArgs("-LO", tgtModFile, tgtModURL))
	if err != nil {
		return err
	}

	tgtModBytes, err := ioutil.ReadFile(tgtModFile)
	if err != nil {
		return err
	}

	tgtMod, err := modfile.Parse(tgtModFile, tgtModBytes, nil)
	if err != nil {
		return err
	}

	tgtReqs := make(map[string]string)
	for _, tgtReq := range tgtMod.Require {
		tgtReqs[tgtReq.Mod.Path] = tgtReq.Mod.Version
	}

	for _, srcReq := range srcMod.Require {
		if tgtReqVersion, ok := tgtReqs[srcReq.Mod.Path]; ok {
			fmt.Fprintf(cmd.OutOrStdout(), "Evaluating common dependency %s\n", srcReq.Mod.Path)
			if srcReq.Mod.Version != tgtReqVersion {
				if checkOnly {
					fmt.Fprintf(cmd.OutOrStderr(), "Detected incompatible dependency %s: %s <> %s\n", srcReq.Mod.Path, srcReq.Mod.Version, tgtReqVersion)
					os.Exit(1)
				} else {
					fmt.Fprintf(cmd.OutOrStderr(), "Updating dependency %s: %s => %s\n", srcReq.Mod.Path, srcReq.Mod.Version, tgtReqVersion)
					srcReq.Mod.Version = tgtReqVersion
				}
			}
		}
	}

	srcModBytes, err = srcMod.Format()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(srcModFile, srcModBytes, 0755)
	if err != nil {
		return err
	}
	return nil
}
