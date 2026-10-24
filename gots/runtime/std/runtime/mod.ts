// Go's runtime package, for a program that is JavaScript on one thread.
import * as $ from "../../rt/mod.ts";

export const GOOS = "js";
export const GOARCH = "wasm";
export const Compiler = "gc";

export function Version(): string {
  return "go1.27.1";
}

export function NumCPU(): number {
  return 1;
}

export function GOMAXPROCS(_n: number): number {
  return 1;
}

export function NumGoroutine(): number {
  return 1;
}

// Gosched lets the goroutines waiting to run go first.
export function Gosched(): Promise<void> {
  return new Promise((resolve) => queueMicrotask(resolve));
}
export function GC(): void {}
export function KeepAlive(_x: any): void {}
export function SetFinalizer(_obj: any, _finalizer: any): void {}
export function LockOSThread(): void {}
export function UnlockOSThread(): void {}
export function Breakpoint(): void {}

export function Goexit(): never {
  throw new $.Panic($.val($.String, "runtime.Goexit"));
}

export function Caller(_skip: number): [number, string, number, boolean] {
  return [0, "", 0, false];
}

export function Callers(_skip: number, _pc: $.Slice<number> | null): number {
  return 0;
}

export function FuncForPC(_pc: number): any {
  return null;
}

export function Stack(_buf: $.Slice<number> | null, _all: boolean): number {
  return 0;
}

export class MemStats extends $.Native {
  Alloc = 0;
  TotalAlloc = 0;
  Sys = 0;
  Mallocs = 0;
  Frees = 0;
  HeapAlloc = 0;
  NumGC = 0;
}
$.nativeType("runtime.MemStats", MemStats);

export function ReadMemStats(_m: MemStats | null): void {}

export interface Error {
  Error(): string;
  RuntimeError(): void;
}
export const Error$type = $.ifaceType("runtime.Error", ["Error", "RuntimeError"]);
