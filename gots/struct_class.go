package gots

import (
	"fmt"
	"go/types"
	"strconv"
	"strings"
)

// structClass is the class holding values of the struct type t.
func (e *emitter) structClass(t types.Type) string {
	t = types.Unalias(t)
	if n, ok := t.(*types.Named); ok {
		return e.namedRef(n)
	}
	return e.anonClass(t.Underlying().(*types.Struct))
}

// anonClass is the class for an anonymous struct type, declared the first
// time the type is met. Identical anonymous struct types share one class,
// as they are one type in Go.
func (e *emitter) anonClass(st *types.Struct) string {
	key := typeString(st)
	if n, ok := e.anon[key]; ok {
		return n
	}
	n := "$S" + strconv.Itoa(len(e.anon)+1)
	e.anon[key] = n
	// Field types may be anonymous structs too, whose classes have to come
	// first: build this one aside and append it after them.
	var w writer
	e.structClassDecl(&w, n, st, false, nil)
	w.line("%s.$type = $.structType(%q, false, %s, null, () => %s);", n, key, n, e.fieldDescs(st))
	e.anonDecls.b.WriteString(w.String())
	return n
}

// fieldJS is the property a struct field is kept in. Blank fields are
// numbered, since a struct may have several.
func fieldJS(st *types.Struct, i int) string {
	f := st.Field(i)
	if f.Name() == "_" {
		return "_$" + strconv.Itoa(i)
	}
	return fieldName(f.Name())
}

// structClassDecl writes the class for a struct: fields, a positional
// constructor defaulting every field to its zero value, and the $clone,
// $set, $eq and $key methods value semantics need. body writes the methods.
func (e *emitter) structClassDecl(w *writer, name string, st *types.Struct, export bool, body func()) {
	kw := "class"
	if export {
		kw = "export class"
	}
	w.open("%s %s {", kw, name)
	var params, assigns, clones, sets, eqs, keys, texts []string
	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		js := e.fieldJS(st, i)
		ts := e.tsType(f.Type())
		w.line("%s: %s;", js, ts)
		// A parameter named like a module-level binding (an embedded
		// field is named after its type) would hide it from the defaults.
		param := mangle(fieldJS(st, i))
		if e.pkgNames[param] || param != fieldJS(st, i) {
			param += "$"
		}
		params = append(params, fmt.Sprintf("%s: %s = %s", param, ts, e.zero(f.Type())))
		if e.textField(st, i) {
			assigns = append(assigns, fmt.Sprintf("this.%s = $.str.toGo(%s);", js, param))
			texts = append(texts, fieldName(f.Name()))
		} else {
			assigns = append(assigns, fmt.Sprintf("this.%s = %s;", js, param))
		}
		clones = append(clones, e.cloneOf("this."+js, f.Type()))
		sets = append(sets, e.setOf("this."+js, "o."+js, f.Type()))
		if f.Name() != "_" {
			eqs = append(eqs, e.equalOf("this."+js, "o."+js, f.Type()))
		}
		keys = append(keys, e.keyExpr("this."+js, f.Type()))
	}
	w.line("static $type: $.Type;")
	w.open("constructor(%s) {", strings.Join(params, ", "))
	for _, a := range assigns {
		w.line("%s", a)
	}
	w.close("}")
	for _, t := range texts {
		// The text a TypeScript caller reads and writes; Go code uses the
		// bytes in t$s.
		w.line("get %s(): string { return $.str.fromGo(this.%s$s); }", t, t)
		w.line("set %s(v: string) { this.%s$s = $.str.toGo(v); }", t, t)
	}
	w.open("$clone(): %s {", name)
	w.line("return new %s(%s);", name, strings.Join(clones, ", "))
	w.close("}")
	w.open("$set(o: %s): this {", name)
	for _, s := range sets {
		w.line("%s;", s)
	}
	w.line("return this;")
	w.close("}")
	w.open("$eq(o: %s): boolean {", name)
	if len(eqs) == 0 {
		w.line("return true;")
	} else {
		w.line("return %s;", strings.Join(eqs, " && "))
	}
	w.close("}")
	w.open("$key(): string {")
	if len(keys) == 0 {
		w.line(`return "{}";`)
	} else {
		w.line(`return "{" + %s + "}";`, strings.Join(keys, ` + "," + `))
	}
	w.close("}")
	if body != nil {
		body()
	}
	w.close("}")
}

// fieldDescs is the thunk body listing a struct's fields for its descriptor.
func (e *emitter) fieldDescs(st *types.Struct) string {
	var fs []string
	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		// Trailing arguments are left out when they are the defaults:
		// embedded false, the property the field's name, no tag.
		// The descriptor has the field's Go name, _ for a blank one, and
		// the property its value is kept in.
		args := []string{fmt.Sprintf("%q", f.Name()), e.descInline(f.Type()), fmt.Sprint(f.Embedded()), fmt.Sprintf("%q", e.fieldJS(st, i)), jsString(st.Tag(i))}
		n := len(args)
		for n > 2 && (args[n-1] == jsString("") && n == 5 || n == 4 && args[3] == args[0] || n == 3 && args[2] == "false") {
			n--
		}
		fs = append(fs, "$.field("+strings.Join(args[:n], ", ")+")")
	}
	return "[" + strings.Join(fs, ", ") + "]"
}

// cloneOf is an expression copying the value x of type t.
func (e *emitter) cloneOf(x string, t types.Type) string {
	switch u := t.Underlying().(type) {
	case *types.Struct:
		return x + ".$clone()"
	case *types.Array:
		if isValue(u.Elem()) {
			return fmt.Sprintf("$.cloneArray(%s, (e: any) => %s)", x, e.cloneOf("e", u.Elem()))
		}
		return "$.cloneArray(" + x + ")"
	}
	return x
}

// setOf is a statement assigning src to the storage dst of type t, in
// place when t is a value type.
func (e *emitter) setOf(dst, src string, t types.Type) string {
	switch u := t.Underlying().(type) {
	case *types.Struct:
		return dst + ".$set(" + src + ")"
	case *types.Array:
		if isValue(u.Elem()) {
			return fmt.Sprintf("$.setArray(%s, %s, (d: any, s: any) => %s)", dst, src, e.setOf("d", "s", u.Elem()))
		}
		return "$.setArray(" + dst + ", " + src + ")"
	}
	return dst + " = " + src
}

// equalOf is an expression for a == b, both of type t.
func (e *emitter) equalOf(a, b string, t types.Type) string {
	switch u := t.Underlying().(type) {
	case *types.Struct:
		return a + ".$eq(" + b + ")"
	case *types.Array:
		return fmt.Sprintf("$.equal(%s, %s, %s)", e.desc(t), a, b)
	case *types.Interface:
		return "$.ifaceEq(" + a + ", " + b + ")"
	case *types.Pointer:
		if !isValue(u.Elem()) {
			return "$.ptrEq(" + a + ", " + b + ")"
		}
	case *types.Basic:
		if u.Info()&types.IsComplex != 0 {
			return fmt.Sprintf("(%s.re === %s.re && %s.im === %s.im)", a, b, a, b)
		}
	}
	return "(" + a + " === " + b + ")"
}
