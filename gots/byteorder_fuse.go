package gots

import (
	"go/ast"
	"go/types"
	"strings"
)

// byteorderFused are the byte-order functions whose calls on a subslice,
// b[lo:] or b[lo:hi], are written as rt/bytes.ts's, on b and the bounds:
// the subslice would be an object made only to read four bytes from.
var byteorderFused = map[string]string{}

func init() {
	for _, pkg := range []string{"internal/byteorder", "crypto/internal/fips140deps/byteorder"} {
		for _, f := range []string{"BEUint16", "BEUint32", "LEUint16", "LEUint32", "BEPutUint16", "BEPutUint32", "LEPutUint16", "LEPutUint32"} {
			name := strings.ToLower(f[:2]) + strings.TrimPrefix(f[2:], "Uint")
			if strings.HasPrefix(f[2:], "Put") {
				name = "put" + f[:2] + strings.TrimPrefix(f[5:], "Uint")
			}
			byteorderFused[pkg+"."+f] = name
		}
	}
	// encoding/binary's, which are methods of stateless values.
	for typ, pre := range map[string]string{"bigEndian": "be", "littleEndian": "le"} {
		for _, bits := range []string{"16", "32"} {
			byteorderFused["encoding/binary."+typ+".Uint"+bits] = pre + bits
			byteorderFused["encoding/binary."+typ+".PutUint"+bits] = "put" + strings.ToUpper(pre) + bits
		}
	}
}

// fusedByteorder is c written as rt/bytes.ts's, if c calls one of
// byteorderFused on a slice.
func (e *emitter) fusedByteorder(c *ast.CallExpr) (string, bool) {
	fn := e.staticCallee(c)
	if fn == nil {
		return "", false
	}
	name, ok := byteorderFused[nativeName(fn)]
	if !ok || !isSliceType(e.typeOf(c.Args[0])) {
		return "", false
	}
	s, bounds := e.fusedBounds(ast.Unparen(c.Args[0]))
	put := strings.HasPrefix(name, "put")
	if put && len(bounds) > 1 {
		return "", false
	}
	args := append([]string{s}, bounds...)
	if put {
		args = append(args, e.expr(c.Args[1]))
	}
	return "$.bytes." + name + "(" + strings.Join(args, ", ") + ")", true
}

// fusedBounds is the slice x reads from, and lo, with hi if x gives one.
func (e *emitter) fusedBounds(x ast.Expr) (string, []string) {
	se, ok := x.(*ast.SliceExpr)
	if !ok || se.Slice3 || !isSliceType(e.typeOf(se.X)) {
		return e.expr(x), []string{"0"}
	}
	bounds := []string{"0"}
	if se.Low != nil {
		bounds[0] = e.indexOf(se.Low)
	}
	if se.High != nil {
		bounds = append(bounds, e.indexOf(se.High))
	}
	return e.expr(se.X), bounds
}

func isSliceType(t types.Type) bool {
	_, ok := t.Underlying().(*types.Slice)
	return ok
}
