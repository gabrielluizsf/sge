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
  static $type = $.structType("exec.Error", true, execError, null, () => [$.field("s", $.String)]);
}

export function LookPath(file: string): [string, any] {
  if (file.includes("/")) return [file, null];
  const path = g.Deno?.env.get("PATH") ?? "";
  for (const dir of path.split(":")) {
    const full = (dir === "" ? "." : dir) + "/" + js(file);
    try {
      const info = g.Deno.statSync(full);
      if (info.isFile && ((info.mode ?? 0o111) & 0o111) !== 0) return [go(full), null];
    } catch {
      // not here
    }
  }
  return ["", new execError(`exec: "${file}": executable file not found in $PATH`)];
}

export class Cmd extends $.Native {
  Path = "";
  Args: $.Slice<string> | null = null;
  Env: $.Slice<string> | null = null;
  Dir = "";
  Stdin: any = null;
  Stdout: any = null;
  Stderr: any = null;

  private run(capture: boolean): [Uint8Array, Uint8Array, any] {
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
      const out = new g.Deno.Command(js(this.Path), opts).outputSync();
      if (!capture) {
        if (this.Stdout !== null) this.Stdout.Write($.slice(out.stdout));
        if (this.Stderr !== null) this.Stderr.Write($.slice(out.stderr));
      }
      const err = out.success ? null : new execError("exit status " + out.code);
      return [out.stdout, out.stderr, err];
    } catch (e) {
      return [new Uint8Array(0), new Uint8Array(0), new execError("exec: " + (e instanceof Error ? e.message : String(e)))];
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

export function Command(name: string, arg: $.Slice<string> | null): Cmd {
  const c = new Cmd();
  const [p, err] = LookPath(name);
  c.Path = err === null ? p : name;
  c.Args = $.slice([name, ...(arg === null ? [] : arg.toArray())]);
  return c;
}
