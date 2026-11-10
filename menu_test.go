package main

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestSwitchFromTrayPanelMenus(t *testing.T) {
	dir := isolatedGit(t)
	root := filepath.Join(dir, "project")
	runGit(t, dir, "init", "-q", root)
	defaultID := identity{"个人", "Test User", "test@example.org"}
	work := identity{"工作", "Work User", "work@example.org"}
	queue := make(chan func(), 4)
	p := &application{accounts: []identity{defaultID, work}, projects: []project{{path: root}}, updateUI: func(f func()) { queue <- f }}
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
	if !p.global.matches(defaultID) || !p.projects[0].current.matches(defaultID) {
		t.Fatal("global switch not reflected in panel")
	}
	choose("project", "工作 · Work User")
	if runGit(t, root, "config", "--local", "user.email") != work.Email+"\n" {
		t.Fatal("project config missing")
	}
	choose("project", "继承全局")
	if !p.projects[0].current.matches(defaultID) || p.projects[0].current.source() != "继承" {
		t.Fatal("inheritance not restored")
	}
	runGit(t, root, "config", "--local", "user.email", "external@example.org")
	p.refresh()
	settle()
	if p.projects[0].current.source() != "混合" || p.projects[0].current.Email != "external@example.org" {
		t.Fatal("external edit not read")
	}
}
