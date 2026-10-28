import * as $ from "../../rt/mod.ts";
import { formatFloat } from "./float.ts";
import { quote, quoteRune } from "./quote.ts";

// Flags are what a verb was written with, between % and the verb letter.
type Flags = {
  // plusV and sharpV are %+v and %#v, which name struct fields and print
  // Go syntax without passing + or # on to the values inside.
  plusV: boolean;
  sharpV: boolean;
  plus: boolean;
  minus: boolean;
  sharp: boolean;
  space: boolean;
  zero: boolean;
  width: number;
  prec: number;
};

function noFlags(): Flags {
  return { plusV: false, sharpV: false, plus: false, minus: false, sharp: false, space: false, zero: false, width: -1, prec: -1 };
}

function byteOf(s: string, i: number): number {
  return $.str.byte(s.charCodeAt(i));
}

function runeCount(s: string): number {
  let n = 0;
  for (let i = 0; i < s.length; i++) if ((byteOf(s, i) & 0xc0) !== 0x80) n++;
  return n;
}

function pad(s: string, f: Flags, zeroOK = false): string {
  const n = runeCount(s);
  if (f.width <= n) return s;
  const fill = f.width - n;
  if (f.minus) return s + " ".repeat(fill);
  if (f.zero && zeroOK) {
    const sign = s[0] === "-" || s[0] === "+" || s[0] === " " ? s[0] : "";
    return sign + "0".repeat(fill) + s.slice(sign.length);
  }
  return " ".repeat(fill) + s;
}

const fakeAddresses = new WeakMap<object, number>();
let nextAddress = 0xc000010000;

function address(o: any): string {
  if (o === null) return "0x0";
  if (typeof o !== "object" && typeof o !== "function") return "0x0";
  let a = fakeAddresses.get(o);
  if (a === undefined) {
    a = nextAddress;
    nextAddress += 0x10;
    fakeAddresses.set(o, a);
  }
  return "0x" + a.toString(16);
}

function isInt(k: $.Kind): boolean {
  return k === "int" || k === "int8" || k === "int16" || k === "int32" || k === "int64" ||
    k === "uint" || k === "uint8" || k === "uint16" || k === "uint32" || k === "uint64" || k === "uintptr";
}

function isUnsigned(k: $.Kind): boolean {
  return k[0] === "u";
}

function fmtInt(v: number | bigint, k: $.Kind, verb: string, f: Flags): string {
  let neg = v < 0;
  let a: bigint = typeof v === "bigint" ? v : BigInt(Math.trunc(v as number));
  if (neg) a = -a;
  let digits: string;
  switch (verb) {
    case "b": digits = a.toString(2); break;
    case "o": case "O": digits = a.toString(8); break;
    case "x": digits = a.toString(16); break;
    case "X": digits = a.toString(16).toUpperCase(); break;
    case "c": return pad($.str.runeStr(Number(v)), f);
    case "q": return pad(quoteRune(Number(v), f.plus), f);
    case "U": {
      let s = "U+" + a.toString(16).toUpperCase().padStart(4, "0");
      if (f.sharp) s += " '" + $.str.runeStr(Number(v)) + "'";
      return pad(s, f);
    }
    default: digits = a.toString(10);
  }
  if (f.prec >= 0) {
    if (f.prec === 0 && a === 0n) digits = "";
    digits = digits.padStart(f.prec, "0");
  }
  let prefix = "";
  if (f.sharp) {
    if (verb === "x") prefix = "0x";
    else if (verb === "X") prefix = "0X";
    else if (verb === "o" && digits[0] !== "0") prefix = "0";
    else if (verb === "b") prefix = "0b";
  }
  if (verb === "O") prefix = "0o";
  const sign = neg ? "-" : f.plus ? "+" : f.space && !isUnsigned(k) ? " " : "";
  return pad(sign + prefix + digits, f, f.prec < 0);
}

function fmtFloatVerb(v: number, bits: number, verb: string, f: Flags): string {
  let fmt = verb, prec = f.prec;
  if (verb === "v") fmt = "g";
  if (verb === "F") fmt = "f";
  if (prec < 0 && (fmt === "f" || fmt === "e" || fmt === "E")) prec = 6;
  let s = formatFloat(v, fmt, prec, bits);
  if (f.plus && s[0] !== "-" && s[0] !== "+") s = "+" + s;
  else if (f.space && s[0] !== "-") s = " " + s;
  else if (s[0] === "+" && !f.plus && s !== "+Inf") s = s.slice(1);
  if (s.endsWith("Inf") || s.endsWith("NaN")) return pad(s, { ...f, zero: false });
  return pad(s, f, true);
}

function hexOfBytes(s: string, upper: boolean, f: Flags): string {
  let out = "";
  for (let i = 0; i < s.length; i++) {
    if (i > 0 && f.space) out += " ";
    if (f.sharp && (i === 0 || f.space)) out += upper ? "0X" : "0x";
    const h = byteOf(s, i).toString(16).padStart(2, "0");
    out += upper ? h.toUpperCase() : h;
  }
  return out;
}

function fmtString(s: string, verb: string, f: Flags): string {
  if (f.prec >= 0) {
    let cut = 0, n = 0;
    while (cut < s.length && n < f.prec) {
      cut += $.str.decodeRune(s, cut)[1];
      n++;
    }
    s = s.slice(0, cut);
  }
  switch (verb) {
    case "q":
      return pad(f.sharp && !s.includes("`") ? "`" + s + "`" : quote(s, f.plus), f);
    case "x":
      return pad(hexOfBytes(s, false, f), f);
    case "X":
      return pad(hexOfBytes(s, true, f), f);
    default:
      return pad(s, f);
  }
}

// methodOf finds the Error or String method fmt prefers over printing a
// value's contents, for an interface value or a plain value of type t. A
// nil pointer has them too, when the interface holding it is given as
// boxed: they are called, as Go calls them, with a nil receiver.
function methodOf(t: $.Type, v: any, boxed: any = null): (() => string) | null {
  let recv: any = null;
  if (t.kind === "ptr") {
    if (v === null && !(boxed instanceof $.NilPtr)) return null;
    recv = v ?? boxed;
  } else if (t.valCls !== null) {
    recv = new t.valCls(t, v);
  }
  if (recv === null) return null;
  // TypeScript-facing methods hand back text; toGo makes it a Go string
  // either way.
  if (typeof recv.Error === "function") return () => $.str.toGo(recv.Error());
  if (typeof recv.String === "function") return () => $.str.toGo(recv.String());
  return null;
}

class Printer {
  out = "";

  badVerb(verb: string, t: $.Type | null, v: any): string {
    if (t === null) return `%!${verb}(<nil>)`;
    return `%!${verb}(${t.str}=${this.value(t, v, "v", noFlags(), 0)})`;
  }

  // arg prints one operand, which is an interface value.
  arg(a: any, verb: string, f: Flags): string {
    if (a === null) {
      if (verb === "T" || verb === "v") return pad("<nil>", f);
      return this.badVerb(verb, null, null);
    }
    const t = $.typeOf(a);
    const v = a instanceof $.Val ? a.v : a;
    if (verb === "T") return pad(t.str, f);
    if (verb === "p") return pad(address(v), f);
    return this.value(t, v, verb, f, 0, a);
  }

  value(t: $.Type, v: any, verb: string, f: Flags, depth: number, boxed: any = null): string {
    if (t.str === "reflect.Value" && v !== null && "p" in v && "f" in v) {
      // fmt prints what a reflect.Value holds, as Go's does.
      if (v.t === null) return pad("<invalid reflect.Value>", f);
      return this.value(v.t, v.get(), verb, f, depth);
    }
    if (verb !== "T" && verb !== "p" && !(verb === "v" && f.sharpV) && "vsxXq".includes(verb)) {
      const m = methodOf(t, v, boxed);
      if (m !== null) {
        try {
          return fmtString(m(), verb, f);
        } catch (e) {
          if (t.kind === "ptr" && v === null) return pad("<nil>", f);
          return `%!${verb}(PANIC=${verb === "v" ? "String" : "Error"} method: ${e instanceof $.Panic ? globalThis.String(e.value?.v ?? e.value) : e})`;
        }
      }
    }
    const k = t.kind;
    if (k === "bool") {
      if (verb === "v" || verb === "t") return pad(v ? "true" : "false", f);
      return this.badVerb(verb, t, v);
    }
    if (isInt(k)) {
      if ("vdboOxXcqU".includes(verb)) {
        if (verb === "v" && f.sharpV && k === "uint8") return fmtInt(v, k, "x", { ...f, sharp: true });
        return fmtInt(v, k, verb === "v" ? "d" : verb, f);
      }
      return this.badVerb(verb, t, v);
    }
    if (k === "float32" || k === "float64") {
      if ("vbgGeEfFxX".includes(verb)) return fmtFloatVerb(v, k === "float32" ? 32 : 64, verb, f);
      return this.badVerb(verb, t, v);
    }
    if (k === "complex64" || k === "complex128") {
      const bits = k === "complex64" ? 32 : 64;
      const im = fmtFloatVerb(v.im, bits, verb, { ...f, plus: true, width: -1 });
      return pad("(" + fmtFloatVerb(v.re, bits, verb, { ...f, width: -1 }) + im + "i)", f);
    }
    if (k === "string") {
      if (verb === "v" && f.sharpV) return pad(quote(v), f);
      if ("vsqxX".includes(verb)) return fmtString(v, verb, f);
      return this.badVerb(verb, t, v);
    }
    switch (k) {
      case "iface":
        if (v === null) return pad(depth > 0 && !f.sharp ? "<nil>" : "<nil>", f);
        return this.value($.typeOf(v), v instanceof $.Val ? v.v : v, verb, f, depth + 1);
      case "struct":
        return this.struct(t, v, verb, f, depth);
      case "slice":
      case "array":
        return this.list(t, v, verb, f, depth);
      case "map":
        return this.map(t, v, verb, f, depth);
      case "ptr":
        if (v === null) return pad(f.sharpV && verb === "v" ? `(${t.str})(nil)` : "<nil>", f);
        if (depth === 0 && verb === "v" || depth === 0 && "sqxX".includes(verb)) {
          const e = t.elem!;
          if (e.kind === "struct" || e.kind === "array" || e.kind === "slice" || e.kind === "map") {
            const inner = e.kind === "struct" || e.kind === "array" ? v : v.v;
            return "&" + this.value(e, inner, verb, f, depth + 1);
          }
        }
        return pad(address(v), f);
      case "func":
      case "chan":
      case "unsafeptr":
        if (v === null && verb === "v") return pad(f.sharpV ? `(${t.str})(nil)` : "<nil>", f);
        return pad(address(v), f);
    }
    return this.badVerb(verb, t, v);
  }

  struct(t: $.Type, v: any, verb: string, f: Flags, depth: number): string {
    const parts: string[] = [];
    for (const fl of t.fields) {
      const name = f.plusV || f.sharpV ? fl.name + ":" : "";
      const inner = fl.exported || !"vsq".includes(verb) ? this.value(fl.type, v[fl.prop], verb, { ...f, width: -1 }, depth + 1)
        : this.plain(fl.type, v[fl.prop], verb, f, depth + 1);
      parts.push(name + inner);
    }
    if (f.sharpV && verb === "v") return t.str + "{" + parts.join(", ") + "}";
    return pad("{" + parts.join(" ") + "}", f);
  }

  // plain prints an unexported field, which fmt cannot call methods on.
  plain(t: $.Type, v: any, verb: string, f: Flags, depth: number): string {
    const saved = t.valCls;
    if (t.kind === "ptr") {
      if (v === null) return "<nil>";
      return address(v);
    }
    try {
      t.valCls = null;
      return this.value(t, v, verb, { ...f, width: -1 }, depth);
    } finally {
      t.valCls = saved;
    }
  }

  list(t: $.Type, v: any, verb: string, f: Flags, depth: number): string {
    const e = t.elem!;
    if (t.kind === "slice" && v === null && f.sharpV && verb === "v") return t.str + "(nil)";
    const items: any[] = t.kind === "slice" ? (v === null ? [] : v.toArray()) : Array.from(v);
    if (e.kind === "uint8" && "sqxX".includes(verb)) {
      return fmtString($.str.bytesToStr($.slice(Uint8Array.from(items))), verb, f);
    }
    const parts = items.map((x) => this.value(e, x, verb, { ...f, width: f.width }, depth + 1));
    if (f.sharpV && verb === "v") return t.str + "{" + parts.join(", ") + "}";
    return "[" + parts.join(" ") + "]";
  }

  map(t: $.Type, v: any, verb: string, f: Flags, depth: number): string {
    if (v === null && f.sharpV && verb === "v") return t.str + "(nil)";
    const entries: [any, any][] = v === null ? [] : Array.from(v as Iterable<[any, any]>);
    entries.sort((a, b) => compareKeys(t.key!, a[0], b[0]));
    const parts = entries.map(([k, x]) =>
      this.value(t.key!, k, verb, { ...f, width: -1 }, depth + 1) + ":" +
      this.value(t.elem!, x, verb, { ...f, width: -1 }, depth + 1)
    );
    if (f.sharpV && verb === "v") return t.str + "{" + parts.join(", ") + "}";
    return "map[" + parts.join(" ") + "]";
  }

  // printf formats like Sprintf.
  printf(format: string, args: any[]): string {
    let out = "";
    let argi = 0;
    let reordered = false;
    for (let i = 0; i < format.length;) {
      const pct = format.indexOf("%", i);
      if (pct < 0) {
        out += format.slice(i);
        break;
      }
      out += format.slice(i, pct);
      i = pct + 1;
      const f = noFlags();
      flags: for (; i < format.length; i++) {
        switch (format[i]) {
          case "+": f.plus = true; break;
          case "-": f.minus = true; f.zero = false; break;
          case "#": f.sharp = true; break;
          case " ": f.space = true; break;
          case "0": f.zero = !f.minus; break;
          default: break flags;
        }
      }
      const argIndex = (): void => {
        if (format[i] !== "[") return;
        const close = format.indexOf("]", i);
        const n = parseInt(format.slice(i + 1, close), 10);
        if (close > 0 && n > 0) {
          argi = n - 1;
          reordered = true;
          i = close + 1;
        }
      };
      const num = (): number => {
        argIndex();
        if (format[i] === "*") {
          i++;
          const a = args[argi++];
          const n = a instanceof $.Val && isInt(a.$t.kind) ? Number(a.v) : -1;
          return n;
        }
        let n = -1;
        while (i < format.length && format[i] >= "0" && format[i] <= "9") {
          n = (n < 0 ? 0 : n * 10) + (format.charCodeAt(i) - 48);
          i++;
        }
        return n;
      };
      f.width = num();
      if (f.width < -1) {
        f.minus = true;
        f.width = -f.width;
      }
      if (format[i] === ".") {
        i++;
        f.prec = num();
        if (f.prec < 0) f.prec = 0;
      }
      argIndex();
      if (format[i] === "v") {
        if (f.plus) f.plus = false, f.plusV = true;
        if (f.sharp) f.sharp = false, f.sharpV = true;
      }
      if (i >= format.length) {
        out += "%!(NOVERB)";
        break;
      }
      const [r, w] = $.str.decodeRune(format, i);
      const verb = w === 1 ? format[i] : $.str.runeStr(r);
      i += w;
      if (verb === "%") {
        out += "%";
        continue;
      }
      if (argi >= args.length) {
        out += `%!${verb}(MISSING)`;
        continue;
      }
      const a = args[argi++];
      out += this.arg(a, verb === "w" ? "v" : verb, f);
    }
    if (!reordered && argi < args.length) {
      out += "%!(EXTRA " + args.slice(argi).map((a) => a === null ? "<nil>" : $.typeOf(a).str + "=" + this.arg(a, "v", noFlags())).join(", ") + ")";
    }
    return out;
  }

  // print formats like Sprint, or Sprintln when ln is set.
  print(args: any[], ln: boolean): string {
    let out = "";
    for (let i = 0; i < args.length; i++) {
      const isStr = (a: any) => a instanceof $.Val && a.$t.kind === "string";
      if (i > 0 && (ln || (!isStr(args[i]) && !isStr(args[i - 1])))) out += " ";
      out += this.arg(args[i], "v", noFlags());
    }
    return ln ? out + "\n" : out;
  }
}

function compareKeys(t: $.Type, a: any, b: any): number {
  if (t.kind === "iface") {
    const ta = a === null ? "" : $.typeOf(a).str, tb = b === null ? "" : $.typeOf(b).str;
    if (ta !== tb) return ta < tb ? -1 : 1;
    return compareKeys(a === null ? t : $.typeOf(a), a instanceof $.Val ? a.v : a, b instanceof $.Val ? b.v : b);
  }
  if (t.kind === "struct") {
    for (const f of t.fields) {
      const c = compareKeys(f.type, a[f.prop], b[f.prop]);
      if (c !== 0) return c;
    }
    return 0;
  }
  if (t.kind === "bool") return a === b ? 0 : a ? 1 : -1;
  return a < b ? -1 : a > b ? 1 : 0;
}

function argsOf(s: $.Slice<any> | null): any[] {
  return s === null ? [] : s.toArray();
}

export function sprintf(format: string, args: $.Slice<any> | null): string {
  return new Printer().printf(format, argsOf(args));
}

export function sprint(args: $.Slice<any> | null, ln: boolean): string {
  return new Printer().print(argsOf(args), ln);
}
