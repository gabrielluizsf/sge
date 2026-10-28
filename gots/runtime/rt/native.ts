import { Type, structType } from "./types.ts";
import { Val } from "./iface.ts";
import { RuntimeError } from "./panic.ts";

// Native is the base of hand-written standard library structs, giving them
// the value methods transpiled code calls on every struct: a copy is a
// fresh value, as for the locks and counters these are.
export class Native {
  static $type: Type;
  $clone(): any {
    const c = new (this.constructor as any)();
    Object.assign(c, this);
    return c;
  }
  $set(o: any): this {
    Object.assign(this, o);
    return this;
  }
  $eq(o: any): boolean {
    return this === o;
  }
  $key(): string {
    return "native";
  }
}

// nativeType gives a hand-written struct class its descriptor, and the
// class an interface holds its values in, forwarding its methods.
export function nativeType(name: string, ctor: any, valCls: any = null): Type {
  const t = structType(name, true, ctor, valCls ?? valClassOf(ctor), () => []);
  ctor.$type = t;
  return t;
}

function valClassOf(ctor: any): any {
  class V extends Val {}
  for (let p = ctor.prototype; p !== null && p !== Object.prototype && p !== Native.prototype; p = Object.getPrototypeOf(p)) {
    for (const m of Object.getOwnPropertyNames(p)) {
      if (m === "constructor" || m.startsWith("$") || m in V.prototype) continue;
      if (typeof Object.getOwnPropertyDescriptor(p, m)?.value !== "function") continue;
      (V.prototype as any)[m] = function (this: Val, ...a: any[]) {
        return this.v[m](...a);
      };
    }
  }
  return V;
}

// A fault recovered from JavaScript is Go's runtime.Error, and an error
// fmt and interfaces can take like any other.
nativeType("runtime.Error", RuntimeError);
