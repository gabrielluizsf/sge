// internal/godebug: settings read once from the GODEBUG environment
// variable, name=value pairs separated by commas.
import * as $ from "../../../rt/mod.ts";

function settings(): Map<string, string> {
  const g: any = globalThis;
  const raw: string = g.Deno?.env.get("GODEBUG") ?? g.process?.env.GODEBUG ?? "";
  const m = new Map<string, string>();
  for (const kv of raw.split(",")) {
    const i = kv.indexOf("=");
    if (i > 0) m.set(kv.slice(0, i), kv.slice(i + 1));
  }
  return m;
}

export class Setting extends $.Native {
  name = "";
  Name(): string {
    return this.name.replace(/^#/, "");
  }
  Undocumented(): boolean {
    return this.name.startsWith("#");
  }
  String(): string {
    return this.Name() + "=" + this.Value();
  }
  Value(): string {
    return $.str.toGo(settings().get($.str.fromGo(this.Name())) ?? "");
  }
  IncNonDefault(): void {}
}
$.nativeType("godebug.Setting", Setting);

export function New(name: string): Setting {
  const s = new Setting();
  s.name = name;
  return s;
}
