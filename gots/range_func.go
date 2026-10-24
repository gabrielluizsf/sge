package gots

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
)

// yieldState is the state of a range over a function whose body leaves it
// by break, continue, return or a jump to an outer loop. The body is the
// yield callback: it returns true to go on and false to stop, setting ctl
// to say why, and after the iterator returns the loop does what the body
// asked for.
type yieldState struct {
	ctl, ret string
	exits    []yieldExit // jumps to statements outside the loop
}

// yieldExit is a break or continue the body made to a statement around
// the loop, done once the iterator has returned.
type yieldExit struct {
	code int
	stmt *ast.BranchStmt
}

const (
	yieldBreak  = 1
	yieldReturn = 2
)

// rangeFunc is range over an iterator function.
func (e *emitter) rangeFunc(w *writer, s *ast.RangeStmt, sig *types.Signature, head string) {
	if e.bodyBlocks(s.Body) {
		w.line("%s;", e.unsupported(s, "a range over a function whose body blocks"))
		return
	}
	params := sig.Params().At(0).Type().(*types.Signature).Params()
	var vars []string
	for i := 0; i < params.Len(); i++ {
		vars = append(vars, e.sc.fresh("$y")+": any")
	}
	y := &yieldState{ctl: e.sc.fresh("$ctl"), ret: e.sc.fresh("$ret")}
	b := &breakable{yield: y}
	if l := strings.TrimSuffix(head, ": "); l != head {
		b.label = l
	}
	w.open("{")
	w.line("let %s = 0, %s: any;", y.ctl, y.ret)
	w.open("%s((%s) => {", e.expr(s.X), strings.Join(vars, ", "))
	// An iterator calling yield again once it said stop is a bug Go
	// panics over.
	w.line("if (%s !== 0) $.panic(new $.RuntimeError(\"range function continued iteration after function for loop body returned false\"));", y.ctl)
	keys := []ast.Expr{s.Key, s.Value}
	for i := range vars {
		e.bindRange(w, keys[i], strings.TrimSuffix(vars[i], ": any"), params.At(i).Type(), s.Tok)
	}
	// The callback cannot wait: what is called in it is not awaited.
	async := e.fn.async
	e.fn.async = false
	e.fn.breaks = append(e.fn.breaks, b)
	e.stmts(w, s.Body.List)
	e.fn.breaks = e.fn.breaks[:len(e.fn.breaks)-1]
	e.fn.async = async
	w.line("return true;")
	w.close("});")
	w.open("if (%s === %d) {", y.ctl, yieldReturn)
	e.finishReturn(w, y.ret, true)
	w.close("}")
	for _, x := range y.exits {
		w.open("if (%s === %d) {", y.ctl, x.code)
		e.branch(w, x.stmt)
		w.close("}")
	}
	w.close("}")
}

// yieldBranch writes a break or continue that crosses the callback of a
// range over a function, and reports whether it did. Meant for that loop,
// it returns from the callback; meant for a statement outside it, it
// returns asking the loop to make the jump once the iterator is done.
func (e *emitter) yieldBranch(w *writer, s *ast.BranchStmt) bool {
	target := -1
	for i := len(e.fn.breaks) - 1; i >= 0; i-- {
		b := e.fn.breaks[i]
		if s.Label != nil {
			if b.label == mangle(s.Label.Name) {
				target = i
				break
			}
			continue
		}
		if b.isSwitch && s.Tok == token.CONTINUE || b.barrier {
			continue
		}
		target = i
		break
	}
	// The innermost range over a function between here and the target.
	for i := len(e.fn.breaks) - 1; i >= 0 && i >= target; i-- {
		y := e.fn.breaks[i].yield
		if y == nil {
			continue
		}
		switch {
		case i == target && s.Tok == token.CONTINUE:
			w.line("return true;")
		case i == target:
			w.line("%s = %d; return false;", y.ctl, yieldBreak)
		default:
			code := yieldReturn + 1 + len(y.exits)
			y.exits = append(y.exits, yieldExit{code, s})
			w.line("%s = %d; return false;", y.ctl, code)
		}
		return true
	}
	return false
}

// innerYield is the range over a function whose callback the code being
// written is in, if any.
func (e *emitter) innerYield() *yieldState {
	for i := len(e.fn.breaks) - 1; i >= 0; i-- {
		if y := e.fn.breaks[i].yield; y != nil {
			return y
		}
	}
	return nil
}

// finishReturn returns the results js, or nothing when it is "", from the
// function: from inside a range-over-function callback, by handing them to
// the loop to return. given says js is already a tuple to return as it is.
func (e *emitter) finishReturn(w *writer, js string, given bool) {
	if y := e.innerYield(); y != nil {
		if js != "" {
			w.line("%s = %s;", y.ret, js)
		}
		w.line("%s = %d; return false;", y.ctl, yieldReturn)
		return
	}
	if js == "" || given && e.fn.sig.Results().Len() == 0 {
		w.line("return;")
		return
	}
	w.line("return %s;", js)
}

// bodyBlocks reports whether body may block, outside the closures in it.
func (e *emitter) bodyBlocks(body *ast.BlockStmt) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.SendStmt:
			found = true
		case *ast.UnaryExpr:
			found = found || n.Op == token.ARROW
		case *ast.SelectStmt:
			found = found || !selectHasDefault(n)
		case *ast.RangeStmt:
			_, isChan := e.typeOf(n.X).Underlying().(*types.Chan)
			found = found || isChan
		case *ast.CallExpr:
			if fn := e.staticCallee(n); fn != nil {
				found = found || e.blocks(fn)
			} else if lit := e.blockingOf().closureCalled(e.info, n); lit != nil {
				b := e.blockingOf()
				found = found || b.blocks[b.key(lit.Pos())]
			}
		}
		return !found
	})
	return found
}
