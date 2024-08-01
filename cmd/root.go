package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/RoseSecurity/terramaid/internal"
	"github.com/caarlos0/env/v11"
	"github.com/spf13/cobra"
)

var Version = "0.0.1"
type opts struct {
	WorkingDir   string `env:"WORKING_DIR" envDefault:"."`
	TFDir        string `env:"TF_DIR" envDefault:"."`
	TFPlan       string `env:"TF_PLAN"`
	TFBinary     string `env:"TF_BINARY"`
	Output       string `env:"OUTPUT" envDefault:"Terramaid.md"`
	Direction    string `env:"DIRECTION" envDefault:"TD"`
	SubgraphName string `env:"SUBGRAPH_NAME" envDefault:"Terraform"`
}

func TerramaidCmd() *cobra.Command {
	opts := &opts{}

	// Parse Envs
	if err := env.ParseWithOptions(opts, env.Options{Prefix: "TERRAMAID_"}); err != nil {
		log.Fatalf("error parsing envs: %w", err)
	}

	cmd := &cobra.Command{
		Use:          "terramaid",
		Short:        "A utility for generating Mermaid diagrams from Terraform",
		SilenceUsage: true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.TFBinary == "" {
				tfBinary, err := exec.LookPath("terraform")
				if err != nil {
					return fmt.Errorf("error finding Terraform binary: %w", err)
				}

				opts.TFBinary = tfBinary
			}

			graph, err := internal.ParseTerraform(opts.WorkingDir, opts.TFBinary, opts.TFPlan)
			if err != nil {
				return fmt.Errorf("error parsing Terraform: %w", err)
			}

			// Convert the graph to a Mermaid diagram
			mermaidDiagram, err := internal.ConvertToMermaid(graph, opts.Direction, opts.SubgraphName)
			if err != nil {
				return fmt.Errorf("error converting to Mermaid: %w\n", err)
			}

			// Write the Mermaid diagram to the specified output file
			if err := os.WriteFile(opts.Output, []byte(mermaidDiagram), 0o644); err != nil {
				return fmt.Errorf("Error writing to file: %w\n", err)
			}

			fmt.Printf("Mermaid diagram successfully written to %s\n", opts.Output)

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.Output, "output", "o", opts.Output, "Output file for Mermaid diagram (env: TERRAMAID_OUTPUT)")
	cmd.Flags().StringVarP(&opts.Direction, "direction", "r", opts.Direction, "Specify the direction of the flowchart (env: TERRAMAID_DIRECTION)")
	cmd.Flags().StringVarP(&opts.SubgraphName, "subgraph-name", "s", opts.SubgraphName, "Specify the subgraph name of the flowchart (env: TERRAMAID_SUBGRAPH_NAME)")
	cmd.Flags().StringVarP(&opts.TFDir, "tf-dir", "d", opts.TFDir, "Path to Terraform directory (env: TERRAMAID_TF_DIR)")
	cmd.Flags().StringVarP(&opts.TFPlan, "tf-plan", "p", opts.TFPlan, "Path to Terraform plan file (env: TERRAMAID_TF_PLAN)")
	cmd.Flags().StringVarP(&opts.TFBinary, "tf-binary", "b", opts.TFBinary, "Path to Terraform binary (env: TERRAMAID_TF_BINARY)")
	cmd.Flags().StringVarP(&opts.WorkingDir, "working-dir", "w", opts.WorkingDir, "Working directory for Terraform (env: TERRAMAID_WORKING_DIR)")

	cmd.AddCommand(docsCmd(), versionCmd())

	return cmd
}

func Execute() error {
	if err := TerramaidCmd().Execute(); err != nil {
		return err
	}

	return nil
}
