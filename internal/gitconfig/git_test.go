package gitconfig

import (
	"os"
	"path/filepath"
	"testing"

	"gus/internal/identity"
	"gus/internal/pathutil"
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
func apply(t *testing.T, dir string, a *identity.Account) {
	t.Helper()
	if err := SwitchIdentity(dir, a); err != nil {
		t.Fatal(err)
	}
}

func TestGitScopesAndSwitching(t *testing.T) {
	dir := isolatedGit(t)
	repo := filepath.Join(dir, "真实 项目")
	runGit(t, dir, "init", "-q", repo)
	other := filepath.Join(dir, "another")
	runGit(t, dir, "init", "-q", other)
	base := identity.Account{Label: "默认", Name: "Lin 王", Email: "lin@example.org"}
	work := identity.Account{Label: "工作", Name: `Lin "Dev" Chen`, Email: "lin+work@example.org"}
	global, _ := ReadIdentity("")
	if global != (identity.Resolved{}) {
		t.Fatal(global)
	}
	apply(t, "", &base)
	if r := InspectProject(repo); r.Problem != "" || !r.Current.Matches(base) || r.Current.Source() != "继承" {
		t.Fatal(r)
	}
	runGit(t, repo, "config", "--local", "user.email", work.Email)
	r := InspectProject(repo)
	if r.Current.Name != base.Name || r.Current.Email != work.Email || r.Current.Source() != "混合" {
		t.Fatal(r)
	}
	runGit(t, repo, "config", "user.signingkey", "keep-me")
	apply(t, repo, &work)
	if r := InspectProject(repo); !r.Current.Matches(work) || r.Current.Source() != "专用" {
		t.Fatal(r)
	}
	if r := InspectProject(other); !r.Current.Matches(base) {
		t.Fatal(r)
	}
	if v, err := ReadIdentity(""); err != nil || !v.Matches(base) {
		t.Fatal(v, err)
	}
	apply(t, repo, nil)
	if r := InspectProject(repo); !r.Current.Matches(base) || r.Current.Source() != "继承" {
		t.Fatal(r)
	}
	if runGit(t, repo, "config", "--local", "--get", "user.signingkey") != "keep-me\n" {
		t.Fatal("unrelated setting changed")
	}
	if _, err := git(repo, "config", "--local", "--get", "user.name"); !exitCode(err, 1) {
		t.Fatal(err)
	}
	apply(t, "", &work)
	if r := InspectProject(repo); !r.Current.Matches(work) || r.Current.Source() != "继承" {
		t.Fatal(r)
	}
}

func TestLockedAndInvalidRepositories(t *testing.T) {
	dir := isolatedGit(t)
	runGit(t, dir, "init", "-q")
	before := identity.Account{Label: "old", Name: "First", Email: "first@example.org"}
	after := identity.Account{Label: "new", Name: "-Second", Email: "second@example.org"}
	apply(t, dir, &before)
	lock := filepath.Join(dir, ".git", "config.lock")
	if err := os.WriteFile(lock, []byte("other writer"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := SwitchIdentity(dir, &after); err == nil {
		t.Fatal("accepted locked repository")
	}
	if r := InspectProject(dir); !r.Current.Matches(before) {
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
	if root, err := RepositoryRoot(sub); err != nil || !pathutil.Same(root, dir) {
		t.Fatal(root, err)
	}
	if err := SwitchIdentity(sub, &before); err == nil {
		t.Fatal("wrote through stale nested repository path")
	}
	if _, err := RepositoryRoot(t.TempDir()); err == nil {
		t.Fatal("accepted non-repository")
	}
}

func TestWorktreeAndConditionalIncludes(t *testing.T) {
	dir := isolatedGit(t)
	repo := filepath.Join(dir, "main")
	runGit(t, dir, "init", "-q", repo)
	base := identity.Account{Label: "base", Name: "Main", Email: "main@example.org"}
	local := identity.Account{Label: "local", Name: "Worktree", Email: "wt@example.org"}
	apply(t, "", &base)
	runGit(t, repo, "-c", "commit.gpgsign=false", "commit", "-q", "--allow-empty", "-m", "initial")
	tree := filepath.Join(dir, "worktree")
	runGit(t, repo, "worktree", "add", "-q", "--detach", tree)
	runGit(t, repo, "config", "extensions.worktreeConfig", "true")
	runGit(t, tree, "config", "--worktree", "user.name", "Original")
	apply(t, tree, &local)
	if r := InspectProject(repo); !r.Current.Matches(base) {
		t.Fatal(r)
	}
	if r := InspectProject(tree); !r.Current.Matches(local) || r.Current.NameScope != "worktree" {
		t.Fatal(r)
	}
	apply(t, tree, nil)
	if r := InspectProject(tree); !r.Current.Matches(base) {
		t.Fatal(r)
	}
	include := filepath.Join(dir, "work.inc")
	runGit(t, dir, "config", "--file", include, "user.email", "included@example.org")
	root, err := RepositoryRoot(repo)
	if err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, "config", "--global", "includeIf.gitdir/i:"+filepath.ToSlash(root)+"/.git.path", filepath.ToSlash(include))
	if r := InspectProject(repo); r.Current.Email != "included@example.org" {
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
	a := identity.Account{Label: "账号", Name: `-王 "Li" \ Dev #1`, Email: "dev@example.org"}
	apply(t, dir, &a)
	if r := InspectProject(dir); !r.Current.Matches(a) {
		t.Fatal(r)
	}
	if s := runGit(t, dir, "config", "--file", includedPath, "--get", "user.name"); s != "Team\n" {
		t.Fatal(s)
	}
	bad := identity.Account{Label: "bad", Name: "Broken\nValue", Email: "fail@example.org"}
	if err := SwitchIdentity(dir, &bad); err == nil {
		t.Fatal("invalid pair accepted")
	}
	if r := InspectProject(dir); !r.Current.Matches(a) {
		t.Fatal("partial update", r)
	}
}
