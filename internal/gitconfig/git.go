// Package gitconfig reads and updates local Git configuration.
package gitconfig

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gus/internal/identity"
	"gus/internal/pathutil"
)

// Project is a repository snapshot; metadata failures do not hide its identity.
type Project struct {
	Path        string
	Current     identity.Resolved
	Problem     string
	Remotes     []Remote
	RemoteError string
	Branch      string
	Detached    bool
	BranchError string
}

func git(dir string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", args...)
	hideGitWindow(cmd)
	if dir == "" {
		dir = os.TempDir()
	}
	cmd.Dir = dir
	// A tray app launched from a Git shell must still operate on the selected repository.
	for _, v := range os.Environ() {
		k, _, _ := strings.Cut(v, "=")
		switch strings.ToUpper(k) {
		case "GIT_DIR", "GIT_COMMON_DIR", "GIT_WORK_TREE", "GIT_CONFIG", "GIT_INDEX_FILE", "GIT_CEILING_DIRECTORIES":
			continue
		}
		cmd.Env = append(cmd.Env, v)
	}
	var out, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &stderr
	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return "", fmt.Errorf("%s: %w", msg, err)
		}
		return "", err
	}
	return out.String(), nil
}

func exitCode(err error, code int) bool {
	var e *exec.ExitError
	return errors.As(err, &e) && e.ExitCode() == code
}

func RepositoryRoot(path string) (string, error) {
	out, err := git(path, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("无法打开工作区 %s：%w", path, err)
	}
	return filepath.EvalSymlinks(filepath.Clean(strings.TrimRight(out, "\r\n")))
}

func ReadIdentity(dir string) (identity.Resolved, error) {
	args := []string{"config", "--includes", "--show-scope", "--null"}
	args = append(args, "--get-regexp", `^user\.(name|email)$`)
	out, err := git(dir, args...)
	var result identity.Resolved
	if exitCode(err, 1) {
		return result, nil
	}
	if err != nil {
		return result, err
	}
	fields := strings.Split(strings.TrimSuffix(out, "\x00"), "\x00")
	for i := 0; i+1 < len(fields); i += 2 {
		// --global reads a single writable file. The combined lookup also
		// sees XDG defaults when both global config files exist.
		if dir == "" && fields[i] != "global" {
			continue
		}
		key, value, _ := strings.Cut(fields[i+1], "\n")
		switch key {
		case "user.name":
			result.Name, result.NameScope = value, fields[i]
		case "user.email":
			result.Email, result.EmailScope = value, fields[i]
		}
	}
	return result, nil
}

func InspectProject(path string) Project {
	r := inspectProjectIdentity(path)
	if r.Problem != "" {
		return r
	}
	var err error
	r.Remotes, err = readRemotes(path)
	if err != nil {
		r.RemoteError = err.Error()
	}
	r.Branch, r.Detached, err = readBranch(path)
	if err != nil {
		r.BranchError = err.Error()
	}
	return r
}

func inspectProjectIdentity(path string) Project {
	r := Project{Path: path}
	root, err := RepositoryRoot(path)
	if err == nil && !pathutil.Same(path, root) {
		err = fmt.Errorf("该目录已不再是仓库根目录")
	}
	if err == nil {
		r.Current, err = ReadIdentity(root)
	}
	if err != nil {
		r.Problem = err.Error()
	}
	return r
}

// git var lists global files in read order. Git writes ~/.gitconfig unless
// only the XDG file exists; GIT_CONFIG_GLOBAL reduces the list to one file.
func globalConfigPath() (string, error) {
	out, err := git("", "var", "GIT_CONFIG_GLOBAL")
	if err != nil {
		return "", fmt.Errorf("无法定位全局配置，请使用支持 git var GIT_CONFIG_GLOBAL 的 Git 版本：%w", err)
	}
	lines := strings.Split(strings.TrimRight(out, "\r\n"), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		name := filepath.Clean(strings.TrimSuffix(lines[i], "\r"))
		if !filepath.IsAbs(name) {
			return "", fmt.Errorf("全局配置路径不是绝对路径：%s", name)
		}
		_, err := os.Stat(name)
		if err == nil {
			return name, nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
	}
	return filepath.Clean(strings.TrimSuffix(lines[len(lines)-1], "\r")), nil
}

// Hold the same lock name as Git while preparing both keys on the copy.
// Readers see either the old pair or the new one, including after a failed edit.
func editConfig(path string, a *identity.Account) error {
	if target, err := filepath.EvalSymlinks(path); err == nil {
		path = target
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	} else if _, lstatErr := os.Lstat(path); lstatErr == nil {
		return fmt.Errorf("配置文件链接的目标不存在：%s", path)
	}
	lock := path + ".lock"
	f, err := os.OpenFile(lock, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	committed := false
	defer func() {
		if !committed {
			_ = os.Remove(lock)
		}
	}()
	old, err := os.Open(path)
	if err == nil {
		info, statErr := old.Stat()
		if statErr == nil && !info.Mode().IsRegular() {
			statErr = fmt.Errorf("不是普通配置文件：%s", path)
		}
		if statErr == nil {
			statErr = f.Chmod(info.Mode().Perm())
		}
		if statErr == nil {
			_, statErr = io.Copy(f, old)
		}
		old.Close()
		err = statErr
	} else if errors.Is(err, os.ErrNotExist) {
		err = nil
	}
	closeErr := f.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	for _, key := range []string{"user.name", "user.email"} {
		_, err = git("", "config", "--file", lock, "--no-includes", "--unset-all", "--", key)
		if err != nil && !exitCode(err, 5) {
			return err
		}
	}
	// Append after any includes so values in included files cannot shadow
	// the selected identity. Only the explicitly selected file is edited.
	if a != nil {
		if strings.ContainsAny(a.Name+a.Email, "\x00\r\n") {
			return errors.New("姓名和邮箱必须为单行文本")
		}
		escape := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\t", `\t`, "\b", `\b`)
		tail, err := os.OpenFile(lock, os.O_RDWR|os.O_APPEND, 0)
		if err != nil {
			return err
		}
		// A hand-written config need not end with a newline. Separate the
		// new section without accumulating blank lines on repeated switches.
		prefix := ""
		end, err := tail.Seek(0, io.SeekEnd)
		if err == nil && end > 0 {
			var last [1]byte
			_, err = tail.ReadAt(last[:], end-1)
			if last[0] != '\n' {
				prefix = "\n"
			}
		}
		if err == nil {
			_, err = fmt.Fprintf(tail, "%s[user]\n\tname = \"%s\"\n\temail = \"%s\"\n", prefix, escape.Replace(a.Name), escape.Replace(a.Email))
		}
		if err == nil {
			err = tail.Sync()
		}
		closeErr := tail.Close()
		if err != nil {
			return err
		}
		if closeErr != nil {
			return closeErr
		}
	}
	if err := os.Rename(lock, path); err != nil {
		return err
	}
	committed = true
	return nil
}

func SwitchIdentity(dir string, a *identity.Account) error {
	var config string
	var err error
	if dir == "" {
		if a == nil {
			return errors.New("请选择全局账号")
		}
		config, err = globalConfigPath()
	} else {
		r := inspectProjectIdentity(dir)
		if r.Problem != "" {
			return errors.New(r.Problem)
		}
		leaf := "config"
		if r.Current.NameScope == "worktree" || r.Current.EmailScope == "worktree" {
			leaf = "config.worktree"
		}
		config, err = git(dir, "rev-parse", "--path-format=absolute", "--git-path", leaf)
		config = filepath.Clean(strings.TrimRight(config, "\r\n"))
	}
	if err != nil {
		return err
	}
	if err = editConfig(config, a); err != nil {
		return err
	}
	effective, err := ReadIdentity(dir)
	if err != nil {
		return err
	}
	if a != nil && !effective.Matches(*a) {
		return errors.New("已写入配置，但其他 Git 配置覆盖了所选身份，请查看实际生效值")
	}
	if a == nil && (identity.RepositoryScope(effective.NameScope) || identity.RepositoryScope(effective.EmailScope)) {
		return errors.New("已移除当前文件的覆盖，仍有其他仓库配置生效")
	}
	return nil
}
