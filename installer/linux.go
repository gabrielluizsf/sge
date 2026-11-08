package installer

// linuxLayout puts the files in Data/<ID>, the desktop entry in
// Data/applications/<ID>.desktop and a link in ~/.local/bin/<ID>, as the
// XDG specs have it; the entry names the icon by its full path.
func linuxLayout(app App, p Place) layout {
	l := layout{system: Linux}
	l.Dir = l.join(p.Data, app.ID)
	l.Exe = l.join(l.Dir, app.Exe)
	l.Uninstaller = "uninstall"
	l.Link = l.join(p.Home, ".local", "bin", app.ID)
	iconPath := ""
	if app.Icon != "" {
		iconPath = l.join(l.Dir, app.Icon)
	}
	entry := desktopEntry(app, l.Exe, l.Dir, iconPath)
	l.Inside = []made{markerFile(app, "")}
	l.Outside = []made{{Path: l.join(p.Data, "applications", app.ID+".desktop"), Data: []byte(entry), Mode: 0o644}}
	return l
}
