// Goroutines: functions run as microtasks, each an async function when it
// may block. A Go program ends when main returns, goroutines or not, and
// dies of a panic no goroutine recovers.
import { Panic } from "./panic.ts";

const g: any = globalThis;
let exited = false;

// go is the go statement: f runs once the current goroutine waits or ends.
export function go(f: () => any): void {
  queueMicrotask(() => {
    if (exited) return;
    try {
      const r = f();
      if (r instanceof Promise) r.catch(crash);
    } catch (e) {
      crash(e);
    }
  });
}

// crash ends the program over a panic, as Go does: the panic on stderr,
// exit status 2.
export function crash(e: any): never {
  exited = true;
  const msg = e instanceof Panic ? e.message : "panic: " + (e instanceof Error ? (e.stack ?? e.message) : g.String(e));
  console.error(msg + "\n\ngoroutine [running]:");
  if (g.Deno !== undefined) g.Deno.exit(2);
  throw e;
}

// runMain runs a program's main. Once it returns the program is over:
// goroutines still running and timers still set stop there, as in Go. If nothing is left that
// could ever let main go on, it is a deadlock.
export function runMain(main: () => any): void {
  let done = false;
  if (g.addEventListener !== undefined) {
    g.addEventListener("beforeunload", () => {
      if (done || exited) return;
      exited = true;
      console.error("fatal error: all goroutines are asleep - deadlock!");
      g.Deno?.exit(2);
    });
  }
  // Timers, tickers and goroutines still pending would keep JavaScript's
  // event loop going.
  const finish = () => {
    done = true;
    exited = true;
    g.Deno?.exit(0);
  };
  let r: any;
  try {
    r = main();
  } catch (e) {
    crash(e);
  }
  if (r instanceof Promise) r.then(finish, crash);
  else finish();
}
