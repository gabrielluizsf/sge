// GHASH's step for crypto/internal/fips140/aes/gcm's overlay (ghash.go):
// Y = (Y ^ block)·H, as Go computes it, on 32-bit words held in numbers,
// with no array or object made along the way. Go's own, on 64-bit words,
// would have each a BigInt.

import type { Slice } from "./slice.ts";

type State = { w0: number; w1: number; w2: number; w3: number };
type Key = { h0: number; h1: number; h2: number; h3: number };

// The last product's words, high and low, set by mul and clmul.
let rh = 0, rl = 0;

function mul(a: number, b: number): void {
  const al = a & 0xffff, ah = a >>> 16, bl = b & 0xffff, bh = b >>> 16;
  const lh = al * bh, hl = ah * bl;
  const mid = ((al * bl) >>> 16) + (lh & 0xffff) + (hl & 0xffff);
  rh = (ah * bh + (lh >>> 16) + (hl >>> 16) + (mid >>> 16)) >>> 0;
  rl = Math.imul(a, b) >>> 0;
}

// clmul is x·y carry-less: integer products of x and y with every fourth
// bit kept, whose carries fall in the bits between, masked off.
function clmul(x: number, y: number): void {
  const x0 = x & 0x11111111, x1 = x & 0x22222222, x2 = x & 0x44444444, x3 = (x & 0x88888888) >>> 0;
  const y0 = y & 0x11111111, y1 = y & 0x22222222, y2 = y & 0x44444444, y3 = (y & 0x88888888) >>> 0;
  let h = 0, l = 0, th = 0, tl = 0;
  mul(x0, y0), th = rh, tl = rl, mul(x1, y3), th ^= rh, tl ^= rl, mul(x2, y2), th ^= rh, tl ^= rl, mul(x3, y1);
  h |= (th ^ rh) & 0x11111111, l |= (tl ^ rl) & 0x11111111;
  mul(x0, y1), th = rh, tl = rl, mul(x1, y0), th ^= rh, tl ^= rl, mul(x2, y3), th ^= rh, tl ^= rl, mul(x3, y2);
  h |= (th ^ rh) & 0x22222222, l |= (tl ^ rl) & 0x22222222;
  mul(x0, y2), th = rh, tl = rl, mul(x1, y1), th ^= rh, tl ^= rl, mul(x2, y0), th ^= rh, tl ^= rl, mul(x3, y3);
  h |= (th ^ rh) & 0x44444444, l |= (tl ^ rl) & 0x44444444;
  mul(x0, y3), th = rh, tl = rl, mul(x1, y2), th ^= rh, tl ^= rl, mul(x2, y1), th ^= rh, tl ^= rl, mul(x3, y0);
  h |= (th ^ rh) & 0x88888888, l |= (tl ^ rl) & 0x88888888;
  rh = h >>> 0, rl = l >>> 0;
}

function be32(a: Uint8Array, i: number): number {
  return ((a[i] << 24) | (a[i + 1] << 16) | (a[i + 2] << 8) | a[i + 3]) >>> 0;
}

export function ghashAbsorb(y: State, k: Key, block: Slice<number>): void {
  const b = block.a as Uint8Array, o = block.o;
  const y3 = (y.w3 ^ be32(b, o)) >>> 0, y2 = (y.w2 ^ be32(b, o + 4)) >>> 0;
  const y1 = (y.w1 ^ be32(b, o + 8)) >>> 0, y0 = (y.w0 ^ be32(b, o + 12)) >>> 0;
  // Karatsuba's three 64-bit products for each half and their sum.
  clmul(y0, k.h0); const lo0h = rh, lo0l = rl;
  clmul(y1, k.h1); const hi0h = rh, hi0l = rl;
  clmul(y0 ^ y1, k.h0 ^ k.h1); const s0h = rh ^ lo0h ^ hi0h, s0l = rl ^ lo0l ^ hi0l;
  clmul(y2, k.h2); const lo1h = rh, lo1l = rl;
  clmul(y3, k.h3); const hi1h = rh, hi1l = rl;
  clmul(y2 ^ y3, k.h2 ^ k.h3); const s1h = rh ^ lo1h ^ hi1h, s1l = rl ^ lo1l ^ hi1l;
  clmul(y0 ^ y2, k.h0 ^ k.h2); const lo2h = rh, lo2l = rl;
  clmul(y1 ^ y3, k.h1 ^ k.h3); const hi2h = rh, hi2l = rl;
  clmul(y0 ^ y2 ^ y1 ^ y3, k.h0 ^ k.h2 ^ k.h1 ^ k.h3); const s2h = rh ^ lo2h ^ hi2h, s2l = rl ^ lo2l ^ hi2l;
  // halves: r[i][0] = lo ^ mid<<32, r[i][1] = hi ^ mid>>32, words high, low.
  const r00h = lo0h ^ s0l, r00l = lo0l, r01h = hi0h, r01l = hi0l ^ s0h;
  let r10h = lo1h ^ s1l, r10l = lo1l;
  const r11h = hi1h, r11l = hi1l ^ s1h;
  const r20h = lo2h ^ s2l ^ r00h ^ r10h, r20l = lo2l ^ r00l ^ r10l;
  const r21h = hi2h ^ r01h ^ r11h, r21l = hi2l ^ s2h ^ r01l ^ r11l;
  const q01h = r01h ^ r20h, q01l = r01l ^ r20l;
  r10h ^= r21h, r10l ^= r21l;
  // The product shifted by one bit, high word then low of each.
  const z0h = (r00h << 1) | (r00l >>> 31), z0l = r00l << 1;
  let z1h = (q01h << 1) | (q01l >>> 31), z1l = (q01l << 1) | (r00h >>> 31);
  let z2h = (r10h << 1) | (r10l >>> 31), z2l = (r10l << 1) | (q01h >>> 31);
  let z3h = (r11h << 1) | (r11l >>> 31), z3l = (r11l << 1) | (r10h >>> 31);
  // fold(z0) into z2 and z1, then fold(z1) into z3 and z2.
  z2h ^= z0h ^ (z0h >>> 1) ^ (z0h >>> 2) ^ (z0h >>> 7);
  z2l ^= z0l ^ ((z0l >>> 1) | (z0h << 31)) ^ ((z0l >>> 2) | (z0h << 30)) ^ ((z0l >>> 7) | (z0h << 25));
  z1h ^= (z0l << 31) ^ (z0l << 30) ^ (z0l << 25);
  z3h ^= z1h ^ (z1h >>> 1) ^ (z1h >>> 2) ^ (z1h >>> 7);
  z3l ^= z1l ^ ((z1l >>> 1) | (z1h << 31)) ^ ((z1l >>> 2) | (z1h << 30)) ^ ((z1l >>> 7) | (z1h << 25));
  z2h ^= (z1l << 31) ^ (z1l << 30) ^ (z1l << 25);
  y.w0 = z2l >>> 0, y.w1 = z2h >>> 0, y.w2 = z3l >>> 0, y.w3 = z3h >>> 0;
}
