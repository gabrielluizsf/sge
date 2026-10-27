// SHA-512's block function for crypto/internal/fips140/sha512's
// blockGeneric: Go's rounds on 64-bit words as pairs of 32-bit ones, high
// and low, in numbers; Go's own would have each word a BigInt. The state,
// Go's [8]uint64, is read and written once a call.

import type { Slice } from "./slice.ts";

// K's words, high then low.
const K = new Int32Array([
  0x428a2f98, 0xd728ae22, 0x71374491, 0x23ef65cd, 0xb5c0fbcf, 0xec4d3b2f, 0xe9b5dba5, 0x8189dbbc,
  0x3956c25b, 0xf348b538, 0x59f111f1, 0xb605d019, 0x923f82a4, 0xaf194f9b, 0xab1c5ed5, 0xda6d8118,
  0xd807aa98, 0xa3030242, 0x12835b01, 0x45706fbe, 0x243185be, 0x4ee4b28c, 0x550c7dc3, 0xd5ffb4e2,
  0x72be5d74, 0xf27b896f, 0x80deb1fe, 0x3b1696b1, 0x9bdc06a7, 0x25c71235, 0xc19bf174, 0xcf692694,
  0xe49b69c1, 0x9ef14ad2, 0xefbe4786, 0x384f25e3, 0xfc19dc6, 0x8b8cd5b5, 0x240ca1cc, 0x77ac9c65,
  0x2de92c6f, 0x592b0275, 0x4a7484aa, 0x6ea6e483, 0x5cb0a9dc, 0xbd41fbd4, 0x76f988da, 0x831153b5,
  0x983e5152, 0xee66dfab, 0xa831c66d, 0x2db43210, 0xb00327c8, 0x98fb213f, 0xbf597fc7, 0xbeef0ee4,
  0xc6e00bf3, 0x3da88fc2, 0xd5a79147, 0x930aa725, 0x6ca6351, 0xe003826f, 0x14292967, 0xa0e6e70,
  0x27b70a85, 0x46d22ffc, 0x2e1b2138, 0x5c26c926, 0x4d2c6dfc, 0x5ac42aed, 0x53380d13, 0x9d95b3df,
  0x650a7354, 0x8baf63de, 0x766a0abb, 0x3c77b2a8, 0x81c2c92e, 0x47edaee6, 0x92722c85, 0x1482353b,
  0xa2bfe8a1, 0x4cf10364, 0xa81a664b, 0xbc423001, 0xc24b8b70, 0xd0f89791, 0xc76c51a3, 0x654be30,
  0xd192e819, 0xd6ef5218, 0xd6990624, 0x5565a910, 0xf40e3585, 0x5771202a, 0x106aa070, 0x32bbd1b8,
  0x19a4c116, 0xb8d2d0c8, 0x1e376c08, 0x5141ab53, 0x2748774c, 0xdf8eeb99, 0x34b0bcb5, 0xe19b48a8,
  0x391c0cb3, 0xc5c95a63, 0x4ed8aa4a, 0xe3418acb, 0x5b9cca4f, 0x7763e373, 0x682e6ff3, 0xd6b2b8a3,
  0x748f82ee, 0x5defb2fc, 0x78a5636f, 0x43172f60, 0x84c87814, 0xa1f0ab72, 0x8cc70208, 0x1a6439ec,
  0x90befffa, 0x23631e28, 0xa4506ceb, 0xde82bde9, 0xbef9a3f7, 0xb2c67915, 0xc67178f2, 0xe372532b,
  0xca273ece, 0xea26619c, 0xd186b8c7, 0x21c0c207, 0xeada7dd6, 0xcde0eb1e, 0xf57d4f7f, 0xee6ed178,
  0x6f067aa, 0x72176fba, 0xa637dc5, 0xa2c898a6, 0x113f9804, 0xbef90dae, 0x1b710b35, 0x131c471b,
  0x28db77f5, 0x23047d84, 0x32caab7b, 0x40c72493, 0x3c9ebe0a, 0x15c9bebc, 0x431d67c4, 0x9c100d4c,
  0x4cc5d4be, 0xcb3e42b6, 0x597f299c, 0xfc657e2a, 0x5fcb6fab, 0x3ad6faec, 0x6c44198c, 0x4a475817,]);

const wh = new Int32Array(80), wl = new Int32Array(80);
const hh = new Int32Array(8), hl = new Int32Array(8);

export function sha512Block(dig: { h: BigUint64Array }, p: Slice<number>): void {
  const h = dig.h;
  for (let i = 0; i < 8; i++) hh[i] = Number(h[i] >> 32n), hl[i] = Number(h[i] & 0xffffffffn);
  const a = p.a as Uint8Array;
  for (let o = p.o, end = p.o + p.n - 127; o < end; o += 128) rounds(a, o);
  for (let i = 0; i < 8; i++) h[i] = (BigInt(hh[i] >>> 0) << 32n) | BigInt(hl[i] >>> 0);
}

// rounds compresses the block of a at o into hh and hl.
function rounds(a: Uint8Array, o: number): void {
  for (let i = 0; i < 16; i++) {
    const j = o + i * 8;
    wh[i] = (a[j] << 24) | (a[j + 1] << 16) | (a[j + 2] << 8) | a[j + 3];
    wl[i] = (a[j + 4] << 24) | (a[j + 5] << 16) | (a[j + 6] << 8) | a[j + 7];
  }
  for (let i = 16; i < 80; i++) {
    // σ1(w[i-2]) = rotr19 ^ rotr61 ^ shr6; σ0(w[i-15]) = rotr1 ^ rotr8 ^ shr7.
    let xh = wh[i - 2], xl = wl[i - 2];
    const s1h = ((xh >>> 19) | (xl << 13)) ^ ((xl >>> 29) | (xh << 3)) ^ (xh >>> 6);
    const s1l = ((xl >>> 19) | (xh << 13)) ^ ((xh >>> 29) | (xl << 3)) ^ ((xl >>> 6) | (xh << 26));
    xh = wh[i - 15], xl = wl[i - 15];
    const s0h = ((xh >>> 1) | (xl << 31)) ^ ((xh >>> 8) | (xl << 24)) ^ (xh >>> 7);
    const s0l = ((xl >>> 1) | (xh << 31)) ^ ((xl >>> 8) | (xh << 24)) ^ ((xl >>> 7) | (xh << 25));
    const lo = (s1l >>> 0) + (wl[i - 7] >>> 0) + (s0l >>> 0) + (wl[i - 16] >>> 0);
    wl[i] = lo;
    wh[i] = s1h + wh[i - 7] + s0h + wh[i - 16] + ((lo / 0x100000000) | 0);
  }
  let ah = hh[0], al = hl[0], bh = hh[1], bl = hl[1], ch = hh[2], cl = hl[2], dh = hh[3], dl = hl[3];
  let eh = hh[4], el = hl[4], fh = hh[5], fl = hl[5], gh = hh[6], gl = hl[6], h_h = hh[7], h_l = hl[7];
  for (let i = 0; i < 80; i++) {
    // Σ1(e) = rotr14 ^ rotr18 ^ rotr41; ch(e, f, g).
    const S1h = ((eh >>> 14) | (el << 18)) ^ ((eh >>> 18) | (el << 14)) ^ ((el >>> 9) | (eh << 23));
    const S1l = ((el >>> 14) | (eh << 18)) ^ ((el >>> 18) | (eh << 14)) ^ ((eh >>> 9) | (el << 23));
    const chh = (eh & fh) ^ (~eh & gh), chl = (el & fl) ^ (~el & gl);
    const t1l = (h_l >>> 0) + (S1l >>> 0) + (chl >>> 0) + (K[i * 2 + 1] >>> 0) + (wl[i] >>> 0);
    const t1h = h_h + S1h + chh + K[i * 2] + wh[i] + ((t1l / 0x100000000) | 0);
    // Σ0(a) = rotr28 ^ rotr34 ^ rotr39; maj(a, b, c).
    const S0h = ((ah >>> 28) | (al << 4)) ^ ((al >>> 2) | (ah << 30)) ^ ((al >>> 7) | (ah << 25));
    const S0l = ((al >>> 28) | (ah << 4)) ^ ((ah >>> 2) | (al << 30)) ^ ((ah >>> 7) | (al << 25));
    const mjh = (ah & bh) ^ (ah & ch) ^ (bh & ch), mjl = (al & bl) ^ (al & cl) ^ (bl & cl);
    const t2l = (S0l >>> 0) + (mjl >>> 0);
    const t2h = S0h + mjh + ((t2l / 0x100000000) | 0);
    h_h = gh, h_l = gl, gh = fh, gl = fl, fh = eh, fl = el;
    let s = (dl >>> 0) + (t1l >>> 0);
    eh = dh + t1h + ((s / 0x100000000) | 0), el = s;
    dh = ch, dl = cl, ch = bh, cl = bl, bh = ah, bl = al;
    s = (t1l >>> 0) + (t2l >>> 0);
    ah = t1h + t2h + ((s / 0x100000000) | 0), al = s;
  }
  add(0, ah, al), add(1, bh, bl), add(2, ch, cl), add(3, dh, dl);
  add(4, eh, el), add(5, fh, fl), add(6, gh, gl), add(7, h_h, h_l);
}

function add(i: number, h: number, l: number): void {
  const s = (hl[i] >>> 0) + (l >>> 0);
  hl[i] = s;
  hh[i] = hh[i] + h + ((s / 0x100000000) | 0);
}
