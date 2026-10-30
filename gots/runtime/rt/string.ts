// Go strings are bytes. Here each byte is one UTF-16 code unit: bytes
// below 0x80 as themselves, the others moved into the private use area,
// byte b as U+F700+b. So len, indexing and slicing stay the JavaScript
// ones, and a Go string is told apart from JavaScript text unit by unit:
// text never has code units in that range (save for private characters),
// and a Go string has nothing else outside ASCII. Converting with toGo is
// therefore safe to repeat, which lets the boundary convert whatever
// comes in without knowing where it came from.

const base = 0xf700;

// byte is the byte a code unit of a Go string stands for. A Latin-1
// unit, as a binary string from atob holds, is taken as the byte too.
export function byte(c: number): number {
  return c >= base ? c - base : c;
}

// byteAt is s[i].
export function byteAt(s: string, i: number): number {
  const c = s.charCodeAt(i);
  if (c !== c) throw new RangeError(`index out of range [${i}] with length ${s.length}`);
  return c >= base ? c - base : c;
}

// unit is the code unit holding byte b.
export function unit(b: number): number {
  return b < 0x80 ? b : base + b;
}

const units: string[] = Array.from({ length: 256 }, (_, b) => globalThis.String.fromCharCode(unit(b)));

// byteStr is the one-byte string of b.
export function byteStr(b: number): string {
  return units[b & 0xff];
}

export function runeStr(r: number): string {
  if (r < 0 || r > 0x10ffff || (r >= 0xd800 && r <= 0xdfff)) r = 0xfffd;
  if (r < 0x80) return units[r];
  if (r < 0x800) return units[0xc0 | (r >> 6)] + units[0x80 | (r & 0x3f)];
  if (r < 0x10000) return units[0xe0 | (r >> 12)] + units[0x80 | ((r >> 6) & 0x3f)] + units[0x80 | (r & 0x3f)];
  return units[0xf0 | (r >> 18)] + units[0x80 | ((r >> 12) & 0x3f)] + units[0x80 | ((r >> 6) & 0x3f)] + units[0x80 | (r & 0x3f)];
}

// decodeRune is the rune starting at byte i of s and its width in bytes,
// or U+FFFD and 1 for an invalid encoding.
export function decodeRune(s: string, i: number): [number, number] {
  const n = s.length - i;
  if (n < 1) return [0xfffd, 0];
  const c0 = byte(s.charCodeAt(i));
  if (c0 < 0x80) return [c0, 1];
  const bad: [number, number] = [0xfffd, 1];
  const cont = (k: number) => {
    const c = byte(s.charCodeAt(i + k));
    return (c & 0xc0) === 0x80 && c < 0x100 ? c & 0x3f : -1;
  };
  if (c0 < 0xc2 || c0 > 0xff) return bad;
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

// bytesToStr is string(b) for a []byte.
export function bytesToStr(b: any): string {
  if (b === null || b.n === 0) return "";
  const a = b.a;
  let s = "";
  let ascii = true;
  for (let i = 0; i < b.n; i++) {
    if (a[b.o + i] >= 0x80) {
      ascii = false;
      break;
    }
  }
  if (ascii) {
    for (let i = 0; i < b.n; i += 8192) {
      s += globalThis.String.fromCharCode.apply(null, a.subarray(b.o + i, b.o + Math.min(b.n, i + 8192)));
    }
    return s;
  }
  const out = new Uint16Array(b.n);
  for (let i = 0; i < b.n; i++) out[i] = unit(a[b.o + i]);
  for (let i = 0; i < out.length; i += 8192) {
    s += globalThis.String.fromCharCode.apply(null, out.subarray(i, Math.min(out.length, i + 8192)) as any);
  }
  return s;
}

// strToBytes is []byte(s): the bytes of a Go string, or of a binary
// string of Latin-1 units.
export function strToBytes(s: string): Uint8Array {
  const b = new Uint8Array(s.length);
  for (let i = 0; i < s.length; i++) b[i] = byte(s.charCodeAt(i));
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

const goOnly = /^[\x00-\x7f-]*$/;
const asciiOnly = /^[\x00-\x7f]*$/;

/**
 * fromGo is the JavaScript text a Go string holds. Text already in it,
 * which a string from TypeScript is, passes through as it is.
 */
export function fromGo(s: string): string {
  if (asciiOnly.test(s)) return s;
  let out = "";
  for (let i = 0; i < s.length;) {
    const c = s.charCodeAt(i);
    if (c >= 0x80 && (c < 0xf780 || c > 0xf7ff)) {
      out += s[i];
      i++;
      continue;
    }
    const [r, w] = decodeRune(s, i);
    out += globalThis.String.fromCodePoint(r);
    i += w;
  }
  return out;
}

/**
 * toGo is s as a Go string: JavaScript text in it becomes its UTF-8
 * bytes, and a Go string stays as it is.
 */
export function toGo(s: string): string {
  if (goOnly.test(s)) return s;
  let out = "";
  for (const ch of s) {
    const c = ch.charCodeAt(0);
    if (c < 0x80 || (c >= 0xf780 && c <= 0xf7ff)) out += ch;
    else out += runeStr(ch.codePointAt(0)!);
  }
  return out;
}

// cmp is strings.Compare: byte order, which code unit order is here.
export function cmp(a: string, b: string): number {
  return a < b ? -1 : a > b ? 1 : 0;
}
