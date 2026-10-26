// TLS for Go's crypto/tls on Deno, and on Node: node:tls speaks the
// protocol over a stream whose bytes Go carries. What the peer sends, Go
// reads from its own net.Conn and feeds in; what TLS sends, Go takes out
// and writes there. The connection stays Go's, dialed by its Dialer with
// its Control hooks, and Go checks what the handshake agreed on.
//
// node:tls is loaded when first wanted, as programs without TLS never need it.

import { failure, type Wake, within, Writing } from "./net.ts";

let mods: Promise<[any, any]> | null = null;
let tlsLoaded: any = null; // node:tls, once loaded

function load(): Promise<[any, any]> {
  return mods ??= Promise.all([import("node:tls"), import("node:stream")]).then((m) => {
    tlsLoaded = m[0];
    return m;
  });
}

// Info is what a handshake agreed on, as the overlay reads it.
export interface Info {
  version: string;
  alpn: string;
  cipher: string;
  // The peer's certificates, its own first: what it sent, with the root
  // node:tls found for it, but not a chain node:tls verified, as it may
  // hold certificates that signed nothing (Go builds the verified ones).
  chain: Uint8Array[];
  authorized: boolean;
  authError: string;
  serverName: string;
}

// A Side is one end of a TLS connection.
export class Side {
  private sock: any = null;
  private duplex: any = null;
  private out: Uint8Array[] = [];
  private outDone = false;
  private outWaiter: (() => void) | null = null;
  private have: Uint8Array[] = [];
  private eof = false;
  private err: unknown = null;
  private secure = false;
  private sni: { name: string; cb: (e: Error | null, ctx?: unknown) => void } | null = null;
  private changed!: Promise<void>;
  private signal!: () => void;
  private wakers = new Set<Wake>();
  private writeWakers = new Set<Wake>();

  constructor() {
    this.renew();
  }

  private renew(): void {
    this.changed = new Promise((resolve) => (this.signal = resolve));
  }

  // notify has whatever waits on the handshake or on reads look again.
  private notify(): void {
    const signal = this.signal;
    this.renew();
    signal();
  }

  // open starts the handshake, as a client or a server with opts.
  async open(server: boolean, opts: Record<string, unknown>): Promise<void> {
    const [tls, stream] = await load();
    if (opts.systemRoots) {
      // The pool holds the system's roots, which are the runtime's.
      opts.ca = [...tls.rootCertificates, ...(opts.ca as string[] ?? [])];
    }
    delete opts.systemRoots;
    this.duplex = new stream.Duplex({
      read() {},
      write: (chunk: Uint8Array, _enc: string, cb: () => void) => {
        this.out.push(new Uint8Array(chunk));
        this.wakeOut();
        cb();
      },
      final: (cb: () => void) => {
        this.outDone = true;
        this.wakeOut();
        cb();
      },
      destroy: (e: Error | null, cb: (e: Error | null) => void) => {
        this.outDone = true;
        this.wakeOut();
        cb(e);
      },
    });
    // node:tls may not see the stream end in the middle of a handshake
    // (after the peer's alert, say): once all that came before is read,
    // the handshake has failed.
    this.duplex.on("end", () => {
      setTimeout(() => {
        if (this.secure) return;
        this.err ??= new Error("tls: connection closed during the handshake");
        this.notify();
      }, 0);
    });
    if (server) {
      if (opts.sni) {
        opts.SNICallback = (name: string, cb: (e: Error | null, ctx?: unknown) => void) => {
          this.sni = { name, cb };
          this.notify();
        };
      }
      this.sock = new tls.TLSSocket(this.duplex, { ...opts, isServer: true });
      this.sock.on("secure", () => this.onSecure());
    } else {
      this.sock = tls.connect({ ...opts, socket: this.duplex });
      this.sock.on("secureConnect", () => this.onSecure());
    }
    this.sock.on("data", (b: Uint8Array) => {
      this.have.push(new Uint8Array(b));
      this.notify();
    });
    this.sock.on("end", () => {
      this.eof = true;
      this.notify();
    });
    this.sock.on("error", (e: unknown) => {
      this.err ??= e;
      this.notify();
    });
    this.sock.on("close", () => {
      this.eof = true;
      this.err ??= new Error("tls: connection closed");
      this.notify();
      this.wakeWriters();
    });
  }

  private onSecure(): void {
    this.secure = true;
    this.notify();
  }

  // step waits for the handshake: its Info once done, the server name to
  // pick a certificate for when a server is asked one, or "wait" when ms
  // milliseconds pass first (ms < 0: no limit) or it is woken.
  async step(ms: number): Promise<Info | { sni: string } | "wait"> {
    for (;;) {
      if (this.secure) return this.info();
      if (this.sni !== null) return { sni: this.sni.name };
      if (this.err !== null) throw failure(this.err);
      if (ms === 0) return "wait";
      if (!(await within(this.changed, ms, this.wakers))) return "wait";
    }
  }

  // answerSni gives the certificate chain and key, in PEM, for the server
  // name step asked about; refuseSni fails the handshake.
  async answerSni(cert: string, key: string): Promise<void> {
    const [tls] = await load();
    const q = this.sni!;
    this.sni = null;
    try {
      q.cb(null, tls.createSecureContext({ cert, key }));
    } catch (e) {
      q.cb(e as Error);
    }
  }

  refuseSni(message: string): void {
    const q = this.sni!;
    this.sni = null;
    q.cb(new Error(message));
  }

  private info(): Info {
    const s = this.sock;
    const chain: Uint8Array[] = [];
    const seen = new Set<unknown>();
    // A certificate its own issuer, a root, ends the chain.
    for (let p = s.getPeerCertificate(true); p && p.raw && !seen.has(p); p = p.issuerCertificate) {
      seen.add(p);
      chain.push(new Uint8Array(p.raw));
    }
    const c = s.getCipher?.();
    return {
      version: String(s.getProtocol?.() ?? ""),
      alpn: typeof s.alpnProtocol === "string" ? s.alpnProtocol : "",
      cipher: String(c?.standardName ?? c?.name ?? ""),
      chain,
      authorized: s.authorized === true,
      authError: String(s.authorizationError ?? ""),
      serverName: typeof s.servername === "string" ? s.servername : "",
    };
  }

  // feed gives TLS bytes the peer sent; feedEnd says it sends no more.
  feed(b: Uint8Array): void {
    this.duplex?.push(b);
  }

  feedEnd(): void {
    this.duplex?.push(null);

  }

  // takeOut gives the bytes TLS has to send, once there are some, or null
  // once it sends no more.
  async takeOut(): Promise<Uint8Array | null> {
    for (;;) {
      if (this.out.length > 0) {
        const parts = this.out;
        this.out = [];
        if (parts.length === 1) return parts[0];
        const b = new Uint8Array(parts.reduce((n, p) => n + p.length, 0));
        let o = 0;
        for (const p of parts) {
          b.set(p, o);
          o += p.length;
        }
        return b;
      }
      if (this.outDone) return null;
      await new Promise<void>((resolve) => (this.outWaiter = resolve));
    }
  }

  private wakeOut(): void {
    const w = this.outWaiter;
    this.outWaiter = null;
    w?.();
  }

  // read gives what has arrived, as Sock.read does.
  async read(max: number, ms: number): Promise<Uint8Array | null | "wait"> {
    for (;;) {
      if (this.have.length > 0) {
        const b = this.have[0];
        if (b.length <= max) {
          this.have.shift();
          return b;
        }
        this.have[0] = b.subarray(max);
        return b.subarray(0, max);
      }
      if (this.eof) return null;
      if (this.err !== null) throw failure(this.err);
      if (ms === 0) return "wait";
      if (!(await within(this.changed, ms, this.wakers))) return "wait";
    }
  }

  // startWrite starts writing all of b, as Sock.startWrite does.
  startWrite(b: Uint8Array): Writing {
    const w = new Writing(this.writeWakers);
    w.promise = new Promise<void>((resolve) => {
      this.sock.write(b, (e?: Error | null) => {
        w.done = true;
        if (e) w.err = failure(e);
        else w.n = b.length;
        resolve();
      });
    });
    return w;
  }

  // end sends close_notify, and no more after it.
  end(): void {
    this.sock.end();
  }

  wake(): void {
    for (const w of [...this.wakers]) w();
  }

  wakeWriters(): void {
    for (const w of [...this.writeWakers]) w();
  }

  destroy(): void {
    this.err ??= new Error("tls: connection closed");
    this.sock?.destroy();
    this.duplex?.destroy();
    this.notify();
    this.wake();
    this.wakeWriters();
  }
}

// rootCertificates are the roots node:tls verifies with by default, in
// PEM, or none when they cannot be read (Deno asks --allow-sys=ca).
function rootCertificates(): string[] {
  try {
    return [...(tlsLoaded?.rootCertificates ?? [])];
  } catch {
    return [];
  }
}

(globalThis as any).gotsTls = { side: () => new Side(), rootCertificates };
