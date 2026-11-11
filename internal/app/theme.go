package app

import "github.com/egoist/mygo/ui"

func panelTheme(base *ui.Theme) ui.Theme {
	t := *base
	if t.Dark {
		t.Background = ui.RGB(27, 27, 30)
		t.Surface = ui.RGB(37, 37, 41)
		t.Text = ui.RGB(242, 242, 245)
		t.TextMuted = ui.RGB(182, 182, 191)
		t.Accent = ui.RGB(226, 224, 240)
		t.AccentText = ui.RGB(31, 30, 37)
	} else {
		t.Background = ui.RGB(245, 245, 246)
		t.Surface = ui.RGB(255, 255, 255)
		t.Text = ui.RGB(29, 29, 32)
		t.TextMuted = ui.RGB(97, 97, 106)
		t.Accent = ui.RGB(39, 38, 46)
		t.AccentText = ui.RGB(255, 255, 255)
	}
	t.Border = t.Surface.Mix(t.Text, 0.13)
	t.SurfaceHover = t.Surface.Mix(t.Accent, 0.10)
	t.SurfacePressed = t.Surface.Mix(t.Accent, 0.20)
	t.AccentHover = t.Accent.Mix(t.Text, 0.12)
	t.AccentPressed = t.Accent.Mix(t.Text, 0.22)
	t.Selection = t.Background.Mix(t.Accent, 0.20)
	t.Focus = t.Accent
	t.Scrollbar = t.Background.Mix(t.TextMuted, 0.38)
	t.ScrollbarWidth = 4
	t.Radius = 6
	return t
}
