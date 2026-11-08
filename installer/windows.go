package installer

import "fmt"

// windowsLayout puts the files in %LOCALAPPDATA%\Programs\<ID>, as
// programs installed for one user are, with the icon as icon.ico; the
// Start menu's shortcut and the installed apps' entry are made by the
// system's own shell (see windowsScript).
func windowsLayout(app App, p Place, icon []byte) (layout, error) {
	if p.Roaming == "" || p.SystemRoot == "" {
		return layout{}, fmt.Errorf("installer: no %%APPDATA%% or %%SystemRoot%% on %q", p.System)
	}
	l := layout{system: Windows}
	l.Dir = l.join(p.Data, "Programs", app.ID)
	l.Exe = l.join(l.Dir, app.Exe)
	l.Uninstaller = "uninstall.exe"
	l.Menu = l.join(p.Roaming, "Microsoft", "Windows", "Start Menu", "Programs")
	l.Shortcut = l.join(l.Menu, app.Name+".lnk")
	l.Key = `HKCU:\Software\Microsoft\Windows\CurrentVersion\Uninstall\` + app.ID
	l.Inside = []made{markerFile(app, "")}
	l.IconLoc = l.Exe + ",0"
	if icon != nil {
		ico, err := icoOf(icon)
		if err != nil {
			return layout{}, err
		}
		l.Inside = append(l.Inside, made{Path: "icon.ico", Data: ico, Mode: 0o644})
		l.IconLoc = l.join(l.Dir, "icon.ico")
	}
	return l, nil
}
