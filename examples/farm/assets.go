package farm

import "embed"

// assets are the farm's sprites, drawn as SVG by cmd/farmart.
//
//go:generate go run ./cmd/farmart -assets assets -preview preview.png
//go:embed assets/*.svg
var assets embed.FS
