// Go's fmt package. Operands arrive as interface values, which carry the
// type descriptors the verbs are chosen by.
import * as $ from "../../rt/mod.ts";
import { sprint, sprintf } from "./print.ts";
import { write } from "./output.ts";

export interface Stringer {
  String(): string;
}

// The descriptors of fmt's interfaces, for type switches, assertions and
// reflect on them.
const stringFunc = () => [$.funcType("func() string", () => [], () => [$.String])];
export const Stringer$type = $.ifaceType("fmt.Stringer", ["String"], stringFunc);
export const GoStringer$type = $.ifaceType("fmt.GoStringer", ["GoString"], stringFunc);
export const Formatter$type = $.ifaceType("fmt.Formatter", ["Format"]);
export const State$type = $.ifaceType("fmt.State", ["Flag", "Precision", "Width", "Write"]);

export function Sprintf(format: string, a: $.Slice<any> | null): string {
  return sprintf(format, a);
}

export function Sprint(a: $.Slice<any> | null): string {
  return sprint(a, false);
}

export function Sprintln(a: $.Slice<any> | null): string {
  return sprint(a, true);
}

export function Printf(format: string, a: $.Slice<any> | null): [number, any] {
  const s = sprintf(format, a);
  write(s);
  return [s.length, null];
}

export function Print(a: $.Slice<any> | null): [number, any] {
  const s = sprint(a, false);
  write(s);
  return [s.length, null];
}

export function Println(a: $.Slice<any> | null): [number, any] {
  const s = sprint(a, true);
  write(s);
  return [s.length, null];
}

// toWriter hands s to an io.Writer as a []byte.
function toWriter(w: any, s: string): [number, any] {
  return $.noWait(w.Write($.slice($.str.strToBytes(s))), "Write");
}

export function Fprintf(w: any, format: string, a: $.Slice<any> | null): [number, any] {
  return toWriter(w, sprintf(format, a));
}

export function Fprint(w: any, a: $.Slice<any> | null): [number, any] {
  return toWriter(w, sprint(a, false));
}

export function Fprintln(w: any, a: $.Slice<any> | null): [number, any] {
  return toWriter(w, sprint(a, true));
}

// The twins of the Fprint family, for code that can wait: they wait for
// a writer whose Write blocks, a socket's.
async function toWriterA(w: any, s: string): Promise<[number, any]> {
  return await $.callA(w, "Write", $.slice($.str.strToBytes(s)));
}

export function Fprintf$a(w: any, format: string, a: $.Slice<any> | null): Promise<[number, any]> {
  return toWriterA(w, sprintf(format, a));
}

export function Fprint$a(w: any, a: $.Slice<any> | null): Promise<[number, any]> {
  return toWriterA(w, sprint(a, false));
}

export function Fprintln$a(w: any, a: $.Slice<any> | null): Promise<[number, any]> {
  return toWriterA(w, sprint(a, true));
}

export function Append(b: $.Slice<number> | null, a: $.Slice<any> | null): $.Slice<number> | null {
  return $.appendString(b, sprint(a, false));
}

export function Appendln(b: $.Slice<number> | null, a: $.Slice<any> | null): $.Slice<number> | null {
  return $.appendString(b, sprint(a, true));
}

export function Appendf(b: $.Slice<number> | null, format: string, a: $.Slice<any> | null): $.Slice<number> | null {
  return $.appendString(b, sprintf(format, a));
}

class wrapError {
  msg: string;
  err: any;
  constructor(msg: string, err: any) {
    this.msg = msg;
    this.err = err;
  }
  Error(): string {
    return this.msg;
  }
  Unwrap(): any {
    return this.err;
  }
  static $type = $.structType("fmt.wrapError", true, wrapError, null, () => [$.field("msg", $.String), $.field("err", $.ifaceType("error", ["Error"]))]);
}

class wrapErrors {
  msg: string;
  errs: $.Slice<any>;
  constructor(msg: string, errs: $.Slice<any>) {
    this.msg = msg;
    this.errs = errs;
  }
  Error(): string {
    return this.msg;
  }
  Unwrap(): $.Slice<any> {
    return this.errs;
  }
  static $type = $.structType("fmt.wrapErrors", true, wrapErrors, null, () => [$.field("msg", $.String)]);
}

class plainError {
  s: string;
  constructor(s: string) {
    this.s = s;
  }
  Error(): string {
    return this.s;
  }
  static $type = $.structType("errors.errorString", true, plainError, null, () => [$.field("s", $.String)]);
}

// Errorf is Sprintf as an error, wrapping the operands of its %w verbs.
export function Errorf(format: string, a: $.Slice<any> | null): any {
  const s = sprintf(format, a);
  const args = a === null ? [] : a.toArray();
  const wrapped: any[] = [];
  let argi = 0;
  for (let i = 0; i < format.length; i++) {
    if (format[i] !== "%") continue;
    i++;
    while (i < format.length && "+-# 0123456789.*[]".includes(format[i])) i++;
    if (format[i] === "%") continue;
    if (format[i] === "w" && typeof args[argi]?.Error === "function") wrapped.push(args[argi]);
    argi++;
  }
  if (wrapped.length === 1) return new wrapError(s, wrapped[0]);
  if (wrapped.length > 1) return new wrapErrors(s, $.slice(wrapped));
  return new plainError(s);
}
export { Fscan, Fscanf, Fscanln, Sscan, Sscanf, Sscanln } from "./scan.ts";
