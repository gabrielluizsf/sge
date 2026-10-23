package gots

import (
	"go/ast"
	"go/token"
)

func (e *emitter) stmts(w *writer, list []ast.Stmt) {
	for _, s := range list {
		e.stmt(w, s)
	}
}

func (e *emitter) stmt(w *writer, s ast.Stmt) {
	switch s := s.(type) {
	case *ast.ExprStmt:
		w.line("%s;", e.expr(s.X))
	case *ast.AssignStmt:
		e.assign(w, s)
	case *ast.IncDecStmt:
		op := token.ADD
		if s.Tok == token.DEC {
			op = token.SUB
		}
		e.opAssign(w, s.X, op, "1")
	case *ast.DeclStmt:
		e.localDecl(w, s.Decl.(*ast.GenDecl))
	case *ast.ReturnStmt:
		e.returnStmt(w, s)
	case *ast.BlockStmt:
		w.open("{")
		e.stmts(w, s.List)
		w.close("}")
	case *ast.IfStmt:
		e.ifStmt(w, s, false)
	case *ast.ForStmt:
		e.forStmt(w, s, "")
	case *ast.RangeStmt:
		e.rangeStmt(w, s, "")
	case *ast.SwitchStmt:
		e.switchStmt(w, s, "")
	case *ast.TypeSwitchStmt:
		e.typeSwitch(w, s, "")
	case *ast.LabeledStmt:
		e.labeled(w, s)
	case *ast.BranchStmt:
		e.branch(w, s)
	case *ast.DeferStmt:
		e.deferStmt(w, s)
	case *ast.EmptyStmt:
	case *ast.GoStmt:
		w.line("%s;", e.unsupported(s, "go statements"))
	case *ast.SendStmt, *ast.SelectStmt:
		w.line("%s;", e.unsupported(s, "channels"))
	default:
		w.line("%s;", e.unsupported(s, "statement"))
	}
}

// ifStmt writes an if and its else chain.
func (e *emitter) ifStmt(w *writer, s *ast.IfStmt, _ bool) {
	if s.Init != nil {
		w.open("{")
		e.stmt(w, s.Init)
		e.ifStmt(w, &ast.IfStmt{If: s.If, Cond: s.Cond, Body: s.Body, Else: s.Else}, false)
		w.close("}")
		return
	}
	w.open("if (%s) {", e.expr(s.Cond))
	for {
		e.stmts(w, s.Body.List)
		el, ok := s.Else.(*ast.IfStmt)
		if !ok || el.Init != nil {
			break
		}
		w.close("} else if (%s) {", e.expr(el.Cond))
		w.indent++
		s = el
	}
	switch el := s.Else.(type) {
	case nil:
		w.close("}")
	case *ast.IfStmt:
		w.close("} else {")
		w.indent++
		e.ifStmt(w, el, false)
		w.close("}")
	case *ast.BlockStmt:
		w.close("} else {")
		w.indent++
		e.stmts(w, el.List)
		w.close("}")
	}
}

func (e *emitter) labeled(w *writer, s *ast.LabeledStmt) {
	label := mangle(s.Label.Name)
	switch inner := s.Stmt.(type) {
	case *ast.ForStmt:
		e.forStmt(w, inner, label)
	case *ast.RangeStmt:
		e.rangeStmt(w, inner, label)
	case *ast.SwitchStmt:
		e.switchStmt(w, inner, label)
	case *ast.TypeSwitchStmt:
		e.typeSwitch(w, inner, label)
	default:
		w.open("%s: {", label)
		e.stmt(w, s.Stmt)
		w.close("}")
	}
}

func (e *emitter) branch(w *writer, s *ast.BranchStmt) {
	switch s.Tok {
	case token.BREAK:
		if s.Label != nil {
			w.line("break %s;", mangle(s.Label.Name))
			return
		}
		if n := len(e.fn.breaks); n > 0 && e.fn.breaks[n-1].isSwitch {
			b := e.fn.breaks[n-1]
			b.used = true
			w.line("break %s;", b.label)
			return
		}
		w.line("break;")
	case token.CONTINUE:
		if s.Label != nil {
			w.line("continue %s;", mangle(s.Label.Name))
			return
		}
		w.line("continue;")
	case token.FALLTHROUGH:
		// Handled by the switch, which joins the clauses.
	case token.GOTO:
		w.line("%s;", e.unsupported(s, "goto"))
	}
}
