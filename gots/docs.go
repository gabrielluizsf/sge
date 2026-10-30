package gots

import (
	"go/ast"
	"go/doc/comment"
	"regexp"
	"strings"
)

// docLink is a Go doc link, [Name] or [pkg.Name], as Markdown escapes it.
var docLink = regexp.MustCompile(`\\\[([A-Za-z_][\w.]*)\]`)

// docs is the Go documentation of a package, by the names its module
// exports: the package's own, each package-level declaration's, and each
// member's of a class, a struct's fields and a type's methods, the
// methods also by the T$$M functions that implement them.
type docs struct {
	pkg     string
	decls   map[string]string
	members map[string]map[string]string
}

// docsOf gathers the documentation of the exported declarations in files.
func docsOf(files []*ast.File) docs {
	d := docs{decls: map[string]string{}, members: map[string]map[string]string{}}
	for _, f := range files {
		if f.Doc != nil && (d.pkg == "" || strings.HasPrefix(f.Doc.Text(), "Package ")) {
			d.pkg = f.Doc.Text()
		}
		for _, decl := range f.Decls {
			switch decl := decl.(type) {
			case *ast.FuncDecl:
				d.funcDoc(decl)
			case *ast.GenDecl:
				d.genDoc(decl)
			}
		}
	}
	return d
}

// funcDoc records a function's doc, or a method's, under its type too.
func (d *docs) funcDoc(f *ast.FuncDecl) {
	if !f.Name.IsExported() || f.Doc == nil {
		return
	}
	if f.Recv == nil {
		d.decls[mangle(f.Name.Name)] = f.Doc.Text()
		return
	}
	t := recvName(f.Recv.List[0].Type)
	if t == "" {
		return
	}
	d.decls[t+"$$"+f.Name.Name] = f.Doc.Text()
	d.member(t, fieldName(f.Name.Name), f.Doc.Text())
}

// genDoc records the docs of a type, const or var declaration, and of a
// struct's fields and an interface's methods.
func (d *docs) genDoc(g *ast.GenDecl) {
	for _, s := range g.Specs {
		switch s := s.(type) {
		case *ast.TypeSpec:
			if !s.Name.IsExported() {
				continue
			}
			if text := specDoc(g, s.Doc, s.Comment); text != "" {
				d.decls[s.Name.Name] = text
			}
			d.fieldDocs(s)
		case *ast.ValueSpec:
			text := specDoc(g, s.Doc, s.Comment)
			for _, n := range s.Names {
				if n.IsExported() && text != "" {
					d.decls[mangle(n.Name)] = text
				}
			}
		}
	}
}

// specDoc is a spec's doc: its own, the declaration's when it stands
// alone, or the comment at the end of its line.
func specDoc(g *ast.GenDecl, doc, line *ast.CommentGroup) string {
	switch {
	case doc != nil:
		return doc.Text()
	case g.Doc != nil && len(g.Specs) == 1 || g.Doc != nil && !g.Lparen.IsValid():
		return g.Doc.Text()
	case line != nil:
		return line.Text()
	case g.Doc != nil:
		return g.Doc.Text()
	}
	return ""
}

// fieldDocs records the docs of a struct type's fields, or an interface
// type's methods, as members of its class or interface.
func (d *docs) fieldDocs(s *ast.TypeSpec) {
	var list *ast.FieldList
	switch t := s.Type.(type) {
	case *ast.StructType:
		list = t.Fields
	case *ast.InterfaceType:
		list = t.Methods
	}
	if list == nil {
		return
	}
	for _, f := range list.List {
		text := ""
		if f.Doc != nil {
			text = f.Doc.Text()
		} else if f.Comment != nil {
			text = f.Comment.Text()
		}
		for _, n := range f.Names {
			if n.IsExported() {
				d.member(s.Name.Name, fieldName(n.Name), text)
			}
		}
	}
}

func (d *docs) member(class, name, text string) {
	if d.members[class] == nil {
		d.members[class] = map[string]string{}
	}
	d.members[class][name] = text
}

// jsdoc is a Go doc comment as a JSDoc comment indented by indent: its
// text as Markdown, which JSR and editors show, and tags after it.
func jsdoc(text, indent string, tags ...string) string {
	var p comment.Parser
	md := strings.TrimSpace(string((&comment.Printer{}).Markdown(p.Parse(text))))
	md = docLink.ReplaceAllString(md, "{@link $1}")
	var b strings.Builder
	b.WriteString(indent + "/**\n")
	lines := []string{}
	if md != "" {
		lines = strings.Split(md, "\n")
	}
	if len(lines) > 0 && len(tags) > 0 {
		lines = append(lines, "")
	}
	lines = append(lines, tags...)
	for _, l := range lines {
		l = strings.ReplaceAll(l, "*/", "*\\/")
		if l == "" {
			b.WriteString(indent + " *\n")
			continue
		}
		b.WriteString(indent + " * " + l + "\n")
	}
	b.WriteString(indent + " */\n")
	return b.String()
}
