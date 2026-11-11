package app

import (
	"path/filepath"

	"github.com/egoist/mygo/ui"

	"gus/internal/gitconfig"
)

func (p *application) projectRow(c *ui.Context, r gitconfig.Project, columns projectColumns) {
	t := c.Theme()
	a := r.Current
	row := ui.Row(c.Key(r.Path)).Height(38).Shrink(0).Padding(0, 12, 0, 8).Gap(8).AlignItems(ui.Center)
	row.Border(1, t.Surface.Mix(t.Text, 0.08)).BorderWidth(0, 0, 1, 0)
	if row.Hovered() {
		row.Background(t.SurfaceHover)
	}
	row.Focusable().Tooltip(p.description(r)).ContextMenu(p.identityMenu(r.Path, r.Current, r.Problem == "")).Children(func() {
		ui.Text(c, filepath.Base(r.Path)).Width(columns.repo).Shrink(0).MaxLines(1).Ellipsis("…").Tooltip(p.description(r)).Bold()
		ui.Text(c, p.displayValue(a.Name)).Width(columns.name).Shrink(0).MaxLines(1).Ellipsis("…")
		ui.Text(c, p.displayValue(a.Email)).Grow(1).MinWidth(0).MaxLines(1).Ellipsis("…").FontSize(13).TextColor(t.TextMuted)
		p.platformStack(c, r)
		p.sourceMark(c, r)
	})
	if row.DoubleClicked() || row.Shortcut(0, ui.KeyEnter) {
		p.openProject(r.Path)
	}
}

func (p *application) sourceMark(c *ui.Context, r gitconfig.Project) {
	t := c.Theme()
	source := r.Current.Source()
	icon, ink := minimizeIcon, t.TextMuted
	switch source {
	case "继承":
		icon, ink = inheritIcon, ui.RGB(41, 108, 77)
	case "专用":
		icon, ink = folderIcon, ui.RGB(110, 82, 166)
	case "混合":
		icon, ink = layersIcon, ui.RGB(159, 109, 34)
	}
	if r.Problem != "" {
		source, icon, ink = "不可用", alertIcon, t.Danger
	}
	if t.Dark {
		ink = ink.Mix(t.Text, 0.50)
	}
	hint := p.tr("身份来源") + " · " + p.tr(source)
	if r.Problem == "" {
		hint += "\n" + p.tr("姓名：") + p.displayValue(r.Current.NameScope) + "\n" + p.tr("邮箱：") + p.displayValue(r.Current.EmailScope)
	} else {
		hint += "\n" + r.Problem
	}
	ui.Box(c.Key("source")).Size(24, 24).Margin(0, 9).Radius(6).Background(t.Surface.Mix(ink, 0.07)).
		Label(p.tr("身份来源") + " · " + p.tr(source)).Tooltip(hint).
		AlignItems(ui.Center).Justify(ui.Center).Children(func() {
		ui.Icon(c, icon).Size(15, 15).TextColor(ink)
	})
}

func platformInk(c *ui.Context) ui.Color {
	ink := ui.RGB(57, 91, 163)
	if c.Theme().Dark {
		ink = ink.Mix(c.Theme().Text, 0.63)
	}
	return ink
}
