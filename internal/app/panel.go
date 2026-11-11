package app

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

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
							if strings.Contains(strings.ToLower(filepath.Base(r.Path)), query) {
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
