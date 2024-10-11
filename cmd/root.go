package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/RoseSecurity/terramaid/internal"
	"github.com/RoseSecurity/terramaid/pkg/utils"
	"github.com/caarlos0/env/v11"
	"github.com/spf13/cobra"
)

var Version string

type options struct {
	WorkingDir   string `env:"WORKING_DIR" envDefault:"."`
	TFDir        string `env:"TF_DIR" envDefault:"."`
	TFPlan       string `env:"TF_PLAN"`
	TFBinary     string `env:"TF_BINARY"`
	Output       string `env:"OUTPUT" envDefault:"Terramaid.md"`
	Direction    string `env:"DIRECTION" envDefault:"TD"`
	SubgraphName string `env:"SUBGRAPH_NAME" envDefault:"Terraform"`
	ChartType    string `env:"CHART_TYPE" envDefault:"flowchart"`
}

func TerramaidCmd() *cobra.Command {
	options := &options{}

	// Parse Envs
	if err := env.ParseWithOptions(options, env.Options{Prefix: "TERRAMAID_"}); err != nil {
		log.Fatalf("error parsing envs: %s", err.Error())
	}

	cmd := &cobra.Command{
		Use:           "terramaid",
		Short:         "A utility for generating Mermaid diagrams from Terraform configurations",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if options.TFDir != "" && !utils.DirExists(options.TFDir) {
				return fmt.Errorf("Terraform directory \"%s\" does not exist", options.TFDir)
			}

			if options.TFDir != "" {
				exists, err := utils.TerraformFilesExist(options.TFDir)
				if err != nil {
					return fmt.Errorf("error checking Terraform files in directory \"%s\": %v", options.TFDir, err)
				}
				if !exists {
					return fmt.Errorf("Terraform files do not exist in directory \"%s\"", options.TFDir)
				}
			}

			if options.WorkingDir != "" && !utils.DirExists(options.WorkingDir) {
				return fmt.Errorf("Working directory \"%s\" does not exist", options.WorkingDir)
			}

			if options.TFPlan != "" && !utils.DirExists(options.TFPlan) {
				return fmt.Errorf("Terraform planfile \"%s\" does not exist", options.TFPlan)
			}

			if options.TFBinary == "" {
				tfBinary, err := exec.LookPath("terraform")
				if err != nil {
					return fmt.Errorf("error finding Terraform binary: %w", err)
				}

				options.TFBinary = tfBinary
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			sp := utils.NewSpinner("Generating Terramaid Diagrams")
			sp.Start()
			graph, err := internal.ParseTerraform(options.WorkingDir, options.TFBinary, options.TFPlan)
			if err != nil {
				return fmt.Errorf("error parsing Terraform: %w", err)
			}

			// Convert the graph to a Mermaid diagram
			var mermaidDiagram string

			switch options.ChartType {
			case "flowchart":
				mermaidDiagram, err = internal.ConvertToMermaidFlowchart(graph, options.Direction, options.SubgraphName)
				if err != nil {
					return fmt.Errorf("error converting to Mermaid flowchart: %w", err)
				}
			default:
				return fmt.Errorf("unsupported chart type: %s", options.ChartType)
			}

			// Write the Mermaid diagram to the specified output file
			if err := os.WriteFile(options.Output, []byte(mermaidDiagram), 0o644); err != nil {
				return fmt.Errorf("error writing to file: %w", err)
			}

			sp.Stop()
			fmt.Printf("Mermaid diagram successfully written to %s\n", options.Output)

			return nil
		},
	}

	cmd.Flags().StringVarP(&options.Output, "output", "o", options.Output, "Output file for Mermaid diagram (env: TERRAMAID_OUTPUT)")
	cmd.Flags().StringVarP(&options.Direction, "direction", "r", options.Direction, "Specify the direction of the diagram (env: TERRAMAID_DIRECTION)")
	cmd.Flags().StringVarP(&options.SubgraphName, "subgraph-name", "s", options.SubgraphName, "Specify the subgraph name of the diagram (env: TERRAMAID_SUBGRAPH_NAME)")
	cmd.Flags().StringVarP(&options.ChartType, "chart-type", "c", options.ChartType, "Specify the type of Mermaid chart to generate (env: TERRAMAID_CHART_TYPE)")
	cmd.Flags().StringVarP(&options.TFDir, "tf-dir", "d", options.TFDir, "Path to Terraform directory (env: TERRAMAID_TF_DIR)")
	cmd.Flags().StringVarP(&options.TFPlan, "tf-plan", "p", options.TFPlan, "Path to Terraform plan file (env: TERRAMAID_TF_PLAN)")
	cmd.Flags().StringVarP(&options.TFBinary, "tf-binary", "b", options.TFBinary, "Path to Terraform binary (env: TERRAMAID_TF_BINARY)")
	cmd.Flags().StringVarP(&options.WorkingDir, "working-dir", "w", options.WorkingDir, "Working directory for Terraform (env: TERRAMAID_WORKING_DIR)")

	cmd.AddCommand(docsCmd(), versionCmd(Version))

	// Disable auto generated string from documentation so that documentation is cleanly built and updated
	cmd.DisableAutoGenTag = true

	return cmd
}

func Execute() error {
	if err := TerramaidCmd().Execute(); err != nil {
		return err
	}

	return nil
}
