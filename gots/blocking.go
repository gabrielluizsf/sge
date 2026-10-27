package gots

import (
	"go/ast"
	"go/token"
	"go/types"
)

// blocking is which functions may block, and so are written as async
// functions their callers await: those that send or receive on a channel,
// select with no default, range over a channel, call a blocking native, or
// call a function that blocks. A go statement does not block its caller.
//
// Functions and closures are known by where their func keyword is, as
// file:line:column, which stays the same when a package is checked again
// with its tests and its objects are new.
type blocking struct {
	// cha counts a call through an interface as a call of every method of
	// that name and signature: the analysis of what may block at all,
	// beside the one of what blocks for sure.
	cha      bool
	fset     *token.FileSet
	blocks   map[string]bool
	callers  map[string][]string // callee → functions calling it
	work     []string
	funcVars map[types.Object]*ast.FuncLit // see closureVars
}

func newBlocking(fset *token.FileSet, cha bool) *blocking {
	b := &blocking{cha: cha, fset: fset, blocks: map[string]bool{}, callers: map[string][]string{}, funcVars: map[types.Object]*ast.FuncLit{}}
	if cha {
		for _, k := range twinIfaceMethods {
			b.blocks[k] = true
		}
	}
	return b
}

func (b *blocking) key(pos token.Pos) string { return b.fset.Position(pos).String() }

// addFiles reads every function and closure of files, checked with info.
func (b *blocking) addFiles(files []*ast.File, info *types.Info) {
	for v, lit := range closureVars(files, info) {
		b.funcVars[v] = lit
	}
	for _, f := range files {
		for _, d := range f.Decls {
			if fd, ok := d.(*ast.FuncDecl); ok && fd.Body != nil {
				b.addFunc(fd.Type, fd.Body, info)
				if fd.Recv != nil && b.cha {
					b.addMethod(fd, info)
				}
			}
			if gd, ok := d.(*ast.GenDecl); ok {
				b.addPackageClosures(gd, info)
			}
		}
		if b.cha {
			b.addFuncValues(f, info)
		}
	}
}

// addFunc reads one function body, and the closures in it, each its own.
func (b *blocking) addFunc(ft *ast.FuncType, body *ast.BlockStmt, info *types.Info) {
	b.addBody(b.key(ft.Pos()), body, info)
}

// addBody reads body as the function known as me.
func (b *blocking) addBody(me string, body *ast.BlockStmt, info *types.Info) {
	nonBlocking := map[ast.Node]bool{} // channel operations of a select with a default
	polls := pollingSelects(body)
	var visit func(n ast.Node) bool
	visit = func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncLit:
			b.addFunc(n.Type, n.Body, info)
			b.funcValue(b.key(n.Type.Pos()), info.TypeOf(n))
			return false
		case *ast.SelectStmt:
			if selectHasDefault(n) {
				for _, c := range n.Body.List {
					if cc := c.(*ast.CommClause); cc.Comm != nil {
						markCommOp(cc.Comm, nonBlocking)
					}
				}
				if polls[n] {
					b.mark(me)
				}
			} else {
				b.mark(me)
			}
		case *ast.SendStmt:
			if !nonBlocking[n] {
				b.mark(me)
			}
		case *ast.UnaryExpr:
			if n.Op == token.ARROW && !nonBlocking[n] {
				b.mark(me)
			}
		case *ast.RangeStmt:
			switch t := info.TypeOf(n.X).Underlying().(type) {
			case *types.Chan:
				b.mark(me)
			case *types.Signature:
				// The body of a range over a function is the yield callback
				// the iterator calls: a function of its own, which the loop
				// waits for.
				ast.Inspect(n.X, visit)
				rk := rangeFuncKey(b.key(n.Pos()))
				b.addBody(rk, n.Body, info)
				b.callers[rk] = append(b.callers[rk], me)
				if b.blocks[rk] {
					b.mark(me)
				}
				b.funcValue(rk, t.Params().At(0).Type())
				return false
			}
		case *ast.GoStmt:
			// The call runs in a goroutine of its own; what it is called
			// with is worked out here and now.
			if lit, ok := ast.Unparen(n.Call.Fun).(*ast.FuncLit); ok {
				// Run, not handed on: not a function value.
				b.addFunc(lit.Type, lit.Body, info)
			} else {
				ast.Inspect(n.Call.Fun, func(m ast.Node) bool { return b.visitOperand(me, m, info) })
			}
			for _, a := range n.Call.Args {
				ast.Inspect(a, func(m ast.Node) bool { return b.visitOperand(me, m, info) })
			}
			return false
		case *ast.CallExpr:
			b.call(me, n, info)
		}
		return true
	}
	ast.Inspect(body, visit)
}

// visitOperand is addFunc's visit for an operand of a go statement, whose
// closures are functions of their own and whose calls happen now.
func (b *blocking) visitOperand(me string, n ast.Node, info *types.Info) bool {
	switch n := n.(type) {
	case *ast.FuncLit:
		b.addFunc(n.Type, n.Body, info)
		b.funcValue(b.key(n.Type.Pos()), info.TypeOf(n))
		return false
	case *ast.UnaryExpr:
		if n.Op == token.ARROW {
			b.mark(me)
		}
	case *ast.CallExpr:
		b.call(me, n, info)
	}
	return true
}

// call notes that me calls what c calls, when that is known.
func (b *blocking) call(me string, c *ast.CallExpr, info *types.Info) {
	fn := staticCalleeIn(info, c)
	if fn == nil && b.cha {
		if m := ifaceMethodCalled(info, c); m != nil {
			k := ifaceKey(m)
			b.callers[k] = append(b.callers[k], me)
			if b.blocks[k] {
				b.mark(me)
			}
			return
		}
	}
	if fn == nil && b.cha && b.closureCalled(info, c) == nil {
		if tv := info.Types[c.Fun]; !tv.IsType() && !tv.IsBuiltin() {
			if sig, ok := tv.Type.Underlying().(*types.Signature); ok {
				// A call through a function value may reach any function of
				// its shape used as a value.
				k := funcShapeKey(sig)
				b.callers[k] = append(b.callers[k], me)
				if b.blocks[k] {
					b.mark(me)
				}
				return
			}
		}
	}
	if fn == nil {
		if lit := b.closureCalled(info, c); lit != nil {
			callee := b.key(lit.Pos())
			b.callers[callee] = append(b.callers[callee], me)
			if b.blocks[callee] {
				b.mark(me)
			}
		}
		return
	}
	if blockingNatives[nativeName(fn)] || b.cha && twinNatives[nativeName(fn)] {
		b.mark(me)
		return
	}
	if s := fn.Origin().Scope(); s != nil && s.Pos().IsValid() {
		callee := b.key(s.Pos())
		b.callers[callee] = append(b.callers[callee], me)
		if b.blocks[callee] {
			b.mark(me)
		}
	}
}

// closureCalled is the closure c calls through a variable holding nothing
// else, or nil.
func (b *blocking) closureCalled(info *types.Info, c *ast.CallExpr) *ast.FuncLit {
	if id, ok := ast.Unparen(c.Fun).(*ast.Ident); ok {
		return b.funcVars[info.Uses[id]]
	}
	return nil
}

// mark records that the function k blocks, and so do its callers.
func (b *blocking) mark(k string) {
	if b.blocks[k] {
		return
	}
	b.blocks[k] = true
	b.work = append(b.work, k)
	for len(b.work) > 0 {
		k := b.work[len(b.work)-1]
		b.work = b.work[:len(b.work)-1]
		for _, c := range b.callers[k] {
			if !b.blocks[c] {
				b.blocks[c] = true
				b.work = append(b.work, c)
			}
		}
	}
}

// clone is a copy to add a package's test files to.
func (b *blocking) clone() *blocking {
	c := newBlocking(b.fset, b.cha)
	for k, v := range b.blocks {
		c.blocks[k] = v
	}
	for k, v := range b.callers {
		c.callers[k] = append([]string(nil), v...)
	}
	for k, v := range b.funcVars {
		c.funcVars[k] = v
	}
	return c
}

func selectHasDefault(s *ast.SelectStmt) bool {
	for _, c := range s.Body.List {
		if c.(*ast.CommClause).Comm == nil {
			return true
		}
	}
	return false
}

// markCommOp marks the send or receive a select case does.
func markCommOp(s ast.Stmt, ops map[ast.Node]bool) {
	switch s := s.(type) {
	case *ast.SendStmt:
		ops[s] = true
	case *ast.ExprStmt:
		ops[ast.Unparen(s.X)] = true
	case *ast.AssignStmt:
		ops[ast.Unparen(s.Rhs[0])] = true
	}
}
