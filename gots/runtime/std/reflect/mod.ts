// Go's reflect, over the transpiler's run-time type descriptors: Type is
// a descriptor, Value a descriptor and a value, and the place it lives in
// when it is addressable. Calling methods and functions through reflect is
// not supported yet.
import * as $ from "../../rt/mod.ts";

export { Kind$String, Kind$type, Kind$Val } from "./kinds.ts";
export type { Kind } from "./kinds.ts";
export { Method, StructField, StructTag$Get, StructTag$Lookup, StructTag$type, StructTag$Val, rtype } from "./rtype.ts";
export type { StructTag } from "./rtype.ts";
export {
  Append, AppendSlice, ArrayOf, Copy, Indirect, MakeChan, MakeMap, MakeMapWithSize, MakeSlice, MapIter, MapOf, New,
  PointerTo, PtrTo, SliceOf, Swapper, TypeAssert, TypeFor, TypeOf, Value, ValueError, ValueOf, Zero,
} from "./value.ts";

export type Type = any;
export const Type$type = $.ifaceType("reflect.Type", [
  "Align", "AssignableTo", "Bits", "ChanDir", "Comparable", "ConvertibleTo", "Elem", "Field", "FieldAlign", "FieldByName",
  "Implements", "In", "IsVariadic", "Key", "Kind", "Len", "Method", "MethodByName", "Name", "NumField", "NumIn",
  "NumMethod", "NumOut", "Out", "PkgPath", "Size", "String",
]);

export const DeepEqual = deepEqual;

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
      for (const f of t.fields) if (!deep(f.type, a[f.prop], b[f.prop])) return false;
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
