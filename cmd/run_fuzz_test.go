package cmd_test

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/RoseSecurity/terramaid/cmd"
)

func FuzzRun(f *testing.F) {
	f.Add("--working-dir=. --direction=TD")
	f.Add("--chart-type=flowchart --verbose")

	f.Fuzz(func(t *testing.T, flagLine string) {
		args := strings.Fields(flagLine)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		os.Setenv("TERRAMAID_TF_BINARY", "/bin/true")

		tmpDir, _ := os.MkdirTemp("", "tf")
		_ = os.WriteFile(tmpDir+"/main.tf", []byte("terraform {}"), 0o644)
		defer os.RemoveAll(tmpDir)

		hasWD := false
		for i, a := range args {
			if strings.HasPrefix(a, "--working-dir") {
				args[i] = "--working-dir=" + tmpDir
				hasWD = true
			}
		}
		if !hasWD {
			args = append(args, "--working-dir="+tmpDir)
		}

		c := cmd.NewRunCmd()
		c.SetArgs(args)
		c.SetOut(io.Discard)
		c.SetErr(io.Discard)
		_ = c.ExecuteContext(ctx)
	})
}
