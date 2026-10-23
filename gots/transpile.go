package gots

import (
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"path/filepath"
	"sort"
)

// Options say what to transpile and where to.
type Options struct {
	// Dir is the directory go list runs in.
	Dir string
	// Out is the directory the TypeScript is written to.
	Out string
	// Patterns are the packages to transpile, as go list takes them. Their
	// dependencies outside the standard library are transpiled too.
	Patterns []string
	// Tests also transpiles the packages' tests, as *_test.ts files that
	// run under deno test.
	Tests bool
}

// Transpile writes the TypeScript for the packages in opts and returns the
// files it wrote, relative to opts.Out.
func Transpile(opts Options) ([]string, error) {
	prog, err := Load(opts.Dir, opts.Tests, opts.Patterns...)
	if err != nil {
		return nil, err
	}
	lay := layout{module: prog.Module}
	var written []string
	var errs []error
	write := func(rel, src string) {
		dst := filepath.Join(opts.Out, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			errs = append(errs, err)
			return
		}
		if err := os.WriteFile(dst, []byte(src), 0o644); err != nil {
			errs = append(errs, err)
			return
		}
		written = append(written, rel)
	}
	needStd := map[string]bool{}
	for _, path := range prog.order {
		p := prog.Packages[path]
		if p.Std && !p.Transpile || p.Path == "unsafe" {
			continue
		}
		e := newEmitter(prog, p, p.Info, p.Types, lay)
		src := e.module(p.Files)
		errs = append(errs, e.errs...)
		write(lay.file(p), src)
		for _, imp := range p.Imports {
			if q := prog.Packages[imp]; q != nil && q.Std && !q.Transpile {
				needStd[imp] = true
			}
		}
	}
	if opts.Tests {
		for _, p := range prog.Roots {
			files, err := emitTests(prog, p, lay, needStd)
			errs = append(errs, err)
			for rel, src := range files {
				write(rel, src)
			}
		}
	}
	for p := range needStd {
		if unsupportedStd[p] != "" {
			write("_gots/std/"+p+"/mod.ts", stubModule(prog.Packages[p]))
		}
	}
	var missing []string
	for p := range needStd {
		if p != "unsafe" && !hasStd(p) && unsupportedStd[p] == "" {
			missing = append(missing, p)
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		errs = append(errs, fmt.Errorf("no TypeScript version of standard package(s) %v yet", missing))
	}
	if err := writeRuntime(opts.Out); err != nil {
		errs = append(errs, err)
	}
	return written, errors.Join(errs...)
}

// emitTests writes a package's tests: its own test files compiled together
// with the package, and its _test package on its own.
func emitTests(prog *Program, p *Package, lay layout, needStd map[string]bool) (map[string]string, error) {
	out := map[string]string{}
	var errs []error
	emit := func(name, path string, goFiles []string, withPkg bool) {
		var files []*ast.File
		if withPkg {
			files = append(files, p.Files...)
		}
		extra, err := prog.parse(p.Dir, goFiles)
		if err != nil {
			errs = append(errs, err)
			return
		}
		files = append(files, extra...)
		info := newInfo()
		tpkg, err := prog.typeCheck(path, files, info, false)
		if err != nil {
			errs = append(errs, err)
			return
		}
		for _, f := range extra {
			for _, imp := range f.Imports {
				if q := prog.Packages[importPath(imp)]; q != nil && q.Std && !q.Transpile {
					needStd[q.Path] = true
				}
			}
		}
		e := newEmitter(prog, p, info, tpkg, lay)
		out[filepath.ToSlash(filepath.Join(lay.dir(p), name))] = e.module(files)
		errs = append(errs, e.errs...)
		for _, s := range e.Skipped {
			fmt.Fprintf(os.Stderr, "gots: skipping a test function: %v\n", s)
		}
	}
	if len(p.TestFiles) > 0 {
		emit("mod_test.ts", p.Path, p.TestFiles, true)
	}
	if len(p.XTestFiles) > 0 {
		emit("mod_x_test.ts", p.Path+"_test", p.XTestFiles, false)
	}
	return out, errors.Join(errs...)
}

func importPath(imp *ast.ImportSpec) string {
	s := imp.Path.Value
	return s[1 : len(s)-1]
}

var _ = types.Universe
