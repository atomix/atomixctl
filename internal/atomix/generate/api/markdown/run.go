// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package markdown

import (
	"fmt"
	"github.com/atomix/cli/internal/exec"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func run(cmd *cobra.Command, args []string) error {
	ctxDir, err := os.Getwd()
	if err != nil {
		return err
	}

	inputDir, err := cmd.Flags().GetString("input")
	if err != nil {
		return err
	}

	pattern, err := cmd.Flags().GetString("pattern")
	if err != nil {
		return err
	}

	outputDir, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	ctxInputDir := filepath.Join(ctxDir, inputDir)
	ctxOutputDir := filepath.Join(ctxDir, outputDir)
	err = doublestar.GlobWalk(os.DirFS(ctxInputDir), pattern, func(path string, info fs.DirEntry) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(info.Name()) != ".proto" {
			return nil
		}

		var pathArgs []string
		pathArgs = append(pathArgs, ctxDir)
		pathArgs = append(pathArgs, ctxInputDir)
		pathArgs = append(pathArgs, filepath.Join(os.Getenv("GOPATH"), "src/github.com/gogo/protobuf"))

		protoPath := strings.Join(pathArgs, ":")

		var args []string
		args = append(args, "-I", protoPath)

		docDir := filepath.Dir(filepath.Join(ctxOutputDir, path))
		if err := os.MkdirAll(docDir, 0755); err != nil {
			return err
		}
		args = append(args, fmt.Sprintf("--doc_out=%s", docDir))

		var optArgs []string
		optArgs = append(optArgs, "markdown")
		optArgs = append(optArgs, fmt.Sprintf("%s.md", info.Name()[:len(info.Name())-len(filepath.Ext(info.Name()))]))
		opt := strings.Join(optArgs, ",")

		args = append(args, fmt.Sprintf("--doc_opt=%s", opt))
		args = append(args, path)
		return exec.Run("protoc", exec.WithEnv(os.Environ()...), exec.WithDir(ctxDir), exec.WithArgs(args...))
	})
	if err != nil {
		return err
	}
	return nil
}
