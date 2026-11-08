package installer

import "strings"

// equalFold is whether two paths in the files are the same on a system
// whose names ignore case, as Windows' and macOS's do.
func equalFold(a, b string) bool {
	return strings.EqualFold(a, b)
}
