package human

import "strings"

// hasFolder is whether file name is in the kit's folder.
func hasFolder(name, folder string) bool { return strings.HasPrefix(name, folder+"/") }
