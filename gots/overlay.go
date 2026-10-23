package gots

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path"
	"strings"
)

// overlay replaces declarations of a standard library package with the
// ones in runtime/overlay/<path>: Go written for the transpiler, for code
// that leans on unsafe, assembly or the runtime. A declaration is replaced
// whole, by name; overlay files may also add declarations of their own.
func (prog *Program) overlay(pkgPath string, files []*ast.File) ([]*ast.File, error) {
	dir := "runtime/overlay/" + pkgPath
	entries, err := fs.ReadDir(runtimeFS, dir)
	if err != nil {
		return files, nil
	}
	var extra []*ast.File
	replaced := map[string]bool{}
	for _, en := range entries {
		if path.Ext(en.Name()) != ".go" {
			continue
		}
		src, err := runtimeFS.ReadFile(dir + "/" + en.Name())
		if err != nil {
			return nil, err
		}
		f, err := parser.ParseFile(prog.Fset, "overlay/"+pkgPath+"/"+en.Name(), src, parser.ParseComments|parser.SkipObjectResolution)
		if err != nil {
			return nil, err
		}
		for _, d := range f.Decls {
			for _, n := range declNames(d) {
				replaced[n] = true
			}
		}
		extra = append(extra, f)
	}
	for _, f := range files {
		var kept []ast.Decl
		for _, d := range f.Decls {
			if g, ok := d.(*ast.GenDecl); ok && len(g.Specs) > 0 {
				var specs []ast.Spec
				for _, s := range g.Specs {
					if !anyReplaced(specNames(s), replaced) {
						specs = append(specs, s)
					}
				}
				if len(specs) == 0 {
					continue
				}
				g.Specs = specs
				kept = append(kept, g)
				continue
			}
			if !anyReplaced(declNames(d), replaced) {
				kept = append(kept, d)
			}
		}
		f.Decls = kept
		prog.pruneImports(f)
	}
	return append(files, extra...), nil
}

// pruneImports drops the imports of f that nothing in it uses any more,
// after an overlay took away the declarations that did.
func (prog *Program) pruneImports(f *ast.File) {
	used := map[string]bool{}
	ast.Inspect(f, func(n ast.Node) bool {
		if sel, ok := n.(*ast.SelectorExpr); ok {
			if id, ok := sel.X.(*ast.Ident); ok {
				used[id.Name] = true
			}
		}
		return true
	})
	var specs []ast.Spec
	var imports []*ast.ImportSpec
	for _, d := range f.Decls {
		g, ok := d.(*ast.GenDecl)
		if !ok || g.Tok != token.IMPORT {
			continue
		}
		specs = specs[:0]
		for _, sp := range g.Specs {
			is := sp.(*ast.ImportSpec)
			path := strings.Trim(is.Path.Value, `"`)
			name := path[strings.LastIndex(path, "/")+1:]
			if p := prog.Packages[path]; p != nil {
				name = p.Name
			}
			if path == "unsafe" {
				name = "unsafe"
			}
			if is.Name != nil {
				name = is.Name.Name
			}
			if name == "_" || name == "." || used[name] {
				specs = append(specs, is)
				imports = append(imports, is)
			}
		}
		g.Specs = append([]ast.Spec(nil), specs...)
	}
	f.Imports = imports
	decls := f.Decls[:0]
	for _, d := range f.Decls {
		if g, ok := d.(*ast.GenDecl); ok && g.Tok == token.IMPORT && len(g.Specs) == 0 {
			continue
		}
		decls = append(decls, d)
	}
	f.Decls = decls
}

func anyReplaced(names []string, replaced map[string]bool) bool {
	for _, n := range names {
		if replaced[n] {
			return true
		}
	}
	return false
}

// declNames are the names a top-level declaration declares, with methods
// written Type.Method.
func declNames(d ast.Decl) []string {
	switch d := d.(type) {
	case *ast.FuncDecl:
		if d.Recv == nil {
			return []string{d.Name.Name}
		}
		return []string{recvName(d.Recv.List[0].Type) + "." + d.Name.Name}
	case *ast.GenDecl:
		var out []string
		for _, s := range d.Specs {
			out = append(out, specNames(s)...)
		}
		return out
	}
	return nil
}

func specNames(s ast.Spec) []string {
	switch s := s.(type) {
	case *ast.TypeSpec:
		return []string{s.Name.Name}
	case *ast.ValueSpec:
		var out []string
		for _, n := range s.Names {
			if n.Name != "_" {
				out = append(out, n.Name)
			}
		}
		return out
	}
	return nil
}

func recvName(t ast.Expr) string {
	switch t := t.(type) {
	case *ast.StarExpr:
		return recvName(t.X)
	case *ast.IndexExpr:
		return recvName(t.X)
	case *ast.IndexListExpr:
		return recvName(t.X)
	case *ast.Ident:
		return t.Name
	}
	return ""
}
