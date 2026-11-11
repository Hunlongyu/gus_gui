package app

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// MyGo compares the complete Run command when reading the login setting.
// Update only our own old command, preserving Windows' StartupApproved state.
func migrateStartupPath() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	if !strings.EqualFold(filepath.Base(exe), "gus.exe") {
		return nil
	}
	const id = "io.github.hunlongyu.gus-gui"
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE|registry.SET_VALUE)
	if errors.Is(err, registry.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	defer key.Close()
	value, _, err := key.GetStringValue(id)
	if errors.Is(err, registry.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	command := func(path string) string { return `"` + path + `" --mygo-opened-at-login` }
	old := filepath.Join(filepath.Dir(exe), "GitUserSwitch.exe")
	if strings.EqualFold(value, command(old)) {
		return key.SetStringValue(id, command(exe))
	}
	return nil
}
