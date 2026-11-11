package app

import (
	"slices"
	"strings"
	"testing"

	"gus/internal/gitconfig"
)

func TestProjectRemotesAndPushDestinations(t *testing.T) {
	dir := isolatedGit(t)
	runGit(t, dir, "init", "-q", "-b", "main")
	runGit(t, dir, "config", "user.name", "Example")
	runGit(t, dir, "config", "user.email", "example@example.org")
	r := gitconfig.InspectProject(dir)
	if r.Problem != "" || r.RemoteError != "" || r.BranchError != "" || len(r.Remotes) != 0 || r.Branch != "main" || r.Detached {
		t.Fatalf("new repository: %+v", r)
	}
	runGit(t, dir, "remote", "add", "origin", "https://github.com/team/project.git")
	runGit(t, dir, "remote", "set-url", "--add", "--push", "origin", "git@github.com:team/project.git")
	runGit(t, dir, "remote", "set-url", "--add", "--push", "origin", "git@gitee.com:team/project.git")
	runGit(t, dir, "remote", "add", "backup", "https://gitlab.com/team/project.git")
	runGit(t, dir, "remote", "add", "rewritten", "work:team/project.git")
	runGit(t, dir, "config", "url.https://codeberg.org/.insteadOf", "work:")
	runGit(t, dir, "config", "url.ssh://git@git.example.org/.pushInsteadOf", "work:")
	r = gitconfig.InspectProject(dir)
	if r.Problem != "" || r.RemoteError != "" || len(r.Remotes) != 3 {
		t.Fatalf("remote lookup: %+v", r)
	}
	for _, remote := range r.Remotes {
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
	if r = gitconfig.InspectProject(dir); !r.Detached || r.Branch == "" || r.BranchError != "" {
		t.Fatalf("detached HEAD: %+v", r)
	}
}

func TestPushPlatformsDeduplicatesHosts(t *testing.T) {
	p := &application{}
	r := gitconfig.Project{Remotes: []gitconfig.Remote{{PushURLs: []string{"git@github.com:a/b", "https://github.com/c/d", "git@gitee.com:a/b"}}}}
	if got := p.pushPlatforms(r); got != "GitHub · Gitee" {
		t.Fatal(got)
	}
}

func TestRemoteNameFallbackKeepsOtherPlatforms(t *testing.T) {
	dir := isolatedGit(t)
	runGit(t, dir, "init", "-q")
	// Git treats a remote without url as a local path named after the remote.
	runGit(t, dir, "config", "remote.backup.fetch", "+refs/heads/*:refs/remotes/backup/*")
	runGit(t, dir, "remote", "add", "origin", "git@gitee.com:team/repo.git")
	r := gitconfig.InspectProject(dir)
	p := &application{}
	if r.Problem != "" || r.RemoteError != "" || p.pushPlatforms(r) != "本地路径 · Gitee" {
		t.Fatalf("local fallback hid a hosted platform: %+v", r)
	}
}
