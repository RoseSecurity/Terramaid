package internal

import (
	"fmt"
	"strings"

	"github.com/awalterschulze/gographviz"
)

// ConvertToMermaid converts a gographviz graph to a Mermaid.js compatible string
func ConvertToMermaid(graph *gographviz.Graph) string {
	var sb strings.Builder

	// Start Mermaid graph definition
	sb.WriteString("```mermaid\n")
	sb.WriteString("flowchart TD;\n")
	sb.WriteString("\tsubgraph Terraform\n")
	// Iterate over nodes to add them to the Mermaid graph
	for _, node := range graph.Nodes.Nodes {
		label := strings.Trim(node.Attrs["label"], "\"")
		nodeName := strings.Trim(node.Name, "\"")
		sb.WriteString(fmt.Sprintf("		%s[\"%s\"]\n", nodeName, label))
	}
	// Iterate over edges to add them to the Mermaid graph
	for _, edge := range graph.Edges.Edges {
		srcName := strings.Trim(edge.Src, "\"")
		dstName := strings.Trim(edge.Dst, "\"")
		sb.WriteString(fmt.Sprintf("		%s --> %s\n", srcName, dstName))
	}
	// End Mermaid graph definition
	sb.WriteString("\tend\n```\n")

	return sb.String()
}
