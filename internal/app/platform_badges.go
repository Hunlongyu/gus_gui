package app

import (
	"strconv"

	"github.com/egoist/mygo/ui"

	"gus/internal/gitconfig"
)

func (p *application) platformStack(c *ui.Context, r gitconfig.Project) {
	t := c.Theme()
	platforms := projectPlatforms(r)
	summary := p.pushPlatforms(r)
	if r.Problem != "" {
		platforms, summary = nil, p.tr("不可用")
	}
	anchor := ui.Box(c.Key("platforms")).Size(56, 26).
		Label(p.tr("推送平台：") + summary).Description(summary)
	anchor.Children(func() {
		if len(platforms) == 0 {
			icon := minimizeIcon
			if r.RemoteError != "" || r.Problem != "" {
				icon = alertIcon
			}
			ui.Box(c).Size(24, 24).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
				ui.Icon(c, icon).Size(14, 14).TextColor(t.TextMuted)
			})
		}
		for i := 0; i < min(3, len(platforms)); i++ {
			name := platforms[i]
			if i == 2 && len(platforms) > 3 {
				name = "+" + strconv.Itoa(len(platforms)-2)
			}
			platformAvatar(c, name).Absolute().Left(float32(i) * 15).Top(1).ZIndex(i)
		}
	})
	ui.TooltipBase(c, anchor, func(tip ui.Element) {
		tip.AttachTo(anchor, ui.AnchorLeft, ui.AnchorRight).Margin(0, 8, 0, 0).
			Width(220).Padding(12).Radius(8).Background(t.Surface).Border(1, t.Border).
			Shadow(0, 4, 14, 0, ui.RGBA(0, 0, 0, 0.18)).Children(func() {
			ui.Column(c).Gap(8).Children(func() {
				ui.Text(c, p.tr("推送平台")).Bold().FontSize(13).TextColor(t.Text)
				for _, platform := range platforms {
					ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
						platformAvatar(c, platform)
						ui.Text(c, p.tr(platform)).Grow(1).MinWidth(0).FontSize(13).TextColor(t.Text)
					})
				}
				if len(platforms) == 0 || r.RemoteError != "" {
					message := summary
					if len(platforms) > 0 {
						message = p.tr("读取远程信息失败")
					}
					ui.Text(c, message).FontSize(12).TextColor(t.TextMuted)
				}
			})
		})
	})
}

func platformAvatar(c *ui.Context, name string) ui.Element {
	fill, text := ui.RGB(97, 105, 123), ""
	var icon *ui.SVG
	switch name {
	case "GitHub":
		fill, icon = ui.RGB(40, 45, 55), githubIcon
	case "GitLab":
		fill, icon = ui.RGB(219, 101, 35), gitlabIcon
	case "Gitee":
		fill, icon = ui.RGB(198, 45, 52), giteeIcon
	case "Bitbucket":
		fill, icon = ui.RGB(36, 103, 211), bitbucketIcon
	case "Codeberg":
		fill, icon = ui.RGB(42, 123, 165), codebergIcon
	case "GitGud":
		fill = ui.RGB(10, 12, 16)
	case "GitCode":
		fill, icon = ui.RGB(210, 69, 57), gitcodeIcon
	case "Azure DevOps":
		fill, icon = ui.RGB(25, 119, 192), azureDevopsIcon
	case "Gitea":
		fill, icon = ui.RGB(75, 132, 42), giteaIcon
	case "Forgejo":
		fill, icon = ui.RGB(168, 84, 30), forgejoIcon
	case "SourceHut":
		fill, icon = ui.RGB(40, 45, 55), sourcehutIcon
	case "SourceForge":
		fill, icon = ui.RGB(201, 83, 0), sourceforgeIcon
	case "Launchpad":
		fill, icon = ui.RGB(139, 98, 12), launchpadIcon
	case "CODING":
		fill, icon = ui.RGB(50, 54, 63), codingIcon
	case "Codeup":
		fill, icon = ui.RGB(212, 77, 0), codeupIcon
	case "本地路径":
		icon = folderIcon
	default:
		if len(name) > 0 && name[0] == '+' {
			text = name
		} else {
			icon = globeIcon
		}
	}
	return ui.Box(c).Size(24, 24).Shrink(0).Radius(12).Background(fill).Border(2, c.Theme().Surface).
		AlignItems(ui.Center).Justify(ui.Center).Children(func() {
		if name == "GitGud" {
			ui.Image(c, gitgudImage).Size(18, 18)
		} else if icon != nil {
			ui.Icon(c, icon).Size(16, 16).TextColor(ui.RGB(255, 255, 255))
		} else {
			ui.Text(c, text).Bold().FontSize(11).TextColor(ui.RGB(255, 255, 255))
		}
	})
}
