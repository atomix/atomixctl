// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"github.com/atomix/runtime-api/pkg/runtime/driver"
	"github.com/spf13/cobra"
	"os"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "driver",
		Args: cobra.MaximumNArgs(1),
		RunE: runCommand,
	}
	cmd.Flags().StringP("name", "n", "", "the driver name")
	cmd.Flags().StringP("version", "v", "", "the driver version")
	cmd.Flags().StringP("api-version", "v", "", "the runtime API version for which to generate the driver")
	cmd.Flags().StringP("registry-host", "r", "atomix-controller.kube-system", "the driver registry host")
	cmd.Flags().IntP("registry-port", "p", 5679, "the driver registry port")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("version")
	_ = cmd.MarkFlagRequired("api-version")
	return cmd
}

func runCommand(cmd *cobra.Command, args []string) error {
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}
	version, err := cmd.Flags().GetString("version")
	if err != nil {
		return err
	}
	apiVersion, err := cmd.Flags().GetString("api-version")
	if err != nil {
		return err
	}
	regHost, err := cmd.Flags().GetString("registry-host")
	if err != nil {
		return err
	}
	regPort, err := cmd.Flags().GetInt("registry-port")
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	repo, err := driver.NewRepository(
		driver.WithRegistryHost(regHost),
		driver.WithRegistryPort(regPort),
		driver.WithPath(wd))
	if err != nil {
		return err
	}

	plugin := driver.PluginInfo{
		Name:       name,
		Version:    version,
		APIVersion: apiVersion,
	}
	_, err = repo.Pull(cmd.Context(), plugin)
	if err != nil {
		return err
	}
	return nil
}
