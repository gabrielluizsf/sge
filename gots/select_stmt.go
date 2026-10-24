package gots

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"
)

// selectStmt is a select: its channels and the values it would send are
// evaluated in order, the runtime picks a ready case, waiting for one when
// there is no default, and the case's body runs with what it received.
func (e *emitter) selectStmt(w *writer, s *ast.SelectStmt, label string) {
	w.open("{")
	clauses := s.Body.List
	var cases []string
	index := make([]int, len(clauses))
	def := false
	for i, c := range clauses {
		index[i] = -1
		switch comm := c.(*ast.CommClause).Comm.(type) {
		case nil:
			def = true
			continue
		case *ast.SendStmt:
			elem := e.typeOf(comm.Chan).Underlying().(*types.Chan).Elem()
			cases = append(cases, "{ c: "+e.expr(comm.Chan)+", send: true, v: "+e.value(comm.Value, elem)+" }")
		default:
			cases = append(cases, "{ c: "+e.expr(commRecv(comm).X)+", send: false }")
		}
		index[i] = len(cases) - 1
	}
	r := e.sc.fresh("$sel")
	list := "[" + strings.Join(cases, ", ") + "]"
	if def {
		if e.fn.async {
			// A poll in a loop lets the other goroutines run now and then.
			p := e.fn.scratch
			w.line("if ((%s = $.ch.poll()) !== undefined) await %s;", p, p)
		}
		w.line("const %s = $.ch.selectNow(%s);", r, list)
	} else if e.fn.async {
		w.line("const %s = %s;", r, e.awaitIfPromise("$.ch.select("+list+")"))
	} else {
		w.line("%s;", e.unsupported(s, "a select that waits outside a function that can wait"))
		w.close("}")
		return
	}
	conds := make([]string, len(clauses))
	for i := range clauses {
		if index[i] >= 0 {
			conds[i] = fmt.Sprintf("%s[0] === %d", r, index[i])
		}
	}
	e.chain(w, label, conds, func(w *writer, i int) {
		cc := clauses[i].(*ast.CommClause)
		if as, ok := cc.Comm.(*ast.AssignStmt); ok {
			elem := e.typeOf(commRecv(as).X).Underlying().(*types.Chan).Elem()
			vals := []string{r + "[1]", r + "[2]"}
			froms := []types.Type{elem, types.Typ[types.Bool]}
			for k, l := range as.Lhs {
				if id, isId := l.(*ast.Ident); isId && id.Name == "_" {
					continue
				}
				e.store(w, l, e.convertValue(vals[k], froms[k], e.typeOf(l), false), as.Tok == token.DEFINE)
			}
		}
		e.stmts(w, cc.Body)
	})
	w.close("}")
}

// commRecv is the receive a select case's statement makes.
func commRecv(s ast.Stmt) *ast.UnaryExpr {
	switch s := s.(type) {
	case *ast.ExprStmt:
		return ast.Unparen(s.X).(*ast.UnaryExpr)
	case *ast.AssignStmt:
		return ast.Unparen(s.Rhs[0]).(*ast.UnaryExpr)
	}
	return nil
}
