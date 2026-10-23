package gots

import (
	"go/ast"
	"go/types"
	"strings"
)

func (e *emitter) call(c *ast.CallExpr) string {
	if tv, ok := e.info.Types[c.Fun]; ok {
		if tv.IsType() {
			return e.convert(c.Args[0], e.T(tv.Type))
		}
		if tv.IsBuiltin() {
			return e.builtin(c)
		}
	}
	sig, ok := e.typeOf(c.Fun).Underlying().(*types.Signature)
	if !ok {
		return e.unsupported(c, "call of a non-function")
	}
	return e.callWith(c, e.args(c, sig))
}

// callWith is the call c with its arguments already written.
func (e *emitter) callWith(c *ast.CallExpr, args []string) string {
	list := strings.Join(args, ", ")
	if sel, ok := ast.Unparen(c.Fun).(*ast.SelectorExpr); ok {
		if s := e.info.Selections[sel]; s != nil && s.Kind() == types.MethodVal {
			fn := s.Obj().(*types.Func)
			recv := e.receiver(sel, s)
			if onClass(fn) {
				return recv + "." + fieldName(fn.Name()) + "(" + list + ")"
			}
			if list != "" {
				list = ", " + list
			}
			return e.methodFunc(fn) + "(" + recv + list + ")"
		}
	}
	fun := e.expr(c.Fun)
	if _, isLit := ast.Unparen(c.Fun).(*ast.FuncLit); isLit {
		fun = "(" + fun + ")"
	}
	return fun + "(" + list + ")"
}

// args writes a call's arguments as the callee's parameters receive them:
// copied, put in interfaces, and with variadic ones gathered in a slice.
func (e *emitter) args(c *ast.CallExpr, sig *types.Signature) []string {
	params := sig.Params()
	if len(c.Args) == 1 {
		if tup, ok := e.typeOf(c.Args[0]).(*types.Tuple); ok && tup.Len() > 1 {
			return e.spreadTuple(c.Args[0], tup, sig)
		}
	}
	var out []string
	n := params.Len()
	if sig.Variadic() {
		n--
	}
	for i := 0; i < n && i < len(c.Args); i++ {
		out = append(out, e.value(c.Args[i], params.At(i).Type()))
	}
	if !sig.Variadic() {
		return out
	}
	last := params.At(n).Type().(*types.Slice)
	if c.Ellipsis.IsValid() {
		arg := c.Args[n]
		if isString(e.typeOf(arg)) {
			return append(out, "$.slice($.str.strToBytes("+e.expr(arg)+"))")
		}
		return append(out, e.expr(arg))
	}
	var rest []string
	for _, a := range c.Args[n:] {
		rest = append(rest, e.value(a, last.Elem()))
	}
	return append(out, e.sliceLit(last.Elem(), rest))
}

// spreadTuple passes the results of a multi-value call as the arguments of
// another, f(g()).
func (e *emitter) spreadTuple(arg ast.Expr, tup *types.Tuple, sig *types.Signature) []string {
	t := e.newTmp()
	var vals []string
	for i := 0; i < tup.Len(); i++ {
		var to types.Type
		switch {
		case sig.Variadic() && i >= sig.Params().Len()-1:
			to = sig.Params().At(sig.Params().Len() - 1).Type().(*types.Slice).Elem()
		default:
			to = sig.Params().At(i).Type()
		}
		vals = append(vals, e.convertValue(t+"["+itoa(i)+"]", tup.At(i).Type(), to, false))
	}
	if sig.Variadic() {
		n := sig.Params().Len() - 1
		elem := sig.Params().At(n).Type().(*types.Slice).Elem()
		vals = append(vals[:n], e.sliceLit(elem, vals[n:]))
	}
	// The tuple is bound by a default parameter so the call stays one
	// expression.
	return []string{"...((" + t + ": any = " + e.expr(arg) + ") => [" + strings.Join(vals, ", ") + "])()"}
}

// sliceLit is a slice holding vals, or nil for none.
func (e *emitter) sliceLit(elem types.Type, vals []string) string {
	if len(vals) == 0 {
		return "null"
	}
	return "$.slice(" + e.arrayLit(elem, vals) + ")"
}

// arrayLit is a backing array holding vals.
func (e *emitter) arrayLit(elem types.Type, vals []string) string {
	if b := backing(elem); b != "" {
		return b + ".of(" + strings.Join(vals, ", ") + ")"
	}
	return "[" + strings.Join(vals, ", ") + "]"
}
