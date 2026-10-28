// Ptr is a pointer to anything but a struct or an array, which are objects
// already and are pointed at by reference. Variables whose address is taken
// are kept in a Cell; fields and elements are reached through accessors.
export interface Ptr<T> {
  v: T;
}

export class Cell<T> implements Ptr<T> {
  v: T;
  constructor(v: T) {
    this.v = v;
  }
}

export class FieldPtr<T> implements Ptr<T> {
  o: any;
  k: string;
  constructor(o: any, k: string) {
    if (o === null) throw new TypeError("invalid memory address or nil pointer dereference");
    this.o = o;
    this.k = k;
  }
  get v(): T {
    return this.o[this.k];
  }
  set v(x: T) {
    this.o[this.k] = x;
  }
}

export class IndexPtr<T> implements Ptr<T> {
  a: any;
  i: number;
  constructor(a: any, i: number) {
    this.a = a;
    this.i = i;
  }
  get v(): T {
    return this.a[this.i];
  }
  set v(x: T) {
    this.a[this.i] = x;
  }
}

// ptrEq compares pointers that may be accessors made at different times for
// the same place.
export function ptrEq(p: any, q: any): boolean {
  if (p === q) return true;
  if (p === null || q === null) return false;
  if (p instanceof FieldPtr && q instanceof FieldPtr) return p.o === q.o && p.k === q.k;
  if (p instanceof IndexPtr && q instanceof IndexPtr) return p.a === q.a && p.i === q.i;
  return false;
}

// nilCheck is p, or the panic Go raises for dereferencing nil.
export function nilCheck<T>(p: T | null): T {
  if (p === null) throw new TypeError("invalid memory address or nil pointer dereference");
  return p;
}

// elemPtr is &s[i] for an element that is not an object.
export function elemPtr<T>(s: any, i: number): IndexPtr<T> {
  if (s === null || i < 0 || i >= s.n) throw new RangeError(`index out of range [${i}]`);
  return new IndexPtr<T>(s.a, s.o + i);
}

// VarPtr is a pointer to a package-level variable, reached through the
// module's getter and setter for it.
export class VarPtr<T> implements Ptr<T> {
  private getter: () => T;
  private setter: (v: T) => void;
  constructor(get: () => T, set: (v: T) => void) {
    this.getter = get;
    this.setter = set;
  }
  get v(): T {
    return this.getter();
  }
  set v(x: T) {
    this.setter(x);
  }
}
