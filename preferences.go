package main

import (
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

type preferences struct {
	Language   string `json:"language,omitempty"`
	Appearance string `json:"appearance,omitempty"`
}

func (s preferences) normalized() preferences {
	if s.Language != "en" {
		s.Language = "zh-CN"
	}
	if s.Appearance != "light" && s.Appearance != "dark" {
		s.Appearance = "system"
	}
	return s
}

func (p *application) tr(zh string) string {
	if p.preferences.Language == "en" {
		if en, ok := english[zh]; ok {
			return en
		}
	}
	return zh
}

// Windows chrome, native menus and built-in controls follow the same choices as the panel.
func (p *application) applyPreferences() {
	prefs := p.preferences.normalized()
	mygo.App.SetLocale(prefs.Language)
	mygo.Theme.SetSource(mygo.ThemeSource(prefs.Appearance))
}

func (p *application) savePreferences(next preferences) {
	old := p.preferences
	p.preferences = next.normalized()
	if err := p.persist(p.accounts, p.projects); err != nil {
		p.preferences = old
		p.notice = err.Error()
		return
	}
	if p.win != nil {
		p.applyPreferences()
	}
}

func (p *application) preferenceControls(c *ui.Context) {
	prefs := p.preferences.normalized()
	language := "简体中文"
	if prefs.Language == "en" {
		language = "English"
	}
	ui.Row(c).Height(38).Padding(0, 8).AlignItems(ui.Center).Children(func() {
		ui.Text(c, p.tr("语言")).Grow(1)
		ui.MenuButton(c, language, func(m *ui.Menu) {
			if m.Item("简体中文").Radio(prefs.Language == "zh-CN").Chosen() {
				prefs.Language = "zh-CN"
				p.savePreferences(prefs)
			}
			if m.Item("English").Radio(prefs.Language == "en").Chosen() {
				prefs.Language = "en"
				p.savePreferences(prefs)
			}
		})
	})
	themes := []struct{ value, label string }{{"system", "跟随系统"}, {"light", "浅色"}, {"dark", "深色"}}
	selected := ""
	for _, theme := range themes {
		if theme.value == prefs.Appearance {
			selected = p.tr(theme.label)
		}
	}
	ui.Row(c).Height(38).Padding(0, 8).AlignItems(ui.Center).Children(func() {
		ui.Text(c, p.tr("外观")).Grow(1)
		ui.MenuButton(c, selected, func(m *ui.Menu) {
			for _, theme := range themes {
				if m.Item(p.tr(theme.label)).Radio(theme.value == prefs.Appearance).Chosen() {
					prefs.Appearance = theme.value
					p.savePreferences(prefs)
				}
			}
		})
	})
}

var english = map[string]string{
	"查看更新 · ":     "View update · ",
	"项目":          "Projects",
	"搜索项目名":       "Search projects",
	"账号管理":        "Accounts",
	"设置":          "Settings",
	"应用设置":        "App settings",
	"添加目录":        "Add folder",
	"添加项目目录":      "Add project folder",
	"添加 Git 项目目录": "Add Git folders",
	"仓库":          "Repository",
	"提交姓名":        "Git name",
	"邮箱":          "Email",
	"身份来源":        "Source",
	"拖入 Git 项目目录，或点击右上角添加": "Drop Git folders here or use Add folder.",
	"没有匹配的项目":              "No matching projects",
	"不可用":                  "Offline",
	"全局":                   "Global",
	"置顶面板":                 "Pin panel",
	"取消置顶":                 "Unpin panel",
	"置顶":                   "Pin",
	"最小化":                  "Minimize",
	"最小化到托盘":               "Minimize to tray",
	"关闭":                   "Close",
	"关闭并退出":                "Close and quit",
	"继承全局":                 "Use global identity",
	"保存当前身份为账号…":           "Save current identity…",
	"新增账号…":                "Add account…",
	"移出目录列表":               "Remove from list",
	"返回项目":                 "Back to projects",
	"＋ 新增账号":               "+ Add account",
	"未使用":                  "Unused",
	" 个项目":                 " projects",
	"全局默认":                 "Global default",
	"编辑账号":                 "Edit account",
	"删除账号":                 "Delete account",
	"账号操作":                 "Account actions",
	"新增账号":                 "Add account",
	"返回账号管理":               "Back to accounts",
	"账号标签":                 "Label",
	"例如：工作":                "e.g. Work",
	"取消":                   "Cancel",
	"保存账号":                 "Save account",
	"请填写账号标签和提交姓名。":     "Enter a label and Git name.",
	"请输入有效的邮箱地址。":       "Enter a valid email address.",
	"该账号标签已存在。":         "This label is already in use.",
	"打开面板":              "Open panel",
	"退出":                "Quit",
	"开机自启":              "Start at login",
	"登录 Windows 后启动到托盘": "Open in the tray at Windows sign-in",
	"无法更改开机自启：":         "Could not change startup setting: ",
	"本地构建":              "Local build",
	"应用更新":              "App updates",
	"检查更新":              "Check for updates",
	"正在检查…":             "Checking…",
	"更新日志":              "Release notes",
	"安装并重启":             "Install and restart",
	"关于":                "About",
	"管理全局和项目的 Git 提交身份。": "Manage global and project Git identities.",
	"此版本未启用在线更新。":        "Online updates are unavailable for this build.",
	"检查失败：":              "Update check failed: ",
	"发现新版本 ":             "New version: ",
	"已是最新版本。":            "You're up to date.",
	"返回设置":               "Back to settings",
	"正在下载 ":              "Downloading ",
	"正在更新 %s · %d%%":     "Updating %s · %d%%",
	"安装失败：":              "Installation failed: ",
	"未设置":                "Unset",
	"混合":                 "Mixed",
	"专用":                 "Local",
	"继承":                 "Global",
	"语言":                 "Language",
	"外观":                 "Appearance",
	"跟随系统":               "System",
	"浅色":                 "Light",
	"深色":                 "Dark",
	"姓名：":                "Name: ",
	"邮箱：":                "Email: ",
	"请选择目录：%s":           "Choose a folder: %s",
	"读取全局身份失败：%v":        "Cannot read global identity: %v",
	"切换失败或未完全生效：":        "Could not apply identity: ",
	"无法读取 %s：%v":         "Cannot load %s: %v",
	"保存失败：%w":            "Could not save: %w",
}
