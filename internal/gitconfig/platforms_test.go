package gitconfig

import "testing"

func TestAdditionalRemotePlatforms(t *testing.T) {
	for _, tc := range []struct {
		remote, platform string
	}{
		{"git@gitea.com:team/repo.git", "Gitea"},
		{"ssh://git@GITEA.COM:2222/team/repo.git", "Gitea"},
		{"https://code.forgejo.org/forgejo/forgejo.git", "Forgejo"},
		{"git@git.sr.ht:~user/repo", "SourceHut"},
		{"ssh://user@git.code.sf.net/p/project/code", "SourceForge"},
		{"https://sourceforge.net/p/project/code", "SourceForge"},
		{"git+ssh://user@git.launchpad.net/~user/project/+git/repo", "Launchpad"},
		{"https://launchpad.net/project", "Launchpad"},
		{"https://coding.net/user/repo.git", "CODING"},
		{"git@e.coding.net:team/project/repo.git", "CODING"},
		{"git@git.coding.net:user/repo.git", "CODING"},
		{"https://TEAM.CODING.NET/p/project/d/repo/git", "CODING"},
		{"git@codeup.aliyun.com:team/repo.git", "Codeup"},
		{"https://code.aliyun.com/team/repo.git", "Codeup"},
		{"https://team.devops.aliyuncs.com/team/repo.git", "Codeup"},
		{"git@gitea.company.example:team/repo.git", "gitea.company.example"},
		{"git@forgejo-work:team/repo.git", "forgejo-work"},
		{"https://gitea.com.example.org/team/repo.git", "gitea.com.example.org"},
		{"https://notcoding.net/team/repo.git", "notcoding.net"},
		{"https://coding.net.example.org/team/repo.git", "coding.net.example.org"},
		{"https://notdevops.aliyuncs.com/team/repo.git", "notdevops.aliyuncs.com"},
		{"https://team.devops.aliyuncs.com.example.org/repo.git", "team.devops.aliyuncs.com.example.org"},
		{"https://example.org/coding.net/repo.git", "example.org"},
	} {
		t.Run(tc.remote, func(t *testing.T) {
			if got := RemotePlatform(tc.remote); got != tc.platform {
				t.Errorf("RemotePlatform(%q) = %q, want %q", tc.remote, got, tc.platform)
			}
		})
	}
}
