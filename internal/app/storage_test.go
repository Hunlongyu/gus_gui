package app

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/egoist/mygo/ui"

	"gus/internal/gitconfig"
	"gus/internal/identity"
	"gus/internal/state"
)

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
	root, err := gitconfig.RepositoryRoot(filepath.Join(dir, "项目 A"))
	if err != nil {
		t.Fatal(err)
	}
	if err = p.addRoots([]string{root, root}); err != nil {
		t.Fatal(err)
	}
	reopened := &application{}
	reopened.restore()
	if reopened.storageError != "" || len(reopened.projects) != 1 || reopened.projects[0].Path != root || len(reopened.accounts) != 1 || reopened.accounts[0] != p.accounts[0] {
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
	if data, err := state.Load(p.stateFile); err != nil || len(data.Accounts) != 0 {
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

func TestProjectCacheRestoresBeforeAnyGitRead(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("APPDATA", dir)
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("HOME", dir)
	p := &application{}
	p.restore()
	root := filepath.Join(dir, "cached-project")
	row := gitconfig.Project{Path: root, Current: identity.Resolved{Name: "Tester", Email: "test@example.org", NameScope: "global", EmailScope: "local"},
		Branch: "feature", Detached: true, Remotes: []gitconfig.Remote{{Name: "origin", FetchURLs: []string{"https://gitgud.io/team/repo.git"}, PushURLs: []string{"git@ssh.gitgud.io:team/repo.git"}}}}
	p.global = identity.Resolved{Name: "Global", Email: "global@example.org", NameScope: "global", EmailScope: "global"}
	if err := p.persist(nil, []gitconfig.Project{row}); err != nil {
		t.Fatal(err)
	}
	reopened := &application{}
	reopened.restore()
	if reopened.storageError != "" || !reflect.DeepEqual(reopened.projects, []gitconfig.Project{row}) || reopened.global != p.global {
		t.Fatal("cached attributes were not restored", reopened)
	}
	v := ui.NewTester(reopened.view, panelWidth, panelHeight)
	if !v.HasText("Tester") || !v.HasText("test@example.org") || !v.HasText("Global") {
		t.Fatal("first frame did not show cached identities", v.Texts())
	}
	if _, ok := v.Find("推送平台：GitGud"); !ok {
		t.Fatal("first frame did not show cached platform")
	}
	// Old state files remain valid; cache entries cannot add unlisted projects.
	if err := state.Store(p.stateFile, state.Data{Directories: []string{root}, ProjectCache: []gitconfig.Project{{Path: filepath.Join(dir, "removed")}}}); err != nil {
		t.Fatal(err)
	}
	legacy := &application{}
	legacy.restore()
	if legacy.storageError != "" || len(legacy.projects) != 1 || legacy.projects[0].Path != root || legacy.projects[0].Current.Name != "" {
		t.Fatal("cache changed the authoritative directory list")
	}
	if err := state.Store(p.stateFile, state.Data{Directories: []string{root}}); err != nil {
		t.Fatal(err)
	}
	legacy = &application{}
	legacy.restore()
	if legacy.storageError != "" || len(legacy.projects) != 1 || legacy.projects[0].Path != root {
		t.Fatal("old directory-only state failed to restore")
	}
}
