package gots

import "go/ast"

// pollingSelects are the selects with a default inside a loop of body, not
// counting closures: a goroutine polling channels so never waits, and only
// takes turns with the others, as Go's scheduler would make it, if the
// poll lets them run now and then.
func pollingSelects(body *ast.BlockStmt) map[*ast.SelectStmt]bool {
	out := map[*ast.SelectStmt]bool{}
	var walk func(n ast.Node, inLoop bool)
	walk = func(n ast.Node, inLoop bool) {
		ast.Inspect(n, func(m ast.Node) bool {
			switch m := m.(type) {
			case *ast.FuncLit:
				return false
			case *ast.ForStmt:
				if m != n {
					walk(m.Body, true)
					return false
				}
			case *ast.RangeStmt:
				if m != n {
					walk(m.Body, true)
					return false
				}
			case *ast.SelectStmt:
				if inLoop && selectHasDefault(m) {
					out[m] = true
				}
			}
			return true
		})
	}
	walk(body, false)
	return out
}
