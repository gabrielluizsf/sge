// Calling through reflect: method sets as reflect sees them, finding a
// method's Go implementation for a receiver, and calling a func Value with
// Values.
import * as $ from "../../rt/mod.ts";
import { box, copyOf, makeSlice } from "./zero.ts";
import { toType } from "./rtype.ts";

function fail(msg: string): never {
  $.panic($.val($.String, $.str.toGo(msg)));
}

// methodsOf is the exported method set of type d, sorted by name: a named
// type's own, or all of them for a pointer to one.
export function methodsOf(d: $.Type): $.Method[] {
  if (d.kind === "ptr" && d.elem!.named) return d.elem!.methodSet;
  if (d.named) return d.methodSet.filter((m) => m.value);
  return [];
}

// bound is the method name of recv, a value of type d, as a function of
// the method's own parameters. Go code's names are tried first: an
// exported method taking or returning strings also has a TypeScript-facing
// one converting them.
export function bound(d: $.Type, recv: any, name: string): (...a: any[]) => any {
  if (d.kind === "iface") {
    if (recv === null) fail("reflect: Method on nil interface value");
    const dyn = $.typeOf(recv);
    return bound(dyn, recv instanceof $.Val ? recv.v : recv, name);
  }
  if (d.kind === "struct" || d.kind === "ptr" && d.elem!.kind === "struct") {
    if (recv === null) fail("reflect: call of method " + name + " on nil pointer");
    return withTwin(recv, (recv[name + "$g"] ?? recv[name]).bind(recv), name);
  }
  // Methods of other named types are functions the Val class holding
  // their values calls.
  const cls = d.valCls;
  if (cls === null || cls === undefined) fail("gots: reflect cannot find method " + name + " of " + d.str);
  const b = new cls(d, recv);
  return withTwin(b, b[name].bind(b), name);
}

// withTwin is f, method name of r bound to it, carrying the method's twin
// when it has one, for Call from code that can wait.
function withTwin(r: any, f: (...a: any[]) => any, name: string): (...a: any[]) => any {
  const a = r[name + "$g$a"] ?? r[name + "$a"];
  return typeof a === "function" ? $.twins(f, a.bind(r)) : f;
}

// argument is the Value v passed as a parameter of type t, which it has to
// be assignable to.
function argument(t: $.Type, v: any): any {
  if (!toType(v.t)!.AssignableTo(toType(t))) fail("reflect: Call using " + v.t.str + " as type " + t.str);
  const x = v.get();
  return t.kind === "iface" && v.t.kind !== "iface" ? box(v.t, copyOf(v.t, x)) : copyOf(t, x);
}

// args are the arguments a call of fn, of the func type t, takes from the
// Values in, packing the extra ones of a variadic function into its last
// parameter unless asSlice says in has the slice already.
function args(t: $.Type, fn: any, ins: any[], asSlice: boolean): any[] {
  if (fn === null) fail("reflect: call of nil function");
  const params = t.ins;
  const n = params.length;
  const fixed = t.variadic ? n - 1 : n;
  if (asSlice ? ins.length !== n : t.variadic ? ins.length < fixed : ins.length !== n) {
    fail("reflect: Call with " + (ins.length < n ? "too few" : "too many") + " input arguments");
  }
  const out: any[] = [];
  for (let i = 0; i < fixed; i++) out.push(argument(params[i], ins[i]));
  if (t.variadic) {
    const st = params[n - 1];
    if (asSlice) out.push(ins[n - 1].get());
    else {
      const rest = ins.slice(fixed);
      const s = makeSlice(st, rest.length, rest.length);
      rest.forEach((v, i) => (s.a[i] = argument(st.elem!, v)));
      out.push(rest.length === 0 ? null : s);
    }
  }
  return out;
}

// results are what a call of the func type t gave, as values of their
// types.
function results(t: $.Type, r: any): any[] {
  const outs = t.outs;
  if (outs.length === 0) return [];
  if (outs.length === 1) return [r];
  return r;
}

// call calls fn, of the func type t, with the Values in. A function that
// blocks is called by callA, from code that can wait.
export function call(t: $.Type, fn: any, ins: any[], asSlice: boolean): any[] {
  const r = fn(...args(t, fn, ins, asSlice));
  if (r instanceof Promise) {
    r.catch(() => {});
    fail("gots: reflect.Value.Call of a function that blocks, from code that does not wait; call it from a goroutine");
  }
  return results(t, r);
}

// callA is call for code that can wait: it calls fn's twin, when it has
// one, and waits for a function that blocks.
export async function callA(t: $.Type, fn: any, ins: any[], asSlice: boolean): Promise<any[]> {
  const a = args(t, fn, ins, asSlice);
  return results(t, await $.fa(fn)(...a));
}
