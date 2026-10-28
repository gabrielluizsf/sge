// What a desktop tells and shows a program outside its window, for the
// SDL window on Deno: notifications, and whether the system's theme is
// dark, through the system's own tools, as SDL2 has neither. Deno runs
// them with --allow-run; without it, or without the tool, there is none.
// Titles and bodies go to the tools as arguments, never into a script.

const D: any = (globalThis as any).Deno;

// notify shows a notification, giving back whether one was sent off.
export function notify(title: string, body: string): boolean {
  if (typeof D?.Command !== "function") return false;
  const [cmd, args, env] = notifier(title, body);
  try {
    new D.Command(cmd, { args, env, stdin: "null", stdout: "null", stderr: "null" }).spawn();
    return true;
  } catch {
    return false;
  }
}

function notifier(title: string, body: string): [string, string[], Record<string, string>] {
  switch (D.build.os) {
    case "darwin":
      return ["osascript", ["-e", "on run argv", "-e", "display notification (item 2 of argv) with title (item 1 of argv)", "-e", "end run", title, body], {}];
    case "windows":
      // The text reaches the script through the environment.
      return ["powershell", ["-NoProfile", "-NonInteractive", "-Command", toast], { GOTS_TITLE: title, GOTS_BODY: body }];
  }
  return ["notify-send", ["--", title, body], {}];
}

const toast = [
  "[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null",
  "$x = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastText02)",
  "$t = $x.GetElementsByTagName('text')",
  "$t.Item(0).AppendChild($x.CreateTextNode($env:GOTS_TITLE)) | Out-Null",
  "$t.Item(1).AppendChild($x.CreateTextNode($env:GOTS_BODY)) | Out-Null",
  "[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier('gots').Show([Windows.UI.Notifications.ToastNotification]::new($x))",
].join("; ");

// The theme, asked at most every few seconds: asking runs a program.
let theme: [boolean, boolean] = [false, false];
let asked = -Infinity;

// systemDark is whether the system's theme is dark, and whether that is
// known.
export function systemDark(): [boolean, boolean] {
  const now = performance.now();
  if (now - asked < 3000) return theme;
  asked = now;
  theme = askTheme();
  return theme;
}

function askTheme(): [boolean, boolean] {
  if (typeof D?.Command !== "function") return [false, false];
  const run = (cmd: string, args: string[]): string | null => {
    try {
      const o = new D.Command(cmd, { args, stdin: "null", stderr: "null" }).outputSync();
      return o.success ? new TextDecoder().decode(o.stdout) : null;
    } catch {
      return null;
    }
  };
  switch (D.build.os) {
    case "darwin": {
      // No such default means light.
      const o = run("defaults", ["read", "-g", "AppleInterfaceStyle"]);
      return [o !== null && o.includes("Dark"), true];
    }
    case "windows": {
      const o = run("reg", ["query", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Themes\\Personalize", "/v", "AppsUseLightTheme"]);
      return o === null ? [false, false] : [/0x0\b/.test(o), true];
    }
  }
  const scheme = run("gsettings", ["get", "org.gnome.desktop.interface", "color-scheme"]);
  if (scheme !== null && scheme.includes("prefer-dark")) return [true, true];
  const gtk = run("gsettings", ["get", "org.gnome.desktop.interface", "gtk-theme"]);
  if (gtk !== null) return [/dark/i.test(gtk), true];
  return scheme === null ? [false, false] : [false, true];
}
