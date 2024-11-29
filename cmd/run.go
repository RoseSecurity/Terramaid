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
	if opts.WorkingDir != "" {
		exists, err := utils.TerraformFilesExist(opts.WorkingDir)
		if err != nil {
			return fmt.Errorf("error checking Terraform files in directory \"%s\": %v", opts.WorkingDir, err)
		}
		if !exists {
			return fmt.Errorf("Terraform files do not exist in directory \"%s\"", opts.WorkingDir)
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
	}

	// Spinner initialization and graph parsing
	sp := utils.NewSpinner("Generating Terramaid Diagrams")
	sp.Start()

	graph, err := internal.ParseTerraform(opts.WorkingDir, opts.TFBinary, opts.TFPlan)
	if err != nil {
		sp.Stop()
		return fmt.Errorf("error parsing Terraform: %w", err)
	}

	// Generate the Mermaid diagram
	mermaidDiagram, err := internal.GenerateMermaidFlowchart(graph, opts.Direction, opts.SubgraphName)
	if err != nil {
		sp.Stop()
		return fmt.Errorf("error generating Mermaid diagram: %w", err)
	}

	// Write the Mermaid diagram to the specified output file
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

	// Disable auto-generated string from documentation so that documentation is cleanly built and updated
	runCmd.DisableAutoGenTag = true
}
