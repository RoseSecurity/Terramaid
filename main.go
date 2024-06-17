package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	u "github.com/RoseSecurity/terramaid/pkg/utils"
	"github.com/awalterschulze/gographviz"
	"github.com/hashicorp/terraform-exec/tfexec"
)

func main() {
	var tfPath, workingDir, planFile, outputFile string
	var output string

	flag.StringVar(&tfPath, "tfPath", "/usr/local/bin/terraform", "Path to Terraform binary")
	flag.StringVar(&workingDir, "workingDir", ".", "Working directory for Terraform")
	flag.StringVar(&planFile, "planFile", "", "Path to Terraform plan file")
	flag.StringVar(&outputFile, "outputFile", "Terramaid.md", "Output file for Mermaid diagram")
	flag.Parse()

	ctx := context.Background()
	tf, err := tfexec.NewTerraform(workingDir, tfPath)
	if err != nil {
		u.LogErrorAndExit(err)
	}

	err = tf.Init(ctx, tfexec.Upgrade(true))
	if err != nil {
		u.LogErrorAndExit(err)
	}

	// Graph Terraform resources
	if planFile != "" {
		output, err = tf.Graph(ctx, tfexec.GraphPlan(planFile))
	} else {
		output, err = tf.Graph(ctx)
	}

	if err != nil {
		u.LogErrorAndExit(err)
	}

	// Parse the DOT output
	dot := string(output)
	graphAst, err := gographviz.ParseString(dot)
	if err != nil {
		fmt.Println("Error parsing DOT:", err)
		return
	}

	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		fmt.Println("Error analyzing graph:", err)
		return
	}

	// Convert to Mermaid format
	mermaidGraph := ConvertToMermaid(graph)
	err = os.WriteFile(outputFile, []byte(mermaidGraph), 0644)
	if err != nil {
		fmt.Println("Error writing to Terramaid file:", err)
		return
	}
}

func ConvertToMermaid(graph *gographviz.Graph) string {
	var sb strings.Builder

	sb.WriteString("```mermaid\n")
	sb.WriteString("flowchart TD;\n")
	sb.WriteString("\tsubgraph Terraform\n")
	for _, node := range graph.Nodes.Nodes {
		label := strings.Trim(node.Attrs["label"], "\"")
		nodeName := strings.Trim(node.Name, "\"")
		sb.WriteString(fmt.Sprintf("		%s[\"%s\"]\n", nodeName, label))
	}

	for _, edge := range graph.Edges.Edges {
		srcName := strings.Trim(edge.Src, "\"")
		dstName := strings.Trim(edge.Dst, "\"")
		sb.WriteString(fmt.Sprintf("		%s --> %s\n", srcName, dstName))
	}
	sb.WriteString("\tend\n```\n")

	return sb.String()
}
