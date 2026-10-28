// AEADs the runtime has, for the overlays of AES-GCM (aes/gcm) and
// ChaCha20-Poly1305 (vendored x/crypto): node:crypto's, Deno's own, seal
// and open records at native speed, where Go's, transpiled, would take a
// hundred times as long. Where there is no node:crypto, a browser, each
// gives back 0 and Go's own does the work. Only the standard shape is
// handed over, a 12-byte nonce and a 16-byte tag; the overlays check
// everything Go checks first.

import type { Slice } from "./slice.ts";

// A Job is an operation's operands, as the overlays' Go structs hold them.
type Job = { cipher: string; key: Slice<number>; nonce: Slice<number>; text: Slice<number>; ad: Slice<number> | null; out: Slice<number> };

let crypto: any = undefined;

function nodeCrypto(): any {
  if (crypto === undefined) {
    const p = (globalThis as any).process;
    try {
      crypto = typeof p?.getBuiltinModule === "function" ? p.getBuiltinModule("node:crypto") ?? null : null;
    } catch {
      crypto = null;
    }
  }
  return crypto;
}

function view(s: Slice<number> | null): Uint8Array {
  return s === null ? new Uint8Array(0) : (s.a as Uint8Array).subarray(s.o, s.o + s.n);
}

// cipherName is the job's cipher as node:crypto names it, by key size for
// AES.
function cipherName(j: Job): string {
  return j.cipher === "aes-gcm" ? `aes-${j.key.n * 8}-gcm` : j.cipher;
}

// seal writes text sealed, and its tag, to out: 1, or 0 when it cannot.
export function aeadSeal(j: Job): number {
  const c = nodeCrypto();
  if (c === null) return 0;
  const s = c.createCipheriv(cipherName(j), view(j.key), view(j.nonce), { authTagLength: 16 });
  s.setAAD(view(j.ad));
  const body = s.update(view(j.text));
  s.final();
  const out = view(j.out);
  out.set(body, 0);
  out.set(s.getAuthTag(), body.length);
  return 1;
}

// open writes text, sealed with its tag last, opened to out: 1, 2 when
// it is not authentic, which writes nothing, or 0 when it cannot.
//
// What is opened is sealed again and its tag compared with the one given,
// in constant time, rather than trusting final() alone to refuse a wrong
// tag: Deno 1.46 to 2.1.6 did not (CVE-2025-24015).
export function aeadOpen(j: Job): number {
  const c = nodeCrypto();
  if (c === null) return 0;
  const text = view(j.text), n = text.length - 16, tag = text.subarray(n);
  const name = cipherName(j), key = view(j.key), nonce = view(j.nonce), ad = view(j.ad);
  const d = c.createDecipheriv(name, key, nonce, { authTagLength: 16 });
  d.setAAD(ad);
  d.setAuthTag(tag);
  const body = d.update(text.subarray(0, n));
  try {
    d.final();
  } catch {
    return 2;
  }
  const s = c.createCipheriv(name, key, nonce, { authTagLength: 16 });
  s.setAAD(ad);
  s.update(body);
  s.final();
  if (!c.timingSafeEqual(s.getAuthTag(), tag)) return 2;
  view(j.out).set(body, 0);
  return 1;
}
