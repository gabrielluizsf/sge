// Twins: a function that blocks only through what it calls by interface
// or function value is written twice, the function itself, which never
// waits, and name$a, which awaits what it calls. Code that can wait calls
// the twin; these find it on objects and function values.
import { panic } from "./panic.ts";
import { val } from "./iface.ts";
import { String as GoString } from "./types.ts";
import { toGo } from "./string.ts";

// callA calls method m of r through its twin when r has one: the name
// gots gives the Go-facing method, or the TypeScript-facing one, with $a.
export function callA(r: any, m: string, ...a: any[]): any {
  const f = r[m + "$g$a"] ?? r[m + "$a"] ?? r[m];
  return f.apply(r, a);
}

// noWait is the result of a call through an interface made where gots
// wrote code that cannot wait: a promise there means a method that blocks
// was reached from code that was not expected to, which goes on no
// further than Go would with a deadlock.
export function noWait<T>(r: T, m: string): T {
  if (r instanceof Promise) {
    r.catch(() => {});
    panic(val(GoString, toGo("gots: " + m + " blocks, but was called where gots wrote code that does not wait for it; call it from a goroutine")));
  }
  return r;
}

// fa is the function value f as code that can wait calls it: its twin
// when it has one.
export function fa(f: any): any {
  return f?.$a ?? f;
}

// twins is f carrying its twin a.
export function twins<F extends (...a: any[]) => any>(f: F, a: (...a: any[]) => any): F {
  (f as any).$a = a;
  return f;
}
