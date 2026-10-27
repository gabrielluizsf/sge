// The runtime's AEADs refuse what is not authentic even where the
// runtime's own check does not: Deno 1.46 to 2.1.6 never refused a wrong
// AES-GCM tag in final() (CVE-2025-24015). final() is made to accept
// everything here, as it did there.

import { aeadOpen, aeadSeal } from "../../runtime/rt/aead.ts";
import { Slice } from "../../runtime/rt/slice.ts";

const crypto = (globalThis as any).process.getBuiltinModule("node:crypto");

function slice(b: Uint8Array): Slice<number> {
  return new Slice<number>(b, 0, b.length, b.length);
}

function job(cipher: string, key: Uint8Array, text: Uint8Array, out: Uint8Array) {
  return { cipher, key: slice(key), nonce: slice(new Uint8Array(12).fill(3)), text: slice(text), ad: slice(new Uint8Array([9, 9])), out: slice(out) };
}

for (const [cipher, size] of [["aes-gcm", 16], ["aes-gcm", 32], ["chacha20-poly1305", 32]] as const) {
  Deno.test(`${cipher}/${size}: a wrong tag is refused though final() accepts it`, () => {
    const key = new Uint8Array(size).fill(1), plain = new Uint8Array(40).fill(7);
    const sealed = new Uint8Array(56);
    if (aeadSeal(job(cipher, key, plain, sealed)) !== 1) throw new Error("seal failed");
    const opened = new Uint8Array(40);
    if (aeadOpen(job(cipher, key, sealed, opened)) !== 1) throw new Error("authentic text refused");
    const proto = Object.getPrototypeOf(crypto.createDecipheriv(cipher === "aes-gcm" ? `aes-${size * 8}-gcm` : cipher, key, new Uint8Array(12), { authTagLength: 16 }));
    const final = proto.final;
    proto.final = function () {
      return new Uint8Array(0);
    };
    try {
      for (const at of [0, 39, 40, 55]) {
        const bad = sealed.slice();
        bad[at] ^= 1;
        const out = new Uint8Array(40);
        if (aeadOpen(job(cipher, key, bad, out)) !== 2) throw new Error(`tampered byte ${at} accepted`);
        if (out.some((b) => b !== 0)) throw new Error(`tampered byte ${at}: plaintext written`);
      }
    } finally {
      proto.final = final;
    }
  });
}
