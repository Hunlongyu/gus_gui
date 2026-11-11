package main

import (
	"github.com/egoist/mygo/ui"
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
	if !samePath(real, alias) || !samePath(alias, real) || samePath(alias, other) {
		t.Fatal("directory aliases must match their target, but not other repositories")
	}
	if samePath(filepath.Join(dir, "missing-a"), filepath.Join(dir, "missing-b")) {
		t.Fatal("missing paths must not compare equal")
	}
}

func TestSavedAccountsAndDirectoriesAfterRestart(t *testing.T) {
	dir := isolatedGit(t)
	t.Setenv("APPDATA", dir)
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("HOME", dir)
	p := &application{}
	p.restore()
	if p.storageError != "" {
		t.Fatal(p.storageError)
	}
	v := ui.NewTester(p.view, panelWidth, panelHeight)
	click := func(s string) {
		t.Helper()
		if err := v.Click(s); err != nil {
			t.Fatal(err)
		}
		v.Frame()
	}
	click("账号管理")
	click("＋ 新增账号")
	for _, field := range []struct{ key, value string }{{"账号标签", "个人"}, {"名字", "Tester"}, {"邮箱", "tester@example.org"}} {
		click(field.key)
		v.Type(field.value)
	}
	click("保存账号")
	if len(p.accounts) != 1 {
		t.Fatal(p.formError)
	}
	runGit(t, dir, "init", "-q", "项目 A")
	root, err := repositoryRoot(filepath.Join(dir, "项目 A"))
	if err != nil {
		t.Fatal(err)
	}
	if err = p.addRoots([]string{root, root}); err != nil {
		t.Fatal(err)
	}
	reopened := &application{}
	reopened.restore()
	if reopened.storageError != "" || len(reopened.projects) != 1 || reopened.projects[0].path != root || len(reopened.accounts) != 1 || reopened.accounts[0] != p.accounts[0] {
		t.Fatalf("restore failed: %+v", reopened)
	}
	reopened.beginAccount(0)
	reopened.draft.Name = "Edited"
	reopened.saveAccount()
	reopened.removeDirectory(root)
	third := &application{}
	third.restore()
	if len(third.projects) != 0 || third.accounts[0].Name != "Edited" {
		t.Fatal("edit/remove not saved")
	}
	third.removeAccount(0)
	if data, err := loadData(p.stateFile); err != nil || len(data.Accounts) != 0 {
		t.Fatal(data, err)
	}
	if err = os.WriteFile(p.stateFile, []byte("invalid"), 0600); err != nil {
		t.Fatal(err)
	}
	broken := &application{}
	broken.restore()
	if broken.storageError == "" || broken.persist(nil, nil) == nil {
		t.Fatal("overwrote damaged data")
	}
}
