// Command gots transpiles Go packages to TypeScript for Deno.
//
//	gots [-o dir] [-test] packages...
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gabrielluizsf/sge/gots"
)

func main() {
	out := flag.String("o", "deno", "directory to write the TypeScript to")
	tests := flag.Bool("test", false, "also transpile the packages' tests")
	flag.Parse()
	patterns := flag.Args()
	if len(patterns) == 0 {
		patterns = []string{"."}
	}
	files, err := gots.Transpile(gots.Options{Dir: ".", Out: *out, Patterns: patterns, Tests: *tests})
	for _, f := range files {
		fmt.Println(f)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
