package gots

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// listed is the part of `go list -json` output the loader reads.
type listed struct {
	ImportPath   string
	Name         string
	Dir          string
	Standard     bool
	GoFiles      []string
	TestGoFiles  []string
	XTestGoFiles []string
	Imports      []string
	Module       *struct{ Path string }
	TestImports  []string
	XTestImports []string
	Error        *struct{ Err string }
}

// Program is every package the requested ones need, checked against each
// other so that a type means the same object wherever it appears.
type Program struct {
	Fset     *token.FileSet
	Packages map[string]*Package
	// Roots are the packages that were asked for, in the order go list
	// gave them.
	Roots []*Package
	// Module is the module the roots belong to.
	Module string
	order  []string
}

// buildEnv is the environment packages are listed and checked under: js/wasm
// is the GOOS/GOARCH whose build tags select the code that has no assembly
// and no cgo.
func buildEnv() []string {
	return append(os.Environ(), "GOOS=js", "GOARCH=wasm", "CGO_ENABLED=0")
}

// Load lists the packages matching patterns from dir with their
// dependencies, then parses and type-checks them all. withTests also loads
// what the packages' test files import.
func Load(dir string, withTests bool, patterns ...string) (*Program, error) {
	roots, err := goList(dir, false, false, patterns)
	if err != nil {
		return nil, err
	}
	all, err := goList(dir, true, withTests, patterns)
	if err != nil {
		return nil, err
	}
	prog := &Program{Fset: token.NewFileSet(), Packages: map[string]*Package{}}
	for _, l := range all {
		if l.Error != nil {
			return nil, fmt.Errorf("%s: %s", l.ImportPath, l.Error.Err)
		}
		if strings.Contains(l.ImportPath, " [") || strings.HasSuffix(l.ImportPath, ".test") {
			continue
		}
		p := &Package{
			Path: l.ImportPath, Name: l.Name, Dir: l.Dir, Std: l.Standard,
			GoFiles: l.GoFiles, TestFiles: l.TestGoFiles, XTestFiles: l.XTestGoFiles,
			Imports: l.Imports, fset: prog.Fset,
		}
		prog.Packages[p.Path] = p
		prog.order = append(prog.order, p.Path)
	}
	for _, path := range prog.order {
		if err := prog.check(prog.Packages[path]); err != nil {
			return nil, err
		}
	}
	for _, r := range roots {
		if p := prog.Packages[r.ImportPath]; p != nil {
			prog.Roots = append(prog.Roots, p)
		}
		if r.Module != nil && prog.Module == "" {
			prog.Module = r.Module.Path
		}
	}
	return prog, nil
}

// goList runs go list on patterns, with their dependencies when deps is set
// and with what their tests import when tests is set too.
func goList(dir string, deps, tests bool, patterns []string) ([]listed, error) {
	args := []string{"list", "-e", "-json"}
	if deps {
		args = append(args, "-deps")
	}
	if tests {
		args = append(args, "-test")
	}
	args = append(args, patterns...)
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	cmd.Env = buildEnv()
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("go list: %v\n%s", err, stderr.String())
	}
	var pkgs []listed
	dec := json.NewDecoder(bytes.NewReader(out))
	for {
		var l listed
		if err := dec.Decode(&l); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		pkgs = append(pkgs, l)
	}
	return pkgs, nil
}

func (prog *Program) parse(dir string, names []string) ([]*ast.File, error) {
	var files []*ast.File
	for _, n := range names {
		f, err := parser.ParseFile(prog.Fset, filepath.Join(dir, n), nil, parser.ParseComments|parser.SkipObjectResolution)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

// check parses and type-checks p's own files. go list orders dependencies
// first, so everything p imports is already checked.
func (prog *Program) check(p *Package) error {
	if p.Path == "unsafe" {
		p.Types = types.Unsafe
		return nil
	}
	files, err := prog.parse(p.Dir, p.GoFiles)
	if err != nil {
		return err
	}
	p.Files = files
	p.Info = newInfo()
	p.Types, err = prog.typeCheck(p.Path, files, p.Info, p.Std)
	return err
}

// typeCheck checks files as the package path. Function bodies are skipped
// for packages that are only needed for their declarations.
func (prog *Program) typeCheck(path string, files []*ast.File, info *types.Info, declsOnly bool) (*types.Package, error) {
	var errs []error
	conf := types.Config{
		Importer:         importerFunc(prog.importPackage),
		Sizes:            types.SizesFor("gc", "wasm"),
		IgnoreFuncBodies: declsOnly,
		Error: func(err error) {
			if len(errs) < 10 {
				errs = append(errs, err)
			}
		},
	}
	pkg, _ := conf.Check(path, prog.Fset, files, info)
	if len(errs) > 0 {
		return nil, fmt.Errorf("type-checking %s: %v", path, errs)
	}
	return pkg, nil
}

func (prog *Program) importPackage(path string) (*types.Package, error) {
	if path == "unsafe" {
		return types.Unsafe, nil
	}
	p := prog.Packages[path]
	if p == nil || p.Types == nil {
		return nil, fmt.Errorf("package %s was not loaded", path)
	}
	return p.Types, nil
}

type importerFunc func(path string) (*types.Package, error)

func (f importerFunc) Import(path string) (*types.Package, error) { return f(path) }
