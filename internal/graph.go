package internal

import (
	"fmt"
	"strings"

	"github.com/awalterschulze/gographviz"
)

// ConvertToMermaid converts a gographviz graph to a Mermaid.js compatible string.
// It accepts a graph, direction, and an optional subgraph name.
func ConvertToMermaid(graph *gographviz.Graph, direction string, subgraphName string) (string, error) {
	var sb strings.Builder

	// Validate the direction of the flowchart. Valid options are: TB, TD, BT, RL, LR
	validDirections := map[string]bool{
		"TB": true, "TD": true, "BT": true, "RL": true, "LR": true,
	}
	if !validDirections[direction] {
		return "", fmt.Errorf("invalid direction %s: valid options are: TB, TD, BT, RL, LR", direction)
	}

	// Start Mermaid graph definition
	sb.WriteString("```mermaid\n")
	sb.WriteString("flowchart " + direction + ";\n")
	if subgraphName != "" {
		sb.WriteString(fmt.Sprintf("\tsubgraph %s\n", subgraphName))
	}

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

	if subgraphName != "" {
		sb.WriteString("\tend\n")
	}
	sb.WriteString("```\n")

	return sb.String(), nil
}
