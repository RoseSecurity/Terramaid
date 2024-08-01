package main

import (
	"fmt"
	"os"

	"github.com/RoseSecurity/terramaid/cmd"
)

var version string

func main() {
	cmd.Version = version 

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)

		os.Exit(1)
	}
}
