// fmt's scanning: Sscan, Sscanln and Sscanf, storing through the pointers
// the operands are, by the kind of what they point at.
import * as $ from "../../rt/mod.ts";

class scanError {
  s: string;
  constructor(s: string) {
    this.s = s;
  }
  Error(): string {
    return this.s;
  }
  static $type = $.structType("fmt.scanError", true, scanError, null, () => [$.field("s", $.String)]);
}

const eof = () => new scanError("unexpected EOF");

function isSpace(c: string): boolean {
  return c === " " || c === "\t" || c === "\r" || c === "\n" || c === "\v" || c === "\f";
}

// target is where an operand points: the kind stored there, and how.
type target = { kind: $.Kind; set: (v: any) => void };

function targetOf(arg: any): target | null {
  if (!(arg instanceof $.Val) || arg.$t.kind !== "ptr" || arg.v === null) return null;
  const p = arg.v;
  return { kind: arg.$t.elem!.kind, set: (v: any) => (p.v = v) };
}

const intWidths: Partial<Record<$.Kind, number>> = { int8: 8, int16: 16, int32: 32, uint8: 8, uint16: 16, uint32: 32 };

function storeInt(t: target, n: bigint): any {
  const w = intWidths[t.kind];
  if (t.kind === "uint64") {
    t.set(BigInt.asUintN(64, n));
    return null;
  }
  if (w !== undefined) {
    const lim = t.kind[0] === "u" ? 1n << BigInt(w) : 1n << BigInt(w - 1);
    if (t.kind[0] === "u" ? n < 0n || n >= lim : n < -lim || n >= lim) {
      return new scanError("integer overflow on token");
    }
  }
  t.set(Number(n));
  return null;
}

class reader {
  s: string;
  i = 0;
  constructor(s: string) {
    this.s = s;
  }
  done(): boolean {
    return this.i >= this.s.length;
  }
  peek(): string {
    return this.s[this.i];
  }
  // skipSpace skips blanks, and newlines too unless they end the input
  // for Sscanln.
  skipSpace(newlines: boolean): void {
    while (!this.done() && isSpace(this.peek()) && (newlines || this.peek() !== "\n")) this.i++;
  }
  token(): string {
    const start = this.i;
    while (!this.done() && !isSpace(this.peek())) this.i++;
    return this.s.slice(start, this.i);
  }
  match(re: RegExp): string | null {
    re.lastIndex = this.i;
    const m = re.exec(this.s);
    if (m === null || m.index !== this.i || m[0] === "") return null;
    this.i += m[0].length;
    return m[0];
  }
}

const intRe: Record<string, RegExp> = {
  d: /[+-]?[0-9_]+/y,
  x: /[+-]?(0[xX])?[0-9a-fA-F_]+/y,
  o: /[+-]?(0[oO])?[0-7_]+/y,
  b: /[+-]?(0[bB])?[01_]+/y,
  v: /[+-]?(0[xX][0-9a-fA-F_]+|0[bB][01_]+|0[oO]?[0-7_]+|[0-9_]+)/y,
};
const floatRe = /[+-]?(inf|infinity|nan|0[xX][0-9a-fA-F_.]+([pP][+-]?[0-9]+)?|[0-9_]*\.?[0-9_]*([eE][+-]?[0-9]+)?)/iy;

function parseInt(tok: string, verb: string): bigint | null {
  let s = tok.replaceAll("_", "");
  let neg = false;
  if (s[0] === "+" || s[0] === "-") {
    neg = s[0] === "-";
    s = s.slice(1);
  }
  let base = { d: 10, x: 16, X: 16, o: 8, b: 2 }[verb] ?? 10;
  if (verb === "v" || verb === "x" || verb === "X" || verb === "o" || verb === "b") {
    if (/^0[xX]/.test(s)) base = 16, s = s.slice(2);
    else if (/^0[bB]/.test(s)) base = 2, s = s.slice(2);
    else if (/^0[oO]/.test(s)) base = 8, s = s.slice(2);
    else if (verb === "v" && s.length > 1 && s[0] === "0") base = 8, s = s.slice(1);
  }
  if (s === "") return null;
  let n = 0n;
  for (const c of s) {
    const d = globalThis.parseInt(c, base);
    if (Number.isNaN(d)) return null;
    n = n * BigInt(base) + BigInt(d);
  }
  return neg ? -n : n;
}

// scanOne reads one operand's text from r by verb and stores it.
function scanOne(r: reader, verb: string, arg: any): any {
  const t = targetOf(arg);
  if (t === null) return new scanError("can't scan type");
  if (verb !== "c") r.skipSpace(false);
  if (r.done()) return eof();
  const k = t.kind;
  if (k === "string") {
    if (verb === "q") {
      const q = r.match(/"(\\.|[^"\\])*"|`[^`]*`/y);
      if (q === null) return new scanError("expected quoted string");
      t.set($.str.toGo(q[0] === "`" ? q.slice(1, -1) : JSON.parse(q)));
      return null;
    }
    t.set(r.token());
    return null;
  }
  if (k === "bool") {
    const tok = r.match(/true|false|TRUE|FALSE|True|False|[tTfF10]/y);
    if (tok === null) return new scanError("syntax error scanning boolean");
    t.set(/^[tT1]/.test(tok));
    return null;
  }
  if (k === "float32" || k === "float64") {
    const tok = r.match(floatRe);
    if (tok === null) return new scanError("expected float");
    const s = tok.replaceAll("_", "").toLowerCase();
    let f = s.endsWith("inf") || s.endsWith("infinity") ? (s[0] === "-" ? -Infinity : Infinity) : s.endsWith("nan") ? NaN : Number(s);
    if (Number.isNaN(f) && !s.endsWith("nan")) return new scanError("strconv.ParseFloat: parsing " + JSON.stringify(tok) + ": invalid syntax");
    if (k === "float32") f = Math.fround(f);
    t.set(f);
    return null;
  }
  if (verb === "c") {
    const [ch, w] = $.str.decodeRune(r.s, r.i);
    r.i += w;
    return storeInt(t, BigInt(ch));
  }
  const tok = r.match(intRe[verb.toLowerCase()] ?? intRe.v);
  if (tok === null) return new scanError("expected integer");
  const n = parseInt(tok, verb);
  if (n === null) return new scanError("expected integer");
  return storeInt(t, n);
}

function argsOf(a: $.Slice<any> | null): any[] {
  return a === null ? [] : a.toArray();
}

export function Sscanf(str: string, format: string, a: $.Slice<any> | null): [number, any] {
  const args = argsOf(a);
  const r = new reader(str);
  let n = 0;
  for (let i = 0; i < format.length;) {
    const c = format[i];
    if (isSpace(c)) {
      const nl = c === "\n";
      while (i < format.length && isSpace(format[i])) i++;
      r.skipSpace(nl);
      continue;
    }
    if (c !== "%" || format[i + 1] === "%") {
      if (c === "%") i++;
      if (r.done()) return [n, eof()];
      if (r.peek() !== format[i]) return [n, new scanError("input does not match format")];
      r.i++;
      i++;
      continue;
    }
    i++;
    while (i < format.length && "+-# 0123456789".includes(format[i])) i++;
    const verb = format[i++];
    if (n >= args.length) return [n, new scanError("too few operands for format '%" + verb + "'")];
    const err = scanOne(r, verb, args[n]);
    if (err !== null) return [n, err];
    n++;
  }
  if (n < args.length) return [n, new scanError("too many operands")];
  return [n, null];
}

function scanAll(str: string, a: $.Slice<any> | null, ln: boolean): [number, any] {
  const args = argsOf(a);
  const r = new reader(str);
  let n = 0;
  for (const arg of args) {
    r.skipSpace(!ln);
    if (r.done() || (ln && r.peek() === "\n")) return [n, ln && !r.done() ? new scanError("unexpected newline") : eof()];
    const err = scanOne(r, "v", arg);
    if (err !== null) return [n, err];
    n++;
  }
  if (ln) {
    r.skipSpace(false);
    if (!r.done() && r.peek() !== "\n") return [n, new scanError("expected newline")];
  }
  return [n, null];
}

export function Sscan(str: string, a: $.Slice<any> | null): [number, any] {
  return scanAll(str, a, false);
}

export function Sscanln(str: string, a: $.Slice<any> | null): [number, any] {
  return scanAll(str, a, true);
}

// readAll drains an io.Reader into a Go string, for the Fscan family.
function readAll(r: any): string {
  let s = "";
  const buf = $.make<number>(Uint8Array, 4096, 4096);
  for (;;) {
    const [n, err] = r.Read(buf);
    s += $.str.bytesToStr($.sub(buf, 0, n));
    if (err !== null || n === 0) return s;
  }
}

export function Fscan(r: any, a: $.Slice<any> | null): [number, any] {
  return Sscan(readAll(r), a);
}

export function Fscanln(r: any, a: $.Slice<any> | null): [number, any] {
  return Sscanln(readAll(r), a);
}

export function Fscanf(r: any, format: string, a: $.Slice<any> | null): [number, any] {
  return Sscanf(readAll(r), format, a);
}
