package utils

import (
	"log"

	"github.com/RoseSecurity/terramaid/cmd"
	"github.com/spf13/cobra/doc"
)

// Generate documentation for the CLI
func generateDocs() {
	err := doc.GenMarkdownTree(cmd.RootCmd, "./docs")
	if err != nil {
		log.Fatal(err)
	}
}

func docs() {
	generateDocs()
}
