package main

import (
	"os"
	"path/filepath"
	"testing"
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
	out, err := git(dir, args...)
	if err != nil {
		t.Fatal(err)
	}
	return out
}
func apply(t *testing.T, dir string, a *identity) {
	t.Helper()
	if err := switchIdentity(dir, a); err != nil {
		t.Fatal(err)
	}
}

func TestGitScopesAndSwitching(t *testing.T) {
	dir := isolatedGit(t)
	repo := filepath.Join(dir, "真实 项目")
	runGit(t, dir, "init", "-q", repo)
	other := filepath.Join(dir, "another")
	runGit(t, dir, "init", "-q", other)
	base := identity{"默认", "Lin 王", "lin@example.org"}
	work := identity{"工作", `Lin "Dev" Chen`, "lin+work@example.org"}
	global, _ := readIdentity("")
	if global != (gitIdentity{}) {
		t.Fatal(global)
	}
	apply(t, "", &base)
	if r := inspectProject(repo); r.problem != "" || !r.current.matches(base) || r.current.source() != "继承" {
		t.Fatal(r)
	}
	runGit(t, repo, "config", "--local", "user.email", work.Email)
	r := inspectProject(repo)
	if r.current.Name != base.Name || r.current.Email != work.Email || r.current.source() != "混合" {
		t.Fatal(r)
	}
	runGit(t, repo, "config", "user.signingkey", "keep-me")
	apply(t, repo, &work)
	if r := inspectProject(repo); !r.current.matches(work) || r.current.source() != "专用" {
		t.Fatal(r)
	}
	if r := inspectProject(other); !r.current.matches(base) {
		t.Fatal(r)
	}
	if v, err := readIdentity(""); err != nil || !v.matches(base) {
		t.Fatal(v, err)
	}
	apply(t, repo, nil)
	if r := inspectProject(repo); !r.current.matches(base) || r.current.source() != "继承" {
		t.Fatal(r)
	}
	if runGit(t, repo, "config", "--local", "--get", "user.signingkey") != "keep-me\n" {
		t.Fatal("unrelated setting changed")
	}
	if _, err := git(repo, "config", "--local", "--get", "user.name"); !exitCode(err, 1) {
		t.Fatal(err)
	}
	apply(t, "", &work)
	if r := inspectProject(repo); !r.current.matches(work) || r.current.source() != "继承" {
		t.Fatal(r)
	}
}

func TestLockedAndInvalidRepositories(t *testing.T) {
	dir := isolatedGit(t)
	runGit(t, dir, "init", "-q")
	before := identity{"old", "First", "first@example.org"}
	after := identity{"new", "-Second", "second@example.org"}
	apply(t, dir, &before)
	lock := filepath.Join(dir, ".git", "config.lock")
	if err := os.WriteFile(lock, []byte("other writer"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := switchIdentity(dir, &after); err == nil {
		t.Fatal("accepted locked repository")
	}
	if r := inspectProject(dir); !r.current.matches(before) {
		t.Fatal(r)
	}
	b, err := os.ReadFile(lock)
	if err != nil || string(b) != "other writer" {
		t.Fatal(err, string(b))
	}
	if err = os.Remove(lock); err != nil {
		t.Fatal(err)
	}
	apply(t, dir, &after)
	sub := filepath.Join(dir, "nested")
	if err := os.Mkdir(sub, 0700); err != nil {
		t.Fatal(err)
	}
	if root, err := repositoryRoot(sub); err != nil || !samePath(root, dir) {
		t.Fatal(root, err)
	}
	if err := switchIdentity(sub, &before); err == nil {
		t.Fatal("wrote through stale nested repository path")
	}
	if _, err := repositoryRoot(t.TempDir()); err == nil {
		t.Fatal("accepted non-repository")
	}
}

func TestWorktreeAndConditionalIncludes(t *testing.T) {
	dir := isolatedGit(t)
	repo := filepath.Join(dir, "main")
	runGit(t, dir, "init", "-q", repo)
	base := identity{"base", "Main", "main@example.org"}
	local := identity{"local", "Worktree", "wt@example.org"}
	apply(t, "", &base)
	runGit(t, repo, "-c", "commit.gpgsign=false", "commit", "-q", "--allow-empty", "-m", "initial")
	tree := filepath.Join(dir, "worktree")
	runGit(t, repo, "worktree", "add", "-q", "--detach", tree)
	runGit(t, repo, "config", "extensions.worktreeConfig", "true")
	runGit(t, tree, "config", "--worktree", "user.name", "Original")
	apply(t, tree, &local)
	if r := inspectProject(repo); !r.current.matches(base) {
		t.Fatal(r)
	}
	if r := inspectProject(tree); !r.current.matches(local) || r.current.NameScope != "worktree" {
		t.Fatal(r)
	}
	apply(t, tree, nil)
	if r := inspectProject(tree); !r.current.matches(base) {
		t.Fatal(r)
	}
	include := filepath.Join(dir, "work.inc")
	runGit(t, dir, "config", "--file", include, "user.email", "included@example.org")
	runGit(t, dir, "config", "--global", "includeIf.gitdir/i:"+filepath.ToSlash(repo)+"/.git.path", filepath.ToSlash(include))
	if r := inspectProject(repo); r.current.Email != "included@example.org" {
		t.Fatal(r)
	}
}

func TestLateIncludeAndFailedPair(t *testing.T) {
	dir := isolatedGit(t)
	runGit(t, dir, "init", "-q")
	includedPath := filepath.Join(dir, "team.inc")
	runGit(t, dir, "config", "user.name", "Original")
	runGit(t, dir, "config", "user.signingkey", "retain")
	runGit(t, dir, "config", "--file", includedPath, "user.name", "Team")
	runGit(t, dir, "config", "include.path", filepath.ToSlash(includedPath))
	a := identity{"账号", `-王 "Li" \ Dev #1`, "dev@example.org"}
	apply(t, dir, &a)
	if r := inspectProject(dir); !r.current.matches(a) {
		t.Fatal(r)
	}
	if s := runGit(t, dir, "config", "--file", includedPath, "--get", "user.name"); s != "Team\n" {
		t.Fatal(s)
	}
	bad := identity{"bad", "Broken\nValue", "fail@example.org"}
	if err := switchIdentity(dir, &bad); err == nil {
		t.Fatal("invalid pair accepted")
	}
	if r := inspectProject(dir); !r.current.matches(a) {
		t.Fatal("partial update", r)
	}
}
