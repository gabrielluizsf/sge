package human

import (
	"bufio"
	"bytes"
	"strings"
)

// Licence is the licence the proxy or material in file name says it is
// under, as MakeHuman's files say it in a comment: "CC0" for those it
// released as CC0 in 2020, as their header says; "" for one that says
// nothing. A pack's name is no promise: some of a "cc0" pack's files
// say AGPL3 or CC BY, or nothing.
func (k *Kit) Licence(name string) (string, error) {
	data, err := k.lim.read(k.fsys, name)
	if err != nil {
		return "", err
	}
	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {
		line := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(sc.Text()), "#"))
		if strings.Contains(line, "explicitly released as CC0") {
			return "CC0", nil
		}
		f := strings.Fields(line)
		if len(f) > 1 && strings.EqualFold(strings.TrimSuffix(f[0], ":"), "license") {
			return normalLicence(strings.Join(f[1:], " ")), nil
		}
	}
	return "", sc.Err()
}

// normalLicence is CC0 however it is written; any other as it is.
func normalLicence(l string) string {
	switch strings.ToUpper(strings.NewReplacer("-", "", "_", "", " ", "").Replace(l)) {
	case "CC0", "CC01.0", "CC0V1.0":
		return "CC0"
	}
	return l
}

// Licences are the licences of every file a person of body b, looking as
// l, is made of that says one: their skin, or the skins blended for b,
// their eyes' colour, and each proxy they wear, by name.
func (k *Kit) Licences(b Body, l Look) (map[string]string, error) {
	names := append([]string(nil), l.Wear...)
	if l.Eyes != "" {
		names = append(names, l.Eyes)
	}
	if l.Skin != "" {
		names = append(names, l.Skin)
	} else {
		for _, s := range skins(b) {
			names = append(names, s.name)
		}
	}
	out := map[string]string{}
	for _, name := range names {
		lic, err := k.Licence(name)
		if err != nil {
			return nil, err
		}
		out[name] = lic
	}
	return out, nil
}
