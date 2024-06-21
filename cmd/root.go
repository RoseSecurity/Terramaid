package cmd

import (
	"fmt"
	"os"

	"github.com/RoseSecurity/terramaid/internal"
	"github.com/spf13/cobra"
)

// Declare variables
var (
	workingDir string
	tfDir      string
	tfPlan     string
	tfBinary   string
	output     string
)

var RootCmd = &cobra.Command{
	Use:          "terramaid",
	Short:        "A utility for generating Mermaid diagrams from Terraform",
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		graph, err := internal.ParseTerraform(workingDir, tfBinary, tfPlan)
		if err != nil {
			fmt.Printf("Error parsing Terraform: %v\n", err)
			os.Exit(1)
		}
		mermaidDiagram := internal.ConvertToMermaid(graph)
		// Write the Mermaid diagram to the specified output file
		err = os.WriteFile(output, []byte(mermaidDiagram), 0644)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Mermaid diagram successfully written to %s\n", output)
	},
}

func Execute() error {
	return RootCmd.Execute()
}

func init() {
	RootCmd.Flags().StringVarP(&output, "output", "o", "Terramaid.md", "Output file for Mermaid diagram")
	RootCmd.Flags().StringVarP(&tfDir, "tfDir", "d", ".", "Path to Terraform directory")
	RootCmd.Flags().StringVarP(&tfPlan, "tfPlan", "p", "", "Path to Terraform plan file")
	RootCmd.Flags().StringVarP(&tfBinary, "tfBinary", "b", "/usr/local/bin/terraform", "Path to Terraform binary")
	RootCmd.Flags().StringVarP(&workingDir, "workingDir", "w", ".", "Working directory for Terraform")
}
