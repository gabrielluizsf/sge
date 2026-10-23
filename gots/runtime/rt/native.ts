import { Type, structType } from "./types.ts";

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

// nativeType gives a hand-written struct class its descriptor.
export function nativeType(name: string, ctor: any, valCls: any = null): Type {
  const t = structType(name, true, ctor, valCls, () => []);
  ctor.$type = t;
  return t;
}
