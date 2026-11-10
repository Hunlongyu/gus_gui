// Older installations relaunch GitUserSwitch.exe after applying a signed
// update. Forward that entry point, including startup arguments, to gus.
package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	self, err := os.Executable()
	if err != nil {
		os.Exit(1)
	}
	cmd := exec.Command(filepath.Join(filepath.Dir(self), "gus.exe"), os.Args[1:]...)
	if err := cmd.Start(); err != nil {
		os.Exit(1)
	}
	_ = cmd.Process.Release()
}
