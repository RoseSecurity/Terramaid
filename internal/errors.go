// SPDX-License-Identifier: Apache-2.0

package internal

import "errors"

var (
	errInvalidDirection     = errors.New("invalid direction")
	errNoTerraformGraphData = errors.New("no output from terraform graph")
)
