// reflect.Type: one rtype for each run-time descriptor, so types compare
// with == and key maps, as Go's do.
import * as $ from "../../rt/mod.ts";
import { Kind, kindOf, Kind$String } from "./kinds.ts";
import { align, size } from "./zero.ts";
import { methodsOf } from "./call.ts";
import { ifaceMethods } from "./value.ts";

const rtypes = new WeakMap<$.Type, rtype>();

// toType is the reflect.Type for the descriptor d.
export function toType(d: $.Type | null): rtype | null {
  if (d === null) return null;
  let r = rtypes.get(d);
  if (r === undefined) {
    r = new rtype(d);
    rtypes.set(d, r);
  }
  return r;
}

function badKind(m: string, t: $.Type): never {
  $.panic($.val($.String, $.str.toGo("reflect: " + m + " of non-" + m.split(".")[0] + " type " + t.str)));
}

// protoOf is the prototype an interface holding a value of type d carries
// its methods on.
function protoOf(d: $.Type): any {
  if (d.kind === "ptr" && d.elem!.kind === "struct") return d.elem!.ctor?.prototype ?? null;
  return d.valCls?.prototype ?? null;
}

export class rtype {
  constructor(public d: $.Type) {}

  Kind(): Kind {
    return kindOf(this.d);
  }
  String(): string {
    return $.str.toGo(this.d.str);
  }
  Name(): string {
    if (!this.d.named) {
      // Only the predeclared types have names without being defined.
      const k = kindOf(this.d);
      return k >= 1 && k <= 16 || k === 24 ? $.str.toGo(this.d.str) : "";
    }
    const s = this.d.str, br = s.indexOf("[");
    const head = br < 0 ? s : s.slice(0, br);
    return $.str.toGo(s.slice(head.lastIndexOf(".") + 1));
  }
  PkgPath(): string {
    return "";
  }
  Elem(): rtype {
    const k = this.d.kind;
    if (k !== "ptr" && k !== "slice" && k !== "array" && k !== "map" && k !== "chan") badKind("Elem", this.d);
    return toType(this.d.elem)!;
  }
  Key(): rtype {
    if (this.d.kind !== "map") badKind("Key", this.d);
    return toType(this.d.key)!;
  }
  Len(): number {
    if (this.d.kind !== "array") badKind("Len", this.d);
    return this.d.len;
  }
  NumField(): number {
    if (this.d.kind !== "struct") badKind("NumField", this.d);
    return this.d.fields.length;
  }
  Field(i: number): StructField {
    if (this.d.kind !== "struct") badKind("Field", this.d);
    const f = this.d.fields[i];
    if (f === undefined) $.panic($.val($.String, $.str.toGo("reflect: Field index out of bounds")));
    const sf = new StructField();
    sf.Name = $.str.toGo(f.name);
    sf.PkgPath = f.exported ? "" : "main";
    sf.Type = toType(f.type);
    sf.Tag = f.tag;
    sf.Index = $.slice(Float64Array.of(i));
    sf.Anonymous = f.embedded;
    let off = 0;
    for (let j = 0; j <= i; j++) {
      const a = align(this.d.fields[j].type);
      off = Math.ceil(off / a) * a;
      if (j < i) off += size(this.d.fields[j].type);
    }
    sf.Offset = off;
    return sf;
  }
  FieldByName(name: string): [StructField, boolean] {
    if (this.d.kind !== "struct") badKind("FieldByName", this.d);
    const i = this.d.fields.findIndex((f) => $.str.toGo(f.name) === name);
    return i < 0 ? [new StructField(), false] : [this.Field(i), true];
  }
  NumMethod(): number {
    if (this.d.kind === "iface") return this.d.methods.length;
    return methodsOf(this.d).length;
  }
  // Method is the i-th exported method, by name, its type taking the
  // receiver first.
  Method(i: number): Method {
    if (this.d.kind === "iface") {
      const im = ifaceMethods(this.d)[i];
      if (im === undefined) $.panic($.val($.String, $.str.toGo("reflect: Method index out of range")));
      const m = new Method();
      m.Name = $.str.toGo(im.name);
      m.Type = toType(im.type);
      m.Index = i;
      return m;
    }
    const ms = methodsOf(this.d);
    const mt = ms[i];
    if (mt === undefined) $.panic($.val($.String, $.str.toGo("reflect: Method index out of range")));
    const m = new Method();
    m.Name = $.str.toGo(mt.name);
    m.Type = toType(funcOf([this.d, ...mt.type.ins], mt.type.outs, mt.type.variadic));
    m.Index = i;
    m.recv = this.d;
    return m;
  }
  MethodByName(name: string): [Method, boolean] {
    const names = this.d.kind === "iface" ? [...this.d.methods].sort() : methodsOf(this.d).map((m) => m.name);
    const i = names.indexOf($.str.fromGo(name));
    return i < 0 ? [new Method(), false] : [this.Method(i), true];
  }
  // Implements reports whether values of this type have every method of
  // the interface u.
  Implements(u: rtype | null): boolean {
    if (u === null) $.panic($.val($.String, $.str.toGo("reflect: nil type passed to Type.Implements")));
    if (u.d.kind !== "iface") $.panic($.val($.String, $.str.toGo("reflect: non-interface type passed to Type.Implements")));
    if (this.d.kind === "iface") return u.d.methods.every((m) => this.d.methods.includes(m));
    const p = protoOf(this.d);
    return u.d.methods.every((m) => p !== null && typeof p[m] === "function");
  }
  AssignableTo(u: rtype | null): boolean {
    if (u === null) return false;
    if (u.d === this.d) return true;
    if (u.d.kind === "iface") return this.Implements(u);
    return !this.d.named || !u.d.named ? this.d.str === u.d.str || sameUnder(this.d, u.d) : false;
  }
  ConvertibleTo(u: rtype | null): boolean {
    if (u === null) return false;
    const a = this.d, b = u.d;
    const num = (t: $.Type) => kindOf(t) >= 2 && kindOf(t) <= 14;
    if (num(a) && num(b)) return true;
    if (sameUnder(a, b)) return true;
    if (b.kind === "string" && (a.kind === "slice" && (a.elem!.kind === "uint8" || a.elem!.kind === "int32") || num(a) && !(a.kind === "float32" || a.kind === "float64"))) return true;
    if (a.kind === "string" && b.kind === "slice" && (b.elem!.kind === "uint8" || b.elem!.kind === "int32")) return true;
    return b.kind === "iface" && this.Implements(u);
  }
  Comparable(): boolean {
    switch (this.d.kind) {
      case "slice": case "map": case "func":
        return false;
      case "array":
        return toType(this.d.elem)!.Comparable();
      case "struct":
        return this.d.fields.every((f) => toType(f.type)!.Comparable());
    }
    return true;
  }
  Size(): number {
    return size(this.d);
  }
  Align(): number {
    return align(this.d);
  }
  FieldAlign(): number {
    return align(this.d);
  }
  Bits(): number {
    const k = kindOf(this.d);
    if (k < 2 || k > 16) badKind("Bits", this.d);
    return size(this.d) * 8;
  }
  ChanDir(): number {
    if (this.d.kind !== "chan") badKind("ChanDir", this.d);
    return this.d.str.startsWith("<-") ? 1 : this.d.str.startsWith("chan<-") ? 2 : 3;
  }
  IsVariadic(): boolean {
    return this.fn("IsVariadic").variadic;
  }
  private fn(what: string): $.Type {
    if (this.d.kind !== "func") badKind(what, this.d);
    return this.d;
  }
  NumIn(): number {
    return this.fn("NumIn").ins.length;
  }
  NumOut(): number {
    return this.fn("NumOut").outs.length;
  }
  In(i: number): rtype {
    return toType(this.fn("In").ins[i])!;
  }
  Out(i: number): rtype {
    return toType(this.fn("Out").outs[i])!;
  }
  // CanSeq and CanSeq2 report whether a value of this type can be ranged
  // over with one or two variables.
  CanSeq(): boolean {
    const d = this.d;
    switch (d.kind) {
      case "func":
        return d.ins.length === 1 && d.outs.length === 0 && d.ins[0].kind === "func" && d.ins[0].ins.length === 1 && d.ins[0].outs.length === 1;
      case "ptr":
        return d.elem!.kind === "array";
      case "slice": case "array": case "string": case "map": case "chan":
        return true;
    }
    const k = kindOf(d);
    return k >= 2 && k <= 12;
  }
  CanSeq2(): boolean {
    const d = this.d;
    switch (d.kind) {
      case "func":
        return d.ins.length === 1 && d.outs.length === 0 && d.ins[0].kind === "func" && d.ins[0].ins.length === 2 && d.ins[0].outs.length === 1;
      case "ptr":
        return d.elem!.kind === "array";
      case "slice": case "array": case "string": case "map":
        return true;
    }
    return false;
  }
  OverflowInt(x: number): boolean {
    return overflowInt(this.d, x);
  }
  OverflowUint(x: bigint): boolean {
    return overflowUint(this.d, x);
  }
  OverflowFloat(x: number): boolean {
    return overflowFloat(this.d, x);
  }
}
$.nativeType("reflect.rtype", rtype);

function sameUnder(a: $.Type, b: $.Type): boolean {
  if (a.kind !== b.kind) return false;
  switch (a.kind) {
    case "slice": case "ptr": case "chan":
      return a.elem === b.elem;
    case "array":
      return a.len === b.len && a.elem === b.elem;
    case "map":
      return a.key === b.key && a.elem === b.elem;
    case "struct":
      return a.fields.length === b.fields.length && a.fields.every((f, i) => f.name === b.fields[i].name && f.type === b.fields[i].type);
  }
  return kindOf(a) >= 1 && kindOf(a) <= 16 || a.kind === "string";
}

export function unsupported(what: string): never {
  $.panic($.val($.String, $.str.toGo("gots: reflect " + what + " is not supported yet")));
}

export function overflowInt(t: $.Type, x: number): boolean {
  const bits: Record<string, number> = { int8: 8, int16: 16, int32: 32 };
  const b = bits[t.kind];
  if (b === undefined) return false;
  return x < -(2 ** (b - 1)) || x > 2 ** (b - 1) - 1;
}

export function overflowUint(t: $.Type, x: bigint): boolean {
  const bits: Record<string, number> = { uint8: 8, uint16: 16, uint32: 32 };
  const b = bits[t.kind];
  if (b === undefined) return false;
  return x >= 1n << BigInt(b);
}

export function overflowFloat(t: $.Type, x: number): boolean {
  if (t.kind !== "float32") return false;
  const a = Math.abs(x);
  return a > 3.4028234663852886e38 && a !== Infinity;
}

// StructField is a struct field as Type.Field describes it.
export class StructField extends $.Native {
  Name = "";
  PkgPath = "";
  Type: rtype | null = null;
  Tag = "";
  Offset = 0;
  Index: $.Slice<number> | null = null;
  Anonymous = false;
  IsExported(): boolean {
    return this.PkgPath === "";
  }
}
$.nativeType("reflect.StructField", StructField);

export type StructTag = string;

// A StructTag in an interface has its methods.
export class StructTag$Val extends $.Val {
  Get(key: string): string {
    return StructTag$Get(this.v, key);
  }
  Lookup(key: string): [string, boolean] {
    return StructTag$Lookup(this.v, key);
  }
}
export const StructTag$type = $.named("reflect.StructTag", $.String, StructTag$Val);

export function StructTag$Get(tag: string, key: string): string {
  return StructTag$Lookup(tag, key)[0];
}

// StructTag$Lookup is StructTag.Lookup: the value under key in the
// conventional key:"value" form, unquoted.
export function StructTag$Lookup(tag: string, key: string): [string, boolean] {
  while (tag !== "") {
    let i = 0;
    while (i < tag.length && tag[i] === " ") i++;
    tag = tag.slice(i);
    if (tag === "") break;
    i = 0;
    while (i < tag.length && tag.charCodeAt(i) > 0x20 && tag[i] !== ":" && tag[i] !== '"' && tag.charCodeAt(i) !== 0x7f) i++;
    if (i === 0 || i + 1 >= tag.length || tag[i] !== ":" || tag[i + 1] !== '"') break;
    const name = tag.slice(0, i);
    tag = tag.slice(i + 1);
    i = 1;
    while (i < tag.length && tag[i] !== '"') {
      if (tag[i] === "\\") i++;
      i++;
    }
    if (i >= tag.length) break;
    const quoted = tag.slice(0, i + 1);
    tag = tag.slice(i + 1);
    if (name === key) {
      try {
        return [$.str.toGo(JSON.parse($.str.fromGo(quoted))), true];
      } catch {
        break;
      }
    }
  }
  return ["", false];
}

// Method is a method as Type.Method describes it. Func, the method as a
// function of its receiver, is made when asked for.
export class Method extends $.Native {
  Name = "";
  PkgPath = "";
  Type: rtype | null = null;
  Index = 0;
  recv: $.Type | null = null;
  IsExported(): boolean {
    return this.PkgPath === "";
  }
}
$.nativeType("reflect.Method", Method);

// funcOf is the func type with those parameters and results, made at run
// time and written as Go writes it.
export function funcOf(ins: $.Type[], outs: $.Type[], variadic: boolean): $.Type {
  const params = ins.map((t, i) => variadic && i === ins.length - 1 ? "..." + t.elem!.str : t.str);
  const res = outs.length === 0 ? "" : outs.length === 1 ? " " + outs[0].str : " (" + outs.map((t) => t.str).join(", ") + ")";
  return $.funcType("func(" + params.join(", ") + ")" + res, () => ins, () => outs, variadic);
}
