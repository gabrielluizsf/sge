package gots

import (
	"go/ast"
	"go/types"
	"strings"
)

// isBenchmark reports whether d is a Go benchmark, func BenchmarkX(b *testing.B).
func isBenchmark(d *ast.FuncDecl, fn *types.Func) bool {
	if !strings.HasPrefix(d.Name.Name, "Benchmark") {
		return false
	}
	sig := fn.Type().(*types.Signature)
	if sig.Params().Len() != 1 || sig.Results().Len() != 0 {
		return false
	}
	return typeString(sig.Params().At(0).Type()) == "*testing.B"
}
