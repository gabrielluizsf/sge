package gots

import (
	"go/ast"
	"go/types"
)

func (e *emitter) typeSwitch(w *writer, s *ast.TypeSwitchStmt, label string) {
	w.open("{")
	if s.Init != nil {
		e.stmt(w, s.Init)
	}
	var x ast.Expr
	switch a := s.Assign.(type) {
	case *ast.ExprStmt:
		x = a.X.(*ast.TypeAssertExpr).X
	case *ast.AssignStmt:
		x = a.Rhs[0].(*ast.TypeAssertExpr).X
	}
	v := e.sc.fresh("$x")
	w.line("const %s = %s;", v, e.expr(x))
	clauses := s.Body.List
	conds := make([]string, len(clauses))
	for i, c := range clauses {
		cc := c.(*ast.CaseClause)
		cond := ""
		for j, t := range cc.List {
			if j > 0 {
				cond += " || "
			}
			cond += e.holds(v, e.typeOf(t))
		}
		conds[i] = cond
	}
	e.chain(w, label, conds, func(w *writer, i int) {
		cc := clauses[i].(*ast.CaseClause)
		if obj, ok := e.info.Implicits[cc].(*types.Var); ok {
			js := v
			if len(cc.List) == 1 {
				js = e.unboxed(v, e.typeOf(cc.List[0]))
				// An interface's value never changes, so a case that only
				// reads its variable can have the value itself.
				body := &ast.BlockStmt{List: cc.Body}
				if t := e.typeOf(cc.List[0]); isValue(t) && !mutates(e.info, body, obj) && !pinned(e.info, body, obj) {
					js = v + ".v"
				}
			}
			e.declareVar(w, obj, js)
		}
		e.stmts(w, cc.Body)
	})
	w.close("}")
}
