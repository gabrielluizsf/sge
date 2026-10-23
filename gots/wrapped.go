package gots

import (
	"go/types"
	"strings"
)

// wrapped reports whether fn gets a TypeScript-facing wrapper: it is
// exported from a public package, takes or returns a string, and, as a
// method, belongs to an exported struct. Its Go body is then named
// fn$g, which Go code calls, and fn converts text on the way in and out.
func (e *emitter) wrapped(fn *types.Func) bool {
	if !fn.Exported() || fn.Pkg() == nil || !e.prog.Public[fn.Pkg().Path()] {
		return false
	}
	// Generic code is left alone; the origin of an instantiated method
	// says whether it is generic, the same at its declaration and its calls.
	origin := fn.Origin().Type().(*types.Signature)
	if origin.TypeParams().Len() > 0 || origin.RecvTypeParams().Len() > 0 {
		return false
	}
	sig := fn.Type().(*types.Signature)
	if recv := sig.Recv(); recv != nil {
		n := namedOf(derefType(recv.Type()))
		if n == nil || !n.Obj().Exported() || !isStruct(n) {
			return false
		}
	}
	return hasText(sig.Params()) || hasText(sig.Results())
}

func hasText(t *types.Tuple) bool {
	for i := 0; i < t.Len(); i++ {
		if isString(t.At(i).Type()) {
			return true
		}
	}
	return false
}

// goName is the name Go code calls fn by.
func (e *emitter) goName(fn *types.Func) string {
	if e.wrapped(fn) {
		return fieldName(fn.Name()) + "$g"
	}
	return fieldName(fn.Name())
}

// wrapper writes the TypeScript-facing function or method head(params)
// calling target with its strings converted.
func (e *emitter) wrapper(w *writer, head, target string, sig *types.Signature) {
	var params, args []string
	for i := 0; i < sig.Params().Len(); i++ {
		p := sig.Params().At(i)
		n := "p" + itoa(i)
		t := p.Type()
		params = append(params, n+": "+e.tsType(t))
		if isString(t) && !(sig.Variadic() && i == sig.Params().Len()-1) {
			n = "$.str.toGo(" + n + ")"
		}
		args = append(args, n)
	}
	call := target + "(" + strings.Join(args, ", ") + ")"
	results := sig.Results()
	w.open("%s(%s): %s {", head, strings.Join(params, ", "), e.tsResults(results))
	switch {
	case results.Len() == 0:
		w.line("%s;", call)
	case results.Len() == 1:
		if isString(results.At(0).Type()) {
			call = "$.str.fromGo(" + call + ")"
		}
		w.line("return %s;", call)
	default:
		w.line("const r = %s;", call)
		var outs []string
		for i := 0; i < results.Len(); i++ {
			o := "r[" + itoa(i) + "]"
			if isString(results.At(i).Type()) {
				o = "$.str.fromGo(" + o + ")"
			}
			outs = append(outs, o)
		}
		w.line("return [%s];", strings.Join(outs, ", "))
	}
	w.close("}")
}

// fromIface is js, the result of calling an interface method that
// returns text, as a Go string: the object behind the interface may be
// TypeScript's, handing back JavaScript text.
func fromIface(js string, sig *types.Signature) string {
	r := sig.Results()
	if !hasText(r) {
		return js
	}
	if r.Len() == 1 {
		return "$.str.toGo(" + js + ")"
	}
	var conv []string
	for i := 0; i < r.Len(); i++ {
		v := "r[" + itoa(i) + "]"
		if isString(r.At(i).Type()) {
			v = "$.str.toGo(" + v + ")"
		}
		conv = append(conv, v)
	}
	return "((r: any) => [" + strings.Join(conv, ", ") + "])(" + js + ")"
}
