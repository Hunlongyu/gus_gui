package gitconfig

import (
	"slices"
	"testing"
)

func TestMultipleURLsFallbackAndCredentialRedaction(t *testing.T) {
	dir := isolatedGit(t)
	runGit(t, dir, "init", "-q")
	runGit(t, dir, "remote", "add", "origin", "https://user:secret@github.com/team/repo.git?token=secret#secret")
	runGit(t, dir, "config", "--add", "remote.origin.url", "ssh://git@gitee.com/team/repo.git")
	r := InspectProject(dir)
	if r.RemoteError != "" || len(r.Remotes) != 1 {
		t.Fatal(r)
	}
	want := []string{"https://github.com/team/repo.git", "ssh://git@gitee.com/team/repo.git"}
	if !slices.Equal(r.Remotes[0].FetchURLs, want) || !slices.Equal(r.Remotes[0].PushURLs, want) {
		t.Fatal(r.Remotes)
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
		"git@ssh.gitgud.io:team/repo.git":             "GitGud",
		"ssh://git@SSH.GITGUD.IO:2222/team/repo.git":  "GitGud",
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
		if got := RemotePlatform(address); got != want {
			t.Errorf("%q: got %q, want %q", address, got, want)
		}
	}
}

func TestRemoteWebURLs(t *testing.T) {
	for address, want := range map[string]string{
		"https://user:token@github.com/team/repo.git?secret=token": "https://github.com/team/repo",
		"https://git.example.org:8443/team/repo.git":               "https://git.example.org:8443/team/repo",
		"git@github.com:team/repo.git":                             "https://github.com/team/repo",
		"ssh://git@altssh.gitlab.com:443/team/repo.git":            "https://gitlab.com/team/repo",
		"git@ssh.gitgud.io:team/repo.git":                          "https://gitgud.io/team/repo",
		"ssh://git@SSH.GITGUD.IO:2222/team/repo.git":               "https://gitgud.io/team/repo",
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
		if got := RemoteWebURL(address); got != want {
			t.Errorf("%q: got %q, want %q", address, got, want)
		}
	}
}
