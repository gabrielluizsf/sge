package farm

import "embed"

// assets are the plain look's sprites, kept for the package's tests.
//
//go:embed assets/*.svg
var assets embed.FS
