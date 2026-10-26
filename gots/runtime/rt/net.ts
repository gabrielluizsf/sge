// TCP sockets for Go's net on Deno: Deno.connect and Deno.listen behind
// the small, promise-based surface the net overlay reaches through
// syscall/js as gotsNet. Hostnames go to Deno as they are, so Deno looks
// them up and its --allow-net checks the name and what it resolves to.
//
// Go's reads take deadlines, and a read that runs out of time must not
// lose what arrives later: a socket keeps at most one read of its own in
// flight, and what it brings waits for the next read.

const chunk = 64 * 1024;

// A waiter is told when something it waits on settles, or that it should
// look again because a deadline changed.
type Wake = () => void;

// within resolves true when p settles, false once ms milliseconds pass
// (never, for ms < 0) or wake is called, whichever is first.
function within(p: Promise<unknown>, ms: number, wakers: Set<Wake>): Promise<boolean> {
  return new Promise((resolve) => {
    let timer: ReturnType<typeof setTimeout> | undefined;
    const done = (v: boolean) => {
      if (timer !== undefined) clearTimeout(timer);
      wakers.delete(wake);
      resolve(v);
    };
    const wake = () => done(false);
    wakers.add(wake);
    if (ms >= 0) timer = setTimeout(() => done(false), ms);
    p.then(() => done(true), () => done(true));
  });
}

// A Failure is a Deno error as the overlay reads it: its name, which it
// turns into an errno (ConnectionRefused, AddrInUse...), and its message.
export interface Failure {
  failed: true;
  name: string;
  message: string;
}

function failure(e: unknown): Failure {
  const err = e instanceof Error ? e : new Error(String(e));
  return { failed: true, name: err.name, message: err.message };
}

function fail(e: unknown): never {
  throw failure(e);
}

function addr(a: Deno.Addr): [string, number] {
  return a.transport === "tcp" || a.transport === "udp" ? [a.hostname, a.port] : ["", 0];
}

export class Sock {
  private buf = new Uint8Array(chunk);
  private have: Uint8Array | null = null;
  private eof = false;
  private err: unknown = null;
  private pending: Promise<void> | null = null;
  private wakers = new Set<Wake>();

  constructor(private conn: Deno.Conn) {}

  local(): [string, number] {
    return addr(this.conn.localAddr);
  }

  remote(): [string, number] {
    return addr(this.conn.remoteAddr);
  }

  // read gives what has arrived, at most max bytes; null at the end of
  // the stream; "wait" when ms milliseconds pass first (ms < 0: no limit)
  // or the reader is woken to look at a new deadline.
  async read(max: number, ms: number): Promise<Uint8Array | null | "wait"> {
    for (;;) {
      if (this.have !== null) {
        const out = this.have.subarray(0, max);
        this.have = max < this.have.length ? this.have.subarray(max) : null;
        return out;
      }
      if (this.err !== null) fail(this.err);
      if (this.eof) return null;
      if (ms === 0) return "wait";
      if (this.pending === null) {
        this.pending = this.conn.read(this.buf).then(
          (n) => {
            if (n === null) this.eof = true;
            else this.have = this.buf.slice(0, n);
            this.pending = null;
          },
          (e) => {
            this.err = e;
            this.pending = null;
          },
        );
      }
      if (!(await within(this.pending, ms, this.wakers))) return "wait";
    }
  }

  // write writes all of b.
  async write(b: Uint8Array): Promise<number> {
    let n = 0;
    try {
      while (n < b.length) n += await this.conn.write(b.subarray(n));
    } catch (e) {
      fail(e);
    }
    return n;
  }

  // wake has every waiting read look at its deadline again.
  wake(): void {
    for (const w of [...this.wakers]) w();
  }

  // closeWrite ends what is sent, giving the error it failed with, if any.
  closeWrite(): Failure | null {
    try {
      this.conn.closeWrite();
      return null;
    } catch (e) {
      return failure(e);
    }
  }

  close(): void {
    try {
      this.conn.close();
    } catch {
      // Closed already.
    }
    this.wake();
  }
}

export class Listener {
  private next: Promise<Deno.TcpConn> | null = null;
  private wakers = new Set<Wake>();

  constructor(private l: Deno.TcpListener) {}

  local(): [string, number] {
    return addr(this.l.addr);
  }

  // accept gives the next connection, or "wait" as Sock.read does. A
  // connection that arrives after a wait is kept for the next accept.
  async accept(ms: number): Promise<Sock | "wait"> {
    if (this.next === null) {
      this.next = this.l.accept();
      this.next.catch(() => {});
    }
    const p = this.next;
    if (ms === 0 || !(await within(p, ms, this.wakers))) return "wait";
    this.next = null;
    let c: Deno.TcpConn;
    try {
      c = await p;
    } catch (e) {
      fail(e);
    }
    c.setNoDelay(true);
    return new Sock(c);
  }

  wake(): void {
    for (const w of [...this.wakers]) w();
  }

  close(): void {
    try {
      this.l.close();
    } catch {
      // Closed already.
    }
    this.wake();
  }
}

// connect dials host:port over TCP, with Nagle off as Go has it.
export async function connect(host: string, port: number): Promise<Sock> {
  let c: Deno.TcpConn;
  try {
    c = await Deno.connect({ hostname: host, port, transport: "tcp" });
  } catch (e) {
    fail(e);
  }
  c.setNoDelay(true);
  return new Sock(c);
}

// listen listens on host:port, all addresses when host is empty, or gives
// the error it failed with: the calls Go makes synchronously never throw.
export function listen(host: string, port: number): Listener | Failure {
  try {
    return new Listener(Deno.listen({ hostname: host === "" ? "0.0.0.0" : host, port, transport: "tcp" }));
  } catch (e) {
    return failure(e);
  }
}

(globalThis as any).gotsNet = { connect, listen };
