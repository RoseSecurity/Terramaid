package internal

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/awalterschulze/gographviz"
)

var labelCleaner = regexp.MustCompile(`\s*\(expand\)|\s*\(close\)|\[root\]\s*|"`)

// CleanLabel removes unnecessary parts from the label
func CleanID(id string) string {
	id = labelCleaner.ReplaceAllString(id, "")
	if strings.HasPrefix(id, "provider[") {
		id = strings.ReplaceAll(id, "provider[", "provider_")
		id = strings.ReplaceAll(id, "]", "")
		id = strings.ReplaceAll(id, "/", "_")
		id = strings.ReplaceAll(id, ".", "_")
		return id
	}
	id = strings.ReplaceAll(id, ".", "_")
	id = strings.ReplaceAll(id, "/", "_")
	return id
}

func CleanLabel(label string) string {
	label = labelCleaner.ReplaceAllString(label, "")
	if strings.HasPrefix(label, "provider[") {
		label = strings.ReplaceAll(label, "[", ": ")
		label = strings.ReplaceAll(label, "]", "")
	}
	label = strings.ReplaceAll(label, "\\", "")
	return label
}

// GenerateMermaidFlowchart generates a Mermaid diagram from a gographviz graph
func GenerateMermaidFlowchart(graph *gographviz.Graph, direction string, subgraphName string) (string, error) {
	validDirections := map[string]bool{"TB": true, "TD": true, "BT": true, "RL": true, "LR": true}
	if !validDirections[direction] {
		return "", fmt.Errorf("invalid direction %s: valid options are TB, TD, BT, RL, LR", direction)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("```mermaid\nflowchart %s\n", direction))

	if subgraphName != "" {
		sb.WriteString(fmt.Sprintf("    subgraph %s\n", subgraphName))
	}

	addedNodes := make(map[string]string)

	addedProviders := make(map[string]bool)

	for _, node := range graph.Nodes.Nodes {
		nodeID := CleanID(node.Name)
		nodeLabel := CleanLabel(node.Attrs["label"])

		if nodeLabel == "" {
			continue
		}

		if strings.HasPrefix(nodeLabel, "provider:") {
			if addedProviders[nodeID] {
				continue
			}
			addedProviders[nodeID] = true
		}

		if _, exists := addedNodes[nodeID]; !exists {
			sb.WriteString(fmt.Sprintf("        %s[\"%s\"]\n", nodeID, nodeLabel))
			addedNodes[nodeID] = nodeLabel
		}
	}

	if subgraphName != "" {
		sb.WriteString("    end\n")
	}

	for _, edge := range graph.Edges.Edges {
		fromID := CleanID(edge.Src)
		toID := CleanID(edge.Dst)

		if _, exists := addedNodes[fromID]; !exists {
			fromLabel := CleanLabel(graph.Nodes.Lookup[edge.Src].Attrs["label"])
			if fromLabel != "" {
				sb.WriteString(fmt.Sprintf("        %s[\"%s\"]\n", fromID, fromLabel))
				addedNodes[fromID] = fromLabel
			}
		}

		if _, exists := addedNodes[toID]; !exists {
			toLabel := CleanLabel(graph.Nodes.Lookup[edge.Dst].Attrs["label"])
			if toLabel != "" {
				sb.WriteString(fmt.Sprintf("        %s[\"%s\"]\n", toID, toLabel))
				addedNodes[toID] = toLabel
			}
		}

		sb.WriteString(fmt.Sprintf("    %s --> %s\n", fromID, toID))
	}

	sb.WriteString("```\n")
	return sb.String(), nil
}
