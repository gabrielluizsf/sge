package gots

import (
	"go/ast"
	"go/token"
	"go/types"
)

// sharedRange reports whether the value variable of the range s over a
// slice or an array may be its element itself, not a copy: it is
// declared by the loop, is only read, kept by no closure or pointer, and
// nothing the body does can change the element while it is read, for the
// body writes only the function's own variables held by value, never the
// one ranged over, and calls only pure functions. A copy each time round
// is an allocation V8 cannot avoid.
func (e *emitter) sharedRange(s *ast.RangeStmt) bool {
	id, ok := s.Value.(*ast.Ident)
	if !ok || s.Tok != token.DEFINE || id.Name == "_" || e.fn == nil || e.fn.body == nil || e.deferring {
		return false
	}
	v, ok := e.info.Defs[id].(*types.Var)
	if !ok || e.boxed[v] || mutates(e.info, s.Body, v) || pinned(e.info, s.Body, v) {
		return false
	}
	ranged := rootVar(e.info, s.X)
	own := func(x ast.Expr) bool {
		if id, ok := ast.Unparen(x).(*ast.Ident); ok && id.Name == "_" {
			return true
		}
		r, ok := rootVar(e.info, x).(*types.Var)
		return ok && r != ranged && !r.IsField() && r.Pos() >= e.fn.ft.Pos() && r.Pos() < e.fn.body.End()
	}
	pure := true
	ast.Inspect(s.Body, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.GoStmt, *ast.DeferStmt, *ast.SendStmt, *ast.SelectStmt, *ast.FuncLit:
			pure = false
		case *ast.UnaryExpr:
			pure = pure && n.Op != token.ARROW
		case *ast.AssignStmt:
			for _, l := range n.Lhs {
				pure = pure && own(l)
			}
		case *ast.IncDecStmt:
			pure = pure && own(n.X)
		case *ast.RangeStmt:
			if n.Tok == token.ASSIGN {
				pure = pure && (n.Key == nil || own(n.Key)) && (n.Value == nil || own(n.Value))
			}
			if _, ok := e.info.TypeOf(n.X).Underlying().(*types.Chan); ok {
				pure = false
			}
		case *ast.CallExpr:
			pure = pure && e.pureCall(e.info, n)
		}
		return pure
	})
	return pure
}
