// Go's os package over Deno's synchronous file system API. Errors come out
// as Go's do: a *fs.PathError wrapping a syscall.Errno.
import * as $ from "../../rt/mod.ts";
import * as fs from "../io/fs/mod.ts";
import * as io from "../io/mod.ts";
import * as time from "../time/mod.ts";
import * as syscall from "../syscall/mod.ts";

const g: any = globalThis;

function deno(): any {
  if (g.Deno === undefined) throw new Error("gots: package os needs Deno");
  return g.Deno;
}

const go = $.str.toGo;
const js = $.str.fromGo;

// errno maps an exception from Deno to the errno Go would have seen.
function errno(e: any): number {
  const D = g.Deno?.errors ?? {};
  if (D.NotFound && e instanceof D.NotFound) return 2;
  if (D.PermissionDenied && e instanceof D.PermissionDenied) return 13;
  if (D.NotCapable && e instanceof D.NotCapable) return 13;
  if (D.AlreadyExists && e instanceof D.AlreadyExists) return 17;
  if (D.IsADirectory && e instanceof D.IsADirectory) return 21;
  if (D.NotADirectory && e instanceof D.NotADirectory) return 20;
  if (D.BadResource && e instanceof D.BadResource) return 9;
  return 22;
}

function pathError(op: string, path: string, e: any): any {
  return new fs.PathError(op, path, $.val(syscall.Errno$type, errno(e)) as any);
}

// attempt runs f, turning a Deno exception into a *PathError.
function attempt<T>(op: string, path: string, zero: T, f: () => T): [T, any] {
  try {
    return [f(), null];
  } catch (e) {
    return [zero, pathError(op, path, e)];
  }
}

class osError {
  s: string;
  constructor(s: string) {
    this.s = s;
  }
  Error(): string {
    return this.s;
  }
  static $type = $.structType("errors.errorString", true, osError, null, () => [$.field("s", $.String)]);
}

export let ErrNotExist: any = fs.ErrNotExist;
export let ErrExist: any = fs.ErrExist;
export let ErrPermission: any = fs.ErrPermission;
export let ErrClosed: any = fs.ErrClosed;
export let ErrInvalid: any = fs.ErrInvalid;

// DeadlineExceededError is internal/poll's, the error a read or write past
// its deadline gives: a timeout, as net.Error has it.
class DeadlineExceededError extends $.Native {
  Error(): string {
    return "i/o timeout";
  }
  Timeout(): boolean {
    return true;
  }
  Temporary(): boolean {
    return true;
  }
}
$.nativeType("poll.DeadlineExceededError", DeadlineExceededError);

export let ErrDeadlineExceeded: any = new DeadlineExceededError();

// SyscallError is an error from a system call, saying which.
export class SyscallError extends $.Native {
  Syscall: string;
  Err: any;
  constructor(syscall = "", err: any = null) {
    super();
    this.Syscall = syscall;
    this.Err = err;
  }
  Error(): string {
    return this.Syscall + ": " + this.Err.Error();
  }
  Unwrap(): any {
    return this.Err;
  }
  // Timeout reports whether the error inside is a timeout.
  Timeout(): boolean {
    return typeof this.Err?.Timeout === "function" && this.Err.Timeout();
  }
}
$.nativeType("os.SyscallError", SyscallError);

// NewSyscallError is err from the system call syscall, or nil for none.
export function NewSyscallError(syscall: string, err: any): any {
  return err === null ? null : new SyscallError(syscall, err);
}
export let ErrNoDeadline: any = new osError("file type does not support deadline");
export const DevNull = "/dev/null";
export const PathSeparator = 47, PathListSeparator = 58;

function isErr(err: any, target: any): boolean {
  for (let e = err; e !== null && e !== undefined;) {
    if (e === target) return true;
    if (e instanceof $.Val && e.$t === syscall.Errno$type) return syscall.Errno$Is(e.v, target);
    e = typeof e.Unwrap === "function" ? e.Unwrap() : null;
  }
  return false;
}

export function IsNotExist(err: any): boolean {
  return isErr(err, fs.ErrNotExist);
}

export function IsExist(err: any): boolean {
  return isErr(err, fs.ErrExist);
}

export function IsPermission(err: any): boolean {
  return isErr(err, fs.ErrPermission);
}

export function IsTimeout(_err: any): boolean {
  return false;
}

// fileStat is the fs.FileInfo of a Deno file.
export class fileStat extends $.Native {
  name = "";
  info: any = null;
  Name(): string {
    return this.name;
  }
  Size(): number {
    return this.info.size;
  }
  Mode(): number {
    let m = (this.info.mode ?? 0o644) & 0o777;
    if (this.info.isDirectory) m |= fs.ModeDir;
    if (this.info.isSymlink) m |= fs.ModeSymlink;
    return m >>> 0;
  }
  ModTime(): any {
    const t: Date | null = this.info.mtime;
    const ms = t === null ? 0 : t.getTime();
    return time.Unix(Math.floor(ms / 1000), (ms % 1000) * 1e6);
  }
  IsDir(): boolean {
    return this.info.isDirectory;
  }
  Sys(): any {
    return null;
  }
}
$.nativeType("os.fileStat", fileStat);

function base(path: string): string {
  const p = path.replace(/\/+$/, "");
  const i = p.lastIndexOf("/");
  return i < 0 ? p : p.slice(i + 1);
}

function stat(info: any, path: string): fileStat {
  const s = new fileStat();
  s.name = base(path);
  s.info = info;
  return s;
}

export function Stat(name: string): [any, any] {
  return attempt("stat", name, null, () => stat(deno().statSync(js(name)), name));
}

export function Lstat(name: string): [any, any] {
  return attempt("lstat", name, null, () => stat(deno().lstatSync(js(name)), name));
}

export const O_RDONLY = 0, O_WRONLY = 1, O_RDWR = 2, O_APPEND = 1024, O_CREATE = 64, O_EXCL = 128, O_SYNC = 1052672, O_TRUNC = 512;

export class File extends $.Native {
  name = "";
  f: any = null;
  // stdio is 0, 1 or 2 for the standard streams, which Deno does not
  // hand out as files.
  stdio = -1;
  closed = false;
  entries: any[] | null = null;

  Name(): string {
    return this.name;
  }

  Fd(): number {
    return this.stdio >= 0 ? this.stdio : this.f?.rid ?? -1;
  }

  private stream(): any {
    if (this.closed) return null;
    if (this.stdio === 0) return deno().stdin;
    if (this.stdio === 1) return deno().stdout;
    if (this.stdio === 2) return deno().stderr;
    return this.f;
  }

  Read(b: $.Slice<number> | null): [number, any] {
    const s = this.stream();
    if (s === null) return [0, pathError("read", this.name, new (g.Deno.errors.BadResource)())];
    if (b === null || b.n === 0) return [0, null];
    try {
      const n = s.readSync(b.a.subarray(b.o, b.o + b.n));
      return n === null ? [0, io.EOF] : [n, null];
    } catch (e) {
      return [0, pathError("read", this.name, e)];
    }
  }

  ReadAt(b: $.Slice<number> | null, off: number): [number, any] {
    const [, err] = this.Seek(off, 0);
    if (err !== null) return [0, err];
    let total = 0;
    while (b !== null && total < b.n) {
      const [n, e] = this.Read($.sub(b, total));
      total += n;
      if (e !== null) return [total, e];
      if (n === 0) break;
    }
    return [total, null];
  }

  Write(b: $.Slice<number> | null): [number, any] {
    const s = this.stream();
    if (s === null) return [0, pathError("write", this.name, new (g.Deno.errors.BadResource)())];
    if (b === null) return [0, null];
    try {
      const data = b.a.subarray(b.o, b.o + b.n);
      for (let n = 0; n < data.length;) n += s.writeSync(data.subarray(n));
      return [b.n, null];
    } catch (e) {
      return [0, pathError("write", this.name, e)];
    }
  }

  WriteString(str: string): [number, any] {
    return this.Write($.slice($.str.strToBytes(str)));
  }

  WriteAt(b: $.Slice<number> | null, off: number): [number, any] {
    const [, err] = this.Seek(off, 0);
    if (err !== null) return [0, err];
    return this.Write(b);
  }

  Seek(offset: number, whence: number): [number, any] {
    const s = this.stream();
    if (s === null || typeof s.seekSync !== "function") return [0, pathError("seek", this.name, null)];
    return attempt("seek", this.name, 0, () => s.seekSync(offset, whence));
  }

  Stat(): [any, any] {
    if (this.f === null) return Stat(this.name);
    return attempt("stat", this.name, null, () => stat(this.f.statSync(), this.name));
  }

  Sync(): any {
    if (this.f !== null) this.f.syncSync();
    return null;
  }

  Truncate(size: number): any {
    return attempt("truncate", this.name, null, () => this.f.truncateSync(size))[1];
  }

  Chmod(mode: number): any {
    return Chmod(this.name, mode);
  }

  Close(): any {
    if (this.closed) return pathError("close", this.name, new (g.Deno?.errors.BadResource ?? Error)());
    this.closed = true;
    if (this.f !== null) this.f.close();
    return null;
  }

  ReadDir(n: number): [any, any] {
    if (this.entries === null) {
      const [es, err] = ReadDir(this.name);
      if (err !== null) return [null, err];
      this.entries = es === null ? [] : es.toArray();
    }
    const take = n <= 0 ? this.entries!.length : Math.min(n, this.entries!.length);
    const out = this.entries!.splice(0, take);
    if (n > 0 && out.length === 0) return [null, io.EOF];
    return [$.slice(out), null];
  }

  Readdirnames(n: number): [any, any] {
    const [es, err] = this.ReadDir(n);
    if (es === null) return [null, err];
    return [$.slice(es.toArray().map((d: any) => d.Name())), err];
  }
}
$.nativeType("os.File", File);

function stdio(fd: number, name: string): File {
  const f = new File();
  f.stdio = fd;
  f.name = name;
  return f;
}

export let Stdin: File | null = stdio(0, "/dev/stdin");
export let Stdout: File | null = stdio(1, "/dev/stdout");
export let Stderr: File | null = stdio(2, "/dev/stderr");

export function OpenFile(name: string, flag: number, perm: number): [File | null, any] {
  const acc = flag & 3;
  const opts = {
    read: acc === O_RDONLY || acc === O_RDWR,
    write: acc === O_WRONLY || acc === O_RDWR,
    append: (flag & O_APPEND) !== 0,
    truncate: (flag & O_TRUNC) !== 0,
    create: (flag & O_CREATE) !== 0 && (flag & O_EXCL) === 0,
    createNew: (flag & O_CREATE) !== 0 && (flag & O_EXCL) !== 0,
    mode: perm & 0o777,
  };
  return attempt("open", name, null, () => {
    const f = new File();
    f.name = name;
    f.f = deno().openSync(js(name), opts);
    return f;
  });
}

export function Open(name: string): [File | null, any] {
  return OpenFile(name, O_RDONLY, 0);
}

export function Create(name: string): [File | null, any] {
  return OpenFile(name, O_RDWR | O_CREATE | O_TRUNC, 0o666);
}

export function ReadFile(name: string): [$.Slice<number> | null, any] {
  return attempt("open", name, null, () => $.slice<number>(deno().readFileSync(js(name))));
}

export function WriteFile(name: string, data: $.Slice<number> | null, perm: number): any {
  const bytes = data === null ? new Uint8Array(0) : data.a.subarray(data.o, data.o + data.n);
  return attempt("open", name, null, () => deno().writeFileSync(js(name), bytes, { mode: perm & 0o777 }))[1];
}

export class dirEntry extends $.Native {
  name = "";
  dir = "";
  entry: any = null;
  Name(): string {
    return this.name;
  }
  IsDir(): boolean {
    return this.entry.isDirectory;
  }
  Type(): number {
    if (this.entry.isDirectory) return fs.ModeDir >>> 0;
    if (this.entry.isSymlink) return fs.ModeSymlink >>> 0;
    return 0;
  }
  Info(): [any, any] {
    return Lstat(this.dir + "/" + this.name);
  }
  String(): string {
    return fs.FormatDirEntry(this);
  }
}
$.nativeType("os.dirEntry", dirEntry);

export function ReadDir(name: string): [any, any] {
  return attempt("open", name, null, () => {
    const out: dirEntry[] = [];
    for (const en of deno().readDirSync(js(name))) {
      const d = new dirEntry();
      d.name = go(en.name);
      d.dir = name;
      d.entry = en;
      out.push(d);
    }
    out.sort((a, b) => (a.name < b.name ? -1 : a.name > b.name ? 1 : 0));
    return $.slice(out);
  });
}

export function Mkdir(name: string, perm: number): any {
  return attempt("mkdir", name, null, () => deno().mkdirSync(js(name), { mode: perm & 0o777 }))[1];
}

export function MkdirAll(name: string, perm: number): any {
  return attempt("mkdir", name, null, () => deno().mkdirSync(js(name), { recursive: true, mode: perm & 0o777 }))[1];
}

export function Remove(name: string): any {
  return attempt("remove", name, null, () => deno().removeSync(js(name)))[1];
}

export function RemoveAll(name: string): any {
  try {
    deno().removeSync(js(name), { recursive: true });
  } catch (e) {
    if (errno(e) !== 2) return pathError("unlinkat", name, e);
  }
  return null;
}

export function Rename(from: string, to: string): any {
  try {
    deno().renameSync(js(from), js(to));
    return null;
  } catch (e) {
    return pathError("rename", from + " " + to, e);
  }
}

export function Chmod(name: string, mode: number): any {
  return attempt("chmod", name, null, () => deno().chmodSync(js(name), mode & 0o777))[1];
}

export function Truncate(name: string, size: number): any {
  return attempt("truncate", name, null, () => deno().truncateSync(js(name), size))[1];
}

export function Symlink(oldname: string, newname: string): any {
  return attempt("symlink", newname, null, () => deno().symlinkSync(js(oldname), js(newname)))[1];
}

export function Readlink(name: string): [string, any] {
  return attempt("readlink", name, "", () => go(deno().readLinkSync(js(name))));
}

export function TempDir(): string {
  return go(deno().env.get("TMPDIR") ?? "/tmp");
}

// errPatternHasSeparator refuses a temporary name pattern that would
// reach outside its directory, as Go does.
const errPatternHasSeparator = new osError("pattern contains path separator");

export function MkdirTemp(dir: string, pattern: string): [string, any] {
  if (pattern.includes("/")) return ["", new fs.PathError("mkdirtemp", pattern, errPatternHasSeparator)];
  const d = dir === "" ? TempDir() : dir;
  const star = pattern.lastIndexOf("*");
  const opts = star < 0 ? { dir: js(d), prefix: js(pattern) } : { dir: js(d), prefix: js(pattern.slice(0, star)), suffix: js(pattern.slice(star + 1)) };
  return attempt("mkdirtemp", d, "", () => go(deno().makeTempDirSync(opts)));
}

export function CreateTemp(dir: string, pattern: string): [File | null, any] {
  if (pattern.includes("/")) return [null, new fs.PathError("createtemp", pattern, errPatternHasSeparator)];
  const d = dir === "" ? TempDir() : dir;
  const star = pattern.lastIndexOf("*");
  const opts = star < 0 ? { dir: js(d), prefix: js(pattern) } : { dir: js(d), prefix: js(pattern.slice(0, star)), suffix: js(pattern.slice(star + 1)) };
  const [name, err] = attempt("createtemp", d, "", () => go(deno().makeTempFileSync(opts)));
  if (err !== null) return [null, err];
  return OpenFile(name, O_RDWR, 0o600);
}

export function Getwd(): [string, any] {
  return attempt("getwd", ".", "", () => go(deno().cwd()));
}

export function Chdir(dir: string): any {
  return attempt("chdir", dir, null, () => deno().chdir(js(dir)))[1];
}

export function Getenv(key: string): string {
  return syscall.Getenv(key)[0];
}

export function LookupEnv(key: string): [string, boolean] {
  return syscall.Getenv(key);
}

export function Setenv(key: string, value: string): any {
  return syscall.Setenv(key, value);
}

export function Unsetenv(key: string): any {
  return syscall.Unsetenv(key);
}

export function Clearenv(): void {
  syscall.Clearenv();
}

export function Environ(): $.Slice<string> {
  return syscall.Environ();
}

export function ExpandEnv(s: string): string {
  return Expand(s, Getenv);
}

export function Expand(s: string, mapping: (k: string) => string): string {
  return s.replace(/\$\{([^}]*)\}|\$([A-Za-z0-9_]+|[*#$@!?0-9-])/g, (_m, a, b) => mapping(a ?? b));
}

export function Exit(code: number): never {
  return syscall.Exit(code);
}

export function Getpid(): number {
  return syscall.Getpid();
}

export function Getppid(): number {
  return g.Deno?.ppid ?? 0;
}

export function Getuid(): number {
  return syscall.Getuid();
}

export function Geteuid(): number {
  return syscall.Getuid();
}

export function Getgid(): number {
  return g.Deno?.gid?.() ?? -1;
}

export function Getpagesize(): number {
  return syscall.Getpagesize();
}

export function Hostname(): [string, any] {
  return attempt("hostname", "", "", () => go(deno().hostname()));
}

export function Executable(): [string, any] {
  return attempt("executable", "", "", () => go(deno().execPath()));
}

export function UserHomeDir(): [string, any] {
  const h = Getenv("HOME");
  return h === "" ? ["", new osError("$HOME is not defined")] : [h, null];
}

export function UserCacheDir(): [string, any] {
  const x = Getenv("XDG_CACHE_HOME");
  if (x !== "") return [x, null];
  const [h, err] = UserHomeDir();
  return err !== null ? ["", err] : [h + "/.cache", null];
}

export function UserConfigDir(): [string, any] {
  const x = Getenv("XDG_CONFIG_HOME");
  if (x !== "") return [x, null];
  const [h, err] = UserHomeDir();
  return err !== null ? ["", err] : [h + "/.config", null];
}

export let Args: $.Slice<string> = $.slice([go(g.Deno?.mainModule ?? "main"), ...(g.Deno?.args ?? []).map(go)]);
export function Args$set(v: $.Slice<string>): void {
  Args = v;
}

export function Stdin$set(v: File | null): void {
  Stdin = v;
}
export function Stdout$set(v: File | null): void {
  Stdout = v;
}
export function Stderr$set(v: File | null): void {
  Stderr = v;
}
