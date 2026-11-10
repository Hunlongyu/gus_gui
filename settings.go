package main

import (
	"context"
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

//go:embed CHANGELOG.md
var releaseHistory string

func (p *application) appSettings(c *ui.Context) {
	t := c.Theme()
	ui.Row(c).Padding(2, 0, 8).Gap(8).AlignItems(ui.Center).Children(func() {
		windowAction(c, ui.ButtonBase(c), backIcon, "返回项目", "返回项目", false).OnClick(func() { p.page = "" })
		ui.Text(c, "设置").Bold().FontSize(18)
	})
	ui.Row(c).Height(54).Padding(0, 8).AlignItems(ui.Center).Children(func() {
		ui.Text(c, "开机自启").Grow(1)
		if ui.Switch(c, &p.loginEnabled).Label("开机自启").Tooltip("登录 Windows 后启动到托盘").Changed() {
			p.loginError = ""
			if err := mygo.App.SetOpenAtLogin(p.loginEnabled); err != nil {
				p.loginError = "无法更改开机自启：" + err.Error()
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
		version = "本地构建"
	}
	ui.Row(c).Padding(12, 8).Gap(8).AlignItems(ui.Center).Children(func() {
		ui.Column(c).Grow(1).Gap(4).Children(func() {
			ui.Text(c, "应用更新")
			ui.Text(c, version).FontSize(12).TextColor(t.TextMuted)
		})
		title := "检查更新"
		if p.checkingUpdate {
			title = "正在检查…"
		}
		ui.Button(c, "更新日志").Padding(5, 10).OnClick(func() { p.page = "changelog" })
		ui.Button(c, title).Disabled(p.checkingUpdate || p.installingUpdate).Padding(5, 10).OnClick(p.checkUpdates)
	})
	if p.updateStatus != "" {
		ui.Text(c, p.updateStatus).Padding(0, 8, 6).FontSize(12).TextColor(t.TextMuted)
	}
	if up := p.availableUpdate; up != nil {
		ui.Scroll(c.Key("update-notes")).Height(94).Padding(6).Background(t.Surface).Children(func() { ui.Text(c, up.Notes).FontSize(13) })
		ui.PrimaryButton(c, "安装并重启").Disabled(p.checkingUpdate || p.busy || p.installingUpdate).OnClick(p.installUpdate)
	}
	ui.Divider(c)
	ui.Column(c).Padding(12, 8).Gap(6).Children(func() {
		ui.Text(c, "关于").FontSize(12).TextColor(t.TextMuted)
		ui.Text(c, "Git User Switch").Bold().FontSize(16)
		ui.Text(c, "管理全局和项目的 Git 提交身份。").FontSize(13).TextColor(t.TextMuted)
	})
}

func (p *application) checkUpdates() {
	if p.checkingUpdate || p.installingUpdate {
		return
	}
	if !mygo.Updater.Enabled() {
		p.updateStatus = "此版本未启用在线更新。"
		return
	}
	p.checkingUpdate = true
	p.lastUpdateCheck = time.Now()
	p.updateStatus = ""
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		update, err := mygo.Updater.Check(ctx)
		p.win.Update(func() {
			p.checkingUpdate = false
			switch {
			case err != nil:
				p.updateStatus = "检查失败：" + err.Error()
			case update != nil:
				p.availableUpdate = update
				p.updateStatus = "发现新版本 " + update.Version
			default:
				p.availableUpdate = nil
				p.updateStatus = "已是最新版本。"
			}
		})
	}()
}

func (p *application) changelogView(c *ui.Context) {
	ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
		windowAction(c, ui.ButtonBase(c), backIcon, "返回设置", "返回设置", false).OnClick(func() { p.page = "settings" })
		ui.Text(c, "更新日志").Bold().FontSize(18)
	})
	ui.Scroll(c).Grow(1).Padding(4).Children(func() {
		ui.Text(c, strings.TrimSpace(strings.TrimPrefix(releaseHistory, "# 更新日志"))).FontSize(14)
	})
}

func (p *application) installUpdate() {
	up := p.availableUpdate
	if up == nil || p.busy || p.checkingUpdate || p.installingUpdate {
		return
	}
	p.installingUpdate, p.busy = true, true
	p.updateStatus = "正在下载 " + up.Version
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
			p.updateUI(func() { p.updateStatus = fmt.Sprintf("正在更新 %s · %d%%", up.Version, percent) })
		})
		p.updateUI(func() {
			p.installingUpdate, p.busy = false, false
			if err != nil {
				p.updateStatus = "安装失败：" + err.Error()
				p.importDirectories(nil)
				return
			}
			mygo.App.Relaunch()
		})
	}()
}
