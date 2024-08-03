package main

import (
	"github.com/RoseSecurity/terramaid/cmd"
	u "github.com/RoseSecurity/terramaid/pkg/utils"
)

var version string

func main() {
	cmd.Version = version

	if err := cmd.Execute(); err != nil {
		u.LogErrorAndExit(err)
	}
}
