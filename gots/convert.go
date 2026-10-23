package gots

import (
	"go/ast"
	"go/types"
	"strings"
)

// convert is the conversion to(x).
func (e *emitter) convert(x ast.Expr, to types.Type) string {
	from := e.typeOf(x)
	if isNil(from) {
		return "null"
	}
	if isInterface(to) {
		return e.value(x, to)
	}
	js := e.expr(x)
	fu, tu := from.Underlying(), to.Underlying()
	if isString(to) {
		switch f := fu.(type) {
		case *types.Basic:
			if f.Info()&types.IsInteger != 0 {
				if isBig(from) {
					js = "Number(" + js + ")"
				}
				return "$.str.runeStr(" + js + ")"
			}
			return js
		case *types.Slice:
			if basicKind(f.Elem()) == types.Int32 {
				return "$.str.runesToStr(" + js + ")"
			}
			return "$.str.bytesToStr(" + js + ")"
		}
	}
	if ts, ok := tu.(*types.Slice); ok && isString(from) {
		if basicKind(ts.Elem()) == types.Int32 {
			return "$.slice($.str.strToRunes(" + js + "))"
		}
		return "$.slice($.str.strToBytes(" + js + "))"
	}
	if isNumeric(to) && isNumeric(from) {
		return numConvert(js, from, to)
	}
	if isStruct(to) && !e.sameClass(to, from) {
		return e.structConvert(js, structOf(from), to)
	}
	return js
}

// structConvert copies a struct into the class of another struct type with
// the same fields.
func (e *emitter) structConvert(js string, st *types.Struct, to types.Type) string {
	var fs []string
	for i := 0; i < st.NumFields(); i++ {
		fs = append(fs, e.cloneOf("o."+fieldJS(st, i), st.Field(i).Type()))
	}
	return "((o: any) => new " + e.structClass(to) + "(" + strings.Join(fs, ", ") + "))(" + js + ")"
}

// numConvert converts between numeric types.
func numConvert(js string, from, to types.Type) string {
	fk, tk := defaultKind(basicKind(from)), defaultKind(basicKind(to))
	if fk == tk {
		return js
	}
	fromFloat, toFloat := isFloat(from), isFloat(to)
	switch {
	case tk == types.Uint64:
		return "$.int.toU64(" + js + ")"
	case fk == types.Uint64:
		switch {
		case toFloat:
			return wrap(tk, "Number("+js+")")
		case small(tk) && unsigned(tk):
			return "Number(BigInt.asUintN(" + itoa(width(tk)) + ", " + js + "))"
		case small(tk):
			return "Number(BigInt.asIntN(" + itoa(width(tk)) + ", " + js + "))"
		}
		return "$.int.fromU64(" + js + ")"
	case toFloat:
		if tk == types.Float32 {
			return "Math.fround(" + js + ")"
		}
		return js
	case fromFloat:
		if small(tk) {
			return wrap(tk, js)
		}
		return "$.int.trunc(" + js + ")"
	}
	return wrap(tk, js)
}
