// Command gots transpiles Go packages to TypeScript for Deno, and, with
// -web, a program to a page for a browser.
//
//	gots [-o dir] [-test] [-index] [-web dir] packages...
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
	index := flag.Bool("index", false, "write a mod.ts exporting the packages by name")
	web := flag.String("web", "", "also write each program to this directory as a page for a browser")
	flag.Parse()
	patterns := flag.Args()
	if len(patterns) == 0 {
		patterns = []string{"."}
	}
	files, err := gots.Transpile(gots.Options{Dir: ".", Out: *out, Patterns: patterns, Tests: *tests, Index: *index, Web: *web})
	for _, f := range files {
		fmt.Println(f)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
