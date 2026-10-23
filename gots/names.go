package gots

import "strconv"

// reserved are words a JavaScript binding cannot be named, and globals the
// emitted code relies on and must not be shadowed.
var reserved = map[string]bool{
	"arguments": true, "await": true, "break": true, "case": true, "catch": true,
	"class": true, "const": true, "continue": true, "debugger": true,
	"default": true, "delete": true, "do": true, "else": true, "enum": true,
	"eval": true, "export": true, "extends": true, "false": true,
	"finally": true, "for": true, "function": true, "if": true,
	"implements": true, "import": true, "in": true, "instanceof": true,
	"interface": true, "let": true, "new": true, "null": true, "package": true,
	"private": true, "protected": true, "public": true, "return": true,
	"static": true, "super": true, "switch": true, "this": true, "throw": true,
	"true": true, "try": true, "typeof": true, "var": true, "void": true,
	"while": true, "with": true, "yield": true, "undefined": true, "NaN": true,
	"Infinity": true, "Math": true, "Array": true, "Object": true,
	"Number": true, "String": true, "Boolean": true, "BigInt": true,
	"Symbol": true, "Error": true, "Map": true, "Set": true, "Promise": true,
	"globalThis": true, "console": true, "Deno": true, "constructor": true,
	"Float64Array": true, "Float32Array": true, "Int32Array": true,
	"Int16Array": true, "Int8Array": true, "Uint32Array": true,
	"Uint16Array": true, "Uint8Array": true, "BigUint64Array": true,
	"BigInt64Array": true, "type": true, "as": true, "of": true,
	// TypeScript keywords that, as an import alias, would be read as an
	// operator where a type is written.
	"unique": true, "keyof": true, "infer": true, "readonly": true,
	"asserts": true, "is": true, "never": true, "unknown": true,
	"symbol": true, "object": true, "declare": true, "abstract": true,
	"namespace": true, "module": true, "global": true, "satisfies": true,
}

// mangle makes a Go identifier safe as a JavaScript binding. Names starting
// with $ belong to the emitter, so Go names never do.
func mangle(name string) string {
	if reserved[name] {
		return name + "$"
	}
	return name
}

// fieldName makes a Go field or method name safe as a property name.
func fieldName(name string) string {
	switch name {
	case "constructor", "__proto__", "prototype":
		return name + "$"
	}
	if len(name) > 0 && name[0] == '$' {
		return "$" + name
	}
	return name
}

// scope hands out binding names for one function and the closures in it,
// renaming a local that would clash with any other in the function or with
// a package-level name, so JavaScript's block scoping and temporal dead zone
// never see two Go variables as one.
type scope struct {
	used map[string]bool
}

func newScope(pkgNames map[string]bool) *scope {
	s := &scope{used: map[string]bool{}}
	for n := range pkgNames {
		s.used[n] = true
	}
	return s
}

func (s *scope) fresh(base string) string {
	base = mangle(base)
	name := base
	for i := 1; s.used[name]; i++ {
		name = base + "$" + strconv.Itoa(i)
	}
	s.used[name] = true
	return name
}
