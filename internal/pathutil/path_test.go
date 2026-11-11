package pathutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSamePathRecognizesDirectoryAliases(t *testing.T) {
	dir := t.TempDir()
	real := filepath.Join(dir, "repository")
	alias := filepath.Join(dir, "alias")
	other := filepath.Join(dir, "other")
	for _, path := range []string{real, other} {
		if err := os.Mkdir(path, 0700); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink(real, alias); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if !Same(real, alias) || !Same(alias, real) || Same(alias, other) {
		t.Fatal("directory aliases must match their target, but not other repositories")
	}
	if Same(filepath.Join(dir, "missing-a"), filepath.Join(dir, "missing-b")) {
		t.Fatal("missing paths must not compare equal")
	}
}
