package gots

import "strconv"

// hoistConst is expr, a value that never changes, as a module-level
// constant made once rather than on every evaluation, shared by every use
// of the same expression.
func (e *emitter) hoistConst(expr string) string {
	if c, ok := e.typeConsts[expr]; ok {
		return c
	}
	c := "$T" + strconv.Itoa(len(e.typeConsts)+1)
	e.typeConsts[expr] = c
	e.typeOrder = append(e.typeOrder, expr)
	return c
}
