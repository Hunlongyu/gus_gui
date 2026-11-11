package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

const (
	appDeveloper     = "Hunlongyu"
	appDeveloperURL  = "https://github.com/" + appDeveloper
	appRepositoryURL = appDeveloperURL + "/gus_gui"
	appFeedbackURL   = appRepositoryURL + "/issues"
)

func (p *application) notes(markdown string) string {
	var lines []string
	visible := true
	for _, line := range strings.Split(strings.ReplaceAll(markdown, "\r\n", "\n"), "\n") {
		switch line {
		case "### 简体中文":
			visible = p.preferences.Language != "en"
			continue
		case "### English":
			visible = p.preferences.Language == "en"
			continue
		}
		if strings.HasPrefix(line, "## ") {
			visible = true
		}
		if visible {
			lines = append(lines, line)
		}
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func (p *application) appSettings(c *ui.Context) {
	t := c.Theme()
	ui.Row(c).Padding(2, 0, 8).Gap(8).AlignItems(ui.Center).Children(func() {
		windowAction(c, ui.ButtonBase(c), backIcon, p.tr("返回项目"), p.tr("返回项目"), false).OnClick(func() { p.page = "" })
		ui.Text(c, p.tr("设置")).Bold().FontSize(18)
	})
	ui.Scroll(c.Key("settings-body")).Grow(1).Gap(6).Children(func() {
		p.preferenceControls(c)
		ui.Divider(c)
		ui.Row(c).Height(48).Padding(0, 8).AlignItems(ui.Center).Children(func() {
			ui.Text(c, p.tr("开机自启")).Grow(1)
			if ui.Switch(c, &p.loginEnabled).Label(p.tr("开机自启")).Tooltip(p.tr("登录系统后启动并最小化")).Changed() {
				p.loginError = ""
				if err := mygo.App.SetOpenAtLogin(p.loginEnabled); err != nil {
					p.loginError = p.tr("无法更改开机自启：") + err.Error()
				}
				p.loginEnabled = mygo.App.OpenAtLogin()
			}
		})
		if p.loginError != "" {
			ui.Text(c, p.loginError).FontSize(12).TextColor(t.Danger)
		}
		ui.Divider(c)
		version := mygo.App.Version()
		if version == "" {
			version = p.tr("本地构建")
		}
		ui.Row(c).Padding(12, 8).Gap(8).AlignItems(ui.Center).Children(func() {
			ui.Column(c).Grow(1).Gap(4).Children(func() {
				ui.Text(c, p.tr("应用更新"))
				ui.Text(c, version).FontSize(12).TextColor(t.TextMuted)
			})
			title := p.tr("检查更新")
			if p.checkingUpdate {
				title = p.tr("正在检查…")
			}
			ui.Button(c, p.tr("更新日志")).Padding(5, 10).OnClick(func() { p.page = "changelog" })
			ui.Button(c, title).Disabled(p.checkingUpdate || p.installingUpdate).Padding(5, 10).OnClick(p.checkUpdates)
		})
		if p.updateStatus != "" {
			ui.Text(c, p.tr(p.updateStatus)+p.updateDetail).Padding(0, 8, 6).FontSize(12).TextColor(t.TextMuted)
		}
		if up := p.availableUpdate; up != nil {
			ui.Scroll(c.Key("update-notes")).Height(94).Padding(6).Background(t.Surface).Children(func() { ui.Text(c, p.notes(up.Notes)).FontSize(13) })
			if canInstallUpdate() {
				ui.PrimaryButton(c, p.tr("安装并重启")).Disabled(p.checkingUpdate || p.busy || p.installingUpdate).OnClick(p.installUpdate)
			} else {
				ui.Text(c, p.tr("单文件版请下载新版替换当前文件。")).FontSize(12).TextColor(t.TextMuted)
				ui.Link(c, p.tr("下载新版"), appRepositoryURL+"/releases/latest")
			}
		}
		ui.Divider(c)
		ui.Column(c).Padding(12, 8).Gap(6).Children(func() {
			ui.Text(c, p.tr("关于")).FontSize(12).TextColor(t.TextMuted)
			ui.Text(c, appName).Bold().FontSize(16)
			ui.Text(c, p.tr("管理全局和项目的 Git 提交身份。")).FontSize(13).TextColor(t.TextMuted)
			ui.Row(c).Gap(16).Children(func() {
				ui.Link(c, p.tr("开发者"), appDeveloperURL).FontSize(12).TextColor(platformInk(c))
				ui.Link(c, p.tr("仓库地址"), appRepositoryURL).FontSize(12).TextColor(platformInk(c))
				ui.Link(c, p.tr("反馈地址"), appFeedbackURL).FontSize(12).TextColor(platformInk(c))
			})
		})
	})
}

func (p *application) checkUpdates() {
	if p.checkingUpdate || p.installingUpdate {
		return
	}
	if !mygo.Updater.Enabled() {
		p.updateStatus, p.updateDetail = "此版本未启用在线更新。", ""
		return
	}
	p.checkingUpdate = true
	p.lastUpdateCheck = time.Now()
	p.updateStatus, p.updateDetail = "", ""
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		update, err := mygo.Updater.Check(ctx)
		p.win.Update(func() {
			p.checkingUpdate = false
			switch {
			case err != nil:
				p.updateStatus, p.updateDetail = "检查失败：", err.Error()
			case update != nil:
				p.availableUpdate = update
				p.updateStatus, p.updateDetail = "发现新版本 ", update.Version
				p.pendingToasts = append(p.pendingToasts, ui.Toast{
					ID: "available-update", Title: p.tr("发现新版本 ") + update.Version,
					Action: p.tr("查看更新"), OnAction: func() { p.page = "settings" }, Timeout: 8 * time.Second,
				})
			default:
				p.availableUpdate = nil
				p.updateStatus, p.updateDetail = "已是最新版本。", ""
			}
		})
	}()
}

func (p *application) changelogView(c *ui.Context) {
	ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
		windowAction(c, ui.ButtonBase(c), backIcon, p.tr("返回设置"), p.tr("返回设置"), false).OnClick(func() { p.page = "settings" })
		ui.Text(c, p.tr("更新日志")).Bold().FontSize(18)
	})
	ui.Scroll(c).Grow(1).Padding(4).Children(func() {
		ui.Text(c, p.notes(strings.TrimPrefix(p.releaseHistory, "# 更新日志"))).FontSize(14)
	})
}

func (p *application) installUpdate() {
	up := p.availableUpdate
	if up == nil || p.busy || p.checkingUpdate || p.installingUpdate || !canInstallUpdate() {
		return
	}
	p.installingUpdate, p.busy = true, true
	p.updateStatus, p.updateDetail = "正在下载 ", up.Version
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		last := -1
		err := up.Install(ctx, func(done, total int64) {
			if total <= 0 {
				return
			}
			percent := int(done * 100 / total)
			if percent == last {
				return
			}
			last = percent
			p.updateUI(func() {
				p.updateStatus, p.updateDetail = "正在下载 ", fmt.Sprintf("%s · %d%%", up.Version, percent)
			})
		})
		p.updateUI(func() {
			p.installingUpdate, p.busy = false, false
			if err != nil {
				p.updateStatus, p.updateDetail = "安装失败：", err.Error()
				p.importDirectories(nil)
				return
			}
			mygo.App.Relaunch()
		})
	}()
}

// MyGo replaces the archive's canonical executable but relaunches the original
// filename. Versioned standalone downloads therefore use manual replacement;
// otherwise an update would restart the old binary and add package sidecars.
func canInstallUpdate() bool {
	exe, err := os.Executable()
	return err == nil && updateExecutableName(filepath.Base(exe), runtime.GOOS)
}

func updateExecutableName(name, goos string) bool {
	switch goos {
	case "windows":
		return strings.EqualFold(name, "gus.exe")
	case "linux", "darwin":
		return name == "gus"
	default:
		return false
	}
}
