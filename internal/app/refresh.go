package app

import (
	"fmt"
	"reflect"
	"runtime"
	"slices"
	"sync"
	"time"

	"github.com/egoist/mygo"

	"gus/internal/gitconfig"
	"gus/internal/identity"
	"gus/internal/pathutil"
)

const projectRefreshInterval = time.Hour

// Keep Git subprocesses off the UI thread and serialize scans with writes.
// Results are applied individually; one slow repository cannot hold up others.
func (p *application) refresh() {
	p.refreshWith(gitconfig.ReadIdentity, gitconfig.InspectProject)
}

func (p *application) refreshWith(readIdentity func(string) (identity.Resolved, error), inspect func(string) gitconfig.Project) {
	if p.busy || p.quitting {
		return
	}
	p.busy, p.refreshing = true, true
	paths := make([]string, len(p.projects))
	for i, row := range p.projects {
		paths[i] = row.Path
	}
	go func() {
		var wg sync.WaitGroup
		wg.Go(func() {
			global, err := readIdentity("")
			p.updateChanged(func() bool {
				previousError := p.globalError
				p.globalError = ""
				if err != nil {
					p.globalError = fmt.Sprintf(p.tr("读取全局身份失败：%v"), err)
					if p.globalError != previousError {
						p.notifyError(p.globalError)
					}
				} else if p.global != global {
					p.global, p.cacheDirty = global, true
					return true
				}
				return p.globalError != previousError
			})
		})
		jobs := make(chan string)
		for range min(len(paths), max(2, min(8, runtime.NumCPU()))) {
			wg.Go(func() {
				for path := range jobs {
					row := inspect(path)
					p.updateChanged(func() bool { return p.updateProject(row) })
				}
			})
		}
		for _, path := range paths {
			jobs <- path
		}
		close(jobs)
		wg.Wait()
		p.updateUI(func() {
			p.busy, p.refreshing = false, false
			p.saveProjectCache()
			p.importDirectories(nil)
		})
	}()
}

// Window.Update redraws even when nothing changed. Compare on the UI thread
// first, and invalidate only changed results, keeping widgets and scroll state.
func (p *application) updateChanged(change func() bool) {
	if p.win == nil {
		p.updateUI(func() { change() })
		return
	}
	mygo.RunOnMain(func() {
		if !p.quitting && change() {
			p.win.Invalidate()
		}
	})
}

func (p *application) updateProject(row gitconfig.Project) bool {
	i := slices.IndexFunc(p.projects, func(r gitconfig.Project) bool { return pathutil.Same(r.Path, row.Path) })
	if i < 0 {
		return false
	}
	previous := p.projects[i]
	row.Path = previous.Path
	// A temporary read failure must not erase the last known properties.
	// Keep the error indicators so stale values are not presented as verified.
	if row.Problem != "" {
		row.Current, row.Remotes = previous.Current, previous.Remotes
		row.Branch, row.Detached = previous.Branch, previous.Detached
		row.RemoteError, row.BranchError = previous.RemoteError, previous.BranchError
	} else {
		if row.RemoteError != "" {
			row.Remotes = previous.Remotes
		}
		if row.BranchError != "" {
			row.Branch, row.Detached = previous.Branch, previous.Detached
		}
	}
	if reflect.DeepEqual(previous, row) {
		return false
	}
	p.projects[i], p.cacheDirty = row, true
	return true
}

func (p *application) saveProjectCache() {
	if !p.cacheDirty || p.stateFile == "" || p.storageError != "" {
		return
	}
	if err := p.persist(p.accounts, p.projects); err != nil {
		p.setNotice(err.Error())
	}
}

func (p *application) scheduleRefresh(delay time.Duration) {
	if p.refreshTimer != nil {
		p.refreshTimer.Stop()
	}
	p.refreshTimer = time.AfterFunc(delay, func() {
		p.updateUI(p.autoRefresh)
	})
}

func (p *application) autoRefresh() {
	if p.quitting {
		return
	}
	if p.busy {
		// Retry shortly after an import/write instead of skipping an entire hour.
		p.scheduleRefresh(time.Minute)
		return
	}
	p.refresh()
	p.scheduleRefresh(projectRefreshInterval)
}
