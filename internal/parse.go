package internal

import (
	"context"

	"github.com/awalterschulze/gographviz"
	"github.com/hashicorp/terraform-exec/tfexec"
)

// ParseTerraform parses the Terraform plan and returns the generated graph
func ParseTerraform(workingDir, tfPath, planFile string) (*gographviz.Graph, error) {
	ctx := context.Background()
	tf, err := tfexec.NewTerraform(workingDir, tfPath)
	if err != nil {
		return nil, err
	}

	err = tf.Init(ctx, tfexec.Upgrade(true))
	if err != nil {
		return nil, err
	}

	var output string
	// Graph Terraform resources
	if planFile != "" {
		output, err = tf.Graph(ctx, tfexec.GraphPlan(planFile))
	} else {
		output, err = tf.Graph(ctx)
	}

	if err != nil {
		return nil, err
	}

	// Parse the DOT output
	dot := string(output)
	graphAst, err := gographviz.ParseString(dot)
	if err != nil {
		return nil, err
	}

	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		return nil, err
	}

	return graph, nil
}
