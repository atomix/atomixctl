// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package _go

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

	outputPath, err := cmd.Flags().GetString("path")
	if err != nil {
		return err
	}

	importOverrides := make(map[string]string)
	importOverrides["google/protobuf/any.proto"] = "github.com/gogo/protobuf/types"
	importOverrides["google/protobuf/timestamp.proto"] = "github.com/gogo/protobuf/types"
	importOverrides["google/protobuf/duration.proto"] = "github.com/gogo/protobuf/types"

	ctxInputDir := filepath.Join(ctxDir, inputDir)
	err = doublestar.GlobWalk(os.DirFS(ctxInputDir), pattern, func(path string, info fs.DirEntry) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(info.Name()) != ".proto" {
			return nil
		}
		importOverrides[path] = filepath.Join(outputPath, filepath.Dir(path))
		return nil
	})

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

		var specArgs []string
		var overrideArgs []string
		for protoPath, goPath := range importOverrides {
			overrideArgs = append(overrideArgs, fmt.Sprintf("M%s=%s", protoPath, goPath))
		}
		overrides := strings.Join(overrideArgs, ",")
		specArgs = append(specArgs, overrides)

		importPath := filepath.Join(outputPath, filepath.Dir(path))
		specArgs = append(specArgs, fmt.Sprintf("import_path=%s", importPath))
		specArgs = append(specArgs, fmt.Sprintf("plugins=grpc:%s", ctxOutputDir))
		spec := strings.Join(specArgs, ",")

		args = append(args, fmt.Sprintf("--gogofaster_out=%s", spec))
		args = append(args, path)
		return exec.Run("protoc", exec.WithEnv(os.Environ()...), exec.WithDir(ctxDir), exec.WithArgs(args...))
	})
	if err != nil {
		return err
	}
	return nil
}
