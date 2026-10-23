// Go's sync package. A transpiled program runs on one thread, so locks
// never wait and a WaitGroup has always finished when it is waited on.
import * as $ from "../../rt/mod.ts";

export interface Locker {
  Lock(): void;
  Unlock(): void;
}
export const Locker$type = $.ifaceType("sync.Locker", ["Lock", "Unlock"]);

export class Mutex extends $.Native {
  locked = false;
  Lock(): void {
    this.locked = true;
  }
  TryLock(): boolean {
    if (this.locked) return false;
    this.locked = true;
    return true;
  }
  Unlock(): void {
    if (!this.locked) $.panic($.val($.String, "sync: unlock of unlocked mutex"));
    this.locked = false;
  }
}
$.nativeType("sync.Mutex", Mutex);

export class RWMutex extends $.Native {
  writer = false;
  readers = 0;
  Lock(): void {
    this.writer = true;
  }
  TryLock(): boolean {
    if (this.writer || this.readers > 0) return false;
    this.writer = true;
    return true;
  }
  Unlock(): void {
    this.writer = false;
  }
  RLock(): void {
    this.readers++;
  }
  TryRLock(): boolean {
    if (this.writer) return false;
    this.readers++;
    return true;
  }
  RUnlock(): void {
    this.readers--;
  }
  RLocker(): Locker {
    return { Lock: () => this.RLock(), Unlock: () => this.RUnlock() };
  }
}
$.nativeType("sync.RWMutex", RWMutex);

export class Once extends $.Native {
  done = false;
  Do(f: () => void): void {
    if (this.done) return;
    this.done = true;
    f();
  }
}
$.nativeType("sync.Once", Once);

export class WaitGroup extends $.Native {
  n = 0;
  Add(delta: number): void {
    this.n += delta;
    if (this.n < 0) $.panic($.val($.String, "sync: negative WaitGroup counter"));
  }
  Done(): void {
    this.Add(-1);
  }
  Wait(): void {}
  // Go runs f and marks it done; without goroutines it runs now.
  Go(f: () => void): void {
    this.Add(1);
    try {
      f();
    } finally {
      this.Done();
    }
  }
}
$.nativeType("sync.WaitGroup", WaitGroup);

export class Pool extends $.Native {
  New: (() => any) | null = null;
  items: any[] = [];
  Get(): any {
    if (this.items.length > 0) return this.items.pop();
    return this.New === null ? null : this.New();
  }
  Put(x: any): void {
    if (x !== null) this.items.push(x);
  }
}
$.nativeType("sync.Pool", Pool);

export class Map extends $.Native {
  m = new $.KeyMap<any, any>($.keyOf);
  Load(key: any): [any, boolean] {
    return this.m.has(key) ? [this.m.get(key), true] : [null, false];
  }
  Store(key: any, value: any): void {
    this.m.set(key, value);
  }
  LoadOrStore(key: any, value: any): [any, boolean] {
    if (this.m.has(key)) return [this.m.get(key), true];
    this.m.set(key, value);
    return [value, false];
  }
  LoadAndDelete(key: any): [any, boolean] {
    const r = this.Load(key);
    this.m.delete(key);
    return r;
  }
  Delete(key: any): void {
    this.m.delete(key);
  }
  Swap(key: any, value: any): [any, boolean] {
    const r = this.Load(key);
    this.m.set(key, value);
    return r;
  }
  CompareAndSwap(key: any, old: any, value: any): boolean {
    if (!this.m.has(key) || !$.ifaceEq(this.m.get(key), old)) return false;
    this.m.set(key, value);
    return true;
  }
  CompareAndDelete(key: any, old: any): boolean {
    if (!this.m.has(key) || !$.ifaceEq(this.m.get(key), old)) return false;
    this.m.delete(key);
    return true;
  }
  Range(f: (key: any, value: any) => boolean): void {
    for (const [k, v] of Array.from(this.m)) if (!f(k, v)) return;
  }
  Clear(): void {
    this.m.clear();
  }
}
$.nativeType("sync.Map", Map);

export class Cond extends $.Native {
  L: Locker | null = null;
  Wait(): void {}
  Signal(): void {}
  Broadcast(): void {}
}
$.nativeType("sync.Cond", Cond);

export function NewCond(l: Locker | null): Cond {
  const c = new Cond();
  c.L = l;
  return c;
}

export function OnceFunc(f: () => void): () => void {
  let done = false;
  return () => {
    if (done) return;
    done = true;
    f();
  };
}

export function OnceValue<T>(f: () => T): () => T {
  let done = false;
  let v: T;
  return () => {
    if (!done) {
      done = true;
      v = f();
    }
    return v;
  };
}

export function OnceValues<T, U>(f: () => [T, U]): () => [T, U] {
  let done = false;
  let v: [T, U];
  return () => {
    if (!done) {
      done = true;
      v = f();
    }
    return v;
  };
}
