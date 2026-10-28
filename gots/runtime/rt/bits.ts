// math/bits's 32-bit sums and products, and uint's, which is 32 bits, on
// numbers: Go's own works it out in uint64, a BigInt here, and math/big,
// bigmod and so RSA and ECDSA spend their time in it. The carries and
// borrows are Go's formulas, with no branches, as Go's are constant time.

import type { Slice } from "./slice.ts";

export function add32(x: number, y: number, carry: number): [number, number] {
  const sum = (x + y + carry) >>> 0;
  return [sum, ((x & y) | ((x | y) & ~sum)) >>> 31];
}

export function sub32(x: number, y: number, borrow: number): [number, number] {
  const diff = (x - y - borrow) >>> 0;
  return [diff, ((~x & y) | (~(x ^ y) & diff)) >>> 31];
}

// mul32 is x*y's high and low words, from the products of their halves,
// each exact in a number.
export function mul32(x: number, y: number): [number, number] {
  const xl = x & 0xffff, xh = x >>> 16, yl = y & 0xffff, yh = y >>> 16;
  const lh = xl * yh, hl = xh * yl;
  const mid = ((xl * yl) >>> 16) + (lh & 0xffff) + (hl & 0xffff);
  return [(xh * yh + (lh >>> 16) + (hl >>> 16) + (mid >>> 16)) >>> 0, Math.imul(x, y) >>> 0];
}

// addMulVVW is bigmod's: z += x*y, word by word, giving back the carry
// word. Each 32-bit product is taken from its 16-bit halves' products,
// exact in a number, so no word goes through an array of two.
export function addMulVVW(z: Slice<number>, x: Slice<number>, y: number): number {
  const za = z.a, zo = z.o, xa = x.a, xo = x.o, n = z.n;
  const yl = y & 0xffff, yh = y >>> 16;
  let carry = 0;
  for (let i = 0; i < n; i++) {
    const xi = xa[xo + i], xl = xi & 0xffff, xh = xi >>> 16;
    const lh = xl * yh, hl = xh * yl;
    // Below 2^35, so exact: the low word and the carry out of it.
    const lo = xl * yl + za[zo + i] + carry + ((lh & 0xffff) + (hl & 0xffff)) * 0x10000;
    const lo32 = lo >>> 0;
    carry = xh * yh + (lh >>> 16) + (hl >>> 16) + (lo - lo32) / 0x100000000;
    za[zo + i] = lo32;
  }
  return carry;
}

// montgomeryWords is bigmod's montgomeryWords (montgomeryMul's loop): T's
// upper half = a·b/R mod m, give or take m, with the carry out of it, for
// R = 2^(32n). It works on 16-bit halves of the words, which keeps R, and
// scans the product a column at a time: a column's sum of products of 32
// bits fits a number exactly, so it takes one carry, not one a product.
export function montgomeryWords(j: { T: Slice<number>; a: Slice<number>; b: Slice<number>; m: Slice<number>; m0inv: number }): number {
  const n = j.m.n, l = 2 * n;
  if (ha.length < l) ha = new Float64Array(l), hb = new Float64Array(l), hm = new Float64Array(l), hq = new Float64Array(l);
  halves(j.a, ha), halves(j.b, hb), halves(j.m, hm);
  const minv = j.m0inv & 0xffff, m0 = hm[0];
  let c = 0;
  for (let i = 0; i < l; i++) {
    let s = c + ha[i] * hb[0], t = 0;
    for (let k = 0; k < i; k++) s += ha[k] * hb[i - k], t += hq[k] * hm[i - k];
    s += t;
    const q = Math.imul(s - floor16(s), minv) & 0xffff;
    hq[i] = q;
    c = (s + q * m0) / 65536;
  }
  const T = j.T.a, o = j.T.o + n;
  let low = 0;
  for (let i = l; i < 2 * l; i++) {
    let s = c, t = 0;
    for (let k = i - l + 1; k < l; k++) s += ha[k] * hb[i - k], t += hq[k] * hm[i - k];
    s += t;
    const h = floor16(s), r = s - h;
    c = h / 65536;
    if ((i & 1) === 0) low = r;
    else T[o + ((i - l) >> 1)] = low + r * 65536;
  }
  return c;
}

let ha = new Float64Array(0), hb = ha, hm = ha, hq = ha;

// floor16 is s rounded down to a multiple of 2^16, for s a whole number
// from 0 to 2^53, by additions alone: the numbers near 3·2^67 are 2^16
// apart, and s - 32767.5, never halfway between two, rounds to the
// nearest. Not %, which fmod computes in a time that depends on s, and so
// on a secret key's words.
const round16 = 3 * 2 ** 67;

function floor16(s: number): number {
  return s - 32767.5 + round16 - round16;
}

// halves writes s's words' 16-bit halves to h, low first.
function halves(s: Slice<number>, h: Float64Array): void {
  const a = s.a, o = s.o;
  for (let i = 0; i < s.n; i++) {
    const w = a[o + i];
    h[2 * i] = w & 0xffff;
    h[2 * i + 1] = w >>> 16;
  }
}

// The rest of math/bits's 32-bit functions, which JavaScript has
// instructions, or near enough, for.

export function rotl32(x: number, k: number): number {
  const s = k & 31;
  return ((x << s) | (x >>> (32 - s))) >>> 0;
}

export function len32(x: number): number {
  return 32 - Math.clz32(x);
}

export function trailingZeros32(x: number): number {
  return x === 0 ? 32 : 31 - Math.clz32(x & -x);
}

export function onesCount32(x: number): number {
  x = x - ((x >>> 1) & 0x55555555);
  x = (x & 0x33333333) + ((x >>> 2) & 0x33333333);
  return (Math.imul((x + (x >>> 4)) & 0x0f0f0f0f, 0x01010101) >>> 24);
}
