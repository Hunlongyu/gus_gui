package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/egoist/mygo"
)

type project struct {
	path    string
	current gitIdentity
	problem string
}

func (p *application) restore() {
	var err error
	p.stateFile, err = dataPath()
	var data savedData
	if err == nil {
		data, err = loadData(p.stateFile)
	}
	if err != nil {
		p.storageError = fmt.Sprintf("无法读取 %s：%v", p.stateFile, err)
		return
	}
	p.accounts = data.Accounts
	for _, path := range data.Directories {
		p.projects = append(p.projects, project{path: path})
	}
}

func (p *application) persist(accounts []identity, projects []project) error {
	if p.storageError != "" {
		return errors.New(p.storageError)
	}
	data := savedData{Accounts: accounts}
	for _, r := range projects {
		data.Directories = append(data.Directories, r.path)
	}
	if err := storeData(p.stateFile, data); err != nil {
		return fmt.Errorf("保存失败：%w", err)
	}
	p.accounts, p.projects = accounts, projects
	return nil
}

func (p *application) addRoots(roots []string) error {
	list := slices.Clone(p.projects)
	for _, root := range roots {
		if !slices.ContainsFunc(list, func(r project) bool { return samePath(root, r.path) }) {
			list = append(list, project{path: root})
		}
	}
	if len(list) == len(p.projects) {
		return nil
	}
	return p.persist(p.accounts, list)
}

func (p *application) removeDirectory(path string) {
	list := slices.DeleteFunc(slices.Clone(p.projects), func(r project) bool { return samePath(r.path, path) })
	if err := p.persist(p.accounts, list); err != nil {
		p.notice = err.Error()
	}
}

func (p *application) pickDirectories() {
	if p.busy {
		return
	}
	p.busy, p.picking = true, true
	go func() {
		paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{Parent: p.win, Directory: true, Multiple: true, Title: "添加 Git 项目目录"})
		p.updateUI(func() {
			p.picking, p.busy = false, false
			if err != nil {
				p.notice = err.Error()
			}
			p.importDirectories(paths)
		})
	}()
}

// The picker and Explorer drops share one queue. A drop arriving during a
// refresh waits for that snapshot to finish before changing the project list.
func (p *application) importDirectories(paths []string) {
	p.pendingDirectories = append(p.pendingDirectories, paths...)
	if p.busy || len(p.pendingDirectories) == 0 {
		return
	}
	paths, p.pendingDirectories = p.pendingDirectories, nil
	p.busy = true
	go func() {
		var roots, issues []string
		for _, path := range paths {
			info, err := os.Stat(path)
			if err != nil {
				issues = append(issues, err.Error())
				continue
			}
			if !info.IsDir() {
				issues = append(issues, fmt.Sprintf("请选择目录：%s", path))
				continue
			}
			root, err := repositoryRoot(path)
			if err != nil {
				issues = append(issues, err.Error())
				continue
			}
			roots = append(roots, root)
		}
		p.updateUI(func() {
			p.busy = false
			if err := p.addRoots(roots); err != nil {
				issues = append(issues, err.Error())
			}
			p.notice = strings.Join(issues, "\n")
			p.refresh()
		})
	}()
}

// Disk access and Git subprocesses stay off the MyGo UI thread. Serialise
// mutations with scans so an older scan cannot overwrite a switch or removal.
func (p *application) refresh() {
	if p.busy {
		return
	}
	p.busy = true
	rows := slices.Clone(p.projects)
	go func() {
		global, err := readIdentity("")
		var wg sync.WaitGroup
		jobs := make(chan int)
		for range min(4, len(rows)) {
			wg.Go(func() {
				for i := range jobs {
					rows[i] = inspectProject(rows[i].path)
				}
			})
		}
		for i := range rows {
			jobs <- i
		}
		close(jobs)
		wg.Wait()
		p.updateUI(func() {
			p.global, p.projects, p.busy = global, rows, false
			p.globalError = ""
			if err != nil {
				p.globalError = fmt.Sprintf("读取全局身份失败：%v", err)
			}
			p.importDirectories(nil)
		})
	}()
}

func (p *application) applyIdentity(path string, account *identity) {
	if p.busy {
		return
	}
	p.busy, p.writing = true, true
	go func() {
		err := switchIdentity(path, account)
		p.updateUI(func() {
			p.busy, p.writing = false, false
			if p.quitAfterWrite && err == nil {
				mygo.App.Quit()
				return
			}
			p.quitAfterWrite = false
			p.notice = ""
			if err != nil {
				p.notice = "切换失败或未完全生效：" + err.Error()
			}
			p.refresh()
		})
	}()
}

func displayValue(value string) string {
	if value == "" {
		return "未设置"
	}
	return value
}

func (r project) description() string {
	if r.problem != "" {
		return r.path + "\n" + r.problem
	}
	return fmt.Sprintf("%s\n姓名：%s\n邮箱：%s", filepath.Clean(r.path), displayValue(r.current.NameScope), displayValue(r.current.EmailScope))
}
