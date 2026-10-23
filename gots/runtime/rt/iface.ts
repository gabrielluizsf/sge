import { Type, ptrTo } from "./types.ts";
import { Cell, ptrEq } from "./ptr.ts";

// Val is an interface holding anything but a pointer to a struct, which an
// interface holds as the object itself. A type with methods holds its values
// in a subclass of Val that has them, so calling through an interface is an
// ordinary method call.
export class Val {
  $t: Type;
  v: any;
  constructor(t: Type, v: any) {
    this.$t = t;
    this.v = v;
  }
}

// val puts v in an interface as a value of type t.
export function val(t: Type, v: any): Val {
  return t.valCls === null ? new Val(t, v) : new t.valCls(t, v);
}

// typeOf is the dynamic type of a non-nil interface value.
export function typeOf(x: any): Type {
  if (x instanceof Val) return x.$t;
  return ptrTo(x.constructor.$type);
}

// is reports whether interface value x holds a value of the non-pointer type t.
export function is(x: any, t: Type): boolean {
  return x instanceof Val && x.$t === t;
}

// implements reports whether x has every method an interface asks for.
export function implementsI(x: any, methods: string[]): boolean {
  if (x === null) return false;
  for (const m of methods) if (typeof x[m] !== "function") return false;
  return true;
}

export class TypeAssertionError extends Error {
  constructor(from: any, want: string) {
    super(`interface conversion: interface {} is ${from === null ? "nil" : typeOf(from).str}, not ${want}`);
  }
}

// assertVal is x.(T) for a non-pointer T: the value, or a panic.
export function assertVal(x: any, t: Type): any {
  if (x instanceof Val && x.$t === t) return x.v;
  throw new TypeAssertionError(x, t.str);
}

// isPtrTo reports whether interface value x holds a *S, for the struct
// class ctor of S. Descriptors are compared rather than classes, since an
// instance of a generic struct has a class in every module that uses it.
export function isPtrTo(x: any, ctor: any): boolean {
  return x !== null && !(x instanceof Val) && x.constructor.$type === ctor.$type;
}

// assertPtr is x.(*S) for a struct S: the object, or a panic.
export function assertPtr(x: any, ctor: any): any {
  if (isPtrTo(x, ctor)) return x;
  throw new TypeAssertionError(x, "*" + ctor.$type.str);
}

// assertIface is x.(I) for an interface I.
export function assertIface(x: any, str: string, methods: string[]): any {
  if (implementsI(x, methods)) return x;
  if (x === null) throw new TypeAssertionError(x, str);
  throw new Error(`interface conversion: ${typeOf(x).str} is not ${str}: missing method`);
}

// ifaceEq is a == b for interface values.
export function ifaceEq(a: any, b: any): boolean {
  if (a === b) return true;
  if (a === null || b === null) return false;
  const av = a instanceof Val, bv = b instanceof Val;
  if (av !== bv) return false;
  if (!av) return false;
  if (a.$t !== b.$t) return false;
  return equal(a.$t, a.v, b.v);
}

// equal is a == b for two values of type t.
export function equal(t: Type, a: any, b: any): boolean {
  switch (t.kind) {
    case "struct":
      for (const f of t.fields) if (!equal(f.type, a[f.name], b[f.name])) return false;
      return true;
    case "array":
      for (let i = 0; i < t.len; i++) if (!equal(t.elem!, a[i], b[i])) return false;
      return true;
    case "iface":
      return ifaceEq(a, b);
    case "ptr":
      return a instanceof Cell || b instanceof Cell ? a === b : ptrEq(a, b);
    case "slice":
    case "map":
    case "func":
      if (a === null && b === null) return true;
      throw new Error(`runtime error: comparing uncomparable type ${t.str}`);
    case "complex64":
    case "complex128":
      return a.re === b.re && a.im === b.im;
    default:
      return a === b;
  }
}
