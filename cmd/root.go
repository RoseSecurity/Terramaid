// Copyright (c) RoseSecurity
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"

	tuiUtils "github.com/RoseSecurity/terramaid/internal/tui/utils"
	u "github.com/RoseSecurity/terramaid/pkg/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "terramaid",
	Short:         "A utility for generating Mermaid diagrams from Terraform configurations",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		fmt.Println()
		err = tuiUtils.PrintStyledText("TERRAMAID")
		if err != nil {
			u.LogErrorAndExit(err)
		}
		cmd.Help()
	},
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(docsCmd)
	rootCmd.AddCommand(versionCmd)

	// Add global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")

	// Disable auto-generated string from documentation so that documentation is cleanly built and updated
	rootCmd.DisableAutoGenTag = true
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		u.LogErrorAndExit(err)
	}
}
