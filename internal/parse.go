// Copyright (c) RoseSecurity
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"context"
	"fmt"

	"github.com/RoseSecurity/terramaid/pkg/utils"
	"github.com/awalterschulze/gographviz"
	"github.com/hashicorp/terraform-exec/tfexec"
)

const emptyGraph = `digraph G {
  rankdir = "RL";
  node [shape = rect, fontname = "sans-serif"];
  /* This configuration does not contain any resources.         */
  /* For a more detailed graph, try: terraform graph -type=plan */
}
`

// ParseTerraform parses the Terraform plan and returns the generated graph
func ParseTerraform(ctx context.Context, workingDir, tfPath, planFile string, verbose bool) (*gographviz.Graph, error) {
	if verbose {
		utils.LogVerbose("Initializing Terraform with working directory: %s", workingDir)
		utils.LogVerbose("Using Terraform binary: %s", tfPath)
	}

	tf, err := tfexec.NewTerraform(workingDir, tfPath)
	if err != nil {
		return nil, err
	}

	if verbose {
		utils.LogVerbose("Running terraform init with upgrade=true")
	}

	if err := tf.Init(ctx, tfexec.Upgrade(true)); err != nil {
		return nil, err
	}

	opts := &tfexec.GraphPlanOption{}

	if planFile != "" {
		if verbose {
			utils.LogVerbose("Using plan file for graph generation: %s", planFile)
		}
		opts = tfexec.GraphPlan(planFile)
	} else if verbose {
		utils.LogVerbose("No plan file specified, using current state")
	}

	if verbose {
		utils.LogVerbose("Running terraform graph command")
	}

	output, err := tf.Graph(ctx, opts)
	if err != nil {
		return nil, err
	}

	if output == emptyGraph {
		return nil, fmt.Errorf("no output from terraform graph")
	}

	if verbose {
		utils.LogVerbose("Successfully retrieved graph output from Terraform")
		utils.LogVerbose("Parsing DOT output")
	}

	// Parse the DOT output
	graphAst, err := gographviz.ParseString(output)
	if err != nil {
		return nil, err
	}

	graph := gographviz.NewGraph()

	if verbose {
		utils.LogVerbose("Analyzing graph structure")
	}

	if err := gographviz.Analyse(graphAst, graph); err != nil {
		return nil, err
	}

	if verbose {
		utils.LogVerbose("Graph analysis complete")
		utils.LogVerbose("Found %d nodes and %d edges", len(graph.Nodes.Nodes), len(graph.Edges.Edges))
	}

	return graph, nil
}
