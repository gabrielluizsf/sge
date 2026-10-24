// reflect.Value: a descriptor, a value, and where the value lives when it
// is addressable. A struct or array is its own place, as a pointer to one
// is the object; any other addressable value is reached through a Ptr (a
// Cell, a FieldPtr or an IndexPtr), so setting it changes the variable,
// field or element it came from, as Go's Value does.
import * as $ from "../../rt/mod.ts";
import { isComplex, isFloat, isInt, isObject, isUint, Kind, kindOf } from "./kinds.ts";
import { overflowFloat, overflowInt, overflowUint, rtype, toType, unsupported } from "./rtype.ts";
import { backing, box, copyOf, makeMap, makeSlice, setIn, zero } from "./zero.ts";

const addressable = 1, readOnly = 2;

function fail(msg: string): never {
  $.panic($.val($.String, $.str.toGo(msg)));
}

export class Value extends $.Native {
  t: $.Type | null = null;
  x: any = undefined; // the value, or the object of a struct or array
  p: $.Ptr<any> | null = null; // where an addressable non-object value lives
  f = 0;

  static of(t: $.Type, x: any, f = 0): Value {
    const v = new Value();
    v.t = t;
    v.x = x;
    v.f = f;
    return v;
  }

  // at is the value at p, addressable: an object when t is a struct or
  // array, the place itself otherwise.
  static at(t: $.Type, p: $.Ptr<any>, f: number): Value {
    if (isObject(t)) return Value.of(t, p.v, f | addressable);
    const v = Value.of(t, undefined, f | addressable);
    v.p = p;
    return v;
  }

  get(): any {
    return this.p !== null ? this.p.v : this.x;
  }

  private must(what: string): $.Type {
    if (this.t === null) fail("reflect: call of reflect.Value." + what + " on zero Value");
    return this.t;
  }

  private kind(what: string, ...kinds: string[]): $.Type {
    const t = this.must(what);
    if (!kinds.includes(t.kind)) fail("reflect: call of reflect.Value." + what + " on " + kindName(t) + " Value");
    return t;
  }

  private settable(what: string): $.Type {
    const t = this.must(what);
    if (this.f & readOnly) fail("reflect: reflect.Value." + what + " using value obtained using unexported field");
    if (!(this.f & addressable)) fail("reflect: reflect.Value." + what + " using unaddressable value");
    return t;
  }

  private store(x: any): void {
    if (isObject(this.t!)) setIn(this.t!, this.x, x);
    else this.p!.v = x;
  }

  IsValid(): boolean {
    return this.t !== null;
  }
  Kind(): Kind {
    return kindOf(this.t);
  }
  Type(): rtype {
    return toType(this.must("Type"))!;
  }
  CanAddr(): boolean {
    return (this.f & addressable) !== 0;
  }
  CanSet(): boolean {
    return (this.f & addressable) !== 0 && (this.f & readOnly) === 0;
  }
  CanInterface(): boolean {
    this.must("CanInterface");
    return (this.f & readOnly) === 0;
  }

  Interface(): any {
    const t = this.must("Interface");
    if (this.f & readOnly) fail("reflect.Value.Interface: cannot return value obtained from unexported field or method");
    return box(t, copyOf(t, this.get()));
  }

  IsNil(): boolean {
    const t = this.kind("IsNil", "ptr", "map", "slice", "func", "iface", "chan", "unsafeptr");
    return t.kind === "ptr" && isObject(t.elem!) ? this.get() === null : this.get() === null;
  }

  IsZero(): boolean {
    return isZero(this.must("IsZero"), this.get());
  }

  Elem(): Value {
    const t = this.kind("Elem", "iface", "ptr");
    const x = this.get();
    if (x === null) return new Value();
    if (t.kind === "iface") {
      const v = ValueOf(x);
      v.f |= this.f & readOnly;
      return v;
    }
    const e = t.elem!;
    if (isObject(e)) return Value.of(e, x, addressable | (this.f & readOnly));
    return Value.at(e, x, this.f & readOnly);
  }

  Addr(): Value {
    const t = this.must("Addr");
    if (!(this.f & addressable)) fail("reflect.Value.Addr of unaddressable value");
    return Value.of($.ptrTo(t), isObject(t) ? this.x : this.p, this.f & readOnly);
  }

  NumField(): number {
    return this.kind("NumField", "struct").fields.length;
  }

  Field(i: number): Value {
    const t = this.kind("Field", "struct");
    const f = t.fields[i];
    if (f === undefined) fail("reflect: Field index out of range");
    const obj = this.get();
    const fl = (this.f & readOnly) | (f.exported ? 0 : readOnly);
    if (isObject(f.type)) return Value.of(f.type, obj[f.prop], fl | (this.f & addressable));
    if (this.f & addressable) return Value.at(f.type, new $.FieldPtr(obj, f.prop), fl);
    return Value.of(f.type, obj[f.prop], fl);
  }

  FieldByName(name: string): Value {
    const t = this.kind("FieldByName", "struct");
    const i = t.fields.findIndex((f) => $.str.toGo(f.name) === name);
    return i < 0 ? new Value() : this.Field(i);
  }

  FieldByIndex(index: $.Slice<number> | null): Value {
    let v: Value = this;
    for (const i of index ?? []) {
      if (v.Kind() === 22 && v.t!.elem!.kind === "struct") {
        if (v.IsNil()) fail("reflect: indirection through nil pointer to embedded struct");
        v = v.Elem();
      }
      v = v.Field(i);
    }
    return v;
  }

  Len(): number {
    const t = this.kind("Len", "slice", "array", "string", "map", "chan", "ptr");
    const x = this.get();
    switch (t.kind) {
      case "slice":
        return $.len(x);
      case "array":
        return t.len;
      case "ptr":
        if (t.elem!.kind !== "array") fail("reflect: call of reflect.Value.Len on ptr Value");
        return t.elem!.len;
      case "string":
        return x.length;
      case "map":
        return $.mapLen(x);
    }
    return $.ch.chanLen(x);
  }

  Cap(): number {
    const t = this.kind("Cap", "slice", "array", "chan");
    if (t.kind === "slice") return $.cap(this.get());
    if (t.kind === "array") return t.len;
    return $.ch.chanCap(this.get());
  }

  Index(i: number): Value {
    const t = this.kind("Index", "slice", "array", "string");
    const x = this.get();
    if (t.kind === "string") {
      if (i < 0 || i >= x.length) fail("reflect: string index out of range");
      return Value.of($.Uint8, $.str.byteAt(x, i), this.f & readOnly);
    }
    const e = t.elem!;
    if (t.kind === "slice") {
      if (x === null || i < 0 || i >= x.n) fail("reflect: slice index out of range");
      return Value.at(e, new $.IndexPtr(x.a, x.o + i), this.f & readOnly);
    }
    if (i < 0 || i >= t.len) fail("reflect: array index out of range");
    if (this.f & addressable) return Value.at(e, new $.IndexPtr(x, i), this.f & readOnly);
    return Value.of(e, x[i], this.f & readOnly);
  }

  Bool(): boolean {
    return this.kind("Bool", "bool") && this.get();
  }
  Int(): number {
    const t = this.must("Int");
    if (!isInt(t)) fail("reflect: call of reflect.Value.Int on " + kindName(t) + " Value");
    return this.get();
  }
  Uint(): bigint {
    const t = this.must("Uint");
    if (!isUint(t)) fail("reflect: call of reflect.Value.Uint on " + kindName(t) + " Value");
    const x = this.get();
    return typeof x === "bigint" ? x : BigInt(x);
  }
  Float(): number {
    return this.kind("Float", "float32", "float64") && this.get();
  }
  Complex(): any {
    return this.kind("Complex", "complex64", "complex128") && this.get();
  }
  String(): string {
    if (this.t === null) return $.str.toGo("<invalid Value>");
    if (this.t.kind === "string") return this.get();
    return $.str.toGo("<" + this.t.str + " Value>");
  }
  Bytes(): $.Slice<number> | null {
    const t = this.must("Bytes");
    if (t.kind === "slice" && t.elem!.kind === "uint8") return this.get();
    if (t.kind === "array" && t.elem!.kind === "uint8" && this.f & addressable) return $.slice(this.x);
    fail("reflect.Value.Bytes of non-byte slice");
  }
  UnsafePointer(): any {
    const t = this.kind("UnsafePointer", "ptr", "map", "slice", "func", "chan", "unsafeptr");
    const x = this.get();
    return t.kind === "slice" ? (x === null ? null : x.a) : x;
  }
  Pointer(): any {
    return this.UnsafePointer();
  }

  Set(x: Value): void {
    const t = this.settable("Set");
    x.must("Set");
    if (x.f & readOnly) fail("reflect: reflect.Value.Set using value obtained using unexported field");
    let v = x.get();
    if (t.kind === "iface" && x.t!.kind !== "iface") v = box(x.t!, copyOf(x.t!, v));
    else v = copyOf(t, v);
    this.store(v);
  }
  SetBool(b: boolean): void {
    this.settable("SetBool");
    this.store(b);
  }
  SetInt(n: number): void {
    const t = this.settable("SetInt");
    if (!isInt(t)) fail("reflect: call of reflect.Value.SetInt on " + kindName(t) + " Value");
    this.store(wrapInt(t, n));
  }
  SetUint(n: bigint): void {
    const t = this.settable("SetUint");
    if (!isUint(t)) fail("reflect: call of reflect.Value.SetUint on " + kindName(t) + " Value");
    this.store(wrapUint(t, n));
  }
  SetFloat(f: number): void {
    const t = this.settable("SetFloat");
    if (!isFloat(t)) fail("reflect: call of reflect.Value.SetFloat on " + kindName(t) + " Value");
    this.store(t.kind === "float32" ? Math.fround(f) : f);
  }
  SetComplex(c: any): void {
    const t = this.settable("SetComplex");
    if (!isComplex(t)) fail("reflect: call of reflect.Value.SetComplex on " + kindName(t) + " Value");
    this.store(t.kind === "complex64" ? { re: Math.fround(c.re), im: Math.fround(c.im) } : c);
  }
  SetString(s: string): void {
    this.settable("SetString");
    this.store(s);
  }
  SetBytes(b: $.Slice<number> | null): void {
    this.settable("SetBytes");
    this.store(b);
  }
  SetZero(): void {
    const t = this.settable("SetZero");
    this.store(zero(t));
  }
  SetLen(n: number): void {
    this.settable("SetLen");
    const s = this.get();
    if (n < 0 || n > $.cap(s)) fail("reflect: slice length out of range in SetLen");
    this.store(s === null ? null : new $.Slice(s.a, s.o, n, s.c));
  }
  SetCap(c: number): void {
    this.settable("SetCap");
    const s = this.get();
    if (c < $.len(s) || c > $.cap(s)) fail("reflect: slice capacity out of range in SetCap");
    this.store(s === null ? null : new $.Slice(s.a, s.o, s.n, c));
  }

  // Grow makes room for n more elements without the slice moving again.
  Grow(n: number): void {
    const t = this.settable("Grow");
    if (t.kind !== "slice") fail("reflect: call of reflect.Value.Grow on " + kindName(t) + " Value");
    if (n < 0) fail("reflect.Value.Grow: negative len");
    const s = this.get();
    const len = $.len(s), c = $.cap(s);
    if (len + n <= c) return;
    const nc = Math.max(len + n, c < 256 ? c * 2 : c + (c >> 2) + 192, 4);
    const grown = makeSlice(t, len, nc);
    for (let i = 0; i < len; i++) grown.a[i] = s.a[s.o + i];
    this.store(grown);
  }

  MapIndex(k: Value): Value {
    const t = this.kind("MapIndex", "map");
    const m = this.get();
    if (m === null) return new Value();
    const x = m.get(mapKey(t, k));
    return x === undefined ? new Value() : Value.of(t.elem!, x, this.f & readOnly);
  }

  MapKeys(): $.Slice<Value> | null {
    const t = this.kind("MapKeys", "map");
    const keys: Value[] = [];
    for (const [k] of $.entries(this.get())) keys.push(keyValue(t, k, this.f));
    return $.slice(keys);
  }

  MapRange(): MapIter {
    this.kind("MapRange", "map");
    const it = new MapIter();
    it.m = this;
    return it;
  }

  SetMapIndex(k: Value, e: Value): void {
    const t = this.kind("SetMapIndex", "map");
    if (this.f & readOnly) fail("reflect: reflect.Value.SetMapIndex using value obtained using unexported field");
    const m = this.get();
    const key = mapKey(t, k);
    if (!e.IsValid()) {
      $.mapDelete(m, key);
      return;
    }
    if (m === null) fail("assignment to entry in nil map");
    exported(k, "reflect.Value.SetMapIndex");
    exported(e, "reflect.Value.SetMapIndex");
    let x = e.get();
    x = t.elem!.kind === "iface" && e.t!.kind !== "iface" ? box(e.t!, copyOf(e.t!, x)) : copyOf(t.elem!, x);
    $.mapSet(m, key, x);
  }

  // SetIterKey and SetIterValue set v to the key or value the iterator is at.
  SetIterKey(it: MapIter): void {
    this.Set(it.Key());
  }
  SetIterValue(it: MapIter): void {
    this.Set(it.Value());
  }

  // Clear empties a map, or zeroes a slice's elements.
  Clear(): void {
    const t = this.kind("Clear", "map", "slice");
    const x = this.get();
    if (x === null) return;
    if (t.kind === "map") x.clear();
    else for (let i = 0; i < x.n; i++) x.a[x.o + i] = zero(t.elem!);
  }

  NumMethod(): number {
    return this.Type().NumMethod();
  }

  OverflowInt(x: number): boolean {
    return overflowInt(this.must("OverflowInt"), x);
  }
  OverflowUint(x: bigint): boolean {
    return overflowUint(this.must("OverflowUint"), x);
  }
  OverflowFloat(x: number): boolean {
    return overflowFloat(this.must("OverflowFloat"), x);
  }

  // Convert is v converted to t, for the conversions between numbers,
  // and between types of one underlying type.
  Convert(u: rtype): Value {
    const v = this.convert(u);
    v.f |= this.f & readOnly;
    return v;
  }

  private convert(u: rtype): Value {
    const from = this.must("Convert"), to = u.d;
    const x = this.get();
    if (isInt(from) || isUint(from) || isFloat(from)) {
      const n = typeof x === "bigint" ? Number(x) : x;
      if (isInt(to)) return Value.of(to, wrapInt(to, Math.trunc(n)));
      if (isUint(to)) return Value.of(to, wrapUint(to, typeof x === "bigint" ? x : BigInt(Math.trunc(n))));
      if (isFloat(to)) return Value.of(to, to.kind === "float32" ? Math.fround(n) : n);
      if (to.kind === "string" && !isFloat(from)) return Value.of(to, $.str.runeStr(n));
    }
    if (from.kind === "string" && to.kind === "slice" && to.elem!.kind === "uint8") return Value.of(to, $.slice($.str.strToBytes(x)));
    if (from.kind === "slice" && from.elem!.kind === "uint8" && to.kind === "string") return Value.of(to, $.str.bytesToStr(x));
    if (to.kind === "iface") return ValueOf(box(from, copyOf(from, x))).retag(to);
    if (from.kind === to.kind) return Value.of(to, copyOf(from, x));
    fail("reflect.Value.Convert: value of type " + from.str + " cannot be converted to type " + to.str);
  }

  private retag(t: $.Type): Value {
    const v = Value.of(t, box(this.t!, this.get()));
    return v;
  }

  Comparable(): boolean {
    return this.t !== null && this.Type().Comparable();
  }

  Equal(u: Value): boolean {
    if (this.t === null || u.t === null) return this.t === u.t;
    return this.t === u.t && $.equal(this.t, this.get(), u.get());
  }

  Method(_i: number): Value {
    return unsupported("Value.Method");
  }
  MethodByName(_name: string): Value {
    return unsupported("Value.MethodByName");
  }
  Call(_in: any): any {
    return unsupported("Value.Call");
  }
}
$.nativeType("reflect.Value", Value);

// MapIter is what MapRange gives: the entries as they were when it began.
export class MapIter extends $.Native {
  m: Value | null = null;
  private entries: [any, any][] | null = null;
  private i = -1;

  Next(): boolean {
    if (this.entries === null) this.entries = [...$.entries(this.m!.get())];
    this.i++;
    return this.i < this.entries.length;
  }
  Key(): Value {
    return keyValue(this.m!.t!, this.entries![this.i][0], this.m!.f);
  }
  Value(): Value {
    return Value.of(this.m!.t!.elem!, this.entries![this.i][1], this.m!.f & readOnly);
  }
  Reset(v: Value): void {
    this.m = v;
    this.entries = null;
    this.i = -1;
  }
}
$.nativeType("reflect.MapIter", MapIter);

// exported panics when v was reached through an unexported field, which
// Go never lets a value from be used to write anything.
function exported(v: Value, what: string): void {
  if (v.f & readOnly) fail("reflect: " + what + " using value obtained using unexported field");
}

function kindName(t: $.Type): string {
  return t.kind === "iface" ? "interface" : t.kind;
}

// mapKey is k as the map of type t keeps its keys: boxed when they are
// interfaces.
function mapKey(t: $.Type, k: Value): any {
  const x = k.get();
  return t.key!.kind === "iface" && k.t!.kind !== "iface" ? box(k.t!, x) : x;
}

function keyValue(t: $.Type, k: any, f: number): Value {
  return Value.of(t.key!, k, f & readOnly);
}

function wrapInt(t: $.Type, n: number): number {
  switch (t.kind) {
    case "int8":
      return (n << 24) >> 24;
    case "int16":
      return (n << 16) >> 16;
    case "int32":
      return n | 0;
  }
  return n;
}

function wrapUint(t: $.Type, n: bigint): any {
  switch (t.kind) {
    case "uint8":
      return Number(n & 0xffn);
    case "uint16":
      return Number(n & 0xffffn);
    case "uint32":
      return Number(n & 0xffffffffn);
    case "uint64":
      return BigInt.asUintN(64, n);
  }
  return Number(BigInt.asUintN(64, n));
}

function isZero(t: $.Type, x: any): boolean {
  switch (t.kind) {
    case "struct":
      return t.fields.every((f) => isZero(f.type, x[f.prop]));
    case "array":
      for (let i = 0; i < t.len; i++) if (!isZero(t.elem!, x[i])) return false;
      return true;
    case "complex64": case "complex128":
      return x.re === 0 && x.im === 0;
    case "float32": case "float64":
      return x === 0 && !Object.is(x, -0);
  }
  return x === zero(t);
}

// ValueOf is the Value an interface value holds.
export function ValueOf(i: any): Value {
  if (i === null || i === undefined) return new Value();
  if (i instanceof $.Val) return Value.of(i.$t, i.v);
  const d = i.constructor?.$type;
  if (d === undefined) fail("gots: reflect.ValueOf of a value from outside Go");
  return Value.of($.ptrTo(d), i);
}

export function Indirect(v: Value): Value {
  return v.Kind() === 22 ? v.Elem() : v;
}

// New is a pointer to a new zero value of type t.
export function New(t: rtype | null): Value {
  if (t === null) fail("reflect: New(nil)");
  const d = t.d;
  const x = isObject(d) ? zero(d) : new $.Cell(zero(d));
  return Value.of($.ptrTo(d), x);
}

export function Zero(t: rtype | null): Value {
  if (t === null) fail("reflect: Zero(nil)");
  return Value.of(t.d, zero(t.d));
}

export function MakeSlice(t: rtype, n: number, c: number): Value {
  if (t.d.kind !== "slice") fail("reflect.MakeSlice of non-slice type");
  if (n < 0 || c < n) fail("reflect.MakeSlice: len out of range");
  return Value.of(t.d, makeSlice(t.d, n, c));
}

export function MakeMap(t: rtype): Value {
  if (t.d.kind !== "map") fail("reflect.MakeMap of non-map type");
  return Value.of(t.d, makeMap(t.d));
}

export function MakeMapWithSize(t: rtype, _n: number): Value {
  return MakeMap(t);
}

export function MakeChan(t: rtype, n: number): Value {
  return Value.of(t.d, new $.ch.Chan(n, () => zero(t.d.elem!)));
}

// Append is append(s, xs...) on Values.
export function Append(s: Value, xs: $.Slice<Value> | null): Value {
  exported(s, "reflect.Append");
  for (const x of xs ?? []) exported(x, "reflect.Append");
  const t = s.Type().d;
  let out = s.get();
  for (const x of xs ?? []) {
    const v = t.elem!.kind === "iface" && x.t!.kind !== "iface" ? box(x.t!, x.get()) : copyOf(t.elem!, x.get());
    out = $.append(out, backing(t.elem!), v);
  }
  return Value.of(t, out);
}

export function AppendSlice(s: Value, t: Value): Value {
  exported(s, "reflect.AppendSlice");
  exported(t, "reflect.AppendSlice");
  return Value.of(s.t!, $.appendSlice(s.get(), backing(s.t!.elem!), t.get()));
}

export function Copy(dst: Value, src: Value): number {
  exported(dst, "reflect.Copy");
  exported(src, "reflect.Copy");
  const s = src.t!.kind === "string" ? $.slice($.str.strToBytes(src.get())) : src.get();
  return $.copy(dst.t!.kind === "array" ? $.slice(dst.get()) : dst.get(), src.t!.kind === "array" ? $.slice(s) : s);
}

export function TypeOf(i: any): rtype | null {
  return i === null || i === undefined ? null : toType($.typeOf(i));
}

export function PointerTo(t: rtype): rtype {
  return toType($.ptrTo(t.d))!;
}
export const PtrTo = PointerTo;

export function SliceOf(t: rtype): rtype {
  return toType($.sliceOf(t.d))!;
}

export function MapOf(k: rtype, e: rtype): rtype {
  return toType($.mapOf(k.d, e.d))!;
}

export function ArrayOf(n: number, e: rtype): rtype {
  return toType($.arrayOf(e.d, n))!;
}

// TypeFor is reflect.TypeFor[T](): gots hands it T's descriptor.
export function TypeFor(d: $.Type): rtype {
  return toType(d)!;
}

// TypeAssert is reflect.TypeAssert[T](v): v's value as a T.
export function TypeAssert(d: $.Type, v: Value): [any, boolean] {
  const x = v.Interface();
  if (d.kind === "iface") return x !== null && $.implementsI(x, d.methods) ? [x, true] : [null, false];
  if (x === null) return [zero(d), false];
  if ($.typeOf(x) !== d) return [zero(d), false];
  return [x instanceof $.Val ? x.v : x, true];
}

// Swapper swaps elements i and j of the slice s.
export function Swapper(s: any): (i: number, j: number) => void {
  const sl = s instanceof $.Val ? s.v : s;
  return (i: number, j: number) => {
    const a = sl.a, o = sl.o;
    const tmp = a[o + i];
    a[o + i] = a[o + j];
    a[o + j] = tmp;
  };
}

// ValueError is the panic a Value method raises on the wrong kind.
export class ValueError extends $.Native {
  Method = "";
  Kind: Kind = 0;
  Error(): string {
    return $.str.toGo("reflect: call of " + $.str.fromGo(this.Method) + " on zero Value");
  }
}
$.nativeType("reflect.ValueError", ValueError);
