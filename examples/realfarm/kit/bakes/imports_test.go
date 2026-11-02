package bakes

import (
	"os"
	"path/filepath"
	"testing"
)

// With none of the user's imports in their folder there is nothing to
// bake, and the farm draws its own; with one there, it is baked, read from
// where it is, nothing copied.
func TestImportsAreBakedOnlyWhereTheUserHasThem(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("SGE_IMPORTS", dir)
	b := &Baker{}
	if n := len(b.importTasks()); n != 0 {
		t.Fatalf("%d imports baked from an empty folder", n)
	}
	well := filepath.Join(dir, "farm", "buildings", "waterwell")
	if err := os.MkdirAll(well, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(well, "scene.gltf"), []byte("not a model"), 0o600); err != nil {
		t.Fatal(err)
	}
	tasks := b.importTasks()
	if len(tasks) != 1 || tasks[0].ID != "well" || tasks[0].Kind != "import" {
		t.Fatalf("tasks %+v", tasks)
	}
	if err := tasks[0].Run(); err == nil {
		t.Error("a file that is no model was baked")
	}
	if got := ImportsDir(); got != dir {
		t.Errorf("imports dir %q, not %q", got, dir)
	}
}
