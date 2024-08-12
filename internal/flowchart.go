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
	ID       string
	Label    string
	Count    int
	Provider string
}

type Edge struct {
	From string
	To   string
}

type Graph struct {
	Nodes   []Node
	Edges   []Edge
	NodeMap map[string]int
}

var labelCleaner = regexp.MustCompile(`\s*\(expand\)|\s*\(close\)|\[root\]\s*|"`)

// CleanLabel removes unnecessary parts from the label
func CleanLabel(label string) string {
	return labelCleaner.ReplaceAllString(label, "")
}

// CleanID removes unnecessary parts from the ID
func CleanID(id string) string {
	return labelCleaner.ReplaceAllString(id, "")
}

// ExtractProvider extracts the provider for separate subgraph
func ExtractProvider(label string) string {
	parts := strings.Split(label, "_")
	if len(parts) > 0 {
		// Remove quotes from the provider name
		return strings.ReplaceAll(parts[0], "\"", "")
	}
	return ""
}

// TransformGraph transforms the parsed graph into cleaned nodes and edges
func TransformGraph(graph *gographviz.Graph) Graph {
	nodes := []Node{}
	edges := []Edge{}
	nodeMap := make(map[string]int)

	for _, node := range graph.Nodes.Nodes {
		cleanedID := CleanID(node.Name)
		cleanedLabel := CleanLabel(node.Attrs["label"])
		provider := ExtractProvider(cleanedLabel)
		if cleanedLabel != "" {
			nodeMap[cleanedLabel]++
			nodes = append(nodes, Node{ID: cleanedID, Label: cleanedLabel, Count: nodeMap[cleanedLabel], Provider: provider})
		}
	}

	for _, edge := range graph.Edges.Edges {
		fromLabel := CleanLabel(graph.Nodes.Lookup[edge.Src].Attrs["label"])
		toLabel := CleanLabel(graph.Nodes.Lookup[edge.Dst].Attrs["label"])
		if fromLabel != "" && toLabel != "" {
			edges = append(edges, Edge{From: CleanID(edge.Src), To: CleanID(edge.Dst)})
		}
	}

	return Graph{Nodes: nodes, Edges: edges, NodeMap: nodeMap}
}

// ConvertToMermaidFlowchart converts a gographviz graph to a Mermaid.js compatible string.
// It accepts a graph, direction, and an optional subgraph name.
func ConvertToMermaidFlowchart(graph *gographviz.Graph, direction string, subgraphName string) (string, error) {
	var sb strings.Builder

	caser := cases.Title(language.English)
	validDirections := map[string]bool{
		"TB": true, "TD": true, "BT": true, "RL": true, "LR": true,
	}
	if !validDirections[direction] {
		return "", fmt.Errorf("invalid direction %s: valid options are: TB, TD, BT, RL, LR", direction)
	}

	sb.WriteString("```mermaid\n")
	sb.WriteString("flowchart " + direction + "\n")

	if subgraphName != "" {
		sb.WriteString(fmt.Sprintf("\tsubgraph %s\n", subgraphName))
	}

	providerSubgraphs := make(map[string]bool)
	for _, n := range graph.Nodes.Nodes {
		provider := ExtractProvider(n.Attrs["label"])
		if provider != "" && !providerSubgraphs[provider] {
			sb.WriteString(fmt.Sprintf("\t\tsubgraph %s\n", caser.String(provider)))
			providerSubgraphs[provider] = true
		}
	}

	nodeMap := make(map[string]int)
	for _, n := range graph.Nodes.Nodes {
		label := CleanLabel(n.Attrs["label"])
		nodeName := CleanID(n.Name)
		if label != "" && nodeName != "" {
			nodeMap[label]++
			count := nodeMap[label]
			if count > 1 {
				sb.WriteString(fmt.Sprintf("\t\t\t%s[\"%s\\nCount: %d\"]\n", nodeName, label, count))
			} else {
				sb.WriteString(fmt.Sprintf("\t\t\t%s[\"%s\"]\n", nodeName, label))
			}
		}
	}

	for range providerSubgraphs {
		sb.WriteString("\t\tend\n")
	}

	for _, edge := range graph.Edges.Edges {
		srcLabel := CleanLabel(graph.Nodes.Lookup[edge.Src].Attrs["label"])
		dstLabel := CleanLabel(graph.Nodes.Lookup[edge.Dst].Attrs["label"])
		srcName := CleanID(edge.Src)
		dstName := CleanID(edge.Dst)
		if srcLabel != "" && dstLabel != "" {
			sb.WriteString(fmt.Sprintf("\t\t%s --> %s\n", srcName, dstName))
		}
	}

	if subgraphName != "" {
		sb.WriteString("\tend\n")
	}

	sb.WriteString("```\n")
	return sb.String(), nil
}
