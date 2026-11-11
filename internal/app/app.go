// Package app implements the gus desktop UI and application lifecycle.
package app

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"gus/internal/gitconfig"
	"gus/internal/identity"
	"gus/internal/state"
)

const panelWidth, panelHeight = 480, 490
const appName = "gus"
const minPanelWidth, minPanelHeight = 440, 320

type application struct {
	releaseHistory                               string
	preferences                                  state.Preferences
	updateUI                                     func(func())
	win                                          *mygo.Window
	tray                                         *mygo.Tray
	editing                                      int
	global                                       identity.Resolved
	stateFile, storageError, globalError, notice string
	busy, picking                                bool
	writing, quitAfterWrite                      bool
	accounts                                     []identity.Account
	draft                                        identity.Account
	formError                                    string
	pinned                                       bool
	loginEnabled                                 bool
	loginError                                   string
	updateStatus                                 string
	updateDetail                                 string
	checkingUpdate                               bool
	lastUpdateCheck                              time.Time
	availableUpdate                              *mygo.Update
	installingUpdate                             bool
	projects                                     []gitconfig.Project
	pendingDirectories                           []string
	projectQuery                                 string
	detailPath                                   string
	page                                         string
	pendingToasts                                []ui.Toast
	windowSize                                   state.WindowSize
	resizeTimer                                  *time.Timer
}

// Run starts the desktop application with the embedded icon and release notes.
func Run(appIcon []byte, releaseHistory string) {
	p := &application{releaseHistory: releaseHistory}
	mygo.App.SetName(appName)
	guard, err := claimInstance()
	if err != nil {
		log.Fatal("single instance: ", err)
	}
	if guard == nil {
		return
	}
	defer guard.close()
	if !mygo.App.RequestSingleInstanceLock() {
		return
	}
	// OnQuit runs before MyGo starts the replacement process after an update.
	mygo.App.OnQuit(guard.close)
	// Keep the saved window geometry of installations from before the rename.
	if dir, err := os.UserConfigDir(); err == nil {
		legacy := filepath.Join(dir, "Git User Switch")
		if info, err := os.Stat(legacy); err == nil && info.IsDir() {
			mygo.App.SetPath(mygo.PathUserData, legacy)
		}
	}
	if err := migrateStartupPath(); err != nil {
		log.Printf("migrate startup path: %v", err)
	}
	mygo.App.OnBeforeQuit(func(e *mygo.QuitEvent) {
		if p.installingUpdate {
			e.PreventDefault()
			return
		}
		if p.writing {
			p.quitAfterWrite = true
			e.PreventDefault()
			return
		}
		p.savePanelSize()
	})
	p.restore()
	mygo.App.SetLanguages("zh-CN", "en")
	p.applyPreferences()
	mygo.App.OnSecondInstance(func(_ []string, _ string) {
		if p.win != nil {
			p.show()
		}
	})
	mygo.App.OnActivate(func(_ bool) {
		if p.win != nil {
			p.show()
		}
	})

	mygo.App.WhenReady(func() {
		var err error
		p.tray, err = mygo.NewTray(mygo.TrayOptions{Icon: trayIcon(appIcon), ToolTip: appName})
		if err != nil {
			// A missing tray must not prevent access to accounts and projects.
			log.Printf("tray unavailable; using a regular window: %v", err)
		}
		p.win = mygo.NewWindow(mygo.WindowOptions{
			Title: appName, Content: ui.View(p.view),
			Width: panelWidth, Height: panelHeight,
			MinWidth: minPanelWidth, MinHeight: minPanelHeight, StateKey: "tray-panel",
			Hidden: true, Frameless: true, SkipTaskbar: runtime.GOOS == "windows" && p.tray != nil,
		})
		if err := p.win.SetIcon(appIcon); err != nil {
			log.Fatal("window icon: ", err)
		}
		p.updateUI = p.win.Update
		p.trackPanelSize()
		guard.watch(func() { p.updateUI(p.show) })
		if mygo.Updater.Enabled() {
			p.checkUpdates()
		}
		p.win.OnBlur(func() {
			if p.hideToTray() && !p.pinned && !p.picking {
				p.win.Hide()
			}
		})
		if p.tray != nil {
			menu := func() *mygo.Menu {
				return mygo.NewMenu([]*mygo.MenuItem{{Label: p.tr("打开面板"), Click: func(*mygo.MenuItem, *mygo.Window) { p.show() }}, {Label: p.tr("退出"), Role: mygo.RoleQuit}})
			}
			if runtime.GOOS == "linux" {
				// AppIndicator only supports persistent menus; PopUpMenu is a no-op.
				p.tray.SetMenu(menu())
			} else {
				p.tray.OnClick(func() {
					if p.win.IsVisible() {
						p.win.Hide()
					} else {
						p.show()
					}
				})
				p.tray.OnRightClick(func() {
					p.tray.SetMenu(menu())
					defer p.tray.SetMenu(nil)
					p.tray.PopUpMenu()
				})
			}
		}
		if runtime.GOOS != "windows" || p.tray == nil {
			p.show()
			if mygo.App.WasOpenedAtLogin() {
				p.minimize()
			}
		}
		mygo.App.OnWillQuit(func(*mygo.QuitEvent) {
			if p.tray != nil {
				p.tray.Destroy()
			}
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
