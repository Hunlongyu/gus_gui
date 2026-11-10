package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestDropDirectoriesOnProjectPanel(t *testing.T) {
	dir := isolatedGit(t)
	first := filepath.Join(dir, "项目 one")
	second := filepath.Join(dir, "two")
	sub := filepath.Join(first, "src")
	runGit(t, dir, "init", "-q", first)
	runGit(t, dir, "init", "-q", second)
	if err := os.Mkdir(sub, 0700); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(sub, "note.txt")
	if err := os.WriteFile(file, nil, 0600); err != nil {
		t.Fatal(err)
	}
	runGit(t, first, "config", "user.email", "drop@example.org")
	queue := make(chan func(), 4)
	p := &application{stateFile: filepath.Join(dir, "state.json"), updateUI: func(f func()) { queue <- f }}
	v := ui.NewTester(p.view, panelWidth, panelHeight)
	drop := func(paths ...string) {
		t.Helper()
		if !v.DropFiles(90, 155, paths...) {
			t.Fatal("project panel rejected native drop")
		}
		v.Frame()
	}
	settle := func() {
		t.Helper()
		deadline := time.After(20 * time.Second)
		for p.busy {
			select {
			case f := <-queue:
				f()
				v.Frame()
			case <-deadline:
				t.Fatal("drop did not complete")
			}
		}
	}
	p.refresh()
	drop(sub, first)
	drop(second, file, dir)
	settle()
	if len(p.projects) != 2 || !samePath(p.projects[0].path, first) || !samePath(p.projects[1].path, second) {
		t.Fatal(p.projects)
	}
	if p.projects[0].current.Email != "drop@example.org" {
		t.Fatal(p.projects[0])
	}
	if !strings.Contains(p.notice, file) || !strings.Contains(p.notice, "无法打开工作区") {
		t.Fatal(p.notice)
	}
	drop(first)
	settle()
	data, err := loadData(p.stateFile)
	if err != nil || len(data.Directories) != 2 {
		t.Fatal(data, err)
	}
	p.page = "settings"
	v.Frame()
	if v.DropFiles(90, 155, dir) {
		t.Fatal("settings page accepted project drop")
	}
}
