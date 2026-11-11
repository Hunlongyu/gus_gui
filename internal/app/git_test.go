package app

import (
	"os/exec"
	"path/filepath"
	"testing"

	"gus/internal/gitconfig"
	"gus/internal/identity"
)

func isolatedGit(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(dir, "global config"))
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_COUNT", "0")
	return dir
}

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}
func apply(t *testing.T, dir string, a *identity.Account) {
	t.Helper()
	if err := gitconfig.SwitchIdentity(dir, a); err != nil {
		t.Fatal(err)
	}
}
