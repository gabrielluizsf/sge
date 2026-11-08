//go:build linux || js

package installer

import "testing"

// testPlace is a new home on Linux (and on Deno, which tests there).
func testPlace(t *testing.T) Place { return homeOn(t, Linux) }
