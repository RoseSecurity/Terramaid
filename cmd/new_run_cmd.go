package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// NewRunCmd returns a pristine copy of runCmd so that each test/fuzz
// iteration starts with clean flag state.
func NewRunCmd() *cobra.Command {
	clone := *runCmd // shallow‑copy the definition
	clone.Flags().VisitAll(func(f *pflag.Flag) { f.Value.Set(f.DefValue) })
	return &clone
}
