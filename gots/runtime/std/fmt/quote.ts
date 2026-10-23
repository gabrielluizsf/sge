import * as $ from "../../rt/mod.ts";

// printable is a close-enough unicode.IsPrint for quoting.
function printable(r: number): boolean {
  if (r < 0x20 || r === 0x7f) return false;
  if (r >= 0x80 && r < 0xa0) return false;
  if (r === 0xad || r === 0xfffd) return r === 0xfffd ? true : false;
  if (r >= 0xd800 && r <= 0xdfff) return false;
  if (r === 0x2028 || r === 0x2029 || (r >= 0x200b && r <= 0x200f) || r === 0xfeff) return false;
  return true;
}

function hex(n: number, w: number): string {
  return n.toString(16).padStart(w, "0");
}

function escape(r: number, quote: number, ascii: boolean): string {
  if (r === quote || r === 0x5c) return "\\" + globalThis.String.fromCharCode(r);
  if (ascii ? r < 0x80 && printable(r) : printable(r)) return $.str.runeStr(r);
  switch (r) {
    case 7: return "\\a";
    case 8: return "\\b";
    case 12: return "\\f";
    case 10: return "\\n";
    case 13: return "\\r";
    case 9: return "\\t";
    case 11: return "\\v";
  }
  if (r < 0x20 || r === 0x7f) return "\\x" + hex(r, 2);
  if (r < 0x10000) return "\\u" + hex(r, 4);
  return "\\U" + hex(r, 8);
}

// quote is strconv.Quote over a Go string, or QuoteToASCII.
export function quote(s: string, ascii = false): string {
  let out = '"';
  for (let i = 0; i < s.length;) {
    const [r, w] = $.str.decodeRune(s, i);
    if (r === 0xfffd && w === 1 && $.str.byte(s.charCodeAt(i)) !== 0xef) {
      out += "\\x" + hex($.str.byte(s.charCodeAt(i)), 2);
    } else {
      out += escape(r, 0x22, ascii);
    }
    i += w;
  }
  return out + '"';
}

// quoteRune is strconv.QuoteRune.
export function quoteRune(r: number, ascii = false): string {
  if (r > 0x10ffff || (r >= 0xd800 && r <= 0xdfff)) r = 0xfffd;
  return "'" + escape(r, 0x27, ascii) + "'";
}
