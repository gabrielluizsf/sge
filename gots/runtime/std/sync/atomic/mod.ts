// Go's sync/atomic package: on one thread every operation is atomic.
import * as $ from "../../../rt/mod.ts";

function addNum(p: $.Ptr<number>, d: number, wrapTo: (x: number) => number): number {
  p.v = wrapTo(p.v + d);
  return p.v;
}

const i32 = (x: number) => x | 0;
const u32 = (x: number) => x >>> 0;
const same = (x: number) => x;

export function AddInt32(p: $.Ptr<number>, d: number): number { return addNum(p, d, i32); }
export function AddUint32(p: $.Ptr<number>, d: number): number { return addNum(p, d, u32); }
export function AddInt64(p: $.Ptr<number>, d: number): number { return addNum(p, d, same); }
export function AddUintptr(p: $.Ptr<number>, d: number): number { return addNum(p, d, same); }
export function AddUint64(p: $.Ptr<bigint>, d: bigint): bigint {
  p.v = BigInt.asUintN(64, p.v + d);
  return p.v;
}

function load<T>(p: $.Ptr<T>): T { return p.v; }
function store<T>(p: $.Ptr<T>, v: T): void { p.v = v; }
function swap<T>(p: $.Ptr<T>, v: T): T {
  const old = p.v;
  p.v = v;
  return old;
}
function cas<T>(p: $.Ptr<T>, old: T, v: T): boolean {
  if (p.v !== old) return false;
  p.v = v;
  return true;
}

export const LoadInt32 = load, LoadInt64 = load, LoadUint32 = load, LoadUint64 = load, LoadUintptr = load, LoadPointer = load;
export const StoreInt32 = store, StoreInt64 = store, StoreUint32 = store, StoreUint64 = store, StoreUintptr = store, StorePointer = store;
export const SwapInt32 = swap, SwapInt64 = swap, SwapUint32 = swap, SwapUint64 = swap, SwapUintptr = swap, SwapPointer = swap;
export const CompareAndSwapInt32 = cas, CompareAndSwapInt64 = cas, CompareAndSwapUint32 = cas, CompareAndSwapUint64 = cas, CompareAndSwapUintptr = cas, CompareAndSwapPointer = cas;

// counter is the shape of the Int32, Int64, Uint32, Uint64 and Uintptr
// types: a value with atomic methods.
function counter(name: string, zero: any, wrapTo: (x: any) => any) {
  class C extends $.Native {
    v: any = zero;
    Load(): any { return this.v; }
    Store(v: any): void { this.v = v; }
    Swap(v: any): any {
      const old = this.v;
      this.v = v;
      return old;
    }
    CompareAndSwap(old: any, v: any): boolean {
      if (this.v !== old) return false;
      this.v = v;
      return true;
    }
    Add(d: any): any {
      this.v = wrapTo(this.v + d);
      return this.v;
    }
    And(m: any): any {
      const old = this.v;
      this.v = wrapTo(this.v & m);
      return old;
    }
    Or(m: any): any {
      const old = this.v;
      this.v = wrapTo(this.v | m);
      return old;
    }
  }
  $.nativeType("atomic." + name, C);
  return C;
}

export const Int32 = counter("Int32", 0, i32);
export const Uint32 = counter("Uint32", 0, u32);
export const Int64 = counter("Int64", 0, same);
export const Uintptr = counter("Uintptr", 0, same);
export const Uint64 = counter("Uint64", 0n, (x: bigint) => BigInt.asUintN(64, x));

export class Bool extends $.Native {
  v = false;
  Load(): boolean { return this.v; }
  Store(v: boolean): void { this.v = v; }
  Swap(v: boolean): boolean {
    const old = this.v;
    this.v = v;
    return old;
  }
  CompareAndSwap(old: boolean, v: boolean): boolean {
    if (this.v !== old) return false;
    this.v = v;
    return true;
  }
}
$.nativeType("atomic.Bool", Bool);

export class Value extends $.Native {
  v: any = null;
  Load(): any { return this.v; }
  Store(v: any): void {
    if (v === null) $.panic($.val($.String, "sync/atomic: store of nil value into Value"));
    this.v = v;
  }
  Swap(v: any): any {
    const old = this.v;
    this.v = v;
    return old;
  }
  CompareAndSwap(old: any, v: any): boolean {
    if (!$.ifaceEq(this.v, old)) return false;
    this.v = v;
    return true;
  }
}
$.nativeType("atomic.Value", Value);

export class Pointer<T> extends $.Native {
  v: T | null = null;
  Load(): T | null { return this.v; }
  Store(v: T | null): void { this.v = v; }
  Swap(v: T | null): T | null {
    const old = this.v;
    this.v = v;
    return old;
  }
  CompareAndSwap(old: T | null, v: T | null): boolean {
    if (this.v !== old) return false;
    this.v = v;
    return true;
  }
}
$.nativeType("atomic.Pointer", Pointer);
