package installer

// macLayout makes the app a bundle, ~/Applications/<Name>.app, the files
// in Contents/MacOS, with its Info.plist and the icon as icon.icns. It
// has no uninstaller: on macOS an app goes by being put in the bin.
func macLayout(app App, p Place, icon []byte) (layout, error) {
	l := layout{system: MacOS}
	l.Dir = l.join(p.Home, "Applications", app.Name+".app")
	l.Files = "Contents/MacOS"
	l.Exe = l.join(l.Dir, l.Files, app.Exe)
	l.Inside = []made{markerFile(app, "Contents")}
	iconFile := ""
	if icon != nil {
		icns, err := icnsOf(icon)
		if err != nil {
			return layout{}, err
		}
		iconFile = "icon.icns"
		l.Inside = append(l.Inside, made{Path: "Contents/Resources/icon.icns", Data: icns, Mode: 0o644})
	}
	plist := made{Path: "Contents/Info.plist", Data: infoPlist(app, iconFile), Mode: 0o644}
	l.Inside = append(l.Inside, plist)
	return l, nil
}
