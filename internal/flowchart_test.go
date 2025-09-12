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
