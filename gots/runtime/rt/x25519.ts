// X25519 by the runtime's node:crypto, for crypto/ecdh's overlay
// (x25519_native.go): RFC 7748's function, the very one Go computes, some
// ten times faster than Go's transpiled. It gives back false, for Go's to
// do it, where there is no node:crypto, or where the runtime refuses, as
// it does an all-zero result, which Go's caller refuses after.

import type { Slice } from "./slice.ts";

type Job = { dst: Slice<number>; scalar: Slice<number>; point: Slice<number> };

// The DER around a raw X25519 key, private (PKCS #8) and public (SPKI).
const pkcs8 = [0x30, 0x2e, 0x02, 0x01, 0x00, 0x30, 0x05, 0x06, 0x03, 0x2b, 0x65, 0x6e, 0x04, 0x22, 0x04, 0x20];
const spki = [0x30, 0x2a, 0x30, 0x05, 0x06, 0x03, 0x2b, 0x65, 0x6e, 0x03, 0x21, 0x00];

let crypto: any = undefined;

function wrap(prefix: number[], s: Slice<number>): Uint8Array {
  const b = new Uint8Array(prefix.length + 32);
  b.set(prefix);
  b.set((s.a as Uint8Array).subarray(s.o, s.o + 32), prefix.length);
  return b;
}

export function x25519(j: Job): boolean {
  if (crypto === undefined) {
    const p = (globalThis as any).process;
    try {
      crypto = p?.getBuiltinModule?.("node:crypto") ?? null;
    } catch {
      crypto = null;
    }
  }
  if (crypto === null || j.scalar.n !== 32 || j.point.n !== 32 || j.dst.n < 32) return false;
  let out: Uint8Array;
  try {
    const privateKey = crypto.createPrivateKey({ key: wrap(pkcs8, j.scalar), format: "der", type: "pkcs8" });
    const publicKey = crypto.createPublicKey({ key: wrap(spki, j.point), format: "der", type: "spki" });
    out = crypto.diffieHellman({ privateKey, publicKey });
  } catch {
    return false;
  }
  (j.dst.a as Uint8Array).set(out.subarray(0, 32), j.dst.o);
  return true;
}
