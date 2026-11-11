package app

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"gus/internal/gitconfig"
	"gus/internal/identity"
)

func TestSwitchFromTrayPanelMenus(t *testing.T) {
	dir := isolatedGit(t)
	root := filepath.Join(dir, "project")
	runGit(t, dir, "init", "-q", root)
	defaultID := identity.Account{Label: "个人", Name: "Test User", Email: "test@example.org"}
	work := identity.Account{Label: "工作", Name: "Work User", Email: "work@example.org"}
	queue := make(chan func(), 4)
	p := &application{accounts: []identity.Account{defaultID, work}, projects: []gitconfig.Project{{Path: root}}, updateUI: func(f func()) { queue <- f }}
	v := ui.NewTester(p.view, panelWidth, panelHeight)
	settle := func() {
		t.Helper()
		deadline := time.After(20 * time.Second)
		for p.busy {
			select {
			case f := <-queue:
				f()
				v.Frame()
			case <-deadline:
				t.Fatal("Git operation stalled")
			}
		}
		if p.notice != "" || p.globalError != "" {
			t.Fatal(p.notice, p.globalError)
		}
	}
	choose := func(row, option string) {
		t.Helper()
		if err := v.RightClick(row); err != nil {
			t.Fatal(err)
		}
		v.Frame()
		if err := v.Click(option); err != nil {
			t.Fatal(err)
		}
		v.Frame()
		settle()
	}
	p.refresh()
	settle()
	choose("全局", "个人 · Test User")
	if !p.global.Matches(defaultID) || !p.projects[0].Current.Matches(defaultID) {
		t.Fatal("global switch not reflected in panel")
	}
	choose("project", "工作 · Work User")
	if runGit(t, root, "config", "--local", "user.email") != work.Email+"\n" {
		t.Fatal("project config missing")
	}
	choose("project", "继承全局")
	if !p.projects[0].Current.Matches(defaultID) || p.projects[0].Current.Source() != "继承" {
		t.Fatal("inheritance not restored")
	}
	runGit(t, root, "config", "--local", "user.email", "external@example.org")
	p.refresh()
	settle()
	if p.projects[0].Current.Source() != "混合" || p.projects[0].Current.Email != "external@example.org" {
		t.Fatal("external edit not read")
	}
}
