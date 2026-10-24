package gots

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strconv"
)

func (e *emitter) binary(x *ast.BinaryExpr) string {
	lt, rt := e.typeOf(x.X), e.typeOf(x.Y)
	switch x.Op {
	case token.LAND, token.LOR:
		return "(" + e.expr(x.X) + " " + x.Op.String() + " " + e.expr(x.Y) + ")"
	case token.EQL, token.NEQ:
		eq := e.equality(x.X, x.Y, lt, rt)
		if x.Op == token.NEQ {
			return "!" + eq
		}
		return eq
	case token.LSS, token.LEQ, token.GTR, token.GEQ:
		return "(" + e.expr(x.X) + " " + x.Op.String() + " " + e.expr(x.Y) + ")"
	}
	if e.hugeTree(x) {
		// Exact, and a number again once the constants are spent.
		return "$Number(" + e.bigTree(x) + ")"
	}
	var count *int
	if x.Op == token.SHL || x.Op == token.SHR {
		if tv := e.info.Types[x.Y]; tv.Value != nil {
			if n, ok := constant.Int64Val(constant.ToInt(tv.Value)); ok {
				c := int(n)
				count = &c
			}
		}
	}
	b := e.expr(x.Y)
	if (x.Op == token.SHL || x.Op == token.SHR) && isBig(rt) {
		b = "$Number(" + b + ")"
	}
	return e.arith(x.Op, e.expr(x.X), b, e.typeOf(x), count)
}

func isNil(t types.Type) bool {
	b, ok := t.(*types.Basic)
	return ok && b.Kind() == types.UntypedNil
}

// equality is a == b.
func (e *emitter) equality(a, b ast.Expr, at, bt types.Type) string {
	switch {
	case isNil(at):
		return "(" + e.expr(b) + " === null)"
	case isNil(bt):
		return "(" + e.expr(a) + " === null)"
	case isInterface(at) || isInterface(bt):
		return "$.ifaceEq(" + e.value(a, types.NewInterfaceType(nil, nil)) + ", " + e.value(b, types.NewInterfaceType(nil, nil)) + ")"
	}
	return e.equalOf(e.expr(a), e.expr(b), at)
}

// arith is a op b for operands of type t. count is the constant count of a
// shift, when it is one.
func (e *emitter) arith(op token.Token, a, b string, t types.Type, count *int) string {
	k := defaultKind(basicKind(t))
	o := op.String()
	bin := "(" + a + " " + o + " " + b + ")"
	switch {
	case isString(t):
		return bin
	case k == types.Float32:
		return "$Math.fround" + bin
	case isFloat(t):
		return bin
	case k == types.Complex64 || k == types.Complex128:
		names := map[token.Token]string{token.ADD: "add", token.SUB: "sub", token.MUL: "mul", token.QUO: "div"}
		js := "$.cplx." + names[op] + "(" + a + ", " + b + ")"
		if k == types.Complex64 {
			return "$.cplx.c64(" + js + ")"
		}
		return js
	case k == types.Uint64:
		return e.arithBig(op, a, b, count)
	case small(k):
		return e.arithSmall(op, k, a, b, count)
	}
	switch op {
	case token.ADD, token.SUB, token.MUL, token.REM:
		return bin
	case token.QUO:
		return "$.int.div(" + a + ", " + b + ")"
	case token.AND:
		return "$.int.and(" + a + ", " + b + ")"
	case token.OR:
		return "$.int.or(" + a + ", " + b + ")"
	case token.XOR:
		return "$.int.xor(" + a + ", " + b + ")"
	case token.AND_NOT:
		return "$.int.andNot(" + a + ", " + b + ")"
	case token.SHL:
		if count != nil && *count < 53 {
			return "(" + a + " * " + strconv.FormatUint(1<<*count, 10) + ")"
		}
		return "$.int.shl(" + a + ", " + b + ")"
	case token.SHR:
		if count != nil && *count < 53 {
			return "$Math.floor(" + a + " / " + strconv.FormatUint(1<<*count, 10) + ")"
		}
		return "$.int.shr(" + a + ", " + b + ")"
	}
	return "$.panic(\"gots: operator " + o + "\")"
}

// arithSmall is arithmetic on integers of 32 bits or fewer.
func (e *emitter) arithSmall(op token.Token, k types.BasicKind, a, b string, count *int) string {
	bin := "(" + a + " " + op.String() + " " + b + ")"
	w := width(k)
	switch op {
	case token.ADD, token.SUB:
		return wrap(k, a+" "+op.String()+" "+b)
	case token.MUL:
		if k == types.Int32 {
			return "$Math.imul(" + a + ", " + b + ")"
		}
		return wrap(k, "$Math.imul("+a+", "+b+")")
	case token.QUO:
		return wrap(k, a+" / "+b)
	case token.REM:
		return bin
	case token.AND, token.OR, token.XOR:
		if k == types.Uint32 {
			return wrap(k, a+" "+op.String()+" "+b)
		}
		return bin
	case token.AND_NOT:
		if k == types.Uint32 {
			return wrap(k, a+" & ~"+b)
		}
		return "(" + a + " & ~" + b + ")"
	case token.SHL:
		if count != nil {
			if *count >= w {
				return "0"
			}
			return wrap(k, a+" << "+b)
		}
		return wrap(k, "$.int.shl32("+a+", "+b+")")
	case token.SHR:
		if unsigned(k) {
			if count != nil {
				if *count >= w {
					return "0"
				}
				return "(" + a + " >>> " + b + ")"
			}
			return "$.int.shrU32(" + a + ", " + b + ")"
		}
		if count != nil {
			if *count >= w {
				return "(" + a + " >> 31)"
			}
			return "(" + a + " >> " + b + ")"
		}
		return "$.int.shr32(" + a + ", " + b + ")"
	}
	return "$.panic(\"gots: operator " + op.String() + "\")"
}

// arithBig is arithmetic on uint64, kept in a bigint.
func (e *emitter) arithBig(op token.Token, a, b string, count *int) string {
	switch op {
	case token.ADD, token.SUB, token.MUL:
		return "$.int.u64(" + a + " " + op.String() + " " + b + ")"
	case token.QUO, token.REM, token.AND, token.OR, token.XOR:
		return "(" + a + " " + op.String() + " " + b + ")"
	case token.AND_NOT:
		return "(" + a + " & ~" + b + ")"
	case token.SHL:
		return "$.int.shlU64(" + a + ", " + b + ")"
	case token.SHR:
		return "$.int.shrU64(" + a + ", " + b + ")"
	}
	return "$.panic(\"gots: operator " + op.String() + "\")"
}
