package main

import (
	"github.com/RoseSecurity/terramaid/cmd"
	u "github.com/RoseSecurity/terramaid/pkg/utils"
)

func main() {
	if err := cmd.Execute(); err != nil {
		u.LogErrorAndExit(err)
	}
}
