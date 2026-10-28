// Kinds as Go numbers them, from the run-time descriptors' kinds.
import * as $ from "../../rt/mod.ts";

export type Kind = number;

const numbers: Record<string, number> = {
  bool: 1, int: 2, int8: 3, int16: 4, int32: 5, int64: 6,
  uint: 7, uint8: 8, uint16: 9, uint32: 10, uint64: 11, uintptr: 12,
  float32: 13, float64: 14, complex64: 15, complex128: 16,
  array: 17, chan: 18, func: 19, iface: 20, map: 21, ptr: 22,
  slice: 23, string: 24, struct: 25, unsafeptr: 26,
};

const names = [
  "invalid", "bool", "int", "int8", "int16", "int32", "int64",
  "uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
  "float32", "float64", "complex64", "complex128",
  "array", "chan", "func", "interface", "map", "ptr", "slice", "string", "struct", "unsafe.Pointer",
];

export function kindOf(t: $.Type | null): Kind {
  return t === null ? 0 : numbers[t.kind] ?? 0;
}

// Kind$String is Kind.String.
export function Kind$String(k: Kind): string {
  return names[k] ?? "kind" + k;
}

// A Kind in an interface has its String method, for fmt.
export class Kind$Val extends $.Val {
  String(): string {
    return Kind$String(this.v);
  }
}
export const Kind$type = $.named("reflect.Kind", $.Uint, Kind$Val);

export function isInt(t: $.Type): boolean {
  const k = kindOf(t);
  return k >= 2 && k <= 6;
}

export function isUint(t: $.Type): boolean {
  const k = kindOf(t);
  return k >= 7 && k <= 12;
}

export function isFloat(t: $.Type): boolean {
  return t.kind === "float32" || t.kind === "float64";
}

export function isComplex(t: $.Type): boolean {
  return t.kind === "complex64" || t.kind === "complex128";
}

// isObject reports whether values of t are objects pointers point at:
// structs and arrays, which a pointer is the very object of.
export function isObject(t: $.Type): boolean {
  return t.kind === "struct" || t.kind === "array";
}
