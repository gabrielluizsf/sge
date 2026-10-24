import { fromGo } from "./string.ts";

// Panic is a Go panic travelling up the JavaScript stack as an exception.
export class Panic extends Error {
  value: any;
  constructor(value: any) {
    super(describe(value));
    this.value = value;
  }
}

function describe(v: any): string {
  if (v === null) return "panic: nil";
  if (typeof v?.Error === "function") return "panic: " + fromGo(v.Error());
  if (v !== null && typeof v === "object" && "v" in v) return "panic: " + globalThis.String(v.v);
  return "panic: " + globalThis.String(v);
}

// panic is Go's panic(v).
export function panic(v: any): never {
  throw new Panic(v);
}

// RuntimeError is what recover returns for a fault JavaScript raised itself,
// such as calling a method on nil: Go's runtime.Error.
export class RuntimeError {
  msg: string;
  constructor(msg: string) {
    this.msg = msg;
  }
  Error(): string {
    return "runtime error: " + this.msg;
  }
  RuntimeError(): void {}
}

// Frame is the defer stack of one call to a function that defers.
export class Frame {
  defers: Array<() => void> = [];
  panicking = false;
  recovered = false;
  value: any = null;
  thrown: any = null;

  defer(f: () => void): void {
    this.defers.push(f);
  }

  // fail records the exception the function body threw.
  fail(e: any): void {
    this.panicking = true;
    this.thrown = e;
    this.value = e instanceof Panic ? e.value : new RuntimeError(e instanceof Error ? e.message : globalThis.String(e));
  }

  // unwind runs the deferred calls, newest first, then rethrows a panic none
  // of them recovered. A panic in a deferred call replaces the one before it.
  unwind(): void {
    while (this.defers.length > 0) {
      const f = this.defers.pop()!;
      frames.push(this);
      try {
        f();
      } catch (e) {
        this.recovered = false;
        this.fail(e);
      } finally {
        frames.pop();
      }
    }
    if (this.panicking && !this.recovered) throw this.thrown;
  }

  // unwindAsync is unwind for a function that may block, whose deferred
  // calls may too. Goroutines run while one waits, pushing and popping
  // frames of their own, so this one leaves the stack by name.
  async unwindAsync(): Promise<void> {
    while (this.defers.length > 0) {
      const f = this.defers.pop()!;
      frames.push(this);
      try {
        const r: any = f();
        if (r instanceof Promise) await r;
      } catch (e) {
        this.recovered = false;
        this.fail(e);
      } finally {
        const i = frames.lastIndexOf(this);
        if (i >= 0) frames.splice(i, 1);
      }
    }
    if (this.panicking && !this.recovered) throw this.thrown;
  }
}

const frames: Frame[] = [];

// recover is Go's recover(): the value of the panic being unwound by the
// innermost deferred call, which then stops being a panic.
export function recover(): any {
  const f = frames[frames.length - 1];
  if (f === undefined || !f.panicking || f.recovered) return null;
  f.recovered = true;
  return f.value;
}
