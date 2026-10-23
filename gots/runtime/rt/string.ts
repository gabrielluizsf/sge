// Go strings are bytes. Here each byte is one UTF-16 code unit from 0 to 255,
// so len, indexing and slicing are the JavaScript ones. Text is decoded only
// where it leaves for JavaScript (fromGo) and encoded where it comes in (toGo).

export function runeStr(r: number): string {
  if (r < 0 || r > 0x10ffff || (r >= 0xd800 && r <= 0xdfff)) r = 0xfffd;
  if (r < 0x80) return globalThis.String.fromCharCode(r);
  if (r < 0x800) return globalThis.String.fromCharCode(0xc0 | (r >> 6), 0x80 | (r & 0x3f));
  if (r < 0x10000) {
    return globalThis.String.fromCharCode(0xe0 | (r >> 12), 0x80 | ((r >> 6) & 0x3f), 0x80 | (r & 0x3f));
  }
  return globalThis.String.fromCharCode(
    0xf0 | (r >> 18),
    0x80 | ((r >> 12) & 0x3f),
    0x80 | ((r >> 6) & 0x3f),
    0x80 | (r & 0x3f),
  );
}

// decodeRune is the rune starting at byte i of s and its width in bytes,
// or U+FFFD and 1 for an invalid encoding.
export function decodeRune(s: string, i: number): [number, number] {
  const n = s.length - i;
  if (n < 1) return [0xfffd, 0];
  const c0 = s.charCodeAt(i);
  if (c0 < 0x80) return [c0, 1];
  const bad: [number, number] = [0xfffd, 1];
  const cont = (k: number) => {
    const c = s.charCodeAt(i + k);
    return (c & 0xc0) === 0x80 ? c & 0x3f : -1;
  };
  if (c0 < 0xc2) return bad;
  if (c0 < 0xe0) {
    if (n < 2) return bad;
    const c1 = cont(1);
    return c1 < 0 ? bad : [((c0 & 0x1f) << 6) | c1, 2];
  }
  if (c0 < 0xf0) {
    if (n < 3) return bad;
    const c1 = cont(1), c2 = cont(2);
    if (c1 < 0 || c2 < 0) return bad;
    const r = ((c0 & 0x0f) << 12) | (c1 << 6) | c2;
    return r < 0x800 || (r >= 0xd800 && r <= 0xdfff) ? bad : [r, 3];
  }
  if (c0 < 0xf5) {
    if (n < 4) return bad;
    const c1 = cont(1), c2 = cont(2), c3 = cont(3);
    if (c1 < 0 || c2 < 0 || c3 < 0) return bad;
    const r = ((c0 & 0x07) << 18) | (c1 << 12) | (c2 << 6) | c3;
    return r < 0x10000 || r > 0x10ffff ? bad : [r, 4];
  }
  return bad;
}

export function bytesToStr(b: any): string {
  if (b === null) return "";
  let s = "";
  for (let i = 0; i < b.n; i += 8192) {
    s += globalThis.String.fromCharCode.apply(null, b.a.subarray(b.o + i, b.o + Math.min(b.n, i + 8192)));
  }
  return s;
}

export function strToBytes(s: string): Uint8Array {
  const b = new Uint8Array(s.length);
  for (let i = 0; i < s.length; i++) b[i] = s.charCodeAt(i);
  return b;
}

export function strToRunes(s: string): Int32Array {
  const out: number[] = [];
  for (let i = 0; i < s.length;) {
    const [r, w] = decodeRune(s, i);
    out.push(r);
    i += w;
  }
  return Int32Array.from(out);
}

export function runesToStr(r: any): string {
  if (r === null) return "";
  let s = "";
  for (let i = 0; i < r.n; i++) s += runeStr(r.a[r.o + i]);
  return s;
}

// fromGo is the JavaScript text a Go string holds.
export function fromGo(s: string): string {
  let ascii = true;
  for (let i = 0; i < s.length; i++) {
    if (s.charCodeAt(i) >= 0x80) {
      ascii = false;
      break;
    }
  }
  if (ascii) return s;
  let out = "";
  for (let i = 0; i < s.length;) {
    const [r, w] = decodeRune(s, i);
    out += globalThis.String.fromCodePoint(r);
    i += w;
  }
  return out;
}

// toGo is JavaScript text as a Go string of its UTF-8 bytes.
export function toGo(s: string): string {
  let ascii = true;
  for (let i = 0; i < s.length; i++) {
    if (s.charCodeAt(i) >= 0x80) {
      ascii = false;
      break;
    }
  }
  if (ascii) return s;
  let out = "";
  for (const ch of s) out += runeStr(ch.codePointAt(0)!);
  return out;
}

// cmp is strings.Compare: byte order, which is code unit order here.
export function cmp(a: string, b: string): number {
  return a < b ? -1 : a > b ? 1 : 0;
}
