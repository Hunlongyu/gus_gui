package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type gitRemote struct {
	Name                string
	FetchURLs, PushURLs []string
}

func gitLines(out string) []string {
	out = strings.TrimRight(out, "\r\n")
	if out == "" {
		return nil
	}
	return strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n")
}

// get-url resolves insteadOf/pushInsteadOf and falls back to fetch URLs when
// no pushurl is configured. --all preserves multiple destinations per remote.
// These commands only read local configuration; they never contact a server.
func readRemotes(dir string) ([]gitRemote, error) {
	out, err := git(dir, "remote")
	if err != nil {
		return nil, err
	}
	var remotes []gitRemote
	var issues []error
	for _, name := range gitLines(out) {
		r := gitRemote{Name: name}
		for _, push := range []bool{false, true} {
			args := []string{"remote", "get-url", "--all"}
			if push {
				args = append(args, "--push")
			}
			out, err := git(dir, append(args, "--", name)...)
			if err != nil {
				issues = append(issues, fmt.Errorf("%s: %w", name, err))
				continue
			}
			urls := gitLines(out)
			for i := range urls {
				urls[i] = displayRemoteURL(urls[i])
			}
			if push {
				r.PushURLs = urls
			} else {
				r.FetchURLs = urls
			}
		}
		remotes = append(remotes, r)
	}
	return remotes, errors.Join(issues...)
}

func readBranch(dir string) (string, bool, error) {
	out, err := git(dir, "symbolic-ref", "--quiet", "--short", "HEAD")
	if err == nil {
		return strings.TrimSpace(out), false, nil
	}
	if !exitCode(err, 1) {
		return "", false, err
	}
	out, err = git(dir, "rev-parse", "--short", "HEAD")
	return strings.TrimSpace(out), true, err
}

// Credentials in HTTPS URLs are not account identities. Keep them out of the
// UI, including tooltips; SSH's transport username (usually git) is retained.
func displayRemoteURL(raw string) string {
	if !strings.Contains(raw, "://") {
		return raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "[invalid URL]"
	}
	if u.User != nil {
		if u.Scheme == "ssh" || u.Scheme == "git+ssh" || u.Scheme == "ssh+git" {
			u.User = url.User(u.User.Username())
		} else {
			u.User = nil
		}
	}
	u.RawQuery, u.Fragment, u.RawFragment, u.ForceQuery = "", "", "", false
	return u.String()
}

// Unknown hosts (including SSH aliases and self-hosted services) keep their
// hostname instead of guessing a hosting provider from a remote's name.
func remotePlatform(raw string) string {
	var host string
	if strings.Contains(raw, "://") {
		u, err := url.Parse(raw)
		if err != nil {
			return "未知平台"
		}
		if u.Scheme == "file" {
			return "本地路径"
		}
		host = u.Hostname()
	} else {
		if len(raw) >= 2 && raw[1] == ':' && ((raw[0] >= 'a' && raw[0] <= 'z') || (raw[0] >= 'A' && raw[0] <= 'Z')) {
			return "本地路径"
		}
		left, _, ok := strings.Cut(raw, ":")
		if !ok || strings.ContainsAny(left, `/\`) {
			return "本地路径"
		}
		if strings.Contains(left, "@") {
			_, left, _ = strings.Cut(left, "@")
		}
		host = strings.Trim(left, "[]")
	}
	host = strings.ToLower(host)
	switch host {
	case "github.com", "ssh.github.com":
		return "GitHub"
	case "gitlab.com", "altssh.gitlab.com":
		return "GitLab"
	case "gitee.com":
		return "Gitee"
	case "bitbucket.org", "altssh.bitbucket.org":
		return "Bitbucket"
	case "codeberg.org":
		return "Codeberg"
	case "gitgud.io", "www.gitgud.io":
		return "GitGud"
	case "gitcode.com":
		return "GitCode"
	case "dev.azure.com", "ssh.dev.azure.com":
		return "Azure DevOps"
	case "":
		return "未知平台"
	default:
		return host
	}
}

// Only browser URLs are returned. Local paths, remote helpers and SSH aliases
// without a known web hostname remain copyable text rather than launching an
// arbitrary protocol handler. SSH ports are transport ports, not web ports.
func remoteWebURL(raw string) string {
	if !strings.Contains(raw, "://") {
		if remotePlatform(raw) == "本地路径" {
			return ""
		}
		host, path, ok := strings.Cut(raw, ":")
		if !ok {
			return ""
		}
		raw = "ssh://" + host + "/" + path
	}
	u, err := url.Parse(displayRemoteURL(raw))
	if err != nil || u.Hostname() == "" {
		return ""
	}
	switch u.Scheme {
	case "http", "https":
		// Keep explicit web ports for self-hosted services.
	case "ssh", "git", "git+ssh", "ssh+git":
		host := strings.ToLower(u.Hostname())
		if !strings.Contains(host, ".") {
			return ""
		}
		switch host {
		case "ssh.github.com":
			host = "github.com"
		case "altssh.gitlab.com":
			host = "gitlab.com"
		case "altssh.bitbucket.org":
			host = "bitbucket.org"
		case "ssh.dev.azure.com":
			parts := strings.Split(strings.Trim(u.Path, "/"), "/")
			if len(parts) != 4 || parts[0] != "v3" {
				return ""
			}
			host = "dev.azure.com"
			u.Path = "/" + parts[1] + "/" + parts[2] + "/_git/" + parts[3]
		}
		u.Scheme, u.Host = "https", host
	default:
		return ""
	}
	u.User = nil
	u.Path = strings.TrimSuffix(strings.TrimRight(u.Path, "/"), ".git")
	u.RawPath = ""
	return u.String()
}
