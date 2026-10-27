// Keccak-f[1600] for crypto/internal/fips140/sha3's keccakF1600Generic,
// on the state's 200 bytes as they are: Go's reads them as 25 uint64
// lanes through unsafe, a cast gots cannot make, and its lanes would be
// BigInts. Here a lane is a pair of 32-bit words, high and low.

const lo = new Int32Array(25), hi = new Int32Array(25);
const bcl = new Int32Array(5), bch = new Int32Array(5);

// The round constants, high then low.
const RC = new Int32Array([
  0x00000000, 0x00000001, 0x00000000, 0x00008082, 0x80000000, 0x0000808a, 0x80000000, 0x80008000,
  0x00000000, 0x0000808b, 0x00000000, 0x80000001, 0x80000000, 0x80008081, 0x80000000, 0x00008009,
  0x00000000, 0x0000008a, 0x00000000, 0x00000088, 0x00000000, 0x80008009, 0x00000000, 0x8000000a,
  0x00000000, 0x8000808b, 0x80000000, 0x0000008b, 0x80000000, 0x00008089, 0x80000000, 0x00008003,
  0x80000000, 0x00008002, 0x80000000, 0x00000080, 0x00000000, 0x0000800a, 0x80000000, 0x8000000a,
  0x80000000, 0x80008081, 0x80000000, 0x00008080, 0x00000000, 0x80000001, 0x80000000, 0x80008008,
]);

// Where rho and pi move each lane, and by how much it is rotated.
const PI = [10, 7, 11, 17, 18, 3, 5, 16, 8, 21, 24, 4, 15, 23, 19, 13, 12, 2, 20, 14, 22, 9, 6, 1];
const ROT = [1, 3, 6, 10, 15, 21, 28, 36, 45, 55, 2, 14, 27, 41, 56, 8, 25, 43, 62, 18, 39, 61, 20, 44];

export function keccakF1600(state: Uint8Array): void {
  for (let i = 0; i < 25; i++) {
    const o = i * 8;
    lo[i] = state[o] | (state[o + 1] << 8) | (state[o + 2] << 16) | (state[o + 3] << 24);
    hi[i] = state[o + 4] | (state[o + 5] << 8) | (state[o + 6] << 16) | (state[o + 7] << 24);
  }
  for (let r = 0; r < 24; r++) round(r);
  for (let i = 0; i < 25; i++) {
    const o = i * 8, l = lo[i], h = hi[i];
    state[o] = l, state[o + 1] = l >>> 8, state[o + 2] = l >>> 16, state[o + 3] = l >>> 24;
    state[o + 4] = h, state[o + 5] = h >>> 8, state[o + 6] = h >>> 16, state[o + 7] = h >>> 24;
  }
}

function round(r: number): void {
  // θ
  for (let x = 0; x < 5; x++) {
    bcl[x] = lo[x] ^ lo[x + 5] ^ lo[x + 10] ^ lo[x + 15] ^ lo[x + 20];
    bch[x] = hi[x] ^ hi[x + 5] ^ hi[x + 10] ^ hi[x + 15] ^ hi[x + 20];
  }
  for (let x = 0; x < 5; x++) {
    const nl = bcl[(x + 1) % 5], nh = bch[(x + 1) % 5];
    const dl = bcl[(x + 4) % 5] ^ ((nl << 1) | (nh >>> 31)), dh = bch[(x + 4) % 5] ^ ((nh << 1) | (nl >>> 31));
    for (let y = 0; y < 25; y += 5) lo[y + x] ^= dl, hi[y + x] ^= dh;
  }
  // ρ and π
  let cl = lo[1], ch = hi[1];
  for (let i = 0; i < 24; i++) {
    const j = PI[i], s = ROT[i], tl = lo[j], th = hi[j];
    if (s < 32) {
      lo[j] = (cl << s) | (ch >>> (32 - s)), hi[j] = (ch << s) | (cl >>> (32 - s));
    } else {
      lo[j] = (ch << (s - 32)) | (cl >>> (64 - s)), hi[j] = (cl << (s - 32)) | (ch >>> (64 - s));
    }
    cl = tl, ch = th;
  }
  // χ
  for (let y = 0; y < 25; y += 5) {
    for (let x = 0; x < 5; x++) bcl[x] = lo[y + x], bch[x] = hi[y + x];
    for (let x = 0; x < 5; x++) {
      lo[y + x] ^= ~bcl[(x + 1) % 5] & bcl[(x + 2) % 5];
      hi[y + x] ^= ~bch[(x + 1) % 5] & bch[(x + 2) % 5];
    }
  }
  // ι
  hi[0] ^= RC[r * 2], lo[0] ^= RC[r * 2 + 1];
}
