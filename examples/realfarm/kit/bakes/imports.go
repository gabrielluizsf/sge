package bakes

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/gltf"
)

// Import is a model of the user's own, from a Sketchfab download kept in
// a folder outside the repository: its licence does not allow passing the
// files on, so they are only read from where the user has them, baked into
// sprites in the user's cache, and never copied into the repository, a
// package or an archive. Dir is the model's folder under ImportsDir, File
// the glTF in it.
type Import struct {
	ID, Dir, File string
}

// Imports are the models of the user's that the farm draws when it has
// them, each in place of what the farm draws without: the village well,
// from "Water Well" by RPSebb, in the folder of the farm diorama's
// buildings.
var Imports = []Import{
	{ID: "well", Dir: "farm/buildings/waterwell", File: "scene.gltf"},
}

// ImportsDir is the folder the user keeps their imports in: SGE_IMPORTS in
// the environment, or sge-imports in their Downloads.
func ImportsDir() string {
	if dir := os.Getenv("SGE_IMPORTS"); dir != "" {
		return dir
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, "Downloads", "sge-imports")
}

// ImportName is the sprite the import id baked at heading i is kept as.
func ImportName(id string, heading int) string {
	return fmt.Sprintf("import_%s_%d", id, heading)
}

// importTasks bake each import whose files the user has; none for those
// they do not, which the farm draws its own way.
func (b *Baker) importTasks() []Task {
	var out []Task
	for _, im := range Imports {
		if _, err := os.Stat(filepath.Join(ImportsDir(), filepath.FromSlash(im.Dir), im.File)); err != nil {
			continue
		}
		out = append(out, Task{Name: "import " + im.ID, Kind: "import", ID: im.ID, Core: true, Run: func() error { return b.importModel(im) }})
	}
	return out
}

// importModel bakes import im at every heading, the model read as
// untrusted from inside its own folder.
func (b *Baker) importModel(im Import) error {
	root, err := os.OpenRoot(filepath.Join(ImportsDir(), filepath.FromSlash(im.Dir)))
	if err != nil {
		return err
	}
	defer root.Close()
	model, err := gltf.Open(root, path.Clean(im.File), gltf.DefaultLimits)
	if err != nil {
		return fmt.Errorf("%s: %w", im.ID, err)
	}
	return b.each(len(Headings), func(i int) error {
		return b.Shelf.Put(ImportName(im.ID, i), bake.Render(model, bake.Shot{Yaw: Headings[i]}, Options()))
	})
}
