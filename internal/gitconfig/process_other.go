//go:build !windows

package gitconfig

import "os/exec"

func hideGitWindow(*exec.Cmd) {}
