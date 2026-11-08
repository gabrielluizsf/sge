package gots

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// A module carries its package's doc as the module's, and each export
// and member Go documents its doc; plumbing Go does not export is
// ignored, and doc links become JSDoc links.
func TestGoDocsBecomeJSDoc(t *testing.T) {
	src := `// Package geo measures shapes.
package geo

// Box is a rectangle, see [Area].
type Box struct {
	// W is its width.
	W float64
	H float64
	name string
}

// Area is how much ground b covers.
func (b *Box) Area() float64 { return b.W * b.H }

// Unit is the box one by one.
var Unit = Box{W: 1, H: 1}
`
	f, err := parser.ParseFile(token.NewFileSet(), "geo.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	ts := `// header
import * as $ from "./rt.ts";
export class Box {
  W: number;
  H: number;
  name$s: string;
  get name(): string { return ""; }
  constructor(W: number = 0) {}
  Area(): number { return 0; }
}
export class Box$Val extends $.Val {
}
export function Box$$Area(b: Box): number { return 0; }
export let Unit: Box = new Box();
export function Unit$set(v: Box): void {}
`
	got := document(ts, "example.com/geo", docsOf([]*ast.File{f}))
	for _, want := range []string{
		"/**\n * Package geo measures shapes.\n *\n * @module\n */\nimport",
		"/**\n * Box is a rectangle, see {@link Area}.\n */\nexport class Box {",
		"  /**\n   * W is its width.\n   */\n  W: number;",
		"  /**\n   * H is a field of Box.\n   */\n  H: number;",
		"  /** @ignore */\n  name$s: string;",
		"  /** @ignore */\n  get name(): string",
		"  /**\n   * Makes a Box of the values",
		"  /**\n   * Area is how much ground b covers.\n   */\n  Area(): number",
		"/** @ignore */\nexport class Box$Val",
		"/**\n * Area is how much ground b covers.\n */\nexport function Box$$Area",
		"/**\n * Unit is the box one by one.\n */\nexport let Unit",
		"/** @ignore */\nexport function Unit$set",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing:\n%s\n\nin:\n%s", want, got)
		}
	}
}

// The engine's packages, transpiled, are documented as deno doc --lint
// wants: a module doc each, and a doc on every exported symbol, as JSR
// scores a package.
func TestEngineIsDocumented(t *testing.T) {
	t.Parallel()
	deno, err := exec.LookPath("deno")
	if err != nil {
		t.Skip("deno is not installed")
	}
	out := engine(t)
	files := []string{filepath.Join(out, "mod.ts")}
	for _, p := range enginePackages {
		files = append(files, filepath.Join(out, p, "mod.ts"))
	}
	got, _ := exec.Command(deno, append([]string{"doc", "--lint"}, files...)...).CombinedOutput()
	for _, line := range strings.Split(string(got), "\n") {
		if strings.Contains(line, "missing-jsdoc") {
			t.Errorf("deno doc --lint:\n%s", got)
			break
		}
	}
	if _, err := os.Stat(files[0]); err != nil {
		t.Fatal(err)
	}
}
