package gots

import (
	"go/ast"
	"go/types"
	"strings"
)

// funcLit is a closure: an arrow function, sharing the enclosing function's
// names so its locals cannot shadow the ones it captures. A closure with a
// twin carries it, as $a, for code that can wait to call.
func (e *emitter) funcLit(x *ast.FuncLit) string {
	sig := e.typeOf(x).(*types.Signature)
	var w writer
	e.function(&w, "", sig, x.Type, x.Body, nil, true, false)
	f := "(" + strings.TrimSuffix(w.String(), "\n") + ")"
	if !e.twinFunc(x.Type) {
		return f
	}
	var a writer
	e.waiting = true
	e.function(&a, "", sig, x.Type, x.Body, nil, true, false)
	return "$.twins(" + f + ", (" + strings.TrimSuffix(a.String(), "\n") + "))"
}
