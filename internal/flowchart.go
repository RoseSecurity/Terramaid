// Copyright (c) RoseSecurity
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/RoseSecurity/terramaid/pkg/utils"
	"github.com/awalterschulze/gographviz"
)

var labelCleaner = regexp.MustCompile(`\s*\(expand\)|\s*\(close\)|\[root\]\s*|"`)

// CleanID removes unnecessary parts from the label and sanitizes for Mermaid compatibility
func CleanID(id string) string {
	id = labelCleaner.ReplaceAllString(id, "")
	if strings.HasPrefix(id, "provider[") {
		id = strings.ReplaceAll(id, "provider[", "provider_")
		id = strings.ReplaceAll(id, "]", "")
		id = strings.ReplaceAll(id, "/", "_")
		id = strings.ReplaceAll(id, ".", "_")
		return sanitizeMermaidID(id)
	}
	id = strings.ReplaceAll(id, ".", "_")
	id = strings.ReplaceAll(id, "/", "_")
	return sanitizeMermaidID(id)
}

// sanitizeMermaidID removes or replaces characters that can cause Mermaid parsing issues
func sanitizeMermaidID(id string) string {
	// Replace problematic characters that can cause Mermaid parsing errors
	id = strings.ReplaceAll(id, "(", "_")
	id = strings.ReplaceAll(id, ")", "_")
	id = strings.ReplaceAll(id, "[", "_")
	id = strings.ReplaceAll(id, "]", "_")
	id = strings.ReplaceAll(id, "{", "_")
	id = strings.ReplaceAll(id, "}", "_")
	id = strings.ReplaceAll(id, "<", "_")
	id = strings.ReplaceAll(id, ">", "_")
	id = strings.ReplaceAll(id, " ", "_")
	id = strings.ReplaceAll(id, "-", "_")
	id = strings.ReplaceAll(id, ":", "_")
	id = strings.ReplaceAll(id, ";", "_")
	id = strings.ReplaceAll(id, ",", "_")
	id = strings.ReplaceAll(id, "!", "_")
	id = strings.ReplaceAll(id, "@", "_")
	id = strings.ReplaceAll(id, "#", "_")
	id = strings.ReplaceAll(id, "$", "_")
	id = strings.ReplaceAll(id, "%", "_")
	id = strings.ReplaceAll(id, "^", "_")
	id = strings.ReplaceAll(id, "&", "_")
	id = strings.ReplaceAll(id, "*", "_")
	id = strings.ReplaceAll(id, "+", "_")
	id = strings.ReplaceAll(id, "=", "_")
	id = strings.ReplaceAll(id, "|", "_")
	id = strings.ReplaceAll(id, "\\", "_")
	id = strings.ReplaceAll(id, "?", "_")
	id = strings.ReplaceAll(id, "'", "_")
	id = strings.ReplaceAll(id, "\"", "_")
	id = strings.ReplaceAll(id, "`", "_")
	id = strings.ReplaceAll(id, "~", "_")

	// Remove multiple consecutive underscores
	for strings.Contains(id, "__") {
		id = strings.ReplaceAll(id, "__", "_")
	}

	// Trim leading and trailing underscores
	id = strings.Trim(id, "_")

	// Ensure the ID is not empty and starts with a letter or underscore
	if id == "" || (!strings.HasPrefix(id, "_") && (id[0] < 'A' || (id[0] > 'Z' && id[0] < 'a') || id[0] > 'z')) {
		id = "node_" + id
	}

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
func GenerateMermaidFlowchart(ctx context.Context, graph *gographviz.Graph, direction string, subgraphName string, verbose bool) (string, error) {
	validDirections := map[string]bool{"TB": true, "TD": true, "BT": true, "RL": true, "LR": true}
	if !validDirections[direction] {
		return "", fmt.Errorf("invalid direction %s: valid options are TB, TD, BT, RL, LR", direction)
	}

	if verbose {
		utils.LogVerbose("Generating Mermaid flowchart with direction: %s", direction)
		if subgraphName != "" {
			utils.LogVerbose("Using subgraph name: %s", subgraphName)
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("```mermaid\nflowchart %s\n", direction))

	if subgraphName != "" {
		sb.WriteString(fmt.Sprintf("    subgraph %s\n", subgraphName))
	}

	addedNodes := make(map[string]string)
	addedProviders := make(map[string]bool)

	if verbose {
		utils.LogVerbose("Processing %d nodes", len(graph.Nodes.Nodes))
	}

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
			if verbose {
				utils.LogVerbose("Added provider node: %s", nodeID)
			}
		}

		if _, exists := addedNodes[nodeID]; !exists {
			sb.WriteString(fmt.Sprintf("        %s[\"%s\"]\n", nodeID, nodeLabel))
			addedNodes[nodeID] = nodeLabel
			if verbose && !strings.HasPrefix(nodeLabel, "provider:") {
				utils.LogVerbose("Added node: %s", nodeID)
			}
		}
	}

	if subgraphName != "" {
		sb.WriteString("    end\n")
	}

	if verbose {
		utils.LogVerbose("Processing %d edges", len(graph.Edges.Edges))
	}

	for _, edge := range graph.Edges.Edges {
		fromID := CleanID(edge.Src)
		toID := CleanID(edge.Dst)

		if _, exists := addedNodes[fromID]; !exists {
			fromLabel := CleanLabel(graph.Nodes.Lookup[edge.Src].Attrs["label"])
			if fromLabel != "" {
				sb.WriteString(fmt.Sprintf("        %s[\"%s\"]\n", fromID, fromLabel))
				addedNodes[fromID] = fromLabel
				if verbose {
					utils.LogVerbose("Added source node from edge: %s", fromID)
				}
			}
		}

		if _, exists := addedNodes[toID]; !exists {
			toLabel := CleanLabel(graph.Nodes.Lookup[edge.Dst].Attrs["label"])
			if toLabel != "" {
				sb.WriteString(fmt.Sprintf("        %s[\"%s\"]\n", toID, toLabel))
				addedNodes[toID] = toLabel
				if verbose {
					utils.LogVerbose("Added destination node from edge: %s", toID)
				}
			}
		}

		sb.WriteString(fmt.Sprintf("    %s --> %s\n", fromID, toID))
		if verbose {
			utils.LogVerbose("Added edge: %s --> %s", fromID, toID)
		}
	}

	sb.WriteString("```\n")

	if verbose {
		nodeCount := len(addedNodes)
		edgeCount := len(graph.Edges.Edges)
		utils.LogVerbose("Mermaid diagram generation complete with %d nodes and %d edges", nodeCount, edgeCount)
	}

	return sb.String(), nil
}
