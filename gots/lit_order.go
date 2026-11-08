package gots

import (
	"go/ast"
	"go/token"
	"strings"
)

// effectful reports whether evaluating x may do something another element
// of the same literal could see or be seen by: a call that is not a
// conversion, or a receive. Making a closure runs nothing.
func (e *emitter) effectful(x ast.Expr) bool {
	found := false
	ast.Inspect(x, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.CallExpr:
			if tv, ok := e.info.Types[n.Fun]; !ok || !tv.IsType() {
				found = true
			}
		case *ast.UnaryExpr:
			found = found || n.Op == token.ARROW
		}
		return !found
	})
	return found
}

// keyed is a struct literal's element: the field it sets, by index, its
// value as written in TypeScript, and whether evaluating it does anything.
type keyed struct {
	field  int
	js     string
	effect bool
}

// outOfOrder reports whether the elements that do something come in
// another order than their fields: Go evaluates a literal's elements in
// the order they are written, a constructor's arguments go in the fields'.
func outOfOrder(els []keyed) bool {
	last, n := -1, 0
	out := false
	for _, k := range els {
		if !k.effect {
			continue
		}
		n++
		out = out || k.field < last
		last = k.field
	}
	return n > 1 && out
}

// inWrittenOrder is the struct of class made from vals, its fields', with
// the elements that do something evaluated first, into an array, in the
// order they are written: the array is built where the literal is, so a
// value that waits may still wait.
func inWrittenOrder(els []keyed, vals []string, class string) string {
	args := append([]string(nil), vals...)
	var first []string
	for _, k := range els {
		if k.effect && k.field < len(args) {
			args[k.field] = "$w[" + itoa(len(first)) + "]"
			first = append(first, k.js)
		}
	}
	return "(($w: any[]) => new " + class + "(" + strings.Join(args, ", ") + "))([" + strings.Join(first, ", ") + "])"
}
