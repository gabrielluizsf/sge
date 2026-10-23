// Maps with keys JavaScript compares by value (numbers, strings, booleans,
// bigints, pointers) are JavaScript Maps. KeyMap is for keys Go compares by
// contents — structs, arrays, interfaces — and keys them by a string made of
// those contents.
export class KeyMap<K, V> {
  private m = new Map<string, [K, V]>();
  private key: (k: K) => string;

  constructor(key: (k: K) => string) {
    this.key = key;
  }

  get size(): number {
    return this.m.size;
  }

  get(k: K): V | undefined {
    return this.m.get(this.key(k))?.[1];
  }

  has(k: K): boolean {
    return this.m.has(this.key(k));
  }

  set(k: K, v: V): this {
    this.m.set(this.key(k), [k, v]);
    return this;
  }

  delete(k: K): boolean {
    return this.m.delete(this.key(k));
  }

  clear(): void {
    this.m.clear();
  }

  *[Symbol.iterator](): IterableIterator<[K, V]> {
    for (const e of this.m.values()) yield e;
  }
}

// GoMap is either kind of map.
export type GoMap<K, V> = Map<K, V> | KeyMap<K, V>;

export function mapLen(m: GoMap<any, any> | null): number {
  return m === null ? 0 : m.size;
}

// mapSet is m[k] = v, panicking on a nil map as Go does.
export function mapSet<K, V>(m: GoMap<K, V> | null, k: K, v: V): void {
  if (m === null) throw new TypeError("assignment to entry in nil map");
  m.set(k, v);
}

export function mapDelete<K>(m: GoMap<K, any> | null, k: K): void {
  if (m !== null) m.delete(k);
}

// keyOf turns a value Go compares by contents into a string with the same
// equality.
export function keyOf(v: any): string {
  if (v === null || v === undefined) return "n";
  switch (typeof v) {
    case "number":
      return Object.is(v, -0) ? "d0" : "d" + v;
    case "string":
      return "s" + v.length + ":" + v;
    case "boolean":
      return v ? "t" : "f";
    case "bigint":
      return "b" + v;
  }
  if (v.$t !== undefined && "v" in v) return "i" + v.$t.str + "(" + keyOf(v.v) + ")";
  if (Array.isArray(v) || ArrayBuffer.isView(v)) {
    let s = "[";
    for (const e of v as any) s += keyOf(e) + ",";
    return s + "]";
  }
  if (typeof v.$key === "function") return v.$key();
  return "p" + identity(v);
}

const ids = new WeakMap<object, number>();
let nextId = 1;

function identity(o: object): number {
  let id = ids.get(o);
  if (id === undefined) {
    id = nextId++;
    ids.set(o, id);
  }
  return id;
}

// keyMapOf is a map literal with keys compared by contents.
export function keyMapOf<K, V>(entries: [K, V][]): KeyMap<K, V> {
  const m = new KeyMap<K, V>(keyOf);
  for (const [k, v] of entries) m.set(k, v);
  return m;
}

// entries is what a range over m walks: a snapshot of its keys, read
// back as the loop reaches them, so an entry deleted before it is reached
// is not produced and one written during the loop is not produced twice.
export function* entries<K, V>(m: GoMap<K, V> | null): IterableIterator<[K, V]> {
  if (m === null || m.size === 0) return;
  const keys: K[] = [];
  for (const [k] of m as Iterable<[K, V]>) keys.push(k);
  for (const k of keys) {
    if (m.has(k)) yield [k, m.get(k) as V];
  }
}
