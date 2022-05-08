// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"
	"github.com/atomix/cli/cmd/atomix/internal/exec"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "config",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := Load()
			bytes, err := yaml.Marshal(config)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), string(bytes))
			return err
		},
	}
	cmd.AddCommand(getAddGeneratorCommand())
	cmd.AddCommand(getInitCommand())
	return cmd
}

func getAddGeneratorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "add-generator",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			image, err := cmd.Flags().GetString("image")
			if err != nil {
				return err
			} else if image == "" {
				image = fmt.Sprintf("atomix/codegen:%s-latest", name)
			}

			// Check if the Docker image exists locally and pull if it not
			if err := exec.Run("docker", "images", "-q", image); err != nil {
				if err := exec.Run("docker", "pull", image); err != nil {
					return err
				}
			}

			// Load the configuration and add the generator configuration
			config := Load()
			config.Generators = append(config.Generators, GeneratorConfig{
				Name:  name,
				Image: image,
			})
			return Store(config)
		},
	}
	cmd.Flags().StringP("image", "i", "", "the generator image")
	return cmd
}

func getInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "init",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := Load()
			return Store(config)
		},
	}
	return cmd
}
