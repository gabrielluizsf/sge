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

// addr is a's host and port, or for a Unix socket its path and -1.
function addr(a: Deno.Addr): [string, number] {
  if (a.transport === "tcp" || a.transport === "udp") return [a.hostname, a.port];
  if (a.transport === "unix" || a.transport === "unixpacket") return [a.path ?? "", -1];
  return ["", 0];
}

export class Sock {
  private buf = new Uint8Array(chunk);
  private have: Uint8Array | null = null;
  private eof = false;
  private err: unknown = null;
  private pending: Promise<void> | null = null;
  private wakers = new Set<Wake>();
  private writeWakers = new Set<Wake>();
  private writing: Promise<void> = Promise.resolve();

  constructor(private conn: Deno.Conn) {}

  // noDelay turns Nagle off on a TCP connection, as Go has it.
  noDelay(): void {
    (this.conn as Deno.TcpConn).setNoDelay?.(true);
  }

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
  // startWrite starts writing all of b, once the write before is done, so
  // the stream keeps its order, and gives the write in progress.
  startWrite(b: Uint8Array): Writing {
    const w = new Writing(this.writeWakers);
    const before = this.writing;
    w.promise = (async () => {
      await before;
      while (w.n < b.length) w.n += await this.conn.write(b.subarray(w.n));
    })().then(
      () => {
        w.done = true;
      },
      (e) => {
        w.done = true;
        w.err = failure(e);
      },
    );
    this.writing = w.promise;
    return w;
  }

  // wakeWriters has every waiting write look at its deadline again.
  wakeWriters(): void {
    for (const w of [...this.writeWakers]) w();
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
    this.wakeWriters();
  }
}

// Writing is a write in progress: how much of it is written, and how it
// ended once done.
export class Writing {
  n = 0;
  done = false;
  err: Failure | null = null;
  promise: Promise<void> = Promise.resolve();

  constructor(private wakers: Set<Wake>) {}

  // wait waits for the write to be done, at most ms milliseconds (ms < 0:
  // no limit) or until woken to look at a new deadline.
  async wait(ms: number): Promise<void> {
    if (!this.done && ms !== 0) await within(this.promise, ms, this.wakers);
  }
}

export class Listener {
  private next: Promise<Deno.Conn> | null = null;
  private wakers = new Set<Wake>();

  constructor(private l: Deno.Listener) {}

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
    let c: Deno.Conn;
    try {
      c = await p;
    } catch (e) {
      fail(e);
    }
    const sock = new Sock(c);
    sock.noDelay();
    return sock;
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

// Dgram is a UDP socket of Deno's. Like Sock's reads, a receive that runs
// out of time loses nothing: the datagram that arrives later is the next
// receive's.
export class Dgram {
  private have: [Uint8Array, Deno.Addr] | null = null;
  private err: unknown = null;
  private pending: Promise<void> | null = null;
  private wakers = new Set<Wake>();

  constructor(private c: Deno.DatagramConn) {}

  local(): [string, number] {
    return addr(this.c.addr);
  }

  // receive gives the next datagram, its sender's host and port, or "wait"
  // as Sock.read does.
  async receive(ms: number): Promise<[Uint8Array, string, number] | "wait"> {
    for (;;) {
      if (this.have !== null) {
        const [b, from] = this.have;
        this.have = null;
        const [h, p] = addr(from);
        return [b, h, p];
      }
      if (this.err !== null) fail(this.err);
      if (ms === 0) return "wait";
      if (this.pending === null) {
        this.pending = this.c.receive().then(
          (r) => {
            this.have = r;
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

  // send sends b to host:port as one datagram.
  async send(b: Uint8Array, host: string, port: number): Promise<number> {
    try {
      return await this.c.send(b, { transport: "udp", hostname: host, port });
    } catch (e) {
      fail(e);
    }
  }

  wake(): void {
    for (const w of [...this.wakers]) w();
  }

  close(): void {
    try {
      this.c.close();
    } catch {
      // Closed already.
    }
    this.wake();
  }
}

// listenUdp opens a UDP socket on host:port, all addresses when host is
// empty. Deno has UDP behind --unstable-net; without it this says so.
export function listenUdp(host: string, port: number): Dgram | Failure {
  const listen = (Deno as any).listenDatagram;
  if (typeof listen !== "function") {
    return { failed: true, name: "NotSupported", message: "UDP on Deno needs --unstable-net" };
  }
  try {
    return new Dgram(listen({ transport: "udp", hostname: host === "" ? "0.0.0.0" : host, port }));
  } catch (e) {
    return failure(e);
  }
}

// connectUnix dials the Unix socket at path.
export async function connectUnix(path: string): Promise<Sock> {
  try {
    return new Sock(await Deno.connect({ transport: "unix", path }));
  } catch (e) {
    fail(e);
  }
}

// listenUnix listens on the Unix socket at path, or gives the error it
// failed with.
export function listenUnix(path: string): Listener | Failure {
  try {
    return new Listener(Deno.listen({ transport: "unix", path }));
  } catch (e) {
    return failure(e);
  }
}

// lookup is the addresses host resolves to, as the system resolves them,
// its hosts file first: node:dns's lookup, loaded when first wanted.
// Deno's --allow-net has to grant the name.
export async function lookup(host: string): Promise<string[]> {
  try {
    const dns = await import("node:dns/promises");
    return (await dns.lookup(host, { all: true })).map((a: { address: string }) => a.address);
  } catch (e) {
    fail(e);
  }
}

(globalThis as any).gotsNet = { connect, listen, connectUnix, listenUnix, lookup, listenUdp };
