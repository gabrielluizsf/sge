package gots

import (
	"go/ast"
	"go/token"
	"go/types"
)

// rangedArray is the array a range statement walks. An array literal is
// written as a plain JavaScript array, not the typed array it would be
// elsewhere: only the loop reads it, the values it holds are already
// those of its element type, and a typed array costs V8 about ten times
// as much to make. So is the array behind &T{...}, which nothing else can
// reach either.
func (e *emitter) rangedArray(x ast.Expr) string {
	inner := ast.Unparen(x)
	if u, ok := inner.(*ast.UnaryExpr); ok && u.Op == token.AND {
		inner = ast.Unparen(u.X)
	}
	lit, ok := inner.(*ast.CompositeLit)
	if !ok {
		return e.expr(x)
	}
	a, ok := e.typeOf(lit).Underlying().(*types.Array)
	if !ok {
		return e.expr(x)
	}
	return e.indexedLitIn(lit, a.Elem(), int(a.Len()), "")
}
