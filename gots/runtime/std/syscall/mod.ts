// Go's syscall package, the little of it that portable packages reach for:
// the environment and errno values. File descriptors belong to os, which is
// written against Deno directly.
import * as $ from "../../rt/mod.ts";

const g: any = globalThis;

export type Errno = number;

const messages: Record<number, string> = {
  1: "operation not permitted",
  2: "no such file or directory",
  9: "bad file descriptor",
  13: "permission denied",
  17: "file exists",
  20: "not a directory",
  21: "is a directory",
  22: "invalid argument",
  38: "function not implemented",
};

export function Errno$Error(e: Errno): string {
  return messages[e] ?? "errno " + e;
}

export function Errno$Is(e: Errno, target: any): boolean {
  const msg = typeof target?.Error === "function" ? $.str.fromGo(target.Error()) : "";
  switch (e) {
    case 1:
    case 13:
      return msg === "permission denied";
    case 17:
      return msg === "file already exists";
    case 2:
      return msg === "file does not exist";
    case 38:
      return msg === "operation not supported";
  }
  return false;
}

export function Errno$Temporary(e: Errno): boolean {
  return e === 4 || e === 24 || e === 11;
}

export function Errno$Timeout(e: Errno): boolean {
  return e === 11 || e === 110;
}

class Errno$Val extends $.Val {
  Error(): string {
    return Errno$Error(this.v);
  }
  Is(t: any): boolean {
    return Errno$Is(this.v, t);
  }
  Temporary(): boolean {
    return Errno$Temporary(this.v);
  }
  Timeout(): boolean {
    return Errno$Timeout(this.v);
  }
}
export const Errno$type = $.named("syscall.Errno", $.Uintptr, Errno$Val);

export function Getenv(key: string): [string, boolean] {
  const v = g.Deno?.env.get($.str.fromGo(key)) ?? g.process?.env[$.str.fromGo(key)];
  return v === undefined ? ["", false] : [$.str.toGo(v), true];
}

export function Setenv(key: string, value: string): any {
  g.Deno?.env.set($.str.fromGo(key), $.str.fromGo(value));
  return null;
}

export function Unsetenv(key: string): any {
  g.Deno?.env.delete($.str.fromGo(key));
  return null;
}

export function Clearenv(): void {
  for (const k of Object.keys(g.Deno?.env.toObject() ?? {})) g.Deno.env.delete(k);
}

export function Environ(): $.Slice<string> {
  const env = g.Deno?.env.toObject() ?? {};
  return $.slice(Object.entries(env).map(([k, v]) => $.str.toGo(k + "=" + v)));
}

export function Exit(code: number): never {
  g.Deno?.exit(code);
  throw new Error("exit " + code);
}

export function Getpid(): number {
  return g.Deno?.pid ?? 0;
}

export function Getuid(): number {
  return g.Deno?.uid?.() ?? -1;
}

export function Getpagesize(): number {
  return 65536;
}

const nosys = (): any => [0, $.val(Errno$type, 38)];
export const Open = nosys, Read = nosys, Write = nosys, Seek = nosys;

export function Close(_fd: number): any {
  return $.val(Errno$type, 38);
}
