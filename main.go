package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/awalterschulze/gographviz"
)

func main() {
	var tfPath string
	flag.StringVar(&tfPath, "tfPath", "/usr/local/bin/terraform", "Path to Terraform binary")
	flag.Parse()

	// Run terraform graph command
	cmd := exec.Command(tfPath, "graph")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running terraform graph command", err)
		return
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
		srcName := strings.Trim(edge.Src, "\"")
		dstName := strings.Trim(edge.Dst, "\"")
		sb.WriteString(fmt.Sprintf("		%s --> %s\n", srcName, dstName))
	}
	sb.WriteString("\tend\n```\n")

	return sb.String()
}
