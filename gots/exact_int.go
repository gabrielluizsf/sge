package gots

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
)

// maxSafe is the largest integer a JavaScript number holds exactly.
var maxSafe = constant.MakeInt64(1<<53 - 1)

// numberInt reports whether t is an integer type gots keeps in a
// JavaScript number though Go gives it 64 bits.
func numberInt(t types.Type) bool {
	switch basicKind(t) {
	case types.Int, types.Int64, types.Uintptr:
		return true
	}
	return false
}

// hugeTree reports whether x is integer arithmetic on numbers with a
// constant in it beyond 2^53. Such constants usually cancel out, as the
// offsets time works out dates with do, but rounded to a number they
// leave an error behind.
func (e *emitter) hugeTree(x ast.Expr) bool {
	if tv, ok := e.info.Types[x]; ok && tv.Value != nil {
		v := constant.ToInt(tv.Value)
		if v.Kind() != constant.Int {
			return false
		}
		return constant.Compare(v, token.GTR, maxSafe) || constant.Compare(constant.UnaryOp(token.SUB, v, 0), token.GTR, maxSafe)
	}
	switch x := ast.Unparen(x).(type) {
	case *ast.BinaryExpr:
		if !isBigArith(x.Op) || !numberInt(e.typeOf(x)) {
			return false
		}
		return e.hugeTree(x.X) || e.hugeTree(x.Y)
	case *ast.UnaryExpr:
		return (x.Op == token.SUB || x.Op == token.XOR) && numberInt(e.typeOf(x)) && e.hugeTree(x.X)
	}
	return false
}

func isUnsigned(t types.Type) bool {
	b, ok := t.Underlying().(*types.Basic)
	return ok && b.Info()&types.IsUnsigned != 0
}

func isBigArith(op token.Token) bool {
	switch op {
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM, token.AND, token.OR, token.XOR, token.SHL, token.SHR, token.AND_NOT:
		return true
	}
	return false
}

// bigTree is x, a tree hugeTree says yes to, as an exact BigInt wrapped to
// x's 64 bits after every operation, as Go wraps.
func (e *emitter) bigTree(x ast.Expr) string {
	t := e.typeOf(x)
	wrap := func(js string) string {
		if isUnsigned(t) {
			return "$BigInt.asUintN(64, " + js + ")"
		}
		return "$BigInt.asIntN(64, " + js + ")"
	}
	if tv, ok := e.info.Types[x]; ok && tv.Value != nil {
		return wrap(constant.ToInt(tv.Value).ExactString() + "n")
	}
	switch y := ast.Unparen(x).(type) {
	case *ast.BinaryExpr:
		if isBigArith(y.Op) && numberInt(t) {
			a, b := e.bigTree(y.X), e.bigTree(y.Y)
			if y.Op == token.AND_NOT {
				return wrap("(" + a + " & ~" + b + ")")
			}
			return wrap("(" + a + " " + y.Op.String() + " " + b + ")")
		}
	case *ast.UnaryExpr:
		if numberInt(t) {
			switch y.Op {
			case token.SUB:
				return wrap("(-" + e.bigTree(y.X) + ")")
			case token.XOR:
				return wrap("(~" + e.bigTree(y.X) + ")")
			}
		}
	case *ast.CallExpr:
		// A conversion from uint64, which is a BigInt already.
		if tv, ok := e.info.Types[y.Fun]; ok && tv.IsType() && len(y.Args) == 1 && isBig(e.typeOf(y.Args[0])) {
			return wrap(e.expr(y.Args[0]))
		}
	}
	if isBig(t) {
		return e.expr(x)
	}
	return "$BigInt(" + e.expr(x) + ")"
}
