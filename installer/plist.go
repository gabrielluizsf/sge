package installer

import (
	"bytes"
	"encoding/xml"
)

// infoPlist is the bundle's Info.plist for app, its icon file in
// Resources iconFile ("" for none). Every value is escaped as XML.
func infoPlist(app App, iconFile string) []byte {
	category := "public.app-category.utilities"
	if app.Game {
		category = "public.app-category.games"
	}
	version := app.Version
	if version == "" {
		version = "1.0"
	}
	keys := [][2]string{
		{"CFBundleDevelopmentRegion", "en"},
		{"CFBundleDisplayName", app.Name},
		{"CFBundleExecutable", app.Exe},
		{"CFBundleIdentifier", app.ID},
		{"CFBundleInfoDictionaryVersion", "6.0"},
		{"CFBundleName", app.Name},
		{"CFBundlePackageType", "APPL"},
		{"CFBundleShortVersionString", version},
		{"CFBundleVersion", version},
		{"LSApplicationCategoryType", category},
	}
	if iconFile != "" {
		keys = append(keys, [2]string{"CFBundleIconFile", iconFile})
	}
	var b bytes.Buffer
	b.WriteString(xml.Header)
	b.WriteString(`<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">` + "\n")
	b.WriteString("<plist version=\"1.0\">\n<dict>\n")
	for _, kv := range keys {
		b.WriteString("\t<key>")
		xml.EscapeText(&b, []byte(kv[0]))
		b.WriteString("</key>\n\t<string>")
		xml.EscapeText(&b, []byte(kv[1]))
		b.WriteString("</string>\n")
	}
	b.WriteString("\t<key>NSHighResolutionCapable</key>\n\t<true/>\n</dict>\n</plist>\n")
	return b.Bytes()
}
