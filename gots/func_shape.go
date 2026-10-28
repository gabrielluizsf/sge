package gots

import (
	"fmt"
	"go/ast"
	"go/types"
)

// funcShapeKey names the functions a call through a function value of
// type sig may reach: every one of its signature used as a value, or, for
// a signature with type parameters, as a generic iterator's yield has, of
// its shape, its counts of parameters and results.
func funcShapeKey(sig *types.Signature) string {
	if hasTypeParam(sig) {
		return shapeKey(sig)
	}
	return exactKey(sig)
}

func shapeKey(sig *types.Signature) string {
	return fmt.Sprintf("fv:%d:%d:%v", sig.Params().Len(), sig.Results().Len(), sig.Variadic())
}

// exactKey names sig by its parameters' and results' types.
func exactKey(sig *types.Signature) string {
	q := func(p *types.Package) string { return p.Path() }
	k := "fs:("
	for i := 0; i < sig.Params().Len(); i++ {
		k += types.TypeString(sig.Params().At(i).Type(), q) + ","
	}
	if sig.Variadic() {
		k += "..."
	}
	k += ")"
	for i := 0; i < sig.Results().Len(); i++ {
		k += types.TypeString(sig.Results().At(i).Type(), q) + ","
	}
	return k
}

// hasTypeParam reports whether t mentions a type parameter.
func hasTypeParam(t types.Type) bool {
	switch t := t.(type) {
	case *types.TypeParam:
		return true
	case *types.Pointer:
		return hasTypeParam(t.Elem())
	case *types.Slice:
		return hasTypeParam(t.Elem())
	case *types.Array:
		return hasTypeParam(t.Elem())
	case *types.Chan:
		return hasTypeParam(t.Elem())
	case *types.Map:
		return hasTypeParam(t.Key()) || hasTypeParam(t.Elem())
	case *types.Signature:
		for _, tup := range []*types.Tuple{t.Params(), t.Results()} {
			for i := 0; i < tup.Len(); i++ {
				if hasTypeParam(tup.At(i).Type()) {
					return true
				}
			}
		}
	case *types.Named:
		if args := t.TypeArgs(); args != nil {
			for i := 0; i < args.Len(); i++ {
				if hasTypeParam(args.At(i)) {
					return true
				}
			}
		}
	case *types.Alias:
		return hasTypeParam(types.Unalias(t))
	}
	return false
}

// rangeFuncKey names the body of a range over a function at pos: the
// yield callback the loop hands its iterator.
func rangeFuncKey(pos string) string { return "rf:" + pos }

// funcValue notes that the function known as k, of type t, is used as a
// value: if it may block, so may a call through a value of its shape.
func (b *blocking) funcValue(k string, t types.Type) {
	if !b.cha || t == nil {
		return
	}
	sig, ok := t.Underlying().(*types.Signature)
	if !ok {
		return
	}
	// Under its signature, for callers that know it, and its shape, for
	// generic ones that do not.
	keys := []string{shapeKey(sig)}
	if !hasTypeParam(sig) {
		keys = append(keys, exactKey(sig))
	}
	for _, fk := range keys {
		b.callers[k] = append(b.callers[k], fk)
		if b.blocks[k] {
			b.mark(fk)
		}
	}
}

// addFuncValues notes the functions and methods f uses as values, not
// calling them where it names them: handlers given to a server, say.
func (b *blocking) addFuncValues(f *ast.File, info *types.Info) {
	called := map[ast.Expr]bool{}
	ast.Inspect(f, func(n ast.Node) bool {
		if c, ok := n.(*ast.CallExpr); ok {
			called[ast.Unparen(c.Fun)] = true
			if ix, ok := ast.Unparen(c.Fun).(*ast.IndexExpr); ok {
				called[ix.X] = true // f[T](...)
			}
		}
		return true
	})
	ast.Inspect(f, func(n ast.Node) bool {
		var id *ast.Ident
		switch x := n.(type) {
		case *ast.SelectorExpr:
			if called[x] {
				return false
			}
			id = x.Sel
		case *ast.Ident:
			if called[x] {
				return true
			}
			id = x
		default:
			return true
		}
		fn, ok := info.Uses[id].(*types.Func)
		if !ok {
			return true
		}
		if s := fn.Origin().Scope(); s != nil && s.Pos().IsValid() {
			b.funcValue(b.key(s.Pos()), fn.Type())
		}
		return true
	})
}
