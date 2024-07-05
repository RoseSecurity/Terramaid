package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsCmd = &cobra.Command{
	Use:          "docs",
	Short:        "Generate documentation for the CLI",
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		err := doc.GenMarkdownTree(RootCmd, "./docs")
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(docsCmd)
}
