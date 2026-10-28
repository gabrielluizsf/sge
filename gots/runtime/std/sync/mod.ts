// Go's sync package. Goroutines take turns on JavaScript's one thread and
// only change hands where one waits, so a lock is free whenever nobody
// holds it across a wait. Code that can wait locks with Lock$a, which
// waits for a lock another goroutine holds, giving back a promise only
// then; code that cannot, say a frame being drawn, locks with Lock, which
// refuses, loudly, a lock it would have to wait for. WaitGroup.Wait and
// Cond.Wait wait, and Go code awaits them.
import * as $ from "../../rt/mod.ts";

export interface Locker {
  Lock(): void;
  Unlock(): void;
}
export const Locker$type = $.ifaceType("sync.Locker", ["Lock", "Unlock"]);

function contended(): never {
  $.panic($.val($.String, $.str.toGo("gots: Lock, in code that cannot wait, of a lock another goroutine holds")));
}

export class Mutex extends $.Native {
  locked = false;
  waiters: Array<() => void> = [];
  Lock(): void {
    if (this.locked) contended();
    this.locked = true;
  }
  // Lock$a locks m, waiting its turn, as Go's does, if it is held.
  Lock$a(): Promise<void> | undefined {
    if (!this.locked) {
      this.locked = true;
      return undefined;
    }
    return new Promise<void>((resolve) => this.waiters.push(resolve));
  }
  TryLock(): boolean {
    if (this.locked) return false;
    this.locked = true;
    return true;
  }
  // Unlock hands the lock to the goroutine waiting longest, if any.
  Unlock(): void {
    if (!this.locked) $.panic($.val($.String, "sync: unlock of unlocked mutex"));
    const next = this.waiters.shift();
    if (next) next();
    else this.locked = false;
  }
}
$.nativeType("sync.Mutex", Mutex);

// RWMutex's waiters take it in turn, a writer alone, readers together; a
// reader waits behind a writer that waits, as Go's do.
type Waiter = { writer: boolean; resolve: () => void };

export class RWMutex extends $.Native {
  writer = false;
  readers = 0;
  waiters: Waiter[] = [];
  Lock(): void {
    if (this.writer || this.readers > 0) contended();
    this.writer = true;
  }
  Lock$a(): Promise<void> | undefined {
    if (!this.writer && this.readers === 0 && this.waiters.length === 0) {
      this.writer = true;
      return undefined;
    }
    return new Promise<void>((resolve) => this.waiters.push({ writer: true, resolve }));
  }
  TryLock(): boolean {
    if (this.writer || this.readers > 0) return false;
    this.writer = true;
    return true;
  }
  Unlock(): void {
    if (!this.writer) $.panic($.val($.String, "sync: Unlock of unlocked RWMutex"));
    this.writer = false;
    this.wake();
  }
  RLock(): void {
    if (this.writer) contended();
    this.readers++;
  }
  RLock$a(): Promise<void> | undefined {
    if (!this.writer && this.waiters.length === 0) {
      this.readers++;
      return undefined;
    }
    return new Promise<void>((resolve) => this.waiters.push({ writer: false, resolve }));
  }
  TryRLock(): boolean {
    if (this.writer) return false;
    this.readers++;
    return true;
  }
  RUnlock(): void {
    if (this.readers <= 0) $.panic($.val($.String, "sync: RUnlock of unlocked RWMutex"));
    this.readers--;
    if (this.readers === 0) this.wake();
  }
  // wake lets the next waiters in: one writer, or every reader up to one.
  private wake(): void {
    if (this.writer || this.readers > 0) return;
    const w = this.waiters[0];
    if (w === undefined) return;
    if (w.writer) {
      this.waiters.shift();
      this.writer = true;
      w.resolve();
      return;
    }
    while (this.waiters.length > 0 && !this.waiters[0].writer) {
      this.readers++;
      this.waiters.shift()!.resolve();
    }
  }
  RLocker(): Locker {
    return { Lock: () => this.RLock(), Lock$a: () => this.RLock$a(), Unlock: () => this.RUnlock() } as Locker;
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
  waiters: Array<() => void> = [];
  Add(delta: number): void {
    this.n += delta;
    if (this.n < 0) $.panic($.val($.String, $.str.toGo("sync: negative WaitGroup counter")));
    if (this.n === 0) {
      const ws = this.waiters;
      this.waiters = [];
      for (const w of ws) w();
    }
  }
  Done(): void {
    this.Add(-1);
  }
  // Wait waits until the counter is back at zero.
  Wait(): Promise<void> {
    if (this.n === 0) return Promise.resolve();
    return new Promise((resolve) => this.waiters.push(resolve));
  }
  // Go runs f in a goroutine, counted until it returns.
  Go(f: () => any): void {
    this.Add(1);
    $.go(async () => {
      try {
        const r = f();
        if (r instanceof Promise) await r;
      } finally {
        this.Done();
      }
    });
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
  waiters: Array<() => void> = [];
  // Wait unlocks L, waits for Signal or Broadcast, and locks L again.
  async Wait(): Promise<void> {
    const woken = new Promise<void>((resolve) => this.waiters.push(resolve));
    this.L!.Unlock();
    await woken;
    const l = this.L as any;
    if (typeof l.Lock$a === "function") await l.Lock$a();
    else l.Lock();
  }
  Signal(): void {
    this.waiters.shift()?.();
  }
  Broadcast(): void {
    const ws = this.waiters;
    this.waiters = [];
    for (const w of ws) w();
  }
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
