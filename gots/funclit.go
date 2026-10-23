package gots

import (
	"go/ast"
	"go/types"
	"strings"
)

// funcLit is a closure: an arrow function, sharing the enclosing function's
// names so its locals cannot shadow the ones it captures.
func (e *emitter) funcLit(x *ast.FuncLit) string {
	var w writer
	e.function(&w, "", e.typeOf(x).(*types.Signature), x.Type, x.Body, nil, true, false)
	return "(" + strings.TrimSuffix(w.String(), "\n") + ")"
}
