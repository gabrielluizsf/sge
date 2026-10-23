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

// B is testing.B. As go test -bench does, a benchmark is run with b.N
// growing until it takes benchtime, and its line is printed in Go's
// format, so the two can be compared.
export class B extends T {
  N = 0;
  private start = 0;
  private elapsed = 0;
  private timing = false;
  private loops = 0;

  ResetTimer(): void {
    if (this.timing) this.start = performance.now();
    this.elapsed = 0;
  }
  StopTimer(): void {
    if (!this.timing) return;
    this.elapsed += performance.now() - this.start;
    this.timing = false;
  }
  StartTimer(): void {
    if (this.timing) return;
    this.start = performance.now();
    this.timing = true;
  }
  ReportAllocs(): void {}
  SetBytes(_n: number): void {}
  ReportMetric(_n: number, _unit: string): void {}
  Elapsed(): number {
    return (this.elapsed + (this.timing ? performance.now() - this.start : 0)) * 1e6;
  }

  // Loop is b.Loop: true b.N times, the time before the first call left
  // out, the time after the last one too.
  Loop(): boolean {
    if (this.loops === 0) this.ResetTimer();
    if (this.loops++ < this.N) return true;
    this.StopTimer();
    return false;
  }

  // Run would need its sub-benchmarks timed apart; the benchmark is
  // skipped instead of timing them together.
  override Run(name: string, _f: (t: any) => void): boolean {
    this.Skip($.slice([$.val($.String, $.str.toGo("gots: sub-benchmark " + name + " needs b.Run, which gots does not support yet"))]));
    return false;
  }
}

const benchtime = 1000; // milliseconds, go test's -benchtime default

// runN runs f once with b.N = n and gives the B, its time spent.
function runN(name: string, f: (b: B) => void, n: number): B {
  const b = new B(name);
  b.N = n;
  b.StartTimer();
  run(b, f as (t: T) => void);
  b.StopTimer();
  return b;
}

// benchmark times f as go test -bench does and prints its line.
function benchmark(name: string, f: (b: B) => void): void {
  let n = 1;
  let b = runN(name, f, n);
  while (!b.failed && !b.skipped && b.Elapsed() < benchtime * 1e6 && n < 1e9) {
    // Go's prediction: aim for benchtime with 20% to spare, growing at
    // most a hundredfold a round and at least by one.
    const per = b.Elapsed() / n;
    const want = per > 0 ? (benchtime * 1e6 * 1.2) / per : n * 100;
    n = Math.max(Math.min(Math.floor(want), n * 100), n + 1);
    b = runN(name, f, n);
  }
  if (b.failed) throw new Error("\n" + report(b));
  if (b.skipped) {
    out(`--- SKIP: ${name}\n${report(b)}\n`);
    return;
  }
  out(`${name}\t${String(n).padStart(8)}\t${(b.Elapsed() / n).toFixed(2).padStart(10)} ns/op\n`);
}

// out writes s to stdout as it is: deno bench puts the bench's name in
// front of what console.log prints, and Go's lines start with it already.
function out(s: string): void {
  const g: any = globalThis;
  if (g.Deno?.stdout !== undefined) g.Deno.stdout.writeSync(new TextEncoder().encode(s));
  else console.log(s.trimEnd());
}

// registerBench is how a transpiled test file announces each of its
// benchmarks; deno bench runs them, deno test leaves them out. Deno calls
// a bench again and again; the first call does the timing.
export function registerBench(name: string, f: (b: B) => void, skipped = false): void {
  const g: any = globalThis;
  if (g.Deno?.bench === undefined) return;
  let done = false;
  g.Deno.bench({
    name,
    ignore: skipped,
    n: 1,
    warmup: 0,
    fn: () => {
      if (done) return;
      done = true;
      benchmark(name, f);
    },
  });
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
