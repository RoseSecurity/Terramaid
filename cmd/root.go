// Copyright RoseSecurity 2024, 2026
// SPDX-License-Identifier: Apache-2.0

// Package cmd provides the CLI utility and commands for Terramaid.
package cmd

import (
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
		_, _ = cmd.OutOrStdout().Write([]byte("\n"))
		err = tuiUtils.PrintStyledText("TERRAMAID")
		if err != nil {
			u.LogErrorAndExit(err)
		}
		_ = cmd.Help()
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
