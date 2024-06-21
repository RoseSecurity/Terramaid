package main

import (
	"github.com/RoseSecurity/terramaid/cmd"
	u "github.com/RoseSecurity/terramaid/pkg/utils"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		u.LogErrorAndExit(err)
	}
}
