package gots

import (
	"encoding/base64"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// embeddedVars are the package variables a //go:embed directive fills,
// with the expression holding the files' contents: a string, a []byte or
// an embed.FS.
func (e *emitter) embeddedVars(files []*ast.File) map[types.Object]string {
	out := map[types.Object]string{}
	for _, f := range files {
		for _, d := range f.Decls {
			g, ok := d.(*ast.GenDecl)
			if !ok || g.Tok != token.VAR {
				continue
			}
			for _, spec := range g.Specs {
				vs := spec.(*ast.ValueSpec)
				patterns := embedPatterns(vs.Doc)
				if len(g.Specs) == 1 && len(patterns) == 0 {
					patterns = embedPatterns(g.Doc)
				}
				if len(patterns) == 0 || len(vs.Names) != 1 {
					continue
				}
				v := e.info.Defs[vs.Names[0]]
				if v == nil {
					continue
				}
				js, err := e.embedValue(v.Type(), e.embedMatches(patterns))
				if err != nil {
					e.errorf(vs.Pos(), "go:embed: %v", err)
					continue
				}
				out[v] = js
			}
		}
	}
	return out
}

// embedPatterns are the patterns of the //go:embed lines in a comment.
func embedPatterns(doc *ast.CommentGroup) []string {
	if doc == nil {
		return nil
	}
	var out []string
	for _, c := range doc.List {
		if rest, ok := strings.CutPrefix(c.Text, "//go:embed "); ok {
			for _, p := range strings.Fields(rest) {
				out = append(out, strings.Trim(p, `"`+"`"))
			}
		}
	}
	return out
}

// embedMatches are the package's embedded files a set of patterns names:
// a file matching a pattern, or under a directory that does.
func (e *emitter) embedMatches(patterns []string) []string {
	var out []string
	for _, f := range e.pkg.EmbedFiles {
		for _, p := range patterns {
			p = strings.TrimPrefix(p, "all:")
			if ok, _ := path.Match(p, f); ok || strings.HasPrefix(f, p+"/") {
				out = append(out, f)
				break
			}
		}
	}
	sort.Strings(out)
	return out
}

func (e *emitter) embedValue(t types.Type, names []string) (string, error) {
	read := func(name string) ([]byte, error) {
		return os.ReadFile(filepath.Join(e.pkg.Dir, filepath.FromSlash(name)))
	}
	if isString(t) || isSliceOfBytes(t) {
		var data []byte
		if len(names) > 0 {
			b, err := read(names[0])
			if err != nil {
				return "", err
			}
			data = b
		}
		if isString(t) {
			return jsString(string(data)), nil
		}
		return "$.slice($.str.strToBytes(atob(" + jsString(base64.StdEncoding.EncodeToString(data)) + ")))", nil
	}
	var entries []string
	for _, n := range names {
		b, err := read(n)
		if err != nil {
			return "", err
		}
		entries = append(entries, "["+jsString(n)+", "+jsString(base64.StdEncoding.EncodeToString(b))+"]")
	}
	return e.importAlias(e.prog.Packages["embed"].Types) + ".files([" + strings.Join(entries, ", ") + "])", nil
}

func isSliceOfBytes(t types.Type) bool {
	s, ok := t.Underlying().(*types.Slice)
	return ok && basicKind(s.Elem()) == types.Uint8
}
