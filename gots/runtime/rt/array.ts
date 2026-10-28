// Go arrays are values: these copy and compare them element by element.

// cloneArray copies a; elemClone copies an element that is itself a value
// with storage (a struct or an array).
export function cloneArray(a: any, elemClone?: (e: any) => any): any {
  if (elemClone === undefined) return a.slice();
  const out = new Array(a.length);
  for (let i = 0; i < a.length; i++) out[i] = elemClone(a[i]);
  return out;
}

// setArray is dst = src, in place so pointers to dst stay valid.
export function setArray(dst: any, src: any, elemSet?: (d: any, s: any) => void): any {
  if (elemSet === undefined) {
    if (Array.isArray(dst)) for (let i = 0; i < src.length; i++) dst[i] = src[i];
    else dst.set(src);
  } else {
    for (let i = 0; i < src.length; i++) elemSet(dst[i], src[i]);
  }
  return dst;
}

export function fillArray(n: number, zero: () => any): any[] {
  const a = new Array(n);
  for (let i = 0; i < n; i++) a[i] = zero();
  return a;
}

// arrayPtr is (*[n]T)(s): the n elements of s as an array that shares
// its storage, which a typed array's subarray does and a proxy does for a
// plain array.
export function arrayPtr(s: any, n: number): any {
  if (s === null) {
    if (n === 0) return null;
    throw new RangeError(`cannot convert slice with length 0 to array or pointer to array with length ${n}`);
  }
  if (s.n < n) throw new RangeError(`cannot convert slice with length ${s.n} to array or pointer to array with length ${n}`);
  if (!Array.isArray(s.a)) return s.a.subarray(s.o, s.o + n);
  if (s.o === 0 && s.a.length === n) return s.a;
  const a = s.a, o = s.o;
  return new Proxy([], {
    get(_t, k) {
      if (k === "length") return n;
      if (typeof k === "string" && /^\d+$/.test(k)) return a[o + Number(k)];
      return (Array.prototype as any)[k];
    },
    set(_t, k, v) {
      if (typeof k === "string" && /^\d+$/.test(k)) a[o + Number(k)] = v;
      return true;
    },
    has(_t, k) {
      if (typeof k === "string" && /^\d+$/.test(k)) return Number(k) < n;
      return k in Array.prototype;
    },
  });
}
