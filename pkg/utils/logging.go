// Copyright (c) RoseSecurity
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/RoseSecurity/terramaid/internal"
	"github.com/RoseSecurity/terramaid/pkg/utils"
	"github.com/caarlos0/env/v11"
	"github.com/spf13/cobra"
)

type options struct {
	WorkingDir   string `env:"WORKING_DIR" envDefault:"."`
	TFPlan       string `env:"TF_PLAN"`
	TFBinary     string `env:"TF_BINARY"`
	Output       string `env:"OUTPUT" envDefault:"Terramaid.md"`
	Direction    string `env:"DIRECTION" envDefault:"TD"`
	SubgraphName string `env:"SUBGRAPH_NAME" envDefault:"Terraform"`
	ChartType    string `env:"CHART_TYPE" envDefault:"flowchart"`
	Verbose      bool   `env:"VERBOSE" envDefault:"false"`
}

var opts options // Global variable for flags and env variables

var runCmd = &cobra.Command{
	Use:           "run",
	Short:         "Generate Mermaid diagrams from Terraform configurations",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// The opts variable is automatically populated with flags here
		return generateDiagrams(&opts)
	},
}

func generateDiagrams(opts *options) error {
	if opts.Verbose {
		utils.LogVerbose("Starting Terramaid with the following options:")
		utils.LogVerbose("- Working Directory: %s", opts.WorkingDir)
		utils.LogVerbose("- Terraform Plan: %s", opts.TFPlan)
		utils.LogVerbose("- Terraform Binary: %s", opts.TFBinary)
		utils.LogVerbose("- Output File: %s", opts.Output)
		utils.LogVerbose("- Direction: %s", opts.Direction)
		utils.LogVerbose("- Subgraph Name: %s", opts.SubgraphName)
		utils.LogVerbose("- Chart Type: %s", opts.ChartType)
	}

	if opts.WorkingDir != "" {
		exists, err := utils.TerraformFilesExist(opts.WorkingDir)
		if err != nil {
			return fmt.Errorf("error checking Terraform files in directory \"%s\": %v", opts.WorkingDir, err)
		}
		if !exists {
			return fmt.Errorf("Terraform files do not exist in directory \"%s\"", opts.WorkingDir)
		}
		if opts.Verbose {
			utils.LogVerbose("Confirmed Terraform files exist in %s", opts.WorkingDir)
		}
	}

	// Validate directories and files
	if opts.WorkingDir != "" && !utils.DirExists(opts.WorkingDir) {
		return fmt.Errorf("terraform directory \"%s\" does not exist", opts.WorkingDir)
	}

	// Check for Terraform binary
	if opts.TFBinary == "" {
		tfBinary, err := exec.LookPath("terraform")
		if err != nil {
			return fmt.Errorf("error finding Terraform binary: %w", err)
		}
		opts.TFBinary = tfBinary
		if opts.Verbose {
			utils.LogVerbose("Terraform binary found at: %s", opts.TFBinary)
		}
	}

	// Spinner initialization and graph parsing
	sp := utils.NewSpinner("Generating Terramaid Diagrams")
	sp.Start()

	if opts.Verbose {
		utils.LogVerbose("Initializing Terraform and building graph...")
	}
	graph, err := internal.ParseTerraform(opts.WorkingDir, opts.TFBinary, opts.TFPlan, opts.Verbose)
	if err != nil {
		sp.Stop()
		return fmt.Errorf("error parsing Terraform: %w", err)
	}

	// Generate the Mermaid diagram
	if opts.Verbose {
		utils.LogVerbose("Generating Mermaid flowchart...")
	}
	mermaidDiagram, err := internal.GenerateMermaidFlowchart(graph, opts.Direction, opts.SubgraphName, opts.Verbose)
	if err != nil {
		sp.Stop()
		return fmt.Errorf("error generating Mermaid diagram: %w", err)
	}

	// Write the Mermaid diagram to the specified output file
	if opts.Verbose {
		utils.LogVerbose("Writing Mermaid diagram to %s", opts.Output)
	}
	if err := os.WriteFile(opts.Output, []byte(mermaidDiagram), 0o644); err != nil {
		sp.Stop()
		return fmt.Errorf("error writing to file: %w", err)
	}

	sp.Stop()
	fmt.Printf("Mermaid diagram successfully written to %s\n", opts.Output)

	return nil
}

func init() {
	// Parse environment variables first, then bind flags to the opts struct
	if err := env.ParseWithOptions(&opts, env.Options{Prefix: "TERRAMAID_"}); err != nil {
		fmt.Printf("Error parsing environment variables: %s\n", err.Error())
	}

	// Bind flags to the opts struct
	runCmd.Flags().StringVarP(&opts.Output, "output", "o", opts.Output, "Output file for Mermaid diagram (env: TERRAMAID_OUTPUT)")
	runCmd.Flags().StringVarP(&opts.Direction, "direction", "r", opts.Direction, "Specify the direction of the diagram (env: TERRAMAID_DIRECTION)")
	runCmd.Flags().StringVarP(&opts.SubgraphName, "subgraph-name", "s", opts.SubgraphName, "Specify the subgraph name of the diagram (env: TERRAMAID_SUBGRAPH_NAME)")
	runCmd.Flags().StringVarP(&opts.ChartType, "chart-type", "c", opts.ChartType, "Specify the type of Mermaid chart to generate (env: TERRAMAID_CHART_TYPE)")
	runCmd.Flags().StringVarP(&opts.TFPlan, "tf-plan", "p", opts.TFPlan, "Path to Terraform plan file (env: TERRAMAID_TF_PLAN)")
	runCmd.Flags().StringVarP(&opts.TFBinary, "tf-binary", "b", opts.TFBinary, "Path to Terraform binary (env: TERRAMAID_TF_BINARY)")
	runCmd.Flags().StringVarP(&opts.WorkingDir, "working-dir", "w", opts.WorkingDir, "Working directory for Terraform (env: TERRAMAID_WORKING_DIR)")
	runCmd.Flags().BoolVarP(&opts.Verbose, "verbose", "v", opts.Verbose, "Enable verbose output (env: TERRAMAID_VERBOSE)")

	// Disable auto-generated string from documentation so that documentation is cleanly built and updated
	runCmd.DisableAutoGenTag = true
}
