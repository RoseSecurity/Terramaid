package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// Generate documentation for Terramaid commands and output to docs directory
func docsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "docs",
		Short:        "Generate documentation for the CLI",
		SilenceUsage: true,
		Hidden:       true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := doc.GenMarkdownTree(cmd.Root(), "./docs")
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
