package gots

import (
	"fmt"
	"go/types"
	"strings"
)

// tsType is the TypeScript annotation for values of t. The emitted modules
// are not type-checked themselves (Go already checked them); the annotations
// are for the TypeScript code that uses them.
func (e *emitter) tsType(t types.Type) string {
	switch t := t.(type) {
	case *types.Basic:
		switch {
		case t.Kind() == types.UntypedNil:
			return "null"
		case t.Info()&types.IsBoolean != 0:
			return "boolean"
		case t.Info()&types.IsString != 0:
			return "string"
		case t.Kind() == types.Uint64:
			return "bigint"
		case t.Info()&types.IsNumeric != 0:
			return "number"
		}
		return "any"
	case *types.Alias:
		return e.tsType(types.Unalias(t))
	case *types.Named:
		if t.Obj().Pkg() == nil {
			if t.Obj().Name() == "error" {
				return "$.error | null"
			}
			return "any"
		}
		if t.TypeArgs().Len() > 0 && isInterface(t) {
			return "any"
		}
		name := e.namedRef(t)
		if isInterface(t) {
			if t.Underlying().(*types.Interface).Empty() {
				return "any"
			}
			return name + " | null"
		}
		return name
	case *types.Pointer:
		if isValue(t.Elem()) {
			return e.tsType(t.Elem()) + " | null"
		}
		return "$.Ptr<" + e.tsType(t.Elem()) + "> | null"
	case *types.Slice:
		return "$.Slice<" + e.tsType(t.Elem()) + "> | null"
	case *types.Array:
		if b := backing(t.Elem()); b != "" {
			return b
		}
		return e.tsType(t.Elem()) + "[]"
	case *types.Map:
		return "$.GoMap<" + e.tsType(t.Key()) + ", " + e.tsType(t.Elem()) + "> | null"
	case *types.Signature:
		return "(" + e.tsSignature(t) + ") | null"
	case *types.Interface:
		if t.Empty() {
			return "any"
		}
		var ms []string
		for i := 0; i < t.NumMethods(); i++ {
			m := t.Method(i)
			ms = append(ms, fieldName(m.Name())+e.tsMethodType(m.Type().(*types.Signature)))
		}
		return "{ " + strings.Join(ms, "; ") + " } | null"
	case *types.Struct:
		return e.anonClass(t)
	case *types.Tuple:
		return e.tsResults(t)
	case *types.TypeParam:
		return t.Obj().Name()
	}
	return "any"
}

// tsSignature is a function type: (a: A, b: B) => R.
func (e *emitter) tsSignature(sig *types.Signature) string {
	return e.tsParams(sig) + " => " + e.tsResults(sig.Results())
}

// tsMethodType is a method's shape in an interface: (a: A): R.
func (e *emitter) tsMethodType(sig *types.Signature) string {
	return e.tsParams(sig) + ": " + e.tsResults(sig.Results())
}

func (e *emitter) tsParams(sig *types.Signature) string {
	var ps []string
	for i := 0; i < sig.Params().Len(); i++ {
		ps = append(ps, fmt.Sprintf("p%d: %s", i, e.tsType(sig.Params().At(i).Type())))
	}
	return "(" + strings.Join(ps, ", ") + ")"
}

func (e *emitter) tsResults(r *types.Tuple) string {
	switch r.Len() {
	case 0:
		return "void"
	case 1:
		return e.tsType(r.At(0).Type())
	}
	var ts []string
	for i := 0; i < r.Len(); i++ {
		ts = append(ts, e.tsType(r.At(i).Type()))
	}
	return "[" + strings.Join(ts, ", ") + "]"
}
