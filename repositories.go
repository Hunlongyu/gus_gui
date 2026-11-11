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
	path        string
	current     gitIdentity
	problem     string
	remotes     []gitRemote
	remoteError string
	branch      string
	detached    bool
	branchError string
}

func (p *application) restore() {
	var err error
	p.stateFile, err = dataPath()
	var data savedData
	if err == nil {
		data, err = loadData(p.stateFile)
	}
	if err != nil {
		p.storageError = fmt.Sprintf(p.tr("无法读取 %s：%v"), p.stateFile, err)
		p.notifyError(p.storageError)
		return
	}
	p.accounts = data.Accounts
	p.preferences = data.Preferences.normalized()
	p.windowSize = data.WindowSize
	for _, path := range data.Directories {
		p.projects = append(p.projects, project{path: path})
	}
}

func (p *application) persist(accounts []identity, projects []project) error {
	if p.storageError != "" {
		return errors.New(p.storageError)
	}
	data := savedData{Accounts: accounts, Preferences: p.preferences, WindowSize: p.windowSize}
	for _, r := range projects {
		data.Directories = append(data.Directories, r.path)
	}
	if err := storeData(p.stateFile, data); err != nil {
		return fmt.Errorf(p.tr("保存失败：%w"), err)
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
		p.setNotice(err.Error())
	}
}

func (p *application) pickDirectories() {
	if p.busy {
		return
	}
	p.busy, p.picking = true, true
	title := p.tr("添加 Git 项目目录")
	go func() {
		paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{Parent: p.win, Directory: true, Multiple: true, Title: title})
		p.updateUI(func() {
			p.picking, p.busy = false, false
			if err != nil {
				p.setNotice(err.Error())
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
	folderError := p.tr("请选择目录：%s")
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
				issues = append(issues, fmt.Sprintf(folderError, path))
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
			p.setNotice(strings.Join(issues, "\n"))
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
			previousError := p.globalError
			p.globalError = ""
			if err != nil {
				p.globalError = fmt.Sprintf(p.tr("读取全局身份失败：%v"), err)
				if p.globalError != previousError {
					p.notifyError(p.globalError)
				}
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
			p.setNotice("")
			if err != nil {
				p.setNotice(p.tr("切换失败或未完全生效：") + err.Error())
			}
			p.refresh()
		})
	}()
}

func (p *application) displayValue(value string) string {
	if value == "" {
		return p.tr("未设置")
	}
	return value
}

func (p *application) description(r project) string {
	if r.problem != "" {
		return r.path + "\n" + r.problem
	}
	return filepath.Clean(r.path) + "\n" +
		p.tr("姓名：") + p.displayValue(r.current.Name) + " (" + p.displayValue(r.current.NameScope) + ")\n" +
		p.tr("邮箱：") + p.displayValue(r.current.Email) + " (" + p.displayValue(r.current.EmailScope) + ")\n" +
		p.tr("推送平台：") + p.pushPlatforms(r) + "\n" + p.tr("双击查看项目详情")
}
