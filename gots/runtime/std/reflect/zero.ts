// Values made from descriptors alone: zeros, copies and backing arrays,
// laid out the way the transpiler lays them out.
import * as $ from "../../rt/mod.ts";

const typed: Record<string, any> = {
  int: Float64Array, int64: Float64Array, uintptr: Float64Array, float64: Float64Array,
  float32: Float32Array, int32: Int32Array, int16: Int16Array, int8: Int8Array,
  uint: Uint32Array, uint32: Uint32Array, uint16: Uint16Array, uint8: Uint8Array, uint64: BigUint64Array,
};

// backing is the typed array class elements of t are kept in, or null for
// a plain array. Named types keep the kind of what they are named over.
export function backing(t: $.Type): any {
  return typed[t.kind] ?? null;
}

export function zero(t: $.Type): any {
  switch (t.kind) {
    case "bool":
      return false;
    case "string":
      return "";
    case "uint64":
      return 0n;
    case "complex64":
    case "complex128":
      return { re: 0, im: 0 };
    case "struct":
      return new t.ctor();
    case "array": {
      const b = backing(t.elem!);
      return b !== null ? new b(t.len) : $.fillArray(t.len, () => zero(t.elem!));
    }
    case "slice": case "map": case "ptr": case "func": case "iface": case "chan": case "unsafeptr":
      return null;
  }
  return 0;
}

// copyOf is a copy of v, of type t, as an assignment makes.
export function copyOf(t: $.Type, v: any): any {
  if (t.kind === "struct") return v.$clone();
  if (t.kind === "array") return t.elem!.kind === "struct" || t.elem!.kind === "array"
    ? $.cloneArray(v, (e: any) => copyOf(t.elem!, e))
    : $.cloneArray(v);
  return v;
}

// setIn copies v into the struct or array dst, in place.
export function setIn(t: $.Type, dst: any, v: any): void {
  if (t.kind === "struct") dst.$set(v);
  else if (t.elem!.kind === "struct" || t.elem!.kind === "array") $.setArray(dst, v, (d: any, s: any) => setIn(t.elem!, d, s));
  else $.setArray(dst, v);
}

// box is v, of type t, as an interface holds it.
export function box(t: $.Type, v: any): any {
  if (t.kind === "iface") return v;
  if (t.kind === "ptr" && t.elem!.kind === "struct") return v;
  return $.val(t, v);
}

// keyedByValue reports whether a map with keys of type t is a JavaScript
// Map, comparing them as Go does, rather than a KeyMap.
export function keyedByValue(t: $.Type): boolean {
  switch (t.kind) {
    case "complex64": case "complex128": case "struct": case "array": case "iface":
      return false;
    case "ptr":
      return t.elem!.kind === "struct" || t.elem!.kind === "array";
  }
  return true;
}

export function makeMap(t: $.Type): any {
  return keyedByValue(t.key!) ? new Map() : new $.KeyMap($.keyFor(t.key!));
}

export function makeSlice(t: $.Type, n: number, c: number): $.Slice<any> {
  const b = backing(t.elem!);
  return b !== null ? $.make(b, n, c) : $.make(null, n, c, () => zero(t.elem!));
}

const sizes: Record<string, number> = {
  bool: 1, int8: 1, uint8: 1, int16: 2, uint16: 2, int32: 4, uint32: 4, float32: 4,
  int: 8, int64: 8, uint: 4, uint64: 8, uintptr: 8, float64: 8, complex64: 8, complex128: 16,
  string: 16, slice: 24, iface: 16, ptr: 8, map: 8, chan: 8, func: 8, unsafeptr: 8,
};

// size is unsafe.Sizeof for t under js/wasm, which lays out as amd64 does
// but for uint, 32 bits under gots (see wasmSizes in gots).
export function size(t: $.Type): number {
  if (t.kind === "array") return t.len * size(t.elem!);
  if (t.kind === "struct") {
    let off = 0;
    for (const f of t.fields) {
      const a = align(f.type);
      off = Math.ceil(off / a) * a + size(f.type);
    }
    const a = align(t);
    return Math.ceil(off / a) * a;
  }
  return sizes[t.kind];
}

export function align(t: $.Type): number {
  if (t.kind === "array") return align(t.elem!);
  if (t.kind === "struct") return t.fields.reduce((m, f) => Math.max(m, align(f.type)), 1);
  if (t.kind === "complex128" || t.kind === "complex64") return t.kind === "complex64" ? 4 : 8;
  return Math.min(8, sizes[t.kind]);
}
