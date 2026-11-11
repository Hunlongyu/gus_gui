package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type identity struct {
	Label string `json:"label"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
type savedData struct {
	Preferences preferences `json:"preferences,omitempty"`
	WindowSize  panelSize   `json:"windowSize,omitzero"`
	Version     int         `json:"version"`
	Accounts    []identity  `json:"accounts"`
	Directories []string    `json:"directories"`
}

func samePath(a, b string) bool {
	a, b = filepath.Clean(a), filepath.Clean(b)
	if runtime.GOOS == "windows" {
		if strings.EqualFold(a, b) {
			return true
		}
		if strings.Contains(a, "~") || strings.Contains(b, "~") {
			first, errA := os.Stat(a)
			second, errB := os.Stat(b)
			return errA == nil && errB == nil && os.SameFile(first, second)
		}
		return false
	}
	return a == b
}

func dataPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "GitUserSwitch", "state.json"), nil
}

func loadData(path string) (savedData, error) {
	var data savedData
	raw, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return data, nil
	}
	if err != nil {
		return data, err
	}
	if err = json.Unmarshal(raw, &data); err != nil {
		return savedData{}, err
	}
	if data.Version != 1 {
		return savedData{}, fmt.Errorf("不支持的数据版本 %d", data.Version)
	}
	labels := map[string]bool{}
	for _, a := range data.Accounts {
		if a.Label == "" || a.Name == "" || a.Email == "" || labels[a.Label] {
			return savedData{}, errors.New("账号数据无效或标签重复")
		}
		labels[a.Label] = true
	}
	for i, path := range data.Directories {
		if !filepath.IsAbs(path) {
			return savedData{}, errors.New("目录列表包含相对路径")
		}
		for _, prior := range data.Directories[:i] {
			if samePath(prior, path) {
				return savedData{}, errors.New("目录列表包含重复路径")
			}
		}
	}
	return data, nil
}

func storeData(path string, data savedData) error {
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
