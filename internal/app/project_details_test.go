package app

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"gus/internal/gitconfig"
	"gus/internal/identity"
	"gus/internal/state"
)

func TestProjectDetailsNavigationAndRefresh(t *testing.T) {
	dir := isolatedGit(t)
	root := filepath.Join(dir, "project")
	runGit(t, dir, "init", "-q", "-b", "main", root)
	runGit(t, root, "config", "user.name", "Example")
	runGit(t, root, "config", "user.email", "example@example.org")
	runGit(t, root, "remote", "add", "origin", "https://github.com/team/project.git")
	runGit(t, root, "remote", "set-url", "--add", "--push", "origin", "git@github.com:team/project.git")
	runGit(t, root, "remote", "set-url", "--add", "--push", "origin", "git@gitee.com:team/project.git")
	queue := make(chan func(), 4)
	p := &application{
		projects: []gitconfig.Project{gitconfig.InspectProject(root)}, projectQuery: "proj",
		accounts: []identity.Account{{Label: "工作", Name: "Example", Email: "example@example.org"}},
		updateUI: func(f func()) { queue <- f },
	}
	v := ui.NewTester(p.view, panelWidth, panelHeight)
	if _, ok := v.Find("托管平台：GitHub · Gitee"); !ok || v.HasText("git@gitee.com:team/project.git") {
		t.Fatal("list should show platform stack without URLs", v.Texts())
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
				t.Fatal("Git refresh stalled")
			}
		}
	}
	click := func(label string) {
		t.Helper()
		if err := v.Click(label); err != nil {
			t.Fatal(err)
		}
		v.Frame()
	}
	// Click an identity cell, not just the repository name: the entire row opens.
	click("example@example.org")
	if p.page != "" {
		t.Fatal("single click opened details")
	}
	click("example@example.org")
	settle()
	if p.page != "project" || p.detailPath != root {
		t.Fatalf("double click did not open project: %q %q", p.page, p.detailPath)
	}
	for _, want := range []string{"项目详情", "工作", "Example (local)", "main", "GitHub · Gitee", "git@gitee.com:team/project.git"} {
		if !v.HasText(want) {
			t.Errorf("details missing %q: %q", want, v.Texts())
		}
	}
	// External edits must replace the cached snapshot while details remain open.
	runGit(t, root, "remote", "set-url", "--push", "origin", "git@gitlab.com:team/project.git", "git@gitee.com:team/project.git")
	click("刷新")
	settle()
	if !v.HasText("GitHub · GitLab") || v.HasText("git@gitee.com:team/project.git") {
		t.Fatal("stale remote details", v.Texts())
	}
	v.Key(0, ui.KeyEscape)
	v.Frame()
	if p.page != "" || p.projectQuery != "proj" {
		t.Fatal("return lost project search")
	}
	if err := v.RightClick("project"); err != nil {
		t.Fatal(err)
	}
	click("项目详情")
	settle()
	click("返回项目")
	if p.page != "" {
		t.Fatal("back button did not return")
	}
	// Platform avatars remain part of the row's double-click target.
	click("托管平台：GitHub · GitLab")
	click("托管平台：GitHub · GitLab")
	settle()
	if p.page != "project" || p.detailPath != root {
		t.Fatal("double click on stack did not open details")
	}
}

func TestRemoteLinksOpenBrowserURLs(t *testing.T) {
	p := &application{}
	for _, tc := range []struct{ address, want string }{
		{"https://github.com/team/project.git", "https://github.com/team/project"},
		{"git@gitee.com:team/project.git", "https://gitee.com/team/project"},
		{"ssh://git@ssh.github.com:443/team/project.git", "https://github.com/team/project"},
	} {
		v := ui.NewTester(func(c *ui.Context) { p.remoteAddress(c, "推送", tc.address) }, panelWidth, panelHeight)
		if err := v.Click(tc.address); err != nil {
			t.Fatal(err)
		}
		if urls := v.OpenedURLs(); len(urls) != 1 || urls[0] != tc.want {
			t.Fatalf("%s: opened %v", tc.address, urls)
		}
	}
}

func TestProjectDetailsEmptyUnavailableAndMinimumSize(t *testing.T) {
	for _, lang := range []string{"zh-CN", "en"} {
		p := &application{
			preferences: state.Preferences{Language: lang}, page: "project", detailPath: "example",
			projects: []gitconfig.Project{{Path: "example", Branch: "main"}},
		}
		v := ui.NewTester(p.view, minPanelWidth, minPanelHeight)
		for _, label := range []string{"项目详情", "返回项目", "刷新"} {
			r, ok := v.Find(p.tr(label))
			if !ok || r.W <= 0 || r.X < 0 || r.Y < 0 || r.X+r.W > minPanelWidth || r.Y+r.H > minPanelHeight {
				t.Errorf("%s %s outside panel: %+v", lang, label, r)
			}
		}
		if !v.HasText(p.tr("未配置远程仓库")) || !v.HasText(p.tr("未匹配已保存账号")) {
			t.Fatal(v.Texts())
		}
		p.projects[0].RemoteError = "remote unavailable"
		v.Frame()
		if !strings.Contains(strings.Join(v.Texts(), "\n"), "remote unavailable") || v.HasText(p.tr("未配置远程仓库")) {
			t.Fatal("remote failure shown as empty", v.Texts())
		}
		p.projects[0].Problem = "repository unavailable"
		v.Frame()
		if !v.HasText("repository unavailable") {
			t.Fatal(v.Texts())
		}
		p.projects = nil
		v.Frame()
		if !v.HasText(p.tr("项目已移出列表")) {
			t.Fatal(v.Texts())
		}
	}
}
