package app

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"gus/internal/gitconfig"
	"gus/internal/identity"
	"gus/internal/state"
)

func drainRefresh(t *testing.T, p *application, queue <-chan func()) {
	t.Helper()
	deadline := time.After(20 * time.Second)
	for p.busy {
		select {
		case f := <-queue:
			f()
		case <-deadline:
			t.Fatal("refresh did not finish")
		}
	}
}

func TestRefreshPublishesBeforeSlowProjectAndGlobalFinish(t *testing.T) {
	dir := t.TempDir()
	queue := make(chan func(), 32)
	p := &application{stateFile: filepath.Join(dir, "state.json"), updateUI: func(f func()) { queue <- f }}
	for i := range 12 {
		p.projects = append(p.projects, gitconfig.Project{Path: filepath.Join(dir, fmt.Sprint(i))})
	}
	first := &p.projects[0]
	slowPath := first.Path
	slow, global := make(chan struct{}), make(chan struct{})
	releaseSlow, releaseGlobal := sync.OnceFunc(func() { close(slow) }), sync.OnceFunc(func() { close(global) })
	defer releaseSlow()
	defer releaseGlobal()
	read := func(string) (identity.Resolved, error) {
		<-global
		return identity.Resolved{Name: "Global", Email: "global@example.org"}, nil
	}
	inspect := func(path string) gitconfig.Project {
		if path == slowPath {
			<-slow
		}
		return gitconfig.Project{Path: path, Current: identity.Resolved{Name: "Updated", Email: "new@example.org"}, Branch: "main"}
	}
	p.refreshWith(read, inspect)
	// A second request must not launch a competing scan.
	p.refreshWith(func(string) (identity.Resolved, error) { panic("overlapping refresh") }, func(string) gitconfig.Project { panic("overlapping refresh") })
	deadline := time.After(5 * time.Second)
	for p.projects[1].Current.Name == "" {
		select {
		case f := <-queue:
			f()
		case <-deadline:
			t.Fatal("fast result waited for blocked repository/global read")
		}
	}
	if !p.busy || !p.refreshing || first.Current.Name != "" || p.global.Name != "" || first != &p.projects[0] {
		t.Fatal("refresh cleared/replaced the list or completed before all reads")
	}
	// Shutdown can save results already received without waiting for slow reads.
	p.saveProjectCache()
	partial, err := state.Load(p.stateFile)
	if err != nil || partial.ProjectCache[1].Current.Name != "Updated" || partial.ProjectCache[0].Current.Name != "" {
		t.Fatal("partial scan could not be saved", err)
	}
	releaseSlow()
	releaseGlobal()
	drainRefresh(t, p, queue)
	if p.refreshing || p.cacheDirty || p.global.Name != "Global" || first.Current.Name != "Updated" {
		t.Fatal("refresh did not finish and persist")
	}
	data, err := state.Load(p.stateFile)
	if err != nil || !reflect.DeepEqual(data.ProjectCache, p.projects) || data.Global != p.global {
		t.Fatal("completed scan was not cached", data, err)
	}
}

func TestProjectChangesPreserveScrollAndUnchangedRows(t *testing.T) {
	p := &application{}
	for i := range 30 {
		p.projects = append(p.projects, gitconfig.Project{Path: fmt.Sprintf("repo-%02d", i), Current: identity.Resolved{Name: "Cached", Email: "old@example.org"}})
	}
	v := ui.NewTester(p.view, minPanelWidth, minPanelHeight)
	v.Scroll(100, 200, 0, 10000)
	before := p.projectScroll.Y
	last, ok := v.Find("repo-29")
	if !ok || before <= 0 || last.Y < 0 || last.Y+last.H > minPanelHeight {
		t.Fatal("did not scroll to the last project", before, last)
	}
	unchanged := p.projects[0]
	if p.updateProject(unchanged) || p.cacheDirty {
		t.Fatal("unchanged properties triggered an update")
	}
	for _, index := range []int{0, 15, 29} {
		row := p.projects[index]
		row.Current = identity.Resolved{Name: "New identity", Email: "longer@example.org", NameScope: "local", EmailScope: "local"}
		row.Remotes = []gitconfig.Remote{{Name: "origin", PushURLs: []string{"git@ssh.gitgud.io:a/b"}}}
		if !p.updateProject(row) {
			t.Fatal("changed properties ignored")
		}
		v.Frame()
		after, _ := v.Find("repo-29")
		if p.projectScroll.Y != before || after.Y != last.Y {
			t.Fatalf("property update moved the list: scroll %v -> %v, row %v -> %v", before, p.projectScroll.Y, last.Y, after.Y)
		}
	}
	// Navigating away while reads finish must also preserve the list position.
	p.page = "accounts"
	v.Frame()
	p.page = ""
	v.Frame()
	if p.projectScroll.Y != before {
		t.Fatal("returning to the list lost its scroll position")
	}
}

func TestRefreshFailuresKeepCachedPropertiesAndRecover(t *testing.T) {
	cached := gitconfig.Project{Path: "example", Current: identity.Resolved{Name: "Cached", Email: "cached@example.org"},
		Branch: "main", Remotes: []gitconfig.Remote{{Name: "origin", PushURLs: []string{"https://example.org/a/b"}}}}
	p := &application{projects: []gitconfig.Project{cached}}
	if !p.updateProject(gitconfig.Project{Path: cached.Path, Problem: "offline"}) {
		t.Fatal("failure not shown")
	}
	if p.projects[0].Current != cached.Current || p.projects[0].Branch != cached.Branch || !reflect.DeepEqual(p.projects[0].Remotes, cached.Remotes) {
		t.Fatal("temporary failure erased cached properties")
	}
	if !p.updateProject(cached) || p.projects[0].Problem != "" {
		t.Fatal("successful scan did not clear error")
	}
	partial := cached
	partial.Current.Email = "changed@example.org"
	partial.Remotes, partial.Branch = nil, ""
	partial.RemoteError, partial.BranchError = "remote failed", "branch failed"
	p.updateProject(partial)
	if p.projects[0].Current.Email != partial.Current.Email || p.projects[0].Branch != cached.Branch || !reflect.DeepEqual(p.projects[0].Remotes, cached.Remotes) {
		t.Fatal("partial result lost good data")
	}
	// A successful empty result is a real deletion, not a failed read.
	p.updateProject(gitconfig.Project{Path: cached.Path})
	if p.projects[0].Current.Name != "" || len(p.projects[0].Remotes) != 0 || p.projects[0].Branch != "" {
		t.Fatal("successful empty properties did not clear cache")
	}
	if p.updateProject(gitconfig.Project{Path: "removed"}) || len(p.projects) != 1 {
		t.Fatal("late result resurrected removed project")
	}
	queue := make(chan func(), 8)
	p.updateUI = func(f func()) { queue <- f }
	p.global = cached.Current
	p.refreshWith(func(string) (identity.Resolved, error) { return identity.Resolved{}, errors.New("unavailable") }, func(string) gitconfig.Project { return cached })
	drainRefresh(t, p, queue)
	if p.global != cached.Current || p.globalError == "" {
		t.Fatal("failed global read erased cache or hid error")
	}
}

func TestManualAndBackgroundRefresh(t *testing.T) {
	dir := isolatedGit(t)
	runGit(t, dir, "init", "-q")
	runGit(t, dir, "config", "user.name", "Before")
	queue := make(chan func(), 8)
	p := &application{projects: []gitconfig.Project{gitconfig.InspectProject(dir)}, updateUI: func(f func()) { queue <- f }}
	v := ui.NewTester(p.view, minPanelWidth, minPanelHeight)
	runGit(t, dir, "config", "user.name", "After")
	if err := v.Click("刷新"); err != nil {
		t.Fatal(err)
	}
	drainRefresh(t, p, queue)
	v.Frame()
	if !v.HasText("After") {
		t.Fatal("manual refresh did not pick up external edit")
	}
	runGit(t, dir, "config", "user.name", "Hourly")
	if projectRefreshInterval != time.Hour {
		t.Fatal("wrong refresh interval")
	}
	// Exercise the real timer callback without waiting an hour or opening a window.
	p.scheduleRefresh(time.Millisecond)
	select {
	case f := <-queue:
		f()
	case <-time.After(5 * time.Second):
		t.Fatal("background timer did not dispatch")
	}
	defer func() { p.quitting = true; p.refreshTimer.Stop() }()
	drainRefresh(t, p, queue)
	if p.projects[0].Current.Name != "Hourly" {
		t.Fatal("background refresh missed edit")
	}
	p.busy = true
	p.autoRefresh()
	if p.refreshing {
		t.Fatal("timer overlapped another operation")
	}
	p.busy, p.quitting = false, true
	p.autoRefresh()
	if p.busy {
		t.Fatal("timer started after shutdown")
	}
}
