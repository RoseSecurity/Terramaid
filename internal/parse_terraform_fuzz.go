//go:build fuzz
// +build fuzz

package internal

import "context"

// Lightweight stubs that satisfy the production signatures but avoid expensive
// disk I/O and external terraform invocations when the "fuzz" build tag is on.

func ParseTerraform(_ context.Context, _ string, _ string, _ string, _ bool) (interface{}, error) {
	return struct{}{}, nil
}

func GenerateMermaidFlowchart(_ context.Context, _ interface{}, _ string, _ string, _ bool) (string, error) {
	return "graph TD;", nil
}
