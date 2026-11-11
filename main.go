package main

import (
	"bytes"
	_ "embed"
	"image"
	"image/png"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"golang.org/x/image/draw"
)

const panelWidth, panelHeight = 480, 490
const appName = "gus"
const minPanelWidth, minPanelHeight = 440, 320

// appIcon is also used by MyGo to embed the executable's Windows icons.
//
//go:embed resources/icon.png
var appIcon []byte

type application struct {
	preferences                                  preferences
	updateUI                                     func(func())
	win                                          *mygo.Window
	tray                                         *mygo.Tray
	editing                                      int
	global                                       gitIdentity
	stateFile, storageError, globalError, notice string
	busy, picking                                bool
	writing, quitAfterWrite                      bool
	accounts                                     []identity
	draft                                        identity
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
	projects                                     []project
	pendingDirectories                           []string
	projectQuery                                 string
	detailPath                                   string
	page                                         string
	pendingToasts                                []ui.Toast
	windowSize                                   panelSize
	resizeTimer                                  *time.Timer
}

func (p *application) view(c *ui.Context) {
	c.SetMenuPresentation(ui.PopoverMenus)
	if c.Shortcut(0, ui.KeyEscape) {
		if p.page == "changelog" {
			p.page = "settings"
		} else if p.page == "account" {
			p.page = "accounts"
		} else if p.page != "" {
			p.page = ""
		} else if p.projectQuery != "" {
			p.projectQuery = ""
		} else {
			p.minimize()
		}
	}
	base := c.Theme()
	switch p.preferences.Appearance {
	case "light":
		base = ui.LightTheme()
	case "dark":
		base = ui.DarkTheme()
	}
	t := panelTheme(base)
	c.SetTheme(&t)
	p.showNotifications(c)
	ui.Column(c).Fill().Padding(4).Background(t.Background).Children(func() {
		panel := ui.Column(c).Grow(1).Radius(8).Background(t.Surface).Border(1, t.Border)
		if p.page == "" && !p.installingUpdate {
			if paths := panel.DroppedFiles(); len(paths) > 0 {
				p.importDirectories(paths)
			}
			if panel.FileDragOver() {
				panel.Border(1, t.Focus)
			}
		}
		panel.Children(func() {
			p.toolbar(c)
			ui.Column(c).Grow(1).Padding(12).Gap(6).Children(func() {
				switch p.page {
				case "accounts":
					p.accountList(c)
				case "settings":
					p.appSettings(c)
				case "changelog":
					p.changelogView(c)
				case "account":
					p.accountForm(c)
				case "project":
					p.projectDetails(c)
				default:
					visible := p.projects
					count := strconv.Itoa(len(p.projects))
					if query := strings.ToLower(strings.TrimSpace(p.projectQuery)); query != "" {
						visible = nil
						for _, r := range p.projects {
							if strings.Contains(strings.ToLower(filepath.Base(r.path)), query) {
								visible = append(visible, r)
							}
						}
						count = strconv.Itoa(len(visible)) + " / " + count
					}
					ui.Row(c).Gap(5).AlignItems(ui.Center).Children(func() {
						ui.Text(c, p.tr("项目")).Bold().FontSize(16)
						ui.Text(c, count).Grow(1).FontSize(13).TextColor(t.TextMuted)
						ui.SearchField(c.Key("project-search"), &p.projectQuery).Width(168).Height(26).Label(p.tr("搜索项目名")).Placeholder(p.tr("搜索项目名"))
						withIcon(c, ui.Button(c, ""), accountIcon, p.tr("账号管理"), p.tr("账号管理")).OnClick(func() { p.page = "accounts" })
						withIcon(c, ui.Button(c, ""), settingsIcon, p.tr("设置"), p.tr("应用设置")).OnClick(func() { p.page = "settings" })
						withIcon(c, ui.Button(c, ""), addDirectoryIcon, p.tr("添加目录"), p.tr("添加项目目录")).Disabled(p.busy).OnClick(p.pickDirectories)
					})
					columns := p.projectColumns(c)
					ui.Row(c).Padding(5, 12, 5, 8).Gap(8).AlignItems(ui.Center).Background(t.Background).Radius(4).Children(func() {
						ui.Text(c, p.tr("仓库")).Width(columns.repo).Shrink(0).FontSize(12).MaxLines(1).Ellipsis("…").TextColor(t.TextMuted)
						ui.Text(c, p.tr("名字")).Width(columns.name).Shrink(0).FontSize(12).MaxLines(1).Ellipsis("…").TextColor(t.TextMuted)
						ui.Text(c, p.tr("邮箱")).Grow(1).MinWidth(0).FontSize(12).TextColor(t.TextMuted)
						ui.Text(c, p.tr("推送")).Width(56).FontSize(12).TextColor(t.TextMuted)
						ui.Text(c, p.tr("来源")).Width(42).FontSize(12).TextAlign(ui.Center).TextColor(t.TextMuted).Label(p.tr("身份来源")).Tooltip(p.tr("身份来源"))
					})
					if len(p.projects) == 0 {
						ui.Text(c, p.tr("拖入 Git 项目目录，或点击右上角添加")).TextColor(t.TextMuted).Padding(12, 8)
					} else if len(visible) == 0 {
						ui.Text(c, p.tr("没有匹配的项目")).TextColor(t.TextMuted).Padding(12, 8)
					}
					ui.Scroll(c.Key("projects")).Grow(1).Children(func() {
						for _, r := range visible {
							p.projectRow(c, r, columns)
						}
					})
				}
			})
		})
	})
}

func (p *application) toolbar(c *ui.Context) {
	t := c.Theme()
	a := p.global
	ui.Row(c).Padding(10, 12).Gap(10).AlignItems(ui.Center).Border(1, t.Border).BorderWidth(0, 0, 1, 0).Children(func() {
		ui.Row(c.Key("global")).Focusable().Grow(1).MinWidth(0).Gap(10).AlignItems(ui.Center).ContextMenu(p.identityMenu("", p.global, true)).Children(func() {
			ui.Row(c).Shrink(0).Padding(6, 8).Gap(5).Radius(6).Border(1, t.Border).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, globeIcon).Size(16, 16)
				ui.Text(c, p.tr("全局")).FontSize(13)
			})
			ui.Text(c, p.displayValue(a.Name)).Bold().FontSize(15).MaxWidth(112).Shrink(0).MaxLines(1).Ellipsis("…")
			ui.Text(c, p.displayValue(a.Email)).Grow(1).MinWidth(0).FontSize(13).TextColor(t.TextMuted).MaxLines(1).Ellipsis("…")
		})
		ui.Row(c).Shrink(0).Padding(2).Gap(1).Radius(7).Background(t.Background).Children(func() {
			hint := p.tr("置顶面板")
			if p.pinned {
				hint = p.tr("取消置顶")
			}
			pin := ui.ToggleBase(c, &p.pinned)
			if windowAction(c, pin, pinIcon, p.tr("置顶"), hint, p.pinned).Changed() {
				p.win.SetAlwaysOnTop(p.pinned)
			}
			minimizeHint := p.tr("最小化")
			if p.hideToTray() {
				minimizeHint = p.tr("最小化到托盘")
			}
			windowAction(c, ui.ButtonBase(c), minimizeIcon, p.tr("最小化"), minimizeHint, false).OnClick(p.minimize)
			windowAction(c, ui.ButtonBase(c), closeIcon, p.tr("关闭"), p.tr("关闭并退出"), false).OnClick(mygo.App.Quit)
		})
	})
}

func (p *application) identityMenu(path string, current gitIdentity, available bool) func(*ui.Menu) {
	return func(m *ui.Menu) {
		disabled := p.busy || !available
		local := repoScope(current.NameScope) || repoScope(current.EmailScope)
		if path != "" {
			if m.Item(p.tr("项目详情")).Chosen() {
				p.openProject(path)
			}
			m.Separator()
			if m.Item(p.tr("继承全局")).Radio(!local).Disabled(disabled || !local).Chosen() {
				p.applyIdentity(path, nil)
			}
			m.Separator()
		}
		for _, a := range p.accounts {
			if m.Item(a.Label + " · " + a.Name).Subtitle(a.Email).Radio(current.matches(a) && (path == "" || local)).Disabled(disabled).Chosen() {
				p.applyIdentity(path, &a)
			}
		}
		if available && current.Name != "" && current.Email != "" && !slices.ContainsFunc(p.accounts, current.matches) {
			m.Separator()
			if m.Item(p.tr("保存当前身份为账号…")).Chosen() {
				p.beginAccount(-1)
				p.draft.Name, p.draft.Email = current.Name, current.Email
			}
		} else if len(p.accounts) == 0 && m.Item(p.tr("新增账号…")).Chosen() {
			p.beginAccount(-1)
		}
		if path != "" {
			m.Separator()
			if m.Item(p.tr("移出目录列表")).Disabled(p.busy).Chosen() {
				p.updateUI(func() {
					if !p.busy {
						p.removeDirectory(path)
					}
				})
			}
		}
	}
}
func (p *application) beginAccount(index int) {
	p.page, p.editing, p.formError = "account", index, ""
	p.draft = identity{}
	if index >= 0 {
		p.draft = p.accounts[index]
	}
}

func (p *application) uses(id int) int {
	count := 0
	for _, r := range p.projects {
		if r.current.matches(p.accounts[id]) && (repoScope(r.current.NameScope) || repoScope(r.current.EmailScope)) {
			count++
		}
	}
	return count
}

func (p *application) removeAccount(id int) {
	if p.global.matches(p.accounts[id]) || p.uses(id) > 0 {
		return
	}
	list := slices.Delete(slices.Clone(p.accounts), id, id+1)
	if err := p.persist(list, p.projects); err != nil {
		p.setNotice(err.Error())
	}
}
func (p *application) accountList(c *ui.Context) {
	t := c.Theme()
	ui.Row(c).Padding(2, 0, 8).Gap(8).AlignItems(ui.Center).Children(func() {
		windowAction(c, ui.ButtonBase(c), backIcon, p.tr("返回项目"), p.tr("返回项目"), false).OnClick(func() { p.page = "" })
		ui.Text(c, p.tr("账号管理")).Grow(1).Bold().FontSize(18)
		ui.PrimaryButton(c, p.tr("＋ 新增账号")).Padding(5, 10).OnClick(func() { p.beginAccount(-1) })
	})
	remove := -1
	ui.Scroll(c.Key("accounts")).Grow(1).Children(func() {
		for i, a := range p.accounts {
			count := p.uses(i)
			usage := p.tr("未使用")
			if count > 0 {
				usage = strconv.Itoa(count) + p.tr(" 个项目")
			}
			if p.global.matches(a) {
				usage = p.tr("全局默认")
			}
			menu := func(m *ui.Menu) {
				if m.Item(p.tr("编辑账号")).Chosen() {
					p.beginAccount(i)
				}
				if m.Item(p.tr("删除账号")).Disabled(p.global.matches(a) || count > 0).Chosen() {
					remove = i
				}
			}
			row := ui.Row(c.Key(a.Label)).Height(73).Shrink(0).Gap(10).Padding(0, 5).AlignItems(ui.Center).
				Border(1, t.Border).BorderWidth(0, 0, 1, 0)
			row.Focusable().ContextMenu(menu).Children(func() {
				ink := ui.RGB(92, 80, 153)
				if p.global.matches(a) {
					ink = ui.RGB(43, 105, 80)
				}
				if t.Dark {
					ink = ink.Mix(t.Text, 0.60)
				}
				ui.Box(c).Size(35, 35).Radius(9).Background(t.Surface.Mix(ink, 0.12)).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
					ui.Icon(c, accountIcon).Size(19, 19).TextColor(ink)
				})
				ui.Column(c).Grow(1).MinWidth(0).Gap(4).Children(func() {
					ui.Row(c).Gap(7).AlignItems(ui.Center).Children(func() {
						ui.Text(c, a.Label).Bold().FontSize(15).MaxWidth(118).MaxLines(1).Ellipsis("…").Tooltip(a.Label)
						ui.Text(c, p.displayValue(a.Name)).Grow(1).MinWidth(0).MaxLines(1).Ellipsis("…").FontSize(13).TextColor(t.TextMuted)
					})
					ui.Text(c, p.displayValue(a.Email)).FontSize(13).TextColor(t.TextMuted).MaxLines(1).Ellipsis("…")
				})
				tag := ui.Text(c, usage).Shrink(0).FontSize(12).Radius(5).Padding(2, 7).TextColor(t.TextMuted).Background(t.Background)
				if p.global.matches(a) {
					tag.TextColor(ink).Background(t.Surface.Mix(ink, 0.08))
				} else {
					ui.Button(c, p.tr("设为全局")).Shrink(0).FontSize(12).Padding(4, 8).
						Label(a.Label + " · " + p.tr("设为全局")).
						Tooltip(p.tr("设为 Git 全局提交身份，未单独配置身份的项目将使用此账号。")).
						Disabled(p.busy).OnClick(func() { p.applyIdentity("", &a) })
				}
				windowAction(c, ui.ButtonBase(c), moreIcon, a.Label+p.tr("账号操作"), p.tr("账号操作"), false).Menu(menu)
			})
		}
	})
	if remove >= 0 {
		p.removeAccount(remove)
	}
}

func (p *application) accountForm(c *ui.Context) {
	title := p.tr("新增账号")
	if p.editing >= 0 {
		title = p.tr("编辑账号")
	}
	t := c.Theme()
	ui.Row(c).Padding(2, 0, 8).Gap(8).AlignItems(ui.Center).Children(func() {
		windowAction(c, ui.ButtonBase(c), backIcon, p.tr("返回账号管理"), p.tr("返回账号管理"), false).OnClick(func() { p.page = "accounts" })
		ui.Text(c, title).Bold().FontSize(18)
	})
	ui.Scroll(c.Key("account-fields")).Grow(1).Padding(4, 5).Gap(16).Children(func() {
		accountField(c, p.tr("账号标签"), p.tr("例如：工作"), &p.draft.Label)
		accountField(c, p.tr("名字"), "Lin Chen", &p.draft.Name)
		accountField(c, p.tr("邮箱"), "lin@example.org", &p.draft.Email)
		if p.formError != "" {
			ui.Text(c, p.formError).TextColor(t.Danger).FontSize(13)
		}
	})
	ui.Row(c).Justify(ui.End).Gap(8).Padding(10, 0, 0).Border(1, t.Border).BorderWidth(1, 0, 0, 0).Children(func() {
		ui.Button(c, p.tr("取消")).Padding(5, 14).OnClick(func() { p.page = "accounts" })
		ui.PrimaryButton(c, p.tr("保存账号")).Padding(5, 14).OnClick(p.saveAccount)
	})
}

func accountField(c *ui.Context, label, example string, value *string) {
	ui.Field(c.Key(label), label, func() {
		ui.TextInput(c, value).Label(label).Placeholder(example).Height(36).Radius(6).Padding(6, 10)
	}).Gap(5)
}

func (p *application) saveAccount() {
	a := identity{strings.TrimSpace(p.draft.Label), strings.TrimSpace(p.draft.Name), strings.TrimSpace(p.draft.Email)}
	if a.Label == "" || a.Name == "" {
		p.formError = p.tr("请填写账号标签和名字。")
		return
	}
	address, err := mail.ParseAddress(a.Email)
	if err != nil || address.Address != a.Email {
		p.formError = p.tr("请输入有效的邮箱地址。")
		return
	}
	for i, existing := range p.accounts {
		if i != p.editing && existing.Label == a.Label {
			p.formError = p.tr("该账号标签已存在。")
			return
		}
	}
	list := slices.Clone(p.accounts)
	if p.editing < 0 {
		list = append(list, a)
	} else {
		list[p.editing] = a
	}
	if err := p.persist(list, p.projects); err != nil {
		p.formError = err.Error()
		return
	}
	p.page = "accounts"
}

func trayIcon() []byte {
	src, err := png.Decode(bytes.NewReader(appIcon))
	if err != nil {
		panic(err)
	}
	// Downsample the original before handing it to the notification area.
	// 64 pixels also covers high-DPI taskbars without a large native bitmap.
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	draw.CatmullRom.Scale(img, img.Bounds(), src, src.Bounds(), draw.Src, nil)
	var out bytes.Buffer
	if err := png.Encode(&out, img); err != nil {
		panic(err)
	}
	return out.Bytes()
}

// Linux desktops may have no StatusNotifier host even when AppIndicator is
// installed. Keep their window reachable from the taskbar at all times.
func (p *application) hideToTray() bool {
	return runtime.GOOS != "linux" && p.tray != nil
}

func (p *application) minimize() {
	p.savePanelSize()
	if p.hideToTray() {
		p.win.Hide()
	} else {
		p.win.Minimize()
	}
}

func (p *application) show() {
	if mygo.Updater.Enabled() && time.Since(p.lastUpdateCheck) >= 6*time.Hour {
		p.checkUpdates()
	}
	p.loginEnabled = mygo.App.OpenAtLogin()
	p.refresh()
	var a mygo.Rectangle
	if p.tray != nil {
		a = p.tray.Bounds()
	}
	if p.win.IsMinimized() {
		p.win.Restore()
	}
	if a.Width > 0 && a.Height > 0 {
		area := mygo.Screen.DisplayNearestPoint(mygo.Point{X: a.X + a.Width/2, Y: a.Y + a.Height/2}).WorkArea
		b := p.win.Bounds()
		x, y := a.X+(a.Width-b.Width)/2, a.Y+a.Height+6
		if y+b.Height > area.Y+area.Height {
			y = a.Y - b.Height - 6
		}
		p.win.SetPosition(max(area.X, min(x, area.X+area.Width-b.Width)), max(area.Y, min(y, area.Y+area.Height-b.Height)))
	} else {
		p.win.Center()
	}
	p.win.Show()
}
func main() {
	p := &application{}
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
		p.tray, err = mygo.NewTray(mygo.TrayOptions{Icon: trayIcon(), ToolTip: appName})
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
