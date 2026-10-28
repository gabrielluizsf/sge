package gots

import "go/ast"

// indexOf is x as an index or a slice bound: Go takes any integer type,
// and a uint64 here is a BigInt, which arrays cannot be indexed with.
func (e *emitter) indexOf(x ast.Expr) string {
	return e.asIndex(x, e.expr(x))
}

// asIndex is js, the value of x, made a number when x is a uint64.
func (e *emitter) asIndex(x ast.Expr, js string) string {
	if isBig(e.typeOf(x)) {
		return "$Number(" + js + ")"
	}
	return js
}
