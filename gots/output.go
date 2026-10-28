package gots

import (
	"path"
	"strings"
)

// layout decides where each package's TypeScript goes inside the output
// directory, as slash-separated paths relative to it.
type layout struct {
	module string
}

const runtimeDir = "_gots/rt"

// dir is the directory a package's modules are written to.
func (l layout) dir(p *Package) string {
	switch {
	case p.Std:
		return "_gots/std/" + p.Path
	case p.Path == l.module:
		return "."
	case strings.HasPrefix(p.Path, l.module+"/"):
		return strings.TrimPrefix(p.Path, l.module+"/")
	default:
		return "_gots/mod/" + p.Path
	}
}

// file is the module a package is imported from.
func (l layout) file(p *Package) string { return path.Join(l.dir(p), "mod.ts") }

// rel is the import specifier for target from a module in fromDir.
func rel(fromDir, target string) string {
	from := strings.Split(path.Clean(fromDir), "/")
	to := strings.Split(path.Clean(target), "/")
	if from[0] == "." {
		from = nil
	}
	i := 0
	for i < len(from) && i < len(to)-1 && from[i] == to[i] {
		i++
	}
	up := strings.Repeat("../", len(from)-i)
	if up == "" {
		up = "./"
	}
	return up + strings.Join(to[i:], "/")
}
