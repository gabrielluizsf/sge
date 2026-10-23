// Go's testing package, enough for tests to run under `deno test` (each Go
// test becomes a Deno test) or, without Deno, as a plain script.
import * as $ from "../../rt/mod.ts";
import { Sprintf, Sprintln } from "../fmt/mod.ts";

class failNow {}
class skipNow {}

export function Short(): boolean {
  return false;
}

export function Verbose(): boolean {
  return false;
}

export function Testing(): boolean {
  return true;
}

export class T {
  name: string;
  failed = false;
  skipped = false;
  logs: string[] = [];
  cleanups: Array<() => void> = [];
  private parent: T | null;

  constructor(name: string, parent: T | null = null) {
    this.name = name;
    this.parent = parent;
  }

  Name(): string {
    return this.name;
  }

  Helper(): void {}

  Parallel(): void {}

  Log(a: $.Slice<any> | null): void {
    this.logs.push(Sprintln(a).replace(/\n$/, ""));
  }

  Logf(format: string, a: $.Slice<any> | null): void {
    this.logs.push(Sprintf(format, a));
  }

  Fail(): void {
    this.failed = true;
    for (let p = this.parent; p !== null; p = p.parent) p.failed = true;
  }

  Failed(): boolean {
    return this.failed;
  }

  FailNow(): never {
    this.Fail();
    throw new failNow();
  }

  Error(a: $.Slice<any> | null): void {
    this.Log(a);
    this.Fail();
  }

  Errorf(format: string, a: $.Slice<any> | null): void {
    this.Logf(format, a);
    this.Fail();
  }

  Fatal(a: $.Slice<any> | null): never {
    this.Log(a);
    return this.FailNow();
  }

  Fatalf(format: string, a: $.Slice<any> | null): never {
    this.Logf(format, a);
    return this.FailNow();
  }

  Skip(a: $.Slice<any> | null): never {
    this.Log(a);
    return this.SkipNow();
  }

  Skipf(format: string, a: $.Slice<any> | null): never {
    this.Logf(format, a);
    return this.SkipNow();
  }

  SkipNow(): never {
    this.skipped = true;
    throw new skipNow();
  }

  Skipped(): boolean {
    return this.skipped;
  }

  Cleanup(f: () => void): void {
    this.cleanups.push(f);
  }

  TempDir(): string {
    const g: any = globalThis;
    return g.Deno?.makeTempDirSync() ?? "/tmp";
  }

  // Run runs f as a subtest, reporting under this test's name.
  Run(name: string, f: (t: T) => void): boolean {
    const sub = new T(this.name + "/" + name.replaceAll(" ", "_"), this);
    run(sub, f);
    for (const l of sub.logs) this.logs.push("    " + sub.name + ": " + l);
    return !sub.failed;
  }
}

// B is enough of testing.B for benchmarks to compile; they do not run.
export class B extends T {
  N = 1;
  ResetTimer(): void {}
  StopTimer(): void {}
  StartTimer(): void {}
  ReportAllocs(): void {}
  Loop(): boolean {
    return this.N-- > 0;
  }
}

function run(t: T, f: (t: T) => void): void {
  try {
    f(t);
  } catch (e) {
    if (!(e instanceof failNow) && !(e instanceof skipNow)) {
      t.logs.push("panic: " + (e instanceof Error ? e.message : globalThis.String(e)));
      if (e instanceof Error && e.stack) t.logs.push(e.stack);
      t.Fail();
    }
  } finally {
    for (let i = t.cleanups.length - 1; i >= 0; i--) t.cleanups[i]();
  }
}

function report(t: T): string {
  return t.logs.map((l) => "    " + $.str.fromGo(l)).join("\n");
}

const queued: Array<[string, (t: T) => void]> = [];

// register is how a transpiled test file announces each of its tests.
export function register(name: string, f: (t: T) => void, skipped = false): void {
  const g: any = globalThis;
  if (g.Deno?.test !== undefined) {
    g.Deno.test({ name, ignore: skipped }, () => {
      const t = new T(name);
      run(t, f);
      if (t.failed) throw new Error("\n" + report(t));
      if (t.skipped) console.log(report(t));
    });
    return;
  }
  if (queued.length === 0) setTimeout(runQueued, 0);
  queued.push([name, f]);
}

function runQueued(): void {
  const g: any = globalThis;
  let failed = 0;
  for (const [name, f] of queued) {
    const t = new T(name);
    run(t, f);
    const status = t.failed ? "FAIL" : t.skipped ? "SKIP" : "ok  ";
    console.log(`${status} ${name}`);
    if (t.failed || t.skipped) console.log(report(t));
    if (t.failed) failed++;
  }
  console.log(failed === 0 ? `PASS (${queued.length} tests)` : `FAIL (${failed} of ${queued.length} tests)`);
  if (failed > 0 && g.process !== undefined) g.process.exitCode = 1;
}
