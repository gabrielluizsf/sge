package gots

import (
	"go/ast"
	"go/types"
	"strings"
)

// builtin is a call to one of Go's predeclared functions.
func (e *emitter) builtin(c *ast.CallExpr) string {
	var name string
	switch f := ast.Unparen(c.Fun).(type) {
	case *ast.Ident:
		name = f.Name
	case *ast.SelectorExpr:
		// unsafe.String and friends.
		name = "unsafe." + f.Sel.Name
	}
	arg := func(i int) string { return e.expr(c.Args[i]) }
	switch name {
	case "unsafe.SliceData":
		return "$.unsafe.sliceData(" + arg(0) + ")"
	case "unsafe.Slice":
		return "$.unsafe.slice(" + arg(0) + ", " + arg(1) + ")"
	case "unsafe.String":
		return "$.unsafe.string(" + arg(0) + ", " + arg(1) + ")"
	case "unsafe.StringData":
		return "$.unsafe.stringData(" + arg(0) + ")"
	case "unsafe.Add":
		return "$.unsafe.add(" + arg(0) + ", " + arg(1) + ")"
	case "len":
		switch u := e.typeOf(c.Args[0]).Underlying().(type) {
		case *types.Basic:
			return arg(0) + ".length"
		case *types.Slice:
			return "$.len(" + arg(0) + ")"
		case *types.Map:
			return "$.mapLen(" + arg(0) + ")"
		case *types.Pointer, *types.Array:
			_ = u
			return arg(0) + ".length"
		case *types.Chan:
			return "$.ch.chanLen(" + arg(0) + ")"
		}
	case "cap":
		switch e.typeOf(c.Args[0]).Underlying().(type) {
		case *types.Slice:
			return "$.cap(" + arg(0) + ")"
		case *types.Chan:
			return "$.ch.chanCap(" + arg(0) + ")"
		}
		return arg(0) + ".length"
	case "append":
		return e.appendCall(c)
	case "copy":
		if isString(e.typeOf(c.Args[1])) {
			return "$.copyString(" + arg(0) + ", " + arg(1) + ")"
		}
		return "$.copy(" + arg(0) + ", " + arg(1) + ")"
	case "delete":
		m := e.typeOf(c.Args[0]).Underlying().(*types.Map)
		return "$.mapDelete(" + arg(0) + ", " + e.value(c.Args[1], m.Key()) + ")"
	case "make":
		return e.makeCall(c)
	case "new":
		t := e.typeOf(c.Args[0])
		if isValue(t) {
			return e.zero(t)
		}
		return "new $.Cell(" + e.zero(t) + ")"
	case "panic":
		return "$.panic(" + e.value(c.Args[0], types.NewInterfaceType(nil, nil)) + ")"
	case "recover":
		return "$.recover()"
	case "print", "println":
		var parts []string
		for i := range c.Args {
			parts = append(parts, arg(i))
		}
		return "console.error(" + strings.Join(parts, ", ") + ")"
	case "min", "max":
		var parts []string
		for i := range c.Args {
			parts = append(parts, arg(i))
		}
		if isString(e.typeOf(c)) || isBig(e.typeOf(c)) {
			return "$." + name + "Of(" + strings.Join(parts, ", ") + ")"
		}
		return "$Math." + name + "(" + strings.Join(parts, ", ") + ")"
	case "clear":
		switch u := e.typeOf(c.Args[0]).Underlying().(type) {
		case *types.Map:
			return "(" + arg(0) + "?.clear())"
		case *types.Slice:
			if z := e.zeroFn(u.Elem()); z != "" {
				return "$.clear(" + arg(0) + ", " + z + ")"
			}
			return "$.clear(" + arg(0) + ", () => 0)"
		}
	case "close":
		return "$.ch.close(" + arg(0) + ")"
	case "complex":
		js := "{ re: " + arg(0) + ", im: " + arg(1) + " }"
		if basicKind(e.typeOf(c)) == types.Complex64 {
			return "$.cplx.c64(" + js + ")"
		}
		return js
	case "real":
		return arg(0) + ".re"
	case "imag":
		return arg(0) + ".im"
	case "unsafe.Sizeof":
		return itoa(int(wasmSizes.Sizeof(e.typeOf(c.Args[0]))))
	case "unsafe.Alignof":
		return itoa(int(wasmSizes.Alignof(e.typeOf(c.Args[0]))))
	}
	return e.unsupported(c, "builtin "+name)
}

func (e *emitter) appendCall(c *ast.CallExpr) string {
	st := e.typeOf(c).Underlying().(*types.Slice)
	s := e.expr(c.Args[0])
	b := backingOrNull(st.Elem())
	if c.Ellipsis.IsValid() {
		if isString(e.typeOf(c.Args[1])) {
			return "$.appendString(" + s + ", " + e.expr(c.Args[1]) + ")"
		}
		return "$.appendSlice(" + s + ", " + b + ", " + e.expr(c.Args[1]) + ")"
	}
	parts := []string{s, b}
	for _, a := range c.Args[1:] {
		parts = append(parts, e.value(a, st.Elem()))
	}
	return "$.append(" + strings.Join(parts, ", ") + ")"
}

func (e *emitter) makeCall(c *ast.CallExpr) string {
	t := e.typeOf(c.Args[0])
	switch u := t.Underlying().(type) {
	case *types.Slice:
		n := e.expr(c.Args[1])
		cp := n
		if len(c.Args) > 2 {
			cp = e.expr(c.Args[2])
		}
		if z := e.zeroFn(u.Elem()); z != "" {
			return "$.make(null, " + n + ", " + cp + ", " + z + ")"
		}
		return "$.make(" + backing(u.Elem()) + ", " + n + ", " + cp + ")"
	case *types.Map:
		return e.newMap(u, "")
	case *types.Chan:
		n := "0"
		if len(c.Args) > 1 {
			n = e.expr(c.Args[1])
		}
		return "new $.ch.Chan(" + n + ", () => " + e.zero(u.Elem()) + ")"
	}
	return e.unsupported(c, "make of this type")
}

// newMap is an empty map of type m, or one holding entries.
func (e *emitter) newMap(m *types.Map, entries string) string {
	if byValue(m.Key()) {
		if entries == "" {
			return "new $Map()"
		}
		return "new $Map(" + entries + ")"
	}
	if entries == "" {
		return "new $.KeyMap($.keyOf)"
	}
	return "$.keyMapOf(" + entries + ")"
}

// byValue reports whether JavaScript's Map compares keys of type t the way
// Go does.
func byValue(t types.Type) bool {
	switch u := t.Underlying().(type) {
	case *types.Basic:
		return u.Info()&types.IsComplex == 0
	case *types.Pointer:
		return isValue(u.Elem())
	case *types.Chan:
		return true
	}
	return false
}

var wasmSizes = types.SizesFor("gc", "wasm")
