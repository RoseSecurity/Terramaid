package internal

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/awalterschulze/gographviz"
	"golang.org/x/text/cases"
  "golang.org/x/text/language"

)

type Node struct {
	ID    string
	Label string
}

type Edge struct {
	From string
	To   string
}

type Graph struct {
	Nodes []Node
	Edges []Edge
}

// Removes unnecessary parts from the label
func CleanLabel(label string) string {
	re := regexp.MustCompile(`\s*\(expand\)|\s*\(close\)|\[root\]\s*|"`)
	return re.ReplaceAllString(label, "")
}

// Removes unnecessary parts from the ID
func CleanID(id string) string {
	re := regexp.MustCompile(`\s*\(expand\)|\s*\(close\)|\[root\]\s*|"`)
	return re.ReplaceAllString(id, "")
}

// Extracts the provider for separate subgraph
func ExtractProvider(label string) string {
	if strings.Contains(label, "provider") {
		parts := strings.Split(label, "/")
		if len(parts) > 2 {
			return parts[len(parts)-2]
		}
	}
	return ""
}

// Transforms the parsed graph into cleaned nodes and edges
func TransformGraph(graph *gographviz.Graph) Graph {
	nodes := []Node{}
	edges := []Edge{}

	for _, node := range graph.Nodes.Nodes {
		cleanedID := CleanID(node.Name)
		cleanedLabel := CleanLabel(node.Attrs["label"])
		if cleanedLabel != "" && !strings.Contains(cleanedLabel, "provider") {
			nodes = append(nodes, Node{ID: cleanedID, Label: cleanedLabel})
		}
	}

	for _, edge := range graph.Edges.Edges {
		fromLabel := CleanLabel(graph.Nodes.Lookup[edge.Src].Attrs["label"])
		toLabel := CleanLabel(graph.Nodes.Lookup[edge.Dst].Attrs["label"])
		if fromLabel != "" && toLabel != "" && !strings.Contains(fromLabel, "provider") && !strings.Contains(toLabel, "provider") {
			edges = append(edges, Edge{From: CleanID(edge.Src), To: CleanID(edge.Dst)})
		}
	}

	return Graph{Nodes: nodes, Edges: edges}
}

// Converts a gographviz graph to a Mermaid.js compatible string.
// It accepts a graph, direction, and an optional subgraph name.
func ConvertToMermaidFlowchart(graph *gographviz.Graph, direction string, subgraphName string) (string, error) {
	var sb strings.Builder

	// Capitalize the provider name
	caser := cases.Title(language.English)
	// Validate the direction of the flowchart. Valid options are: TB, TD, BT, RL, LR
	validDirections := map[string]bool{
		"TB": true, "TD": true, "BT": true, "RL": true, "LR": true,
	}
	if !validDirections[direction] {
		return "", fmt.Errorf("invalid direction %s: valid options are: TB, TD, BT, RL, LR", direction)
	}
	// Start Mermaid graph definition
	sb.WriteString("```mermaid\n")
	sb.WriteString("flowchart " + direction + "\n")
	
	// Add subgraph for providers
	providerSubgraphs := make(map[string]bool)
	for _, n := range graph.Nodes.Nodes {
		label := n.Attrs["label"]
		provider := ExtractProvider(label)
		if provider != "" && !providerSubgraphs[provider] {
			sb.WriteString(fmt.Sprintf("\tsubgraph %s\n", caser.String(provider)))
			providerSubgraphs[provider] = true
		}
	}
	
	if subgraphName != "" {
		sb.WriteString(fmt.Sprintf("\tsubgraph %s\n", subgraphName))
	}
	
	// Iterate over nodes to add them to the Mermaid graph
	for _, n := range graph.Nodes.Nodes {
		label := CleanLabel(n.Attrs["label"])
		nodeName := CleanID(n.Name)
		if label != "" && nodeName != "" && !strings.Contains(label, "provider") {
			sb.WriteString(fmt.Sprintf("\t\t%s[\"%s\"]\n", nodeName, label))
		}
	}
	
	// Iterate over edges to add them to the Mermaid graph
	for _, edge := range graph.Edges.Edges {
		srcLabel := CleanLabel(graph.Nodes.Lookup[edge.Src].Attrs["label"])
		dstLabel := CleanLabel(graph.Nodes.Lookup[edge.Dst].Attrs["label"])
		srcName := CleanID(edge.Src)
		dstName := CleanID(edge.Dst)
		if srcLabel != "" && dstLabel != "" && !strings.Contains(srcLabel, "provider") && !strings.Contains(dstLabel, "provider") {
			sb.WriteString(fmt.Sprintf("\t\t%s --> %s\n", srcName, dstName))
		}
	}
	
	// Close all open subgraphs
	for range providerSubgraphs {
		sb.WriteString("\tend\n")
	}
	if subgraphName != "" {
		sb.WriteString("\tend\n")
	}
	
	sb.WriteString("```\n")
	return sb.String(), nil
}
