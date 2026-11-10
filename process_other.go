//go:build !windows

package main

import "os/exec"

func hideGitWindow(*exec.Cmd) {}
