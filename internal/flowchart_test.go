// Copyright (c) RoseSecurity
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"testing"
)

func TestCleanID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic node name",
			input:    "module_example",
			expected: "module_example",
		},
		{
			name:     "node with parentheses",
			input:    "module_eks_public_web_module_self_managed_node_group_var_platform (validation)",
			expected: "module_eks_public_web_module_self_managed_node_group_var_platform_validation",
		},
		{
			name:     "node with brackets",
			input:    "provider[registry.terraform.io/hashicorp/aws]",
			expected: "provider_registry_terraform_io_hashicorp_aws",
		},
		{
			name:     "node with dots and slashes",
			input:    "module.foo/bar.baz",
			expected: "module_foo_bar_baz",
		},
		{
			name:     "node with special characters",
			input:    "node-name:with@special#chars",
			expected: "node_name_with_special_chars",
		},
		{
			name:     "node with multiple consecutive spaces",
			input:    "node    with    spaces",
			expected: "node_with_spaces",
		},
		{
			name:     "empty node",
			input:    "",
			expected: "node_",
		},
		{
			name:     "node starting with number",
			input:    "123node",
			expected: "node_123node",
		},
		{
			name:     "complex case from bug report",
			input:    "module_eks_public_web_module_self_managed_node_group_var_platform        (validation)    mod",
			expected: "module_eks_public_web_module_self_managed_node_group_var_platform_validation_mod",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanID(tt.input)
			if result != tt.expected {
				t.Errorf("CleanID(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeMermaidID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "remove problematic characters",
			input:    "node(with)special[chars]{and}more<stuff>",
			expected: "node_with_special_chars_and_more_stuff",
		},
		{
			name:     "handle multiple underscores",
			input:    "node___with___many___underscores",
			expected: "node_with_many_underscores",
		},
		{
			name:     "trim underscores",
			input:    "___node_name___",
			expected: "node_name",
		},
		{
			name:     "handle empty string",
			input:    "",
			expected: "node_",
		},
		{
			name:     "handle string that becomes empty after sanitization",
			input:    "()[]{}",
			expected: "node_",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeMermaidID(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeMermaidID(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseLabelComponents(t *testing.T) {
	tests := []struct {
		name             string
		label            string
		wantModulePath   string
		wantResourceType string
		wantProvider     string
	}{
		{
			name:             "empty label",
			label:            "",
			wantModulePath:   "",
			wantResourceType: "",
			wantProvider:     "",
		},
		{
			name:             "single part label",
			label:            "invalid",
			wantModulePath:   "",
			wantResourceType: "",
			wantProvider:     "",
		},
		{
			name:             "simple resource",
			label:            "aws_instance.web",
			wantModulePath:   "",
			wantResourceType: "aws_instance",
			wantProvider:     "aws",
		},
		{
			name:             "resource with single module",
			label:            "module.vpc.aws_subnet.private",
			wantModulePath:   "vpc",
			wantResourceType: "aws_subnet",
			wantProvider:     "aws",
		},
		{
			name:             "resource with nested modules",
			label:            "module.network.module.subnets.aws_subnet.main",
			wantModulePath:   "network.subnets",
			wantResourceType: "aws_subnet",
			wantProvider:     "aws",
		},
		{
			name:             "data source",
			label:            "data.aws_ami.latest",
			wantModulePath:   "",
			wantResourceType: "aws_ami",
			wantProvider:     "aws",
		},
		{
			name:             "data source with module",
			label:            "module.images.data.aws_ami.ubuntu",
			wantModulePath:   "images",
			wantResourceType: "aws_ami",
			wantProvider:     "aws",
		},
		{
			name:             "google provider resource",
			label:            "google_compute_instance.default",
			wantModulePath:   "",
			wantResourceType: "google_compute_instance",
			wantProvider:     "google",
		},
		{
			name:             "azurerm provider resource",
			label:            "azurerm_resource_group.main",
			wantModulePath:   "",
			wantResourceType: "azurerm_resource_group",
			wantProvider:     "azurerm",
		},
		{
			name:             "resource type without underscore",
			label:            "random.value",
			wantModulePath:   "",
			wantResourceType: "random",
			wantProvider:     "",
		},
		{
			name:             "module only (incomplete)",
			label:            "module.foo",
			wantModulePath:   "foo",
			wantResourceType: "",
			wantProvider:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotModule, gotType, gotProvider := parseLabelComponents(tt.label)
			if gotModule != tt.wantModulePath {
				t.Errorf("parseLabelComponents(%q) modulePath = %q, want %q", tt.label, gotModule, tt.wantModulePath)
			}
			if gotType != tt.wantResourceType {
				t.Errorf("parseLabelComponents(%q) resourceType = %q, want %q", tt.label, gotType, tt.wantResourceType)
			}
			if gotProvider != tt.wantProvider {
				t.Errorf("parseLabelComponents(%q) provider = %q, want %q", tt.label, gotProvider, tt.wantProvider)
			}
		})
	}
}

func TestMatchesGlobPattern(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		pattern string
		want    bool
	}{
		{
			name:    "exact match",
			s:       "aws_instance",
			pattern: "aws_instance",
			want:    true,
		},
		{
			name:    "wildcard prefix",
			s:       "aws_instance",
			pattern: "aws_*",
			want:    true,
		},
		{
			name:    "wildcard suffix",
			s:       "aws_instance",
			pattern: "*_instance",
			want:    true,
		},
		{
			name:    "wildcard middle",
			s:       "aws_security_group",
			pattern: "aws_*_group",
			want:    true,
		},
		{
			name:    "no match",
			s:       "google_compute_instance",
			pattern: "aws_*",
			want:    false,
		},
		{
			name:    "question mark wildcard",
			s:       "aws_s3_bucket",
			pattern: "aws_s?_bucket",
			want:    true,
		},
		{
			name:    "invalid pattern returns false",
			s:       "test",
			pattern: "[invalid",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesGlobPattern(tt.s, tt.pattern)
			if got != tt.want {
				t.Errorf("matchesGlobPattern(%q, %q) = %v, want %v", tt.s, tt.pattern, got, tt.want)
			}
		})
	}
}

func TestMatchesAnyPattern(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		patterns []string
		want     bool
	}{
		{
			name:     "matches first pattern",
			s:        "aws_instance",
			patterns: []string{"aws_*", "google_*"},
			want:     true,
		},
		{
			name:     "matches second pattern",
			s:        "google_compute_instance",
			patterns: []string{"aws_*", "google_*"},
			want:     true,
		},
		{
			name:     "matches none",
			s:        "azurerm_resource_group",
			patterns: []string{"aws_*", "google_*"},
			want:     false,
		},
		{
			name:     "empty patterns",
			s:        "aws_instance",
			patterns: []string{},
			want:     false,
		},
		{
			name:     "single pattern match",
			s:        "aws_instance",
			patterns: []string{"aws_instance"},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesAnyPattern(tt.s, tt.patterns)
			if got != tt.want {
				t.Errorf("matchesAnyPattern(%q, %v) = %v, want %v", tt.s, tt.patterns, got, tt.want)
			}
		})
	}
}

func TestFilterConfig_IsEmpty(t *testing.T) {
	tests := []struct {
		name   string
		filter FilterConfig
		want   bool
	}{
		{
			name:   "empty filter",
			filter: FilterConfig{},
			want:   true,
		},
		{
			name:   "with include types",
			filter: FilterConfig{IncludeTypes: []string{"aws_*"}},
			want:   false,
		},
		{
			name:   "with exclude types",
			filter: FilterConfig{ExcludeTypes: []string{"aws_*"}},
			want:   false,
		},
		{
			name:   "with include providers",
			filter: FilterConfig{IncludeProviders: []string{"aws"}},
			want:   false,
		},
		{
			name:   "with exclude modules",
			filter: FilterConfig{ExcludeModules: []string{"vpc"}},
			want:   false,
		},
		{
			name: "with all filters",
			filter: FilterConfig{
				IncludeTypes:     []string{"aws_*"},
				ExcludeTypes:     []string{"aws_iam_*"},
				IncludeProviders: []string{"aws"},
				ExcludeModules:   []string{"legacy"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.IsEmpty()
			if got != tt.want {
				t.Errorf("FilterConfig.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterConfig_ShouldInclude(t *testing.T) {
	tests := []struct {
		name   string
		filter FilterConfig
		label  string
		want   bool
	}{
		{
			name:   "empty filter includes everything",
			filter: FilterConfig{},
			label:  "aws_instance.web",
			want:   true,
		},
		{
			name:   "include types - matching",
			filter: FilterConfig{IncludeTypes: []string{"aws_instance", "aws_s3_bucket"}},
			label:  "aws_instance.web",
			want:   true,
		},
		{
			name:   "include types - not matching",
			filter: FilterConfig{IncludeTypes: []string{"aws_s3_bucket"}},
			label:  "aws_instance.web",
			want:   false,
		},
		{
			name:   "include types with glob - matching",
			filter: FilterConfig{IncludeTypes: []string{"aws_*"}},
			label:  "aws_instance.web",
			want:   true,
		},
		{
			name:   "include types with glob - not matching",
			filter: FilterConfig{IncludeTypes: []string{"google_*"}},
			label:  "aws_instance.web",
			want:   false,
		},
		{
			name:   "exclude types - matching",
			filter: FilterConfig{ExcludeTypes: []string{"aws_instance"}},
			label:  "aws_instance.web",
			want:   false,
		},
		{
			name:   "exclude types - not matching",
			filter: FilterConfig{ExcludeTypes: []string{"aws_s3_bucket"}},
			label:  "aws_instance.web",
			want:   true,
		},
		{
			name:   "exclude types with glob",
			filter: FilterConfig{ExcludeTypes: []string{"aws_iam_*"}},
			label:  "aws_iam_role.admin",
			want:   false,
		},
		{
			name:   "include providers - matching",
			filter: FilterConfig{IncludeProviders: []string{"aws"}},
			label:  "aws_instance.web",
			want:   true,
		},
		{
			name:   "include providers - not matching",
			filter: FilterConfig{IncludeProviders: []string{"google"}},
			label:  "aws_instance.web",
			want:   false,
		},
		{
			name:   "include providers - case insensitive",
			filter: FilterConfig{IncludeProviders: []string{"AWS"}},
			label:  "aws_instance.web",
			want:   true,
		},
		{
			name:   "exclude modules - matching",
			filter: FilterConfig{ExcludeModules: []string{"vpc"}},
			label:  "module.vpc.aws_subnet.private",
			want:   false,
		},
		{
			name:   "exclude modules - not matching",
			filter: FilterConfig{ExcludeModules: []string{"database"}},
			label:  "module.vpc.aws_subnet.private",
			want:   true,
		},
		{
			name:   "exclude modules - nested module match",
			filter: FilterConfig{ExcludeModules: []string{"subnets"}},
			label:  "module.network.module.subnets.aws_subnet.main",
			want:   false,
		},
		{
			name:   "exclude modules - glob pattern",
			filter: FilterConfig{ExcludeModules: []string{"legacy*"}},
			label:  "module.legacy_vpc.aws_subnet.old",
			want:   false,
		},
		{
			name:   "data source with include types",
			filter: FilterConfig{IncludeTypes: []string{"aws_ami"}},
			label:  "data.aws_ami.latest",
			want:   true,
		},
		{
			name:   "data source with exclude types",
			filter: FilterConfig{ExcludeTypes: []string{"aws_ami"}},
			label:  "data.aws_ami.latest",
			want:   false,
		},
		{
			name:   "combined filters - include type and provider both match",
			filter: FilterConfig{IncludeTypes: []string{"aws_*"}, IncludeProviders: []string{"aws"}},
			label:  "aws_instance.web",
			want:   true,
		},
		{
			name:   "combined filters - exclude takes priority",
			filter: FilterConfig{IncludeTypes: []string{"aws_*"}, ExcludeTypes: []string{"aws_instance"}},
			label:  "aws_instance.web",
			want:   false,
		},
		{
			name:   "resource without module not affected by module exclusion",
			filter: FilterConfig{ExcludeModules: []string{"vpc"}},
			label:  "aws_instance.web",
			want:   true,
		},
		{
			name:   "empty label with type filter excludes",
			filter: FilterConfig{IncludeTypes: []string{"aws_*"}},
			label:  "",
			want:   false,
		},
		{
			name:   "no resource type detected with type filter active",
			filter: FilterConfig{IncludeTypes: []string{"aws_*"}},
			label:  "invalid",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.ShouldInclude(tt.label, false)
			if got != tt.want {
				t.Errorf("FilterConfig.ShouldInclude(%q) = %v, want %v", tt.label, got, tt.want)
			}
		})
	}
}
