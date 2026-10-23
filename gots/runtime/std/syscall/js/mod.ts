// Go's syscall/js. There is no wasm memory to cross: a Value holds the
// JavaScript value itself, and Go strings are converted at the boundary.
import * as $ from "../../../rt/mod.ts";

export const TypeUndefined = 0, TypeNull = 1, TypeBoolean = 2, TypeNumber = 3, TypeString = 4, TypeSymbol = 5, TypeObject = 6, TypeFunction = 7;
export type Type = number;
const typeNames = ["undefined", "null", "boolean", "number", "string", "symbol", "object", "function"];

export function Type$String(t: Type): string {
  return typeNames[t] ?? "invalid";
}

export function Type$isObject(t: Type): boolean {
  return t === TypeObject || t === TypeFunction;
}

class Type$Val extends $.Val {
  String(): string {
    return Type$String(this.v);
  }
}
export const Type$type = $.named("js.Type", $.Int, Type$Val);

// toJS converts a Go value, as it arrives in an interface, to JavaScript.
export function toJS(x: any): any {
  if (x === null || x === undefined) return null;
  if (x instanceof Value) return x.v;
  if (!(x instanceof $.Val)) return x;
  const v = x.v;
  if (v instanceof Value) return v.v;
  switch (x.$t.kind) {
    case "string":
      return $.str.fromGo(v);
    case "uint64":
      return Number(v);
    case "slice": {
      if (v === null) return null;
      if (!Array.isArray(v.a)) return v.a.subarray(v.o, v.o + v.n);
      return v.toArray().map((e: any) => toJS($.val(x.$t.elem!, e)));
    }
    case "map": {
      const o: any = {};
      if (v !== null) for (const [k, e] of v) o[$.str.fromGo(k)] = toJS(x.$t.elem!.kind === "iface" ? e : $.val(x.$t.elem!, e));
      return o;
    }
  }
  return v;
}

function argsJS(args: $.Slice<any> | null): any[] {
  return args === null ? [] : args.toArray().map(toJS);
}

export class Value extends $.Native {
  v: any = undefined;

  override $eq(o: Value): boolean {
    return this.v === o.v;
  }
  override $key(): string {
    return "js" + $.keyOf(this.v);
  }

  Equal(w: Value): boolean {
    return this.v === w.v;
  }
  IsUndefined(): boolean {
    return this.v === undefined;
  }
  IsNull(): boolean {
    return this.v === null;
  }
  IsNaN(): boolean {
    return Number.isNaN(this.v);
  }
  Type(): Type {
    const v = this.v;
    if (v === undefined) return TypeUndefined;
    if (v === null) return TypeNull;
    switch (typeof v) {
      case "boolean": return TypeBoolean;
      case "number": case "bigint": return TypeNumber;
      case "string": return TypeString;
      case "symbol": return TypeSymbol;
      case "function": return TypeFunction;
    }
    return TypeObject;
  }
  Get(p: string): Value {
    return wrap(this.v[$.str.fromGo(p)]);
  }
  Set(p: string, x: any): void {
    this.v[$.str.fromGo(p)] = toJS(x);
  }
  Delete(p: string): void {
    delete this.v[$.str.fromGo(p)];
  }
  Index(i: number): Value {
    return wrap(this.v[i]);
  }
  SetIndex(i: number, x: any): void {
    this.v[i] = toJS(x);
  }
  Length(): number {
    return this.v.length;
  }
  Call(m: string, args: $.Slice<any> | null): Value {
    const name = $.str.fromGo(m);
    const f = this.v[name];
    if (typeof f !== "function") {
      $.panic(new ValueError("Value.Call", this.Get(m).Type()));
    }
    return wrap(f.apply(this.v, argsJS(args)));
  }
  Invoke(args: $.Slice<any> | null): Value {
    return wrap(this.v(...argsJS(args)));
  }
  New(args: $.Slice<any> | null): Value {
    return wrap(new this.v(...argsJS(args)));
  }
  Float(): number {
    return Number(this.v);
  }
  Int(): number {
    return Math.trunc(Number(this.v));
  }
  Bool(): boolean {
    return this.v === true;
  }
  Truthy(): boolean {
    return !!this.v;
  }
  String(): string {
    const v = this.v;
    switch (this.Type()) {
      case TypeString: return $.str.toGo(v);
      case TypeUndefined: return "<undefined>";
      case TypeNull: return "<null>";
      case TypeBoolean: return `<boolean: ${v}>`;
      case TypeNumber: return `<number: ${v}>`;
      case TypeSymbol: return "<symbol>";
      case TypeObject: return "<object>";
    }
    return "<function>";
  }
  InstanceOf(t: Value): boolean {
    return this.v instanceof t.v;
  }
}
$.nativeType("js.Value", Value);

export function wrap(x: any): Value {
  const v = new Value();
  v.v = x;
  return v;
}

export function Undefined(): Value {
  return wrap(undefined);
}

export function Null(): Value {
  return wrap(null);
}

export function Global(): Value {
  return wrap(globalThis);
}

export function ValueOf(x: any): Value {
  return wrap(toJS(x));
}

export class Func extends Value {
  id = 0;
  Release(): void {}
}
$.nativeType("js.Func", Func);

// FuncOf wraps a Go function as a JavaScript one. It is called directly:
// there is no event loop of Go's to hand the call to.
export function FuncOf(fn: (this_: Value, args: $.Slice<Value> | null) => any): Func {
  const f = new Func();
  f.v = function (this: any, ...args: any[]) {
    return toJS(fn(wrap(this), $.slice(args.map(wrap))));
  };
  return f;
}

export class Error$ extends $.Native {
  Value: Value = new Value();
  Error(): string {
    return "JavaScript error: " + $.str.toGo(String(this.Value.v?.message ?? this.Value.v));
  }
}
$.nativeType("js.Error", Error$);

export class ValueError extends $.Native {
  Method: string;
  Type: Type;
  constructor(method = "", type: Type = 0) {
    super();
    this.Method = method;
    this.Type = type;
  }
  Error(): string {
    return "syscall/js: call of " + this.Method + " on " + Type$String(this.Type);
  }
}
$.nativeType("js.ValueError", ValueError);

export function CopyBytesToGo(dst: $.Slice<number> | null, src: Value): number {
  if (dst === null) return 0;
  const s = src.v as Uint8Array;
  const n = Math.min(dst.n, s.length);
  dst.a.set(s.subarray(0, n), dst.o);
  return n;
}

export function CopyBytesToJS(dst: Value, src: $.Slice<number> | null): number {
  if (src === null) return 0;
  const d = dst.v as Uint8Array;
  const n = Math.min(d.length, src.n);
  d.set(src.a.subarray(src.o, src.o + n));
  return n;
}
