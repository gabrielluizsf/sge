// Go's reflect package is not supported by gots yet. What is here lets
// packages that import it load; using it panics.
import * as $ from "../../rt/mod.ts";

function unsupported(): never {
  return $.panic($.val($.String, "gots: package reflect is not supported yet"));
}

export const ValueOf = unsupported, TypeOf = unsupported, Indirect = unsupported, New = unsupported;
export const MakeSlice = unsupported, MakeMap = unsupported, Zero = unsupported, DeepEqual = deepEqual;
export const Copy = unsupported, Append = unsupported, Swapper = unsupported, PointerTo = unsupported;

export class Value extends $.Native {}
$.nativeType("reflect.Value", Value);

export type Kind = number;
export const Kind$type = $.named("reflect.Kind", $.Uint, null);

// deepEqual is reflect.DeepEqual for the values gots represents: what
// tests reach for most.
function deepEqual(a: any, b: any): boolean {
  if (a === b) return true;
  if (a === null || b === null) return false;
  const ta = $.typeOf(a), tb = $.typeOf(b);
  if (ta !== tb) return false;
  return deep(ta, a instanceof $.Val ? a.v : a, b instanceof $.Val ? b.v : b);
}

function deep(t: $.Type, a: any, b: any): boolean {
  if (a === b) return true;
  switch (t.kind) {
    case "slice": {
      if (a === null || b === null) return a === b;
      if (a.n !== b.n) return false;
      for (let i = 0; i < a.n; i++) if (!deep(t.elem!, a.a[a.o + i], b.a[b.o + i])) return false;
      return true;
    }
    case "array":
      for (let i = 0; i < t.len; i++) if (!deep(t.elem!, a[i], b[i])) return false;
      return true;
    case "struct":
      for (const f of t.fields) if (!deep(f.type, a[f.name], b[f.name])) return false;
      return true;
    case "map": {
      if (a === null || b === null) return a === b;
      if (a.size !== b.size) return false;
      for (const [k, v] of a) {
        if (!b.has(k) || !deep(t.elem!, v, b.get(k))) return false;
      }
      return true;
    }
    case "ptr":
      if (a === null || b === null) return false;
      return t.elem!.kind === "struct" || t.elem!.kind === "array" ? deep(t.elem!, a, b) : deep(t.elem!, a.v, b.v);
    case "iface":
      return deepEqual(a, b);
    case "func":
      return false;
  }
  return $.equal(t, a, b);
}
