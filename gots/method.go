package gots

import (
	"go/ast"
	"go/types"
)

// method writes a method: inside its class when inClass, otherwise as a
// function taking the receiver first, named as methodFunc says.
// A pointer method of a struct is both: the function, which Go code
// calls, written to funcs, and a class method forwarding to it.
func (e *emitter) method(w, funcs *writer, fd *ast.FuncDecl, inClass bool) {
	fn := e.info.Defs[fd.Name].(*types.Func)
	if e.methodSubst != nil {
		e.methodSubst(fn)
	}
	sig := e.T(fn.Type()).(*types.Signature)
	if _, ptr := sig.Recv().Type().(*types.Pointer); ptr && inClass {
		e.forward(w, fn, sig)
		w, inClass = funcs, false
	}
	recvField := fd.Recv.List[0]
	var recvObj types.Object
	if len(recvField.Names) > 0 {
		recvObj = e.info.Defs[recvField.Names[0]]
	}
	outer := e.sc
	e.sc = newScope(e.pkgNames)
	defer func() { e.sc = outer }()

	if inClass {
		var prologue []string
		if recvObj != nil && recvObj.Name() != "_" {
			n := e.declare(recvObj)
			_, ptr := sig.Recv().Type().(*types.Pointer)
			self := "this"
			if !ptr && mutates(e.info, fd.Body, recvObj) {
				self = "this.$clone()"
			}
			prologue = append(prologue, "let "+n+" = "+self+";")
		}
		if e.wrapped(fn) {
			e.wrapper(w, fieldName(fn.Name()), "this."+e.goName(fn), sig)
		}
		e.functionIn(w, e.goName(fn), sig, fd.Type, fd.Body, prologue)
		return
	}
	// A receiver parameter, first; Cell-boxed if the body takes its address.
	recvParam := types.NewParam(recvField.Pos(), e.tpkg, "", sig.Recv().Type())
	params := types.NewTuple(append([]*types.Var{recvParam}, tupleVars(sig.Params())...)...)
	withRecv := types.NewSignatureType(nil, nil, nil, params, sig.Results(), sig.Variadic())
	ft := &ast.FuncType{Params: &ast.FieldList{List: append([]*ast.Field{recvField}, fd.Type.Params.List...)}, Results: fd.Type.Results}
	e.functionIn(w, "export function "+e.methodFunc(fn), withRecv, ft, fd.Body, nil)
}

func tupleVars(t *types.Tuple) []*types.Var {
	var out []*types.Var
	for i := 0; i < t.Len(); i++ {
		out = append(out, t.At(i))
	}
	return out
}

// functionIn is function declared with the local names the caller has
// started, which already hold the receiver's.
func (e *emitter) functionIn(w *writer, prefix string, sig *types.Signature, ft *ast.FuncType, body *ast.BlockStmt, prologue []string) {
	e.function(w, prefix, sig, ft, body, prologue, false, false)
}

// mutates reports whether body may change the variable obj in place:
// assigning to it or its fields, incrementing them, taking an address in
// it, or calling a pointer method on it.
func mutates(info *types.Info, body *ast.BlockStmt, obj types.Object) bool {
	if body == nil {
		return false
	}
	rooted := func(x ast.Expr) bool {
		for {
			switch y := ast.Unparen(x).(type) {
			case *ast.Ident:
				return info.ObjectOf(y) == obj
			case *ast.SelectorExpr:
				if s := info.Selections[y]; s == nil || isPointer(s.Recv()) && y.X != nil {
					if s == nil {
						return false
					}
				}
				x = y.X
			case *ast.IndexExpr:
				if _, isArr := info.TypeOf(y.X).Underlying().(*types.Array); !isArr {
					return false
				}
				x = y.X
			default:
				return false
			}
		}
	}
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.AssignStmt:
			for _, l := range n.Lhs {
				found = found || rooted(l)
			}
		case *ast.IncDecStmt:
			found = found || rooted(n.X)
		case *ast.UnaryExpr:
			found = found || n.Op.String() == "&" && rooted(n.X)
		case *ast.RangeStmt:
			found = found || n.Key != nil && rooted(n.Key) || n.Value != nil && rooted(n.Value)
		case *ast.CallExpr:
			if sel, ok := ast.Unparen(n.Fun).(*ast.SelectorExpr); ok {
				if s := info.Selections[sel]; s != nil && s.Kind() == types.MethodVal {
					if _, ptr := s.Obj().Type().(*types.Signature).Recv().Type().(*types.Pointer); ptr {
						found = found || rooted(sel.X)
					}
				}
			}
		}
		return !found
	})
	return found
}
