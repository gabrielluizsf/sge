package gots

import (
	"regexp"
	"strings"
)

var (
	// exportLine is a top-level export of a module, and its name.
	exportLine = regexp.MustCompile(`^export (?:async )?(?:function|class|const|let|type|interface) ([A-Za-z_$][\w$]*)`)
	// memberLine is a member of a class or interface, two spaces in, and
	// its name; a getter or setter is named as the field it stands for.
	memberLine = regexp.MustCompile(`^  (?:static )?(?:async )?(?:get |set )?([A-Za-z_$][\w$]*)[?]?[(:<]`)
)

// document writes Go's documentation into the module src of a package as
// JSDoc: the package's doc as the module's, before its imports, then each
// export's and each class or interface member's. What Go does not export,
// the plumbing of the transpiled code, is marked @ignore, so the
// TypeScript's documentation is the Go package's.
func document(src, path string, d docs) string {
	lines := strings.SplitAfter(src, "\n")
	var b strings.Builder
	class := ""
	moduleDone := false
	for _, l := range lines {
		if !moduleDone && strings.HasPrefix(l, "import ") {
			pkg := d.pkg
			if pkg == "" {
				pkg = "The Go package " + path + ", transpiled to TypeScript."
			}
			b.WriteString(jsdoc(pkg, "", "@module"))
			moduleDone = true
		}
		switch {
		case exportLine.MatchString(l):
			name := exportLine.FindStringSubmatch(l)[1]
			b.WriteString(exportDoc(d, name))
			class = ""
			if (strings.Contains(l, " class ") || strings.Contains(l, " interface ")) && strings.HasSuffix(strings.TrimSpace(l), "{") {
				class = name
			}
		case class != "" && l == "}\n":
			class = ""
		case class != "" && strings.HasPrefix(l, "  [more: string]"):
			b.WriteString(jsdoc("Any other methods the value has.", "  "))
		case class != "" && memberLine.MatchString(l):
			b.WriteString(memberDoc(d, class, memberLine.FindStringSubmatch(l)[1]))
		}
		b.WriteString(l)
	}
	if !moduleDone {
		return jsdoc(d.pkg, "", "@module") + b.String()
	}
	return b.String()
}

// exportDoc is the JSDoc of the module's export name.
func exportDoc(d docs, name string) string {
	if text, ok := d.decls[name]; ok {
		return jsdoc(text, "")
	}
	return "/** @ignore */\n"
}

// memberDoc is the JSDoc of member name of class, documented by Go, made
// up for a struct class's constructor, and ignored if Go has no such
// member.
func memberDoc(d docs, class, name string) string {
	if text, ok := d.members[class][name]; ok {
		if text == "" {
			text = name + " is a field of " + class + "." // Go says nothing more
		}
		return jsdoc(text, "  ")
	}
	if name == "constructor" && d.decls[class] != "" {
		return jsdoc("Makes a "+class+" of the values of its fields, in their order; those left out are their zero values.", "  ")
	}
	return "  /** @ignore */\n"
}
