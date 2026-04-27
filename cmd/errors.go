// Copyright RoseSecurity 2024, 2026
// SPDX-License-Identifier: Apache-2.0

package cmd

import "errors"

var (
	errCheckTerraformFiles       = errors.New("error checking Terraform files in directory")
	errTerraformFilesDoNotExist  = errors.New("terraform files do not exist in directory")
	errTerraformDirectoryMissing = errors.New("terraform directory does not exist")
	errFetchVersionHTTPStatus    = errors.New("failed to fetch version")
)
