package gots

import (
	"go/ast"
	"go/types"
)

// ifaceKey names a method by its name and signature, without the names of
// its parameters: a call through an interface may reach any method of
// that name and signature.
func ifaceKey(fn *types.Func) string {
	sig := fn.Origin().Type().(*types.Signature)
	q := func(p *types.Package) string { return p.Path() }
	k := "iface:" + fn.Name() + "("
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

// ifaceMethodCalled is the interface method c calls, or nil for a call of
// anything else.
func ifaceMethodCalled(info *types.Info, c *ast.CallExpr) *types.Func {
	sel, ok := ast.Unparen(c.Fun).(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	s := info.Selections[sel]
	if s == nil || s.Kind() != types.MethodVal || !types.IsInterface(s.Recv()) {
		return nil
	}
	m, _ := s.Obj().(*types.Func)
	return m
}
