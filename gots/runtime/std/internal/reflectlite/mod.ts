// The parts of internal/reflectlite the standard library uses outside
// errors (which is hand-written): sort.Slice's swapping and length.
import * as $ from "../../../rt/mod.ts";

export function Swapper(slice: any): (i: number, j: number) => void {
  const s: $.Slice<any> | null = slice instanceof $.Val ? slice.v : null;
  if (s === null) return () => {};
  return (i: number, j: number) => {
    const a = s.a, x = s.o + i, y = s.o + j;
    const t = a[x];
    a[x] = a[y];
    a[y] = t;
  };
}

export class Value {
  v: any;
  constructor(v: any) {
    this.v = v;
  }
  Len(): number {
    const x = this.v instanceof $.Val ? this.v.v : this.v;
    if (x === null) return 0;
    if (x instanceof $.Slice) return x.n;
    if (typeof x === "string" || Array.isArray(x) || ArrayBuffer.isView(x)) return (x as any).length;
    if (x instanceof Map || x instanceof $.KeyMap) return x.size;
    return 0;
  }
  $clone(): Value {
    return this;
  }
}

export function ValueOf(i: any): Value {
  return new Value(i);
}

export class rtype {
  t: $.Type;
  constructor(t: $.Type) {
    this.t = t;
  }
  String(): string {
    return this.t.str;
  }
  // Comparable reports whether values of the type can be compared with ==,
  // as context.WithValue asks of its keys.
  Comparable(): boolean {
    return comparable(this.t);
  }
}

function comparable(t: $.Type): boolean {
  switch (t.kind) {
    case "slice":
    case "map":
    case "func":
      return false;
    case "array":
      return comparable(t.elem!);
    case "struct":
      return t.fields.every((f) => comparable(f.type));
  }
  return true;
}

export function TypeOf(i: any): rtype | null {
  return i === null ? null : new rtype($.typeOf(i));
}
