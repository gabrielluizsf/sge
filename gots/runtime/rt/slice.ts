import { byte } from "./string.ts";

// A Go slice: a window [o, o+n) over a backing array that can grow to c.
// The backing array is a typed array for numeric element types, so stores
// wrap the way Go's do, and a plain array for everything else.
export class Slice<T> {
  a: any;
  o: number;
  n: number;
  c: number;

  constructor(a: any, o: number, n: number, c: number) {
    this.a = a;
    this.o = o;
    this.n = n;
    this.c = c;
  }

  get(i: number): T {
    if (i < 0 || i >= this.n) throw outOfRange(i, this.n);
    return this.a[this.o + i];
  }

  set(i: number, v: T): void {
    if (i < 0 || i >= this.n) throw outOfRange(i, this.n);
    this.a[this.o + i] = v;
  }

  toArray(): T[] {
    const out: T[] = new Array(this.n);
    for (let i = 0; i < this.n; i++) out[i] = this.a[this.o + i];
    return out;
  }

  *[Symbol.iterator](): IterableIterator<T> {
    for (let i = 0; i < this.n; i++) yield this.a[this.o + i];
  }
}

function outOfRange(i: number, n: number): Error {
  return new RangeError(`index out of range [${i}] with length ${n}`);
}

// slice wraps a whole array, for slice literals and conversions.
export function slice<T>(a: any): Slice<T> {
  return new Slice<T>(a, 0, a.length, a.length);
}

// fromArray is the slice of a TypeScript array, for callers that are not
// transpiled Go.
export function fromArray<T>(a: T[]): Slice<T> {
  return slice<T>(a.slice());
}

export function len(s: Slice<any> | null): number {
  return s === null ? 0 : s.n;
}

export function cap(s: Slice<any> | null): number {
  return s === null ? 0 : s.c;
}

// make is make([]T, n, c): zero comes from the backing array for typed
// arrays and from zero() for the rest.
export function make<T>(backing: any, n: number, c: number, zero?: () => T): Slice<T> {
  if (n < 0 || c < n) throw new RangeError("makeslice: len out of range");
  const a = typeof backing === "function" ? new backing(c) : new Array(c);
  if (zero !== undefined) for (let i = 0; i < c; i++) a[i] = zero();
  return new Slice<T>(a, 0, n, c);
}

// grow is a backing array of the same kind as a with room for c elements,
// holding the first n elements of s.
function grow(s: Slice<any>, c: number): any {
  const Kind = s.a.constructor;
  const a = Kind === Array ? new Array(c) : new Kind(c);
  if (Kind === Array) {
    for (let i = 0; i < s.n; i++) a[i] = s.a[s.o + i];
  } else {
    a.set(s.a.subarray(s.o, s.o + s.n));
  }
  return a;
}

function newCap(old: number, need: number): number {
  let c = old < 256 ? old * 2 : old + (old >> 2) + 192;
  if (c < need) c = need;
  return c < 4 ? 4 : c;
}

// append is append(s, vals...). backing is the typed array constructor for
// a nil s, which has no backing array to copy the kind of.
export function append<T>(s: Slice<T> | null, backing: any, ...vals: T[]): Slice<T> {
  return appendArray(s, backing, vals, 0, vals.length);
}

// appendSlice is append(s, t...).
export function appendSlice<T>(s: Slice<T> | null, backing: any, t: Slice<T> | null): Slice<T> | null {
  if (t === null || t.n === 0) return s;
  return appendArray(s, backing, t.a, t.o, t.n);
}

// appendString is append(b, str...) for a []byte.
export function appendString(s: Slice<number> | null, str: string): Slice<number> | null {
  if (str.length === 0) return s;
  const b = new Uint8Array(str.length);
  for (let i = 0; i < str.length; i++) b[i] = byte(str.charCodeAt(i));
  return appendArray(s, Uint8Array, b, 0, b.length);
}

function appendArray<T>(s: Slice<T> | null, backing: any, src: any, o: number, k: number): Slice<T> {
  if (s === null) {
    s = new Slice<T>(typeof backing === "function" ? new backing(0) : [], 0, 0, 0);
  }
  const n = s.n + k;
  if (n <= s.c) {
    if (src === s.a && o < s.o + n && s.o + s.n < o + k) src = src.slice(o, o + k), o = 0;
    for (let i = 0; i < k; i++) s.a[s.o + s.n + i] = src[o + i];
    return new Slice<T>(s.a, s.o, n, s.c);
  }
  const c = newCap(s.c, n);
  const a = grow(s, c);
  for (let i = 0; i < k; i++) a[s.n + i] = src[o + i];
  return new Slice<T>(a, 0, n, c);
}

// copy is copy(dst, src), safe when the two share a backing array.
export function copy(dst: Slice<any> | null, src: Slice<any> | null): number {
  if (dst === null || src === null) return 0;
  const k = Math.min(dst.n, src.n);
  if (k === 0) return 0;
  if (dst.a === src.a) {
    if (Array.isArray(dst.a)) {
      const tmp = src.a.slice(src.o, src.o + k);
      for (let i = 0; i < k; i++) dst.a[dst.o + i] = tmp[i];
    } else {
      dst.a.copyWithin(dst.o, src.o, src.o + k);
    }
    return k;
  }
  if (Array.isArray(dst.a) || Array.isArray(src.a)) {
    for (let i = 0; i < k; i++) dst.a[dst.o + i] = src.a[src.o + i];
  } else {
    dst.a.set(src.a.subarray(src.o, src.o + k), dst.o);
  }
  return k;
}

// copyString is copy(dst, str) for a []byte.
export function copyString(dst: Slice<number> | null, str: string): number {
  if (dst === null) return 0;
  const k = Math.min(dst.n, str.length);
  for (let i = 0; i < k; i++) dst.a[dst.o + i] = byte(str.charCodeAt(i));
  return k;
}

// sub is s[lo:hi:max] on a slice.
export function sub<T>(s: Slice<T> | null, lo: number, hi?: number, max?: number): Slice<T> | null {
  if (s === null) {
    if (lo === 0 && (hi ?? 0) === 0) return null;
    throw new RangeError("slice bounds out of range");
  }
  const h = hi ?? s.n;
  const m = max ?? s.c;
  if (lo < 0 || h < lo || m < h || m > s.c) {
    throw new RangeError(`slice bounds out of range [${lo}:${h}] with capacity ${s.c}`);
  }
  return new Slice<T>(s.a, s.o + lo, h - lo, m - lo);
}

// subArray is a[lo:hi:max] on an array, sharing its storage.
export function subArray<T>(a: any, lo: number, hi?: number, max?: number): Slice<T> {
  return sub<T>(new Slice<T>(a, 0, a.length, a.length), lo, hi, max) as Slice<T>;
}

// clear is clear(s): every element back to its zero value.
export function clear<T>(s: Slice<T> | null, zero: () => T): void {
  if (s === null) return;
  for (let i = 0; i < s.n; i++) s.a[s.o + i] = zero();
}

// at is s[i] for a slice expression that cannot be written twice.
export function at<T>(s: Slice<T> | null, i: number): T {
  if (s === null) throw new RangeError(`index out of range [${i}] with length 0`);
  return s.get(i);
}
