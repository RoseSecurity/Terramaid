// Copyright (c) RoseSecurity
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

const (
	LogLevelTrace   = "Trace"
	LogLevelDebug   = "Debug"
	LogLevelInfo    = "Info"
	LogLevelWarning = "Warning"
)

// LogErrorAndExit logs errors to std.Error and exits with an error code
func LogErrorAndExit(err error) {
	if err != nil {
		LogError(err)

		// Find the executed command's exit code from the error
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exitCode := exitError.ExitCode()
			os.Exit(exitCode)
		}
	}
}

// LogError logs errors to std.Error
func LogError(err error) {
	if err != nil {
		c := color.New(color.FgRed)
		_, err2 := c.Fprintln(color.Error, err.Error()+"\n")
		if err2 != nil {
			color.Red("Error logging the error:")
			color.Red("%s\n", err2)
			color.Red("Original error:")
			color.Red("%s\n", err)
		}
	}
}

// LogVerbose logs messages in verbose mode
func LogVerbose(format string, a ...any) {
	c := color.New(color.FgBlue)
	message := fmt.Sprintf(format, a...)
	c.Fprintf(color.Output, "[VERBOSE] %s\n", message)
}
