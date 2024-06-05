package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/hashicorp/terraform-exec/tfexec"
)

func main() {
	var tfPath, workingDir string
	flag.StringVar(&tfPath, "tfPath", "/usr/local/bin/terraform", "Path to Terraform binary")
	flag.StringVar(&workingDir, "workingDir", ".", "Working directory for Terraform")
	flag.Parse()

	ctx := context.Background()
	tf, err := tfexec.NewTerraform(workingDir, tfPath)
	if err != nil {
		log.Fatalf("error creating new Terraform: %s", err)
	}

	err = tf.Init(ctx, tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error initializing Terraform: %s", err)
	}

	output, err := tf.Graph(ctx)
	if err != nil {
		log.Fatalf("error running tf.Graph: %s", err)
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
	err = os.WriteFile("Terramaid.md", []byte(mermaidGraph), 0644)
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
		srcName := strings.ReplaceAll(edge.Src, ".", "_")
		dstName := strings.ReplaceAll(edge.Dst, ".", "_")
		sb.WriteString(fmt.Sprintf("		%s --> %s\n", srcName, dstName))
	}
	sb.WriteString("\tend\n```\n")

	return sb.String()
}
