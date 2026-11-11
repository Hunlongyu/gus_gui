package gitconfig

import (
	"os/exec"
	"syscall"
)

func hideGitWindow(cmd *exec.Cmd) { cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} }
