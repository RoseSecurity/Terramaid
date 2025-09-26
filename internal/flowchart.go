// SPDX-License-Identifier: Apache-2.0

package internal

import (
    "context"
    "fmt"
    "os"
    "regexp"
    "strings"

    "github.com/RoseSecurity/terramaid/pkg/utils"
    "github.com/awalterschulze/gographviz"
)

var (
    labelCleaner = regexp.MustCompile(`\s*\(expand\)|\s*\(close\)|\[root\]\s*|"`)
    // Regex to match all problematic characters for Mermaid IDs in one pass
    mermaidUnsafeChars = regexp.MustCompile(`[()\[\]{}<>\s\-:;,!@#$%^&*+=|\\?\'"` + "`" + `~]+`)
    // Regex to match multiple consecutive underscores
    multipleUnderscores = regexp.MustCompile(`_+`)
    // Optional, user-provided regex (via env TERRAMAID_RESOURCE_TYPE_REGEX) to classify resource types
    customResourceTypeMatcher *regexp.Regexp
    // Optional, user-provided prefixes (via env TERRAMAID_RESOURCE_TYPE_PREFIXES, comma-separated)
    customResourceTypePrefixes []string
    // Built-in defaults for common/major provider resource type prefixes
    defaultResourceTypePrefixes = []string{
        "aws", "azurerm", "google", "kubernetes", "helm",
        "cloudflare", "datadog", "github", "gitlab", "digitalocean",
        "linode", "openstack", "alicloud", "oci", "heroku",
        "pagerduty", "random", "null", "tls",
    }
)

func init() {
    if pattern := os.Getenv("TERRAMAID_RESOURCE_TYPE_REGEX"); pattern != "" {
        // Best-effort compile; if it fails, we silently ignore and proceed with defaults
        if re, err := regexp.Compile(pattern); err == nil {
            customResourceTypeMatcher = re
        }
    }

    if prefixes := os.Getenv("TERRAMAID_RESOURCE_TYPE_PREFIXES"); prefixes != "" {
        // Normalize: trim spaces and drop empties
        for _, p := range strings.Split(prefixes, ",") {
            p = strings.TrimSpace(p)
            if p != "" {
                customResourceTypePrefixes = append(customResourceTypePrefixes, p)
            }
        }
    }
}

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
    // Replace all problematic characters with underscores in one pass
    id = mermaidUnsafeChars.ReplaceAllString(id, "_")

    // Replace multiple consecutive underscores with single underscore
    id = multipleUnderscores.ReplaceAllString(id, "_")

    // Trim leading and trailing underscores
    id = strings.Trim(id, "_")

    // Ensure the ID is not empty and starts with a letter or underscore
    if id == "" {
        id = "node_"
    } else if len(id) > 0 && !strings.HasPrefix(id, "_") && (id[0] < 'A' || (id[0] > 'Z' && id[0] < 'a') || id[0] > 'z') {
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

// isResourceLabel attempts to determine if a Terraform DOT node label represents a resource
// rather than a provider/module/variable/meta node. This is heuristic-based but works well
// across common providers and module addressing.
func isResourceLabel(label string) bool {
    if label == "" {
        return false
    }

    // Exclude providers explicitly handled elsewhere
    if strings.HasPrefix(label, "provider:") {
        return false
    }

    // Terraform graph labels generally follow addressing like:
    // - <type>.<name>
    // - module.<mod>.<type>.<name>
    // - data.<type>.<name> (optionally with module prefixes)
    parts := strings.Split(label, ".")
    if len(parts) < 2 {
        return false
    }

    // Strip any number of leading module segments: module.<name>.
    for len(parts) >= 2 && parts[0] == "module" {
        parts = parts[2:]
    }

    if len(parts) < 2 {
        return false
    }

    // Determine the resource type segment
    var typeSeg string
    if parts[0] == "data" {
        // data.<type>.<name>
        if len(parts) < 3 {
            return false
        }
        typeSeg = parts[1]
    } else {
        // <type>.<name> (possibly with extra addressing suffixes)
        typeIdx := len(parts) - 2
        if typeIdx < 0 {
            return false
        }
        typeSeg = parts[typeIdx]
    }

    // Optional custom matcher provided by user
    if customResourceTypeMatcher != nil && customResourceTypeMatcher.MatchString(typeSeg) {
        return true
    }

    // Optional custom prefixes provided by user
    for _, pref := range customResourceTypePrefixes {
        if strings.HasPrefix(typeSeg, pref) {
            return true
        }
    }

    // Built-in default prefixes for major providers
    for _, pref := range defaultResourceTypePrefixes {
        if strings.HasPrefix(typeSeg, pref) {
            return true
        }
    }

    // Generic heuristic: most Terraform resource types contain an underscore
    // (e.g., aws_s3_bucket, azurerm_resource_group, google_compute_instance, customprov_widget)
    if strings.Contains(typeSeg, "_") {
        return true
    }

    return false
}

// GenerateMermaidFlowchart generates a Mermaid diagram from a gographviz graph
func GenerateMermaidFlowchart(ctx context.Context, graph *gographviz.Graph, direction string, subgraphName string, resourcesOnly bool, verbose bool) (string, error) {
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

		// If resourcesOnly is enabled, skip non-resource nodes
		if resourcesOnly && !isResourceLabel(nodeLabel) {
			if verbose {
				utils.LogVerbose("Skipping non-resource node due to resourcesOnly: %s (%s)", nodeID, nodeLabel)
			}
			continue
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
				if resourcesOnly && !isResourceLabel(fromLabel) {
					if verbose {
						utils.LogVerbose("Skipping edge source (non-resource) due to resourcesOnly: %s", fromID)
					}
					// Skip adding this node; if either endpoint is non-resource, skip the edge later
				} else {
					sb.WriteString(fmt.Sprintf("        %s[\"%s\"]\n", fromID, fromLabel))
					addedNodes[fromID] = fromLabel
					if verbose {
						utils.LogVerbose("Added source node from edge: %s", fromID)
					}
				}
			}
		}

		if _, exists := addedNodes[toID]; !exists {
			toLabel := CleanLabel(graph.Nodes.Lookup[edge.Dst].Attrs["label"])
			if toLabel != "" {
				if resourcesOnly && !isResourceLabel(toLabel) {
					if verbose {
						utils.LogVerbose("Skipping edge destination (non-resource) due to resourcesOnly: %s", toID)
					}
					// Skip adding this node; if either endpoint is non-resource, skip the edge later
				} else {
					sb.WriteString(fmt.Sprintf("        %s[\"%s\"]\n", toID, toLabel))
					addedNodes[toID] = toLabel
					if verbose {
						utils.LogVerbose("Added destination node from edge: %s", toID)
					}
				}
			}
		}

		// If resourcesOnly, include the edge only if both endpoints are resources
		if resourcesOnly {
			fromLabel := CleanLabel(graph.Nodes.Lookup[edge.Src].Attrs["label"])
			toLabel := CleanLabel(graph.Nodes.Lookup[edge.Dst].Attrs["label"])
			if !isResourceLabel(fromLabel) || !isResourceLabel(toLabel) {
				if verbose {
					utils.LogVerbose("Skipping edge due to non-resource endpoint(s): %s --> %s", fromID, toID)
				}
				continue
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
