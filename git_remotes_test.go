package main

import (
	"slices"
	"strings"
	"testing"
)

func TestProjectRemotesAndPushDestinations(t *testing.T) {
	dir := isolatedGit(t)
	runGit(t, dir, "init", "-q", "-b", "main")
	runGit(t, dir, "config", "user.name", "Example")
	runGit(t, dir, "config", "user.email", "example@example.org")
	r := inspectProject(dir)
	if r.problem != "" || r.remoteError != "" || r.branchError != "" || len(r.remotes) != 0 || r.branch != "main" || r.detached {
		t.Fatalf("new repository: %+v", r)
	}
	runGit(t, dir, "remote", "add", "origin", "https://github.com/team/project.git")
	runGit(t, dir, "remote", "set-url", "--add", "--push", "origin", "git@github.com:team/project.git")
	runGit(t, dir, "remote", "set-url", "--add", "--push", "origin", "git@gitee.com:team/project.git")
	runGit(t, dir, "remote", "add", "backup", "https://gitlab.com/team/project.git")
	runGit(t, dir, "remote", "add", "rewritten", "work:team/project.git")
	runGit(t, dir, "config", "url.https://codeberg.org/.insteadOf", "work:")
	runGit(t, dir, "config", "url.ssh://git@git.example.org/.pushInsteadOf", "work:")
	r = inspectProject(dir)
	if r.problem != "" || r.remoteError != "" || len(r.remotes) != 3 {
		t.Fatalf("remote lookup: %+v", r)
	}
	for _, remote := range r.remotes {
		var fetch, push []string
		switch remote.Name {
		case "origin":
			fetch = []string{"https://github.com/team/project.git"}
			push = []string{"git@github.com:team/project.git", "git@gitee.com:team/project.git"}
		case "backup":
			fetch = []string{"https://gitlab.com/team/project.git"}
			push = fetch
		case "rewritten":
			fetch = []string{"https://codeberg.org/team/project.git"}
			push = []string{"ssh://git@git.example.org/team/project.git"}
		default:
			t.Fatal(remote.Name)
		}
		if !slices.Equal(remote.FetchURLs, fetch) || !slices.Equal(remote.PushURLs, push) {
			t.Fatalf("%s: %+v", remote.Name, remote)
		}
	}
	p := &application{}
	description := p.description(r)
	for _, want := range []string{"Example (local)", "example@example.org (local)", "GitHub", "Gitee", "GitLab", "git.example.org"} {
		if !strings.Contains(description, want) {
			t.Errorf("tooltip missing %q: %s", want, description)
		}
	}
	if strings.Contains(p.pushPlatforms(r), "Codeberg") {
		t.Fatal("fetch-only platform included in push summary")
	}
	runGit(t, dir, "-c", "commit.gpgsign=false", "commit", "-q", "--allow-empty", "-m", "initial")
	runGit(t, dir, "checkout", "-q", "--detach")
	if r = inspectProject(dir); !r.detached || r.branch == "" || r.branchError != "" {
		t.Fatalf("detached HEAD: %+v", r)
	}
}

func TestMultipleURLsFallbackAndCredentialRedaction(t *testing.T) {
	dir := isolatedGit(t)
	runGit(t, dir, "init", "-q")
	runGit(t, dir, "remote", "add", "origin", "https://user:secret@github.com/team/repo.git?token=secret#secret")
	runGit(t, dir, "config", "--add", "remote.origin.url", "ssh://git@gitee.com/team/repo.git")
	r := inspectProject(dir)
	if r.remoteError != "" || len(r.remotes) != 1 {
		t.Fatal(r)
	}
	want := []string{"https://github.com/team/repo.git", "ssh://git@gitee.com/team/repo.git"}
	if !slices.Equal(r.remotes[0].FetchURLs, want) || !slices.Equal(r.remotes[0].PushURLs, want) {
		t.Fatal(r.remotes)
	}
	if got := displayRemoteURL("ssh://git:secret@example.org/project.git"); got != "ssh://git@example.org/project.git" {
		t.Fatal(got)
	}
	if got := displayRemoteURL("https://token@github.com/team/repo.git"); got != want[0] {
		t.Fatal(got)
	}
}

func TestRemotePlatforms(t *testing.T) {
	for address, want := range map[string]string{
		"git@github.com:team/repo.git":                "GitHub",
		"ssh://git@ssh.github.com:443/team/repo.git":  "GitHub",
		"https://GITLAB.COM/team/repo.git":            "GitLab",
		"git@gitee.com:team/repo.git":                 "Gitee",
		"git@gitgud.io:team/repo.git":                 "GitGud",
		"https://www.gitgud.io/team/repo.git":         "GitGud",
		"https://gitgud.io.example.org/team/repo.git": "gitgud.io.example.org",
		"https://git.example.org/team/repo.git":       "git.example.org",
		"git@work-alias:team/repo.git":                "work-alias",
		"https://github.com.example.org/repo.git":     "github.com.example.org",
		`E:\repos\project.git`:                        "本地路径",
		"E:/repos/project.git":                        "本地路径",
		"../backup.git":                               "本地路径",
		"/srv/git/backup.git":                         "本地路径",
		"file:///srv/git/backup.git":                  "本地路径",
	} {
		if got := remotePlatform(address); got != want {
			t.Errorf("%q: got %q, want %q", address, got, want)
		}
	}
	p := &application{}
	r := project{remotes: []gitRemote{{PushURLs: []string{"git@github.com:a/b", "https://github.com/c/d", "git@gitee.com:a/b"}}}}
	if got := p.pushPlatforms(r); got != "GitHub · Gitee" {
		t.Fatal(got)
	}
}

func TestRemoteWebURLs(t *testing.T) {
	for address, want := range map[string]string{
		"https://user:token@github.com/team/repo.git?secret=token": "https://github.com/team/repo",
		"https://git.example.org:8443/team/repo.git":               "https://git.example.org:8443/team/repo",
		"git@github.com:team/repo.git":                             "https://github.com/team/repo",
		"ssh://git@altssh.gitlab.com:443/team/repo.git":            "https://gitlab.com/team/repo",
		"git@ssh.dev.azure.com:v3/team/project/repo":               "https://dev.azure.com/team/project/_git/repo",
		"ssh://git@git.example.org:2222/team/repo.git":             "https://git.example.org/team/repo",
		"git://git.example.org/team/repo.git":                      "https://git.example.org/team/repo",
		"git@work-alias:team/repo.git":                             "",
		`E:\repos\project.git`:                                     "",
		"../backup.git":                                            "",
		"file:///srv/git/backup.git":                               "",
		"ext::some-command":                                        "",
		"javascript:alert(1)":                                      "",
		"custom://example.org/repo":                                "",
	} {
		if got := remoteWebURL(address); got != want {
			t.Errorf("%q: got %q, want %q", address, got, want)
		}
	}
}

func TestRemoteNameFallbackKeepsOtherPlatforms(t *testing.T) {
	dir := isolatedGit(t)
	runGit(t, dir, "init", "-q")
	// Git treats a remote without url as a local path named after the remote.
	runGit(t, dir, "config", "remote.backup.fetch", "+refs/heads/*:refs/remotes/backup/*")
	runGit(t, dir, "remote", "add", "origin", "git@gitee.com:team/repo.git")
	r := inspectProject(dir)
	p := &application{}
	if r.problem != "" || r.remoteError != "" || p.pushPlatforms(r) != "本地路径 · Gitee" {
		t.Fatalf("local fallback hid a hosted platform: %+v", r)
	}
}
