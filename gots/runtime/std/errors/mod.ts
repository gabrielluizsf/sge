// Go's errors package.
import * as $ from "../../rt/mod.ts";

export interface error {
  Error(): string;
}

class errorString {
  s: string;
  constructor(s: string) {
    this.s = s;
  }
  Error(): string {
    return this.s;
  }
  static $type = $.structType("errors.errorString", true, errorString, null, () => [$.field("s", $.String)]);
}

export function New(text: string): error {
  return new errorString(text);
}

// unwrapAll is what err wraps: one error, several, or none.
function unwrapAll(err: any): any[] {
  if (typeof err?.Unwrap !== "function") return [];
  const u = err.Unwrap();
  if (u instanceof $.Slice) return u.toArray();
  return u === null ? [] : [u];
}

export function Unwrap(err: any): any {
  if (typeof err?.Unwrap !== "function") return null;
  const u = err.Unwrap();
  return u instanceof $.Slice ? null : u;
}

export function Is(err: any, target: any): boolean {
  if (err === null || target === null) return err === target;
  for (const todo = [err]; todo.length > 0;) {
    const e = todo.pop();
    if ($.ifaceEq(e, target)) return true;
    if (typeof e?.Is === "function" && e.Is(target)) return true;
    todo.push(...unwrapAll(e));
  }
  return false;
}

// As finds the first error in err's tree whose type is the one target
// points at, and stores it there. target is any pointer, as in Go.
export function As(err: any, target: any): boolean {
  if (target === null) $.panic(New("errors: target cannot be nil"));
  // A pointer to a struct is the struct object itself; anything else is a
  // Val holding a Ptr.
  const ptrToStruct = !(target instanceof $.Val);
  const want: $.Type = ptrToStruct ? target.constructor.$type : target.$t.elem;
  for (const todo = [err]; todo.length > 0;) {
    const e = todo.shift();
    if (e === null) continue;
    const ok = want.kind === "iface" ? $.implementsI(e, want.methods) : $.typeOf(e) === want;
    if (ok) {
      if (ptrToStruct) target.$set(e.v);
      else target.v.v = e instanceof $.Val && want.kind !== "iface" ? e.v : e;
      return true;
    }
    if (typeof e?.As === "function" && e.As(target)) return true;
    todo.push(...unwrapAll(e));
  }
  return false;
}

class joinError {
  errs: $.Slice<any>;
  constructor(errs: $.Slice<any>) {
    this.errs = errs;
  }
  Error(): string {
    return this.errs.toArray().map((e) => $.str.toGo(e.Error())).join("\n");
  }
  Unwrap(): $.Slice<any> {
    return this.errs;
  }
  static $type = $.structType("errors.joinError", true, joinError, null, () => []);
}

export function Join(errs: $.Slice<any> | null): error | null {
  const kept = errs === null ? [] : errs.toArray().filter((e) => e !== null);
  return kept.length === 0 ? null : new joinError($.slice(kept));
}
