// Type describes a Go type at run time: what fmt prints, what a type switch
// compares, and how two values of it are equal.
export type Kind =
  | "bool" | "int" | "int8" | "int16" | "int32" | "int64"
  | "uint" | "uint8" | "uint16" | "uint32" | "uint64" | "uintptr"
  | "float32" | "float64" | "complex64" | "complex128" | "string" | "unsafeptr"
  | "slice" | "array" | "map" | "ptr" | "func" | "iface" | "struct" | "chan" | "nil";

export type Field = { name: string; type: Type; embedded: boolean; exported: boolean };

export class Type {
  kind: Kind;
  str: string;
  named = false;
  elem: Type | null = null;
  key: Type | null = null;
  len = 0;
  // methods are the method names an interface type asks for.
  methods: string[] = [];
  // ctor is the class of a struct type's values.
  ctor: any = null;
  // valCls is the class an interface holds values of this type in, carrying
  // the type's methods.
  valCls: any = null;
  private fieldsOf: (() => Field[]) | null = null;
  private fieldList: Field[] | null = null;

  constructor(kind: Kind, str: string) {
    this.kind = kind;
    this.str = str;
  }

  get fields(): Field[] {
    if (this.fieldList === null) this.fieldList = this.fieldsOf === null ? [] : this.fieldsOf();
    return this.fieldList;
  }

  setFields(f: () => Field[]): void {
    this.fieldsOf = f;
    this.fieldList = null;
  }
}

const memo = new Map<string, Type>();

function memoized(key: string, make: () => Type): Type {
  let t = memo.get(key);
  if (t === undefined) {
    t = make();
    memo.set(key, t);
  }
  return t;
}

function basic(kind: Kind): Type {
  return new Type(kind, kind === "unsafeptr" ? "unsafe.Pointer" : kind);
}

export const Bool = basic("bool");
export const Int = basic("int");
export const Int8 = basic("int8");
export const Int16 = basic("int16");
export const Int32 = basic("int32");
export const Int64 = basic("int64");
export const Uint = basic("uint");
export const Uint8 = basic("uint8");
export const Uint16 = basic("uint16");
export const Uint32 = basic("uint32");
export const Uint64 = basic("uint64");
export const Uintptr = basic("uintptr");
export const Float32 = basic("float32");
export const Float64 = basic("float64");
export const Complex64 = basic("complex64");
export const Complex128 = basic("complex128");
export const String = basic("string");
export const UnsafePointer = basic("unsafeptr");

export function sliceOf(elem: Type): Type {
  return memoized("[]" + elem.str, () => {
    const t = new Type("slice", "[]" + elem.str);
    t.elem = elem;
    return t;
  });
}

export function arrayOf(elem: Type, n: number): Type {
  return memoized(`[${n}]${elem.str}`, () => {
    const t = new Type("array", `[${n}]${elem.str}`);
    t.elem = elem;
    t.len = n;
    return t;
  });
}

export function ptrTo(elem: Type): Type {
  return memoized("*" + elem.str, () => {
    const t = new Type("ptr", "*" + elem.str);
    t.elem = elem;
    return t;
  });
}

export function mapOf(key: Type, elem: Type): Type {
  return memoized(`map[${key.str}]${elem.str}`, () => {
    const t = new Type("map", `map[${key.str}]${elem.str}`);
    t.key = key;
    t.elem = elem;
    return t;
  });
}

export function chanOf(elem: Type, dir: string): Type {
  return memoized(dir + elem.str, () => {
    const t = new Type("chan", dir + elem.str);
    t.elem = elem;
    return t;
  });
}

export function funcType(str: string): Type {
  return memoized(str, () => new Type("func", str));
}

export function ifaceType(str: string, methods: string[]): Type {
  return memoized(str, () => {
    const t = new Type("iface", str);
    t.methods = methods;
    return t;
  });
}

// structType describes a struct, named or not. A named one is not memoized:
// two named types are different types even with the same name.
export function structType(str: string, named: boolean, ctor: any, valCls: any, fields: () => Field[]): Type {
  const make = () => {
    const t = new Type("struct", str);
    t.named = named;
    t.ctor = ctor;
    t.valCls = valCls;
    t.setFields(fields);
    return t;
  };
  return named ? make() : memoized(str, make);
}

// named is a defined type over a non-struct underlying type.
export function named(str: string, under: Type, valCls: any): Type {
  const t = new Type(under.kind, str);
  t.named = true;
  t.elem = under.elem;
  t.key = under.key;
  t.len = under.len;
  t.methods = under.methods;
  t.valCls = valCls;
  return t;
}

const canon = new Map<string, Type>();

// canonical is the one descriptor for the generic instance key, made by
// the first module to ask for it.
export function canonical(key: string, make: () => Type): Type {
  let t = canon.get(key);
  if (t === undefined) {
    t = make();
    canon.set(key, t);
  }
  return t;
}

export function field(name: string, type: Type, embedded = false): Field {
  const c = name.charCodeAt(0);
  return { name, type, embedded, exported: c >= 65 && c <= 90 };
}
