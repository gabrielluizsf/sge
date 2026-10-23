// Go's os/exec over Deno.Command, run synchronously.
import * as $ from "../../../rt/mod.ts";

const g: any = globalThis;
const js = $.str.fromGo, go = $.str.toGo;

class execError {
  s: string;
  constructor(s: string) {
    this.s = s;
  }
  Error(): string {
    return this.s;
  }
  static $type = $.structType("errors.errorString", true, execError, null, () => [$.field("s", $.String)]);
}

export let ErrNotFound: any = new execError("executable file not found in $PATH");
export let ErrDot: any = new execError("cannot run executable found relative to current directory");

// Error is Go's *exec.Error: the name looked up, and why it failed.
export class Error extends $.Native {
  Name = "";
  Err: any = null;
  Error(): string {
    return "exec: " + $.str.toGo(JSON.stringify($.str.fromGo(this.Name))) + ": " + $.str.toGo(this.Err.Error());
  }
  Unwrap(): any {
    return this.Err;
  }
}
$.nativeType("exec.Error", Error);

function lookErr(name: string, err: any): Error {
  const e = new Error();
  e.Name = name;
  e.Err = err;
  return e;
}

function executable(full: string): boolean {
  try {
    const info = g.Deno.statSync(full);
    return info.isFile && ((info.mode ?? 0o111) & 0o111) !== 0;
  } catch {
    return false;
  }
}

// LookPath searches PATH as Go does on Unix, empty entries meaning the
// current directory, and refuses, with ErrDot, an executable found
// through a relative entry.
export function LookPath(file: string): [string, any] {
  if (file.includes("/")) return executable(js(file)) ? [file, null] : ["", lookErr(file, ErrNotFound)];
  const path = g.Deno?.env.get("PATH") ?? "";
  for (let dir of path.split(":")) {
    if (dir === "") dir = ".";
    const full = dir === "." ? js(file) : dir.replace(/\/+$/, "") + "/" + js(file);
    if (executable(full)) {
      if (!full.startsWith("/")) return [go(full), lookErr(file, ErrDot)];
      return [go(full), null];
    }
  }
  return ["", lookErr(file, ErrNotFound)];
}

export class Cmd extends $.Native {
  Path = "";
  Args: $.Slice<string> | null = null;
  Env: $.Slice<string> | null = null;
  Dir = "";
  Stdin: any = null;
  Stdout: any = null;
  Stderr: any = null;
  // Err is why the command cannot run, as LookPath said in Command; run
  // gives it back without starting anything.
  Err: any = null;

  private run(capture: boolean): [Uint8Array, Uint8Array, any] {
    if (this.Err !== null) return [new Uint8Array(0), new Uint8Array(0), this.Err];
    if (this.Path === "") return [new Uint8Array(0), new Uint8Array(0), new execError("exec: no command")];
    try {
      const args = this.Args === null ? [] : this.Args.toArray().slice(1).map(js);
      const opts: any = { args, stdout: "piped", stderr: "piped", stdin: "null" };
      if (this.Dir !== "") opts.cwd = js(this.Dir);
      if (this.Env !== null) {
        opts.clearEnv = true;
        opts.env = Object.fromEntries(this.Env.toArray().map((kv: string) => {
          const i = kv.indexOf("=");
          return [js(kv.slice(0, i)), js(kv.slice(i + 1))];
        }));
      }
      // A bare name is the file in the current directory, as execve takes
      // it: Deno would search PATH for it, finding what LookPath refused.
      const path = this.Path.includes("/") ? this.Path : "./" + this.Path;
      const out = new g.Deno.Command(js(path), opts).outputSync();
      if (!capture) {
        if (this.Stdout !== null) this.Stdout.Write($.slice(out.stdout));
        if (this.Stderr !== null) this.Stderr.Write($.slice(out.stderr));
      }
      const err = out.success ? null : new execError("exit status " + out.code);
      return [out.stdout, out.stderr, err];
    } catch (e) {
      return [new Uint8Array(0), new Uint8Array(0), new execError("exec: " + (e instanceof globalThis.Error ? e.message : String(e)))];
    }
  }

  Run(): any {
    return this.run(false)[2];
  }

  Output(): [$.Slice<number> | null, any] {
    const [out, , err] = this.run(true);
    return [$.slice(out), err];
  }

  CombinedOutput(): [$.Slice<number> | null, any] {
    const [out, errOut, err] = this.run(true);
    const all = new Uint8Array(out.length + errOut.length);
    all.set(out);
    all.set(errOut, out.length);
    return [$.slice(all), err];
  }

  Start(): any {
    return this.Run();
  }

  Wait(): any {
    return null;
  }
}
$.nativeType("exec.Cmd", Cmd);

// Command is Go's exec.Command: a name without a slash is looked up in
// PATH, and a lookup that fails, ErrDot included, is kept in Err so the
// command refuses to run rather than running something else.
export function Command(name: string, arg: $.Slice<string> | null): Cmd {
  const c = new Cmd();
  c.Path = name;
  c.Args = $.slice([name, ...(arg === null ? [] : arg.toArray())]);
  if (!name.includes("/")) {
    const [p, err] = LookPath(name);
    if (p !== "") c.Path = p;
    if (err !== null) c.Err = err;
  }
  return c;
}
