package gots

import (
	"go/ast"
	"go/token"
)

func (e *emitter) stmts(w *writer, list []ast.Stmt) {
	if backwardGoto(list) {
		e.stmtsStateMachine(w, list)
		return
	}
	if regions := gotoRegions(list); len(regions) > 0 {
		e.stmtsWithGotos(w, list, regions)
		return
	}
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
		e.goStmt(w, s)
	case *ast.SendStmt:
		e.sendStmt(w, s)
	case *ast.SelectStmt:
		e.selectStmt(w, s, "")
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
	case *ast.SelectStmt:
		e.selectStmt(w, inner, label)
	default:
		w.open("%s: {", label)
		e.stmt(w, s.Stmt)
		w.close("}")
	}
}

// target finds the statement a bare break (loops and switches) or continue
// (loops only) leaves, and whether a goto state machine's loop lies in
// between, in which case the jump has to name it.
func (e *emitter) target(forContinue bool) (b *breakable, crossed bool) {
	for i := len(e.fn.breaks) - 1; i >= 0; i-- {
		t := e.fn.breaks[i]
		switch {
		case t.barrier:
			crossed = true
		case t.isSwitch && forContinue:
		default:
			return t, crossed
		}
	}
	return nil, crossed
}

func (e *emitter) branch(w *writer, s *ast.BranchStmt) {
	switch s.Tok {
	case token.BREAK, token.CONTINUE:
		if e.yieldBranch(w, s) {
			return
		}
		kw := "break"
		if s.Tok == token.CONTINUE {
			kw = "continue"
		}
		if s.Label != nil {
			w.line("%s %s;", kw, mangle(s.Label.Name))
			return
		}
		t, crossed := e.target(s.Tok == token.CONTINUE)
		if t != nil && (t.isSwitch || crossed) {
			t.used = true
			w.line("%s %s;", kw, t.label)
			return
		}
		w.line("%s;", kw)
	case token.FALLTHROUGH:
		// Handled by the switch, which joins the clauses.
	case token.GOTO:
		if gc, ok := e.fn.gotos[s.Label.Name]; ok {
			w.line("{ %s = %d; continue %s; }", gc.state, gc.index, gc.loop)
			return
		}
		w.line("break %s;", gotoLabel(s.Label.Name))
	}
}
