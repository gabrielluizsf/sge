// Go channels over promises. A goroutine that has to wait for a channel
// parks in its queue and is resumed by the goroutine that makes it ready;
// everything runs on JavaScript's one thread, so a channel's state only
// changes between the steps a goroutine takes.
import { panic } from "./panic.ts";
import { val } from "./iface.ts";
import { String as GoString } from "./types.ts";

// A waiter is a parked goroutine, alone or one case of a select: once any
// case of its select fires, the others are stale and skipped.
interface Waiter {
  sel: { done: boolean } | null;
  index: number; // the select case this is, or 0
  v: any; // the value a sender offers
  wake(v: any, ok: boolean, index: number): void;
  fail(e: any): void;
}

export class Chan<T> {
  buf: T[] = [];
  closed = false;
  recvq: Waiter[] = [];
  sendq: Waiter[] = [];
  constructor(public cap: number, public zero: () => T) {
    if (cap < 0) panic(val(GoString, "makechan: size out of range"));
  }
}

function sendClosed(): never {
  panic(val(GoString, "send on closed channel"));
}

// first takes the first waiter of q whose select has not fired yet.
function first(q: Waiter[]): Waiter | null {
  while (q.length > 0) {
    const w = q.shift()!;
    if (w.sel === null || !w.sel.done) {
      if (w.sel !== null) w.sel.done = true;
      return w;
    }
  }
  return null;
}

function hasLive(q: Waiter[]): boolean {
  for (const w of q) if (w.sel === null || !w.sel.done) return true;
  return false;
}

const forever = (): Promise<never> => new Promise(() => {});

// trySend sends v if that needs no waiting.
function trySend<T>(c: Chan<T>, v: T): boolean {
  if (c.closed) sendClosed();
  const r = first(c.recvq);
  if (r !== null) {
    r.wake(v, true, r.index);
    return true;
  }
  if (c.buf.length < c.cap) {
    c.buf.push(v);
    return true;
  }
  return false;
}

// tryRecv receives if that needs no waiting: the value and ok, or null.
function tryRecv<T>(c: Chan<T>): [T, boolean] | null {
  if (c.buf.length > 0) {
    const v = c.buf.shift()!;
    // A sender waiting on a full buffer takes the room that made.
    const s = first(c.sendq);
    if (s !== null) {
      c.buf.push(s.v);
      s.wake(undefined, true, s.index);
    }
    return [v, true];
  }
  const s = first(c.sendq);
  if (s !== null) {
    s.wake(undefined, true, s.index);
    return [s.v, true];
  }
  if (c.closed) return [c.zero(), false];
  return null;
}

// send is c <- v.
export function send<T>(c: Chan<T> | null, v: T): Promise<void> | void {
  if (c === null) return forever();
  if (trySend(c, v)) return;
  return new Promise((resolve, reject) => {
    c.sendq.push({ sel: null, index: 0, v, wake: () => resolve(), fail: reject });
  });
}

// recv is v, ok := <-c.
export function recv<T>(c: Chan<T> | null): Promise<[T, boolean]> | [T, boolean] {
  if (c === null) return forever();
  const r = tryRecv(c);
  if (r !== null) return r;
  return new Promise((resolve, reject) => {
    c.recvq.push({ sel: null, index: 0, v: undefined, wake: (v, ok) => resolve([v, ok]), fail: reject });
  });
}

// close is close(c): waiting receivers get the zero value, waiting
// senders panic.
export function close<T>(c: Chan<T> | null): void {
  if (c === null) panic(val(GoString, "close of nil channel"));
  if (c.closed) panic(val(GoString, "close of closed channel"));
  c.closed = true;
  let w: Waiter | null;
  while ((w = first(c.recvq)) !== null) w.wake(c.zero(), false, w.index);
  while ((w = first(c.sendq)) !== null) {
    try {
      sendClosed();
    } catch (e) {
      w.fail(e);
    }
  }
}

export function chanLen(c: Chan<any> | null): number {
  return c === null ? 0 : c.buf.length;
}

export function chanCap(c: Chan<any> | null): number {
  return c === null ? 0 : c.cap;
}

// A select case: a send of v on c, or a receive from c.
export interface Case {
  c: Chan<any> | null;
  send: boolean;
  v?: any;
}

// The case a select took, -1 for its default, and what it received.
export type Chosen = [number, any, boolean];

// trySelect takes a ready case, in random order as Go does, or none.
function trySelect(cases: Case[]): Chosen | null {
  const n = cases.length;
  const start = n > 1 ? Math.floor(Math.random() * n) : 0;
  for (let k = 0; k < n; k++) {
    const i = (start + k) % n;
    const cs = cases[i];
    if (cs.c === null) continue;
    if (cs.send) {
      if (trySend(cs.c, cs.v)) return [i, undefined, false];
    } else {
      const r = tryRecv(cs.c);
      if (r !== null) return [i, r[0], r[1]];
    }
  }
  return null;
}

let polls = 0;

// poll is called before a select with a default in a goroutine that may
// spin on it: every 64th call it hands back a promise the goroutine waits
// on, letting the others run, which Go's preemptive scheduler would.
export function poll(): Promise<void> | undefined {
  if (++polls % 64 !== 0) return undefined;
  return new Promise((resolve) => setTimeout(resolve, 0));
}

// selectNow is a select with a default: a ready case, or -1.
export function selectNow(cases: Case[]): Chosen {
  return trySelect(cases) ?? [-1, undefined, false];
}

// select is a select with no default: it waits on every case, and the
// first to be ready is taken.
export function select(cases: Case[]): Promise<Chosen> | Chosen {
  const now = trySelect(cases);
  if (now !== null) return now;
  if (cases.every((cs) => cs.c === null)) return forever();
  return new Promise((resolve, reject) => {
    const sel = { done: false };
    cases.forEach((cs, index) => {
      if (cs.c === null) return;
      const w: Waiter = {
        sel,
        index,
        v: cs.v,
        wake: (v, ok, i) => resolve([i, v, ok]),
        fail: reject,
      };
      if (cs.send) {
        // A closed channel makes a send panic once it is chosen.
        if (cs.c.closed) {
          sel.done = true;
          try {
            sendClosed();
          } catch (e) {
            reject(e);
          }
          return;
        }
        cs.c.sendq.push(w);
      } else cs.c.recvq.push(w);
    });
  });
}

// waiting reports whether anything waits on c, for tests of the runtime.
export function waiting(c: Chan<any>): boolean {
  return hasLive(c.recvq) || hasLive(c.sendq);
}
