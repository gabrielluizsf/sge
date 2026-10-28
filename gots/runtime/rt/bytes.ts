// The byte-order functions on a subslice, b[lo:] or b[lo:hi], fused with
// it (see byteorder_fuse.go): they read or write b's array in place, where
// the subslice would be a new object made only to be read four bytes of.
// The bounds are checked as Go checks the subslice's and the read's.

import type { Slice } from "./slice.ts";

// at is the index in s's array of s[lo:hi]'s first byte, which has at
// least k bytes; hi is s's length when not given.
function at(s: Slice<number> | null, lo: number, hi: number | undefined, k: number): number {
  const n = s === null ? 0 : s.n, c = s === null ? 0 : s.c, h = hi ?? n;
  if (lo < 0 || h < lo || h > (hi === undefined ? n : c)) {
    throw new RangeError(`slice bounds out of range [${lo}:${h}] with capacity ${c}`);
  }
  if (h - lo < k) throw new RangeError(`index out of range [${k - 1}] with length ${h - lo}`);
  return s!.o + lo;
}

export function be16(s: Slice<number> | null, lo: number, hi?: number): number {
  const a = s?.a, i = at(s, lo, hi, 2);
  return (a[i] << 8) | a[i + 1];
}

export function be32(s: Slice<number> | null, lo: number, hi?: number): number {
  const a = s?.a, i = at(s, lo, hi, 4);
  return ((a[i] << 24) | (a[i + 1] << 16) | (a[i + 2] << 8) | a[i + 3]) >>> 0;
}

export function le16(s: Slice<number> | null, lo: number, hi?: number): number {
  const a = s?.a, i = at(s, lo, hi, 2);
  return a[i] | (a[i + 1] << 8);
}

export function le32(s: Slice<number> | null, lo: number, hi?: number): number {
  const a = s?.a, i = at(s, lo, hi, 4);
  return (a[i] | (a[i + 1] << 8) | (a[i + 2] << 16) | (a[i + 3] << 24)) >>> 0;
}

export function putBE16(s: Slice<number> | null, lo: number, v: number): void {
  const a = s?.a, i = at(s, lo, undefined, 2);
  a[i] = v >>> 8, a[i + 1] = v;
}

export function putBE32(s: Slice<number> | null, lo: number, v: number): void {
  const a = s?.a, i = at(s, lo, undefined, 4);
  a[i] = v >>> 24, a[i + 1] = v >>> 16, a[i + 2] = v >>> 8, a[i + 3] = v;
}

export function putLE16(s: Slice<number> | null, lo: number, v: number): void {
  const a = s?.a, i = at(s, lo, undefined, 2);
  a[i] = v, a[i + 1] = v >>> 8;
}

export function putLE32(s: Slice<number> | null, lo: number, v: number): void {
  const a = s?.a, i = at(s, lo, undefined, 4);
  a[i] = v, a[i + 1] = v >>> 8, a[i + 2] = v >>> 16, a[i + 3] = v >>> 24;
}
