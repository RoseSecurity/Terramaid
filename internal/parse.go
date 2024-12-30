// Copyright (c) RoseSecurity
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"context"
	"fmt"

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
func ParseTerraform(workingDir, tfPath, planFile string) (*gographviz.Graph, error) {
	ctx := context.Background()
	tf, err := tfexec.NewTerraform(workingDir, tfPath)
	if err != nil {
		return nil, err
	}

	if err := tf.Init(ctx, tfexec.Upgrade(true)); err != nil {
		return nil, err
	}

	opts := &tfexec.GraphPlanOption{}

	if planFile != "" {
		opts = tfexec.GraphPlan(planFile)
	}

	output, err := tf.Graph(ctx, opts)
	if err != nil {
		return nil, err
	}

	if output == emptyGraph {
		return nil, fmt.Errorf("no output from terraform graph")
	}

	// Parse the DOT output
	graphAst, err := gographviz.ParseString(string(output))
	if err != nil {
		return nil, err
	}

	graph := gographviz.NewGraph()

	if err := gographviz.Analyse(graphAst, graph); err != nil {
		return nil, err
	}

	return graph, nil
}
