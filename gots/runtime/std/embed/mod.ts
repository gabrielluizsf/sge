// Go's embed package: an FS of the files a //go:embed directive named,
// which the transpiler writes into the module as base64.
import * as $ from "../../rt/mod.ts";
import * as fs from "../io/fs/mod.ts";
import * as io from "../io/mod.ts";
import * as time from "../time/mod.ts";

const go = $.str.toGo;

class info extends $.Native {
  name = "";
  size = 0;
  dir = false;
  Name(): string {
    return this.name;
  }
  Size(): number {
    return this.size;
  }
  Mode(): number {
    return this.dir ? (fs.ModeDir | 0o555) >>> 0 : 0o444;
  }
  ModTime(): any {
    return new time.Time();
  }
  IsDir(): boolean {
    return this.dir;
  }
  Sys(): any {
    return null;
  }
  Type(): number {
    return this.dir ? fs.ModeDir >>> 0 : 0;
  }
  Info(): [any, any] {
    return [this, null];
  }
  String(): string {
    return fs.FormatFileInfo(this);
  }
}
$.nativeType("embed.file", info);

function base(p: string): string {
  return p.slice(p.lastIndexOf("/") + 1);
}

class openFile extends $.Native {
  data: Uint8Array = new Uint8Array(0);
  off = 0;
  fi: info = new info();
  fsys: FS | null = null;
  path = "";
  Stat(): [any, any] {
    return [this.fi, null];
  }
  Read(b: $.Slice<number> | null): [number, any] {
    if (this.fi.dir) return [0, new fs.PathError("read", this.path, fs.ErrInvalid)];
    if (this.off >= this.data.length) return [0, io.EOF];
    const n = Math.min(b === null ? 0 : b.n, this.data.length - this.off);
    if (b !== null) b.a.set(this.data.subarray(this.off, this.off + n), b.o);
    this.off += n;
    return [n, null];
  }
  Close(): any {
    return null;
  }
  ReadDir(n: number): [any, any] {
    const [all, err] = this.fsys!.ReadDir(this.path);
    if (err !== null) return [null, err];
    const items = all.toArray();
    const rest = items.slice(this.off);
    const take = n <= 0 ? rest.length : Math.min(n, rest.length);
    this.off += take;
    if (n > 0 && take === 0) return [null, io.EOF];
    return [$.slice(rest.slice(0, take)), null];
  }
}
$.nativeType("embed.openFile", openFile);

export class FS extends $.Native {
  data: Map<string, Uint8Array> = new Map();

  private dirs(): Set<string> {
    const out = new Set<string>(["."]);
    for (const p of this.data.keys()) {
      const parts = p.split("/");
      for (let i = 1; i < parts.length; i++) out.add(parts.slice(0, i).join("/"));
    }
    return out;
  }

  Open(name: string): [any, any] {
    const f = new openFile();
    f.fsys = this;
    f.path = name;
    f.fi.name = name === "." ? "." : base(name);
    const data = this.data.get(name);
    if (data !== undefined) {
      f.data = data;
      f.fi.size = data.length;
      return [f, null];
    }
    if (this.dirs().has(name)) {
      f.fi.dir = true;
      return [f, null];
    }
    return [null, new fs.PathError("open", name, fs.ErrNotExist)];
  }

  ReadFile(name: string): [$.Slice<number> | null, any] {
    const data = this.data.get(name);
    if (data === undefined) {
      return [null, new fs.PathError(this.dirs().has(name) ? "read" : "open", name, this.dirs().has(name) ? fs.ErrInvalid : fs.ErrNotExist)];
    }
    return [$.slice<number>(data.slice()), null];
  }

  ReadDir(name: string): [any, any] {
    if (!this.dirs().has(name)) return [null, new fs.PathError("open", name, fs.ErrNotExist)];
    const prefix = name === "." ? "" : name + "/";
    const seen = new Map<string, info>();
    for (const p of [...this.data.keys(), ...this.dirs()]) {
      if (p === "." || !p.startsWith(prefix) || p === name) continue;
      const rest = p.slice(prefix.length);
      const entry = rest.split("/")[0];
      if (seen.has(entry)) continue;
      const i = new info();
      i.name = go(entry);
      i.dir = rest.includes("/") || !this.data.has(p);
      i.size = this.data.get(prefix + entry)?.length ?? 0;
      seen.set(entry, i);
    }
    const out = [...seen.values()].sort((a, b) => (a.name < b.name ? -1 : a.name > b.name ? 1 : 0));
    return [$.slice(out), null];
  }
}
$.nativeType("embed.FS", FS);

// files is an FS of the given names and base64 contents.
export function files(entries: [string, string][]): FS {
  const f = new FS();
  for (const [name, b64] of entries) {
    const bin = atob(b64);
    const data = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; i++) data[i] = bin.charCodeAt(i);
    f.data.set(name, data);
  }
  return f;
}
