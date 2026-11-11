package app

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"gus/internal/gitconfig"
	"gus/internal/pathutil"
)

func projectPlatforms(r gitconfig.Project) []string {
	var platforms []string
	for _, remote := range r.Remotes {
		for _, address := range remote.PushURLs {
			platform := gitconfig.RemotePlatform(address)
			if !slices.Contains(platforms, platform) {
				platforms = append(platforms, platform)
			}
		}
	}
	return platforms
}

func (p *application) pushPlatforms(r gitconfig.Project) string {
	platforms := projectPlatforms(r)
	for i := range platforms {
		platforms[i] = p.tr(platforms[i])
	}
	if r.RemoteError != "" {
		platforms = append(platforms, p.tr("读取远程信息失败"))
	}
	if len(platforms) == 0 {
		return p.tr("未配置远程仓库")
	}
	return strings.Join(platforms, " · ")
}

func (p *application) projectAccount(r gitconfig.Project) string {
	var labels []string
	for _, a := range p.accounts {
		if r.Current.Matches(a) {
			labels = append(labels, a.Label)
		}
	}
	if len(labels) == 0 {
		return p.tr("未匹配已保存账号")
	}
	return strings.Join(labels, " · ")
}

func (p *application) openProject(path string) {
	p.detailPath, p.page = path, "project"
	p.refresh()
}

func (p *application) projectDetails(c *ui.Context) {
	t := c.Theme()
	ui.Row(c).Padding(2, 0, 8).Gap(8).AlignItems(ui.Center).Children(func() {
		windowAction(c, ui.ButtonBase(c), backIcon, p.tr("返回项目"), p.tr("返回项目"), false).OnClick(func() { p.page = "" })
		ui.Text(c, p.tr("项目详情")).Grow(1).Bold().FontSize(18)
		label := p.tr("刷新")
		if p.busy {
			label = p.tr("读取中…")
		}
		ui.Button(c, label).Disabled(p.busy).OnClick(p.refresh)
	})
	index := slices.IndexFunc(p.projects, func(r gitconfig.Project) bool { return pathutil.Same(r.Path, p.detailPath) })
	if index < 0 {
		ui.Text(c, p.tr("项目已移出列表")).TextColor(t.TextMuted)
		return
	}
	r := p.projects[index]
	ui.Scroll(c.Key("project-details-" + r.Path)).Grow(1).Gap(12).Children(func() {
		ui.Column(c).Gap(4).Children(func() {
			ui.Text(c, filepath.Base(r.Path)).Bold().FontSize(16).Selectable()
			ui.Text(c, filepath.Clean(r.Path)).FontSize(12).TextColor(t.TextMuted).Selectable()
		})
		if r.Problem != "" {
			ui.Text(c, r.Problem).FontSize(13).TextColor(t.Danger).Selectable()
			return
		}
		ui.Column(c).Padding(10).Gap(5).Radius(6).Background(t.Surface.Mix(platformInk(c), 0.07)).Children(func() {
			ui.Text(c, p.tr("推送平台")).FontSize(12).TextColor(t.TextMuted)
			ui.Text(c, p.pushPlatforms(r)).Bold().FontSize(16).TextColor(platformInk(c)).Selectable()
		})
		ui.Column(c).Padding(10).Gap(7).Radius(6).Background(t.Background).Children(func() {
			p.detailField(c, "提交账号", p.projectAccount(r))
			p.detailField(c, "名字", p.displayValue(r.Current.Name)+" ("+p.displayValue(r.Current.NameScope)+")")
			p.detailField(c, "邮箱", p.displayValue(r.Current.Email)+" ("+p.displayValue(r.Current.EmailScope)+")")
			branch := p.displayValue(r.Branch)
			if r.Detached {
				branch += " (" + p.tr("分离 HEAD") + ")"
			}
			if r.BranchError != "" {
				branch = p.tr("读取分支失败") + ": " + r.BranchError
			}
			p.detailField(c, "当前分支", branch)
		})
		ui.Text(c, p.tr("远程仓库")).Bold().FontSize(14)
		if r.RemoteError != "" {
			ui.Text(c, p.tr("读取远程信息失败")+": "+r.RemoteError).FontSize(12).TextColor(t.Danger).Selectable()
		} else if len(r.Remotes) == 0 {
			ui.Text(c, p.tr("未配置远程仓库")).FontSize(13).TextColor(t.TextMuted)
		}
		for _, remote := range r.Remotes {
			ui.Column(c.Key(remote.Name)).Padding(10).Gap(6).Radius(6).Border(1, t.Border).Children(func() {
				ui.Text(c, remote.Name).Bold().Selectable()
				for _, address := range remote.FetchURLs {
					p.remoteAddress(c, "拉取", address)
				}
				for _, address := range remote.PushURLs {
					p.remoteAddress(c, "推送", address)
				}
			})
		}
	})
}

func (p *application) detailField(c *ui.Context, label, value string) {
	ui.Row(c).Gap(8).Children(func() {
		ui.Text(c, p.tr(label)).Width(94).Shrink(0).FontSize(12).TextColor(c.Theme().TextMuted)
		ui.Text(c, value).Grow(1).MinWidth(0).FontSize(13).Selectable()
	})
}

func (p *application) remoteAddress(c *ui.Context, kind, address string) {
	ui.Column(c).Gap(2).Children(func() {
		ui.Text(c, p.tr(kind)+" · "+p.tr(gitconfig.RemotePlatform(address))).FontSize(12).TextColor(c.Theme().TextMuted)
		if webURL := gitconfig.RemoteWebURL(address); webURL != "" {
			ui.Link(c, address, webURL).FontSize(13).TextColor(platformInk(c)).Tooltip(p.tr("在浏览器中打开") + "\n" + webURL)
		} else {
			ui.Text(c, address).FontSize(13).Selectable().Tooltip(address)
		}
	})
}
