package installer

import (
	"fmt"
	"path"
	"strings"
)

// layout is where an app goes on a place: Dir is its own folder, made
// whole and taken away whole; Files the folder inside it that the files
// go into ("" for Dir itself); Exe the program's full path; Inside the
// files of the system's own made inside Dir, by path in it, the marker
// first; Outside the ones made elsewhere, by full path; Link a program
// link to Exe (Linux); Uninstaller the path in Dir of the program that
// takes it away ("" for none); Menu the Start menu's folder, Shortcut
// the link in it and Key the installed apps' entry (Windows).
type layout struct {
	system              string
	Dir, Files, Exe     string
	Inside, Outside     []made
	Link, Uninstaller   string
	Menu, Shortcut, Key string
	IconLoc             string
}

// made is one file the installer writes itself, and its mode.
type made struct {
	Path string
	Data []byte
	Mode uint32
}

// marker is the file that tells a folder is an installed app's, by ID.
const marker = ".sge-installer"

// layoutOf is where app goes on p, its icon data icon (a PNG or nil).
func layoutOf(app App, p Place, icon []byte) (layout, error) {
	if err := app.check(); err != nil {
		return layout{}, err
	}
	if p.Home == "" || strings.ContainsAny(p.Home+p.Data+p.Roaming+p.SystemRoot, "\x00\n\r") {
		return layout{}, fmt.Errorf("installer: no usable home on %q", p.System)
	}
	switch p.System {
	case Linux:
		return linuxLayout(app, p), nil
	case Windows:
		return windowsLayout(app, p, icon)
	case MacOS:
		return macLayout(app, p, icon)
	}
	return layout{}, fmt.Errorf("installer: no installing on %q", p.System)
}

// join joins a path on the layout's system, with its own separator.
func (l layout) join(elem ...string) string {
	return joinOn(l.system, elem...)
}

func joinOn(system string, elem ...string) string {
	if system != Windows {
		return path.Join(elem...)
	}
	parts := make([]string, 0, len(elem))
	for _, e := range elem {
		if e = strings.Trim(strings.ReplaceAll(e, "/", `\`), `\`); e != "" {
			parts = append(parts, e)
		}
	}
	return strings.Join(parts, `\`)
}

// markerFile is the marker for app, at dir in the app's folder: its ID.
func markerFile(app App, dir string) made {
	return made{Path: path.Join(dir, marker), Data: []byte(app.ID + "\n"), Mode: 0o644}
}
