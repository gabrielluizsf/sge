// The parts of package unsafe that work on JavaScript values: a pointer to
// an element is an IndexPtr, its array and index, so pointer arithmetic over
// a slice's backing array still means something.
import { IndexPtr } from "./ptr.ts";
import { Slice } from "./slice.ts";
import { byte, byteStr } from "./string.ts";

export function sliceData(s: Slice<any> | null): IndexPtr<any> | null {
  return s === null ? null : new IndexPtr(s.a, s.o);
}

export function slice(p: any, n: number): Slice<any> | null {
  if (p === null) return null;
  if (p instanceof IndexPtr) return new Slice(p.a, p.i, n, p.a.length - p.i);
  throw new Error("gots: unsafe.Slice of a pointer that is not to an element");
}

export function string(p: any, n: number): string {
  if (n === 0) return "";
  if (!(p instanceof IndexPtr)) throw new Error("gots: unsafe.String of a pointer that is not to bytes");
  let s = "";
  for (let i = 0; i < n; i++) s += byteStr(p.a[p.i + i]);
  return s;
}

export function stringData(s: string): IndexPtr<number> | null {
  if (s.length === 0) return null;
  const b = new Uint8Array(s.length);
  for (let i = 0; i < s.length; i++) b[i] = byte(s.charCodeAt(i));
  return new IndexPtr(b, 0);
}

export function add(p: any, n: number): any {
  if (p instanceof IndexPtr) return new IndexPtr(p.a, p.i + n);
  throw new Error("gots: unsafe.Add of a pointer that is not to an element");
}

const view = new DataView(new ArrayBuffer(8));

const writers: Record<string, (v: any) => void> = {
  float64: (v) => view.setFloat64(0, v, true),
  float32: (v) => view.setFloat32(0, v, true),
  uint64: (v) => view.setBigUint64(0, v, true),
  int64: (v) => view.setBigInt64(0, BigInt(v), true),
  int: (v) => view.setBigInt64(0, BigInt(v), true),
  uint: (v) => view.setUint32(0, v, true),
  uintptr: (v) => view.setBigUint64(0, BigInt(v), true),
  uint32: (v) => view.setUint32(0, v, true),
  int32: (v) => view.setInt32(0, v, true),
  uint16: (v) => view.setUint16(0, v, true),
  int16: (v) => view.setInt16(0, v, true),
  uint8: (v) => view.setUint8(0, v),
  int8: (v) => view.setInt8(0, v),
};

const readers: Record<string, () => any> = {
  float64: () => view.getFloat64(0, true),
  float32: () => view.getFloat32(0, true),
  uint64: () => view.getBigUint64(0, true),
  int64: () => Number(view.getBigInt64(0, true)),
  int: () => Number(view.getBigInt64(0, true)),
  uint: () => view.getUint32(0, true),
  uintptr: () => Number(view.getBigUint64(0, true)),
  uint32: () => view.getUint32(0, true),
  int32: () => view.getInt32(0, true),
  uint16: () => view.getUint16(0, true),
  int16: () => view.getInt16(0, true),
  uint8: () => view.getUint8(0),
  int8: () => view.getInt8(0),
};

// reinterpret is *(*to)(unsafe.Pointer(&v)) for a v of type from: the same
// bytes, little-endian as on wasm, read as the other type.
export function reinterpret(v: any, from: string, to: string): any {
  writers[from](v);
  return readers[to]();
}
