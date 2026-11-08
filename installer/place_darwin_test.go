//go:build darwin

package installer

import "testing"

// testPlace is a new home on macOS.
func testPlace(t *testing.T) Place { return homeOn(t, MacOS) }
