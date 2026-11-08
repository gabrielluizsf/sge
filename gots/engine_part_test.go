package gots

import (
	"os"
	"slices"
	"strings"
)

// enginePart is the engine packages this run tests on Deno: all of them,
// or those GOTS_ENGINE names, comma separated, or, when it starts with
// "-", all but those. CI splits the engine's slowest packages among jobs
// running at once this way.
func enginePart() []string {
	spec := os.Getenv("GOTS_ENGINE")
	if spec == "" {
		return enginePackages
	}
	but := strings.HasPrefix(spec, "-")
	named := strings.Split(strings.TrimPrefix(spec, "-"), ",")
	var part []string
	for _, p := range enginePackages {
		if slices.Contains(named, p) != but {
			part = append(part, p)
		}
	}
	return part
}
