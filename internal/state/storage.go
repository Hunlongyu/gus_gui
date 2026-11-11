// Package state persists application data without depending on the UI.
package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gus/internal/gitconfig"
	"gus/internal/identity"
	"gus/internal/pathutil"
)

type Data struct {
	Preferences  Preferences         `json:"preferences,omitempty"`
	WindowSize   WindowSize          `json:"windowSize,omitzero"`
	Version      int                 `json:"version"`
	Accounts     []identity.Account  `json:"accounts"`
	Directories  []string            `json:"directories"`
	ProjectCache []gitconfig.Project `json:"projectCache,omitempty"`
	Global       identity.Resolved   `json:"global,omitzero"`
}

func Path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "GitUserSwitch", "state.json"), nil
}

func Load(path string) (Data, error) {
	var data Data
	raw, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return data, nil
	}
	if err != nil {
		return data, err
	}
	if err = json.Unmarshal(raw, &data); err != nil {
		return Data{}, err
	}
	if data.Version != 1 {
		return Data{}, fmt.Errorf("不支持的数据版本 %d", data.Version)
	}
	labels := map[string]bool{}
	for _, a := range data.Accounts {
		if a.Label == "" || a.Name == "" || a.Email == "" || labels[a.Label] {
			return Data{}, errors.New("账号数据无效或标签重复")
		}
		labels[a.Label] = true
	}
	for i, path := range data.Directories {
		if !filepath.IsAbs(path) {
			return Data{}, errors.New("目录列表包含相对路径")
		}
		for _, prior := range data.Directories[:i] {
			if pathutil.Same(prior, path) {
				return Data{}, errors.New("目录列表包含重复路径")
			}
		}
	}
	return data, nil
}

func Store(path string, data Data) error {
	data.Version = 1
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err = os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, ".state-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err = f.Write(raw); err == nil {
		err = f.Sync()
	}
	closeErr := f.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	return os.Rename(f.Name(), path)
}
