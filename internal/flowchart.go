// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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

// init initializes package-level resource type customization from environment variables.
// It sets `customResourceTypeMatcher` from `TERRAMAID_RESOURCE_TYPE_REGEX` if the value is a valid
// regular expression (invalid patterns are ignored). It also populates
// `customResourceTypePrefixes` from `TERRAMAID_RESOURCE_TYPE_PREFIXES` by splitting on commas,
// trimming whitespace, and discarding empty entries.
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

// FilterConfig holds the configuration for filtering resources in the diagram
type FilterConfig struct {
	IncludeTypes     []string // Include only these resource types (supports glob patterns)
	ExcludeTypes     []string // Exclude these resource types (supports glob patterns)
	IncludeProviders []string // Include only resources from these providers
	ExcludeModules   []string // Exclude resources from these modules
}

// IsEmpty returns true if no filters are configured
func (f *FilterConfig) IsEmpty() bool {
	return len(f.IncludeTypes) == 0 &&
		len(f.ExcludeTypes) == 0 &&
		len(f.IncludeProviders) == 0 &&
		len(f.ExcludeModules) == 0
}

// parseLabelComponents extracts the module path, resource type, and provider from a Terraform label
// Label formats:
// - <type>.<name> (e.g., aws_instance.web)
// - module.<mod>.<type>.<name> (e.g., module.vpc.aws_subnet.private)
// parseLabelComponents parses a Terraform-style node label and returns the module path,
// resource type, and provider.
//
// The function recognizes labels with leading module segments (e.g. "module.foo.module.bar.resource_type.name")
// and returns the module path as dot-joined module names ("foo.bar"). For data sources of the form
// "data.<type>.<name>" it returns the `<type>` as the resource type. For regular resources it returns
// the first segment as the resource type (e.g. "aws_instance"). The provider is inferred as the prefix
// before the first underscore in the resource type (e.g. "aws" from "aws_instance").
//
// Returns:
//   - modulePath: dot-joined module names extracted from leading "module" segments, or empty if none.
//   - resourceType: resource type or data source type if determinable, otherwise empty.
//   - provider: prefix before '_' in resourceType, or empty if not present.
func parseLabelComponents(label string) (modulePath string, resourceType string, provider string) {
	if label == "" {
		return "", "", ""
	}

	parts := strings.Split(label, ".")
	if len(parts) < 2 {
		return "", "", ""
	}

	// Collect module path segments
	var moduleSegments []string
	for len(parts) >= 2 && parts[0] == "module" {
		moduleSegments = append(moduleSegments, parts[1])
		parts = parts[2:]
	}
	modulePath = strings.Join(moduleSegments, ".")

	if len(parts) < 2 {
		return modulePath, "", ""
	}

	// Handle data sources: data.<type>.<name>
	if parts[0] == "data" {
		if len(parts) >= 3 {
			resourceType = parts[1]
		}
	} else {
		// Regular resource: <type>.<name>
		resourceType = parts[0]
	}

	// Extract provider from resource type (e.g., "aws" from "aws_instance")
	if idx := strings.Index(resourceType, "_"); idx > 0 {
		provider = resourceType[:idx]
	}

	return modulePath, resourceType, provider
}

// matchesGlobPattern reports whether s matches the glob pattern.
// It returns false when the pattern is invalid or when no match is found.
func matchesGlobPattern(s string, pattern string) bool {
	// Use filepath.Match for glob matching
	matched, err := filepath.Match(pattern, s)
	if err != nil {
		return false
	}
	return matched
}

// matchesAnyPattern reports whether s matches any of the provided glob patterns.
// It returns true if at least one pattern matches s, false otherwise.
func matchesAnyPattern(s string, patterns []string) bool {
	for _, pattern := range patterns {
		if matchesGlobPattern(s, pattern) {
			return true
		}
	}
	return false
}

// ShouldInclude determines if a resource label should be included based on the filter configuration
func (f *FilterConfig) ShouldInclude(label string, verbose bool) bool {
	// If no filters are configured, include everything
	if f.IsEmpty() {
		return true
	}

	modulePath, resourceType, provider := parseLabelComponents(label)

	// Check module exclusions first (exclusions take priority)
	if len(f.ExcludeModules) > 0 && modulePath != "" {
		for _, excludeModule := range f.ExcludeModules {
			// Check if the module path contains or matches the excluded module
			if strings.Contains(modulePath, excludeModule) || matchesGlobPattern(modulePath, excludeModule) {
				if verbose {
					utils.LogVerbose("Excluding %s: module %s matches exclude pattern %s", label, modulePath, excludeModule)
				}
				return false
			}
		}
	}

	// Check type exclusions
	if len(f.ExcludeTypes) > 0 && resourceType != "" {
		if matchesAnyPattern(resourceType, f.ExcludeTypes) {
			if verbose {
				utils.LogVerbose("Excluding %s: type %s matches exclude pattern", label, resourceType)
			}
			return false
		}
	}

	// Check provider inclusions (if specified, only these providers are allowed)
	if len(f.IncludeProviders) > 0 {
		if provider == "" {
			if verbose {
				utils.LogVerbose("Excluding %s: no provider detected and provider filter is active", label)
			}
			return false
		}
		providerMatch := false
		for _, includeProvider := range f.IncludeProviders {
			if strings.EqualFold(provider, includeProvider) {
				providerMatch = true
				break
			}
		}
		if !providerMatch {
			if verbose {
				utils.LogVerbose("Excluding %s: provider %s not in include list", label, provider)
			}
			return false
		}
	}

	// Check type inclusions (if specified, only these types are allowed)
	if len(f.IncludeTypes) > 0 {
		if resourceType == "" {
			if verbose {
				utils.LogVerbose("Excluding %s: no resource type detected and type filter is active", label)
			}
			return false
		}
		if !matchesAnyPattern(resourceType, f.IncludeTypes) {
			if verbose {
				utils.LogVerbose("Excluding %s: type %s not in include list", label, resourceType)
			}
			return false
		}
	}

	return true
}

// CleanID removes inline annotations and provider wrappers from an identifier, replaces dot and path separators with underscores, and returns a sanitized Mermaid-compatible identifier.
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

// sanitizeMermaidID cleans and normalizes a string for use as a Mermaid node ID.
// It replaces characters matched by mermaidUnsafeChars with '_' , collapses
// consecutive underscores, trims leading/trailing underscores, and ensures the
// result is non-empty and begins with a letter or underscore by prefixing
// "node_" when necessary.
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

// CleanLabel returns a human-friendly node label suitable for rendering.
// It removes inline annotations and surrounding quotes, converts `provider[...]` to `provider: ...`, and strips backslashes.
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
// isResourceLabel reports whether the given Terraform-style label represents a resource node.
//
// It returns true for labels that look like Terraform resource or data addresses and false for provider,
// module or other non-resource labels. The function ignores empty labels and labels beginning with
// "provider:". Leading module segments ("module.<name>") are stripped before analysis. For data
// addresses ("data.<type>.<name>") the type segment is inspected; for other addresses the type is taken
// from the canonical position in the address. Detection is guided by an optional user-supplied regex
// matcher and optional user prefixes, falls back to a set of built-in provider prefixes, and finally
// treats any type segment containing an underscore ('_') as a resource.
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

// GenerateMermaidFlowchart generates a Mermaid flowchart diagram from a gographviz graph.
// It validates the layout direction (must be one of TB, TD, BT, RL, LR) and returns an error for invalid directions.
// The output may include an optional named subgraph, can be limited to Terraform resource-like nodes when resourcesOnly is true, and is filtered by the provided FilterConfig (a nil filter is treated as empty).
// When verbose is true the function emits progress messages via the utils logger.
// It returns the complete Mermaid diagram as a string or an error if validation fails.
func GenerateMermaidFlowchart(ctx context.Context, graph *gographviz.Graph, direction string, subgraphName string, resourcesOnly bool, filter *FilterConfig, verbose bool) (string, error) {
	validDirections := map[string]bool{"TB": true, "TD": true, "BT": true, "RL": true, "LR": true}
	if !validDirections[direction] {
		return "", fmt.Errorf("invalid direction %s: valid options are TB, TD, BT, RL, LR", direction)
	}

	// Initialize filter if nil to avoid nil pointer checks
	if filter == nil {
		filter = &FilterConfig{}
	}

	if verbose {
		utils.LogVerbose("Generating Mermaid flowchart with direction: %s", direction)
		if subgraphName != "" {
			utils.LogVerbose("Using subgraph name: %s", subgraphName)
		}
		if !filter.IsEmpty() {
			utils.LogVerbose("Filter configuration active:")
			if len(filter.IncludeTypes) > 0 {
				utils.LogVerbose("  - Include types: %v", filter.IncludeTypes)
			}
			if len(filter.ExcludeTypes) > 0 {
				utils.LogVerbose("  - Exclude types: %v", filter.ExcludeTypes)
			}
			if len(filter.IncludeProviders) > 0 {
				utils.LogVerbose("  - Include providers: %v", filter.IncludeProviders)
			}
			if len(filter.ExcludeModules) > 0 {
				utils.LogVerbose("  - Exclude modules: %v", filter.ExcludeModules)
			}
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

		// Apply custom filters
		if !filter.ShouldInclude(nodeLabel, verbose) {
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
		fromLabel := CleanLabel(graph.Nodes.Lookup[edge.Src].Attrs["label"])
		toLabel := CleanLabel(graph.Nodes.Lookup[edge.Dst].Attrs["label"])

		// Check if source node should be included
		fromIncluded := true
		if fromLabel != "" {
			if resourcesOnly && !isResourceLabel(fromLabel) {
				fromIncluded = false
				if verbose {
					utils.LogVerbose("Skipping edge source (non-resource) due to resourcesOnly: %s", fromID)
				}
			} else if !filter.ShouldInclude(fromLabel, verbose) {
				fromIncluded = false
			}
		}

		// Check if destination node should be included
		toIncluded := true
		if toLabel != "" {
			if resourcesOnly && !isResourceLabel(toLabel) {
				toIncluded = false
				if verbose {
					utils.LogVerbose("Skipping edge destination (non-resource) due to resourcesOnly: %s", toID)
				}
			} else if !filter.ShouldInclude(toLabel, verbose) {
				toIncluded = false
			}
		}

		// Add source node if not already added and passes filters
		if _, exists := addedNodes[fromID]; !exists && fromIncluded && fromLabel != "" {
			sb.WriteString(fmt.Sprintf("        %s[\"%s\"]\n", fromID, fromLabel))
			addedNodes[fromID] = fromLabel
			if verbose {
				utils.LogVerbose("Added source node from edge: %s", fromID)
			}
		}

		// Add destination node if not already added and passes filters
		if _, exists := addedNodes[toID]; !exists && toIncluded && toLabel != "" {
			sb.WriteString(fmt.Sprintf("        %s[\"%s\"]\n", toID, toLabel))
			addedNodes[toID] = toLabel
			if verbose {
				utils.LogVerbose("Added destination node from edge: %s", toID)
			}
		}

		// Skip edge if either endpoint is filtered out
		if !fromIncluded || !toIncluded {
			if verbose {
				utils.LogVerbose("Skipping edge due to filtered endpoint(s): %s --> %s", fromID, toID)
			}
			continue
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