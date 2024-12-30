// Copyright (c) RoseSecurity
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/mattn/go-colorable"
)

const (
	ColorReset     = "\033[0m"
	ColorGreen     = "\033[32m"
	ColorYellow    = "\033[33m"
	ColorRed       = "\033[31m"
	ColorBold      = "\033[1m"
	ColorUnderline = "\033[4m"
)

type Spinner struct {
	s *spinner.Spinner
}

// Check if a directory exists
func DirExists(dir string) bool {
	_, err := os.Stat(dir)
	return !os.IsNotExist(err)
}

// Check if Terraform files exist in a directory
func TerraformFilesExist(dir string) (bool, error) {
	validExtensions := []string{".tf", ".tf.json", ".tftest.hcl", ".tftest.json", "terraform.tfvars", "terraform.tfvars.json"}

	var found bool
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		for _, ext := range validExtensions {
			if filepath.Ext(path) == ext || filepath.Base(path) == ext {
				found = true
				return filepath.SkipDir
			}
		}
		return nil
	})
	if err != nil {
		return false, err
	}

	return found, nil
}

// Initialize a new spinner
func NewSpinner(text string) *Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Color("blue")
	s.Writer = colorable.NewColorableStdout() // Ensure colors are supported on Windows
	s.Suffix = " " + text
	return &Spinner{s: s}
}

// Start the spinner
func (sp *Spinner) Start() {
	fmt.Printf("%s%s%s ", ColorBold+ColorGreen, sp.s.Suffix, ColorReset)
	sp.s.Start()
}

// Stop the spinner
func (sp *Spinner) Stop() {
	sp.s.Stop()
}
