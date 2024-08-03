package utils

import (
	"os"
	"path/filepath"
)

// Check if a directory exists
func DirExists(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}

	return true
}

// Check if Terraform files exist in a directory
func TerraformFilesExist(dir string) bool {
	validExtensions := []string{".tf", ".tf.json", ".tftest.hcl", ".tftest.json", "terraform.tfvars", "terraform.tfvars.json"}

	found := false
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
		return false
	}

	return found
}
