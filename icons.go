package main

import "github.com/egoist/mygo/ui"

func lineIcon(shapes string) *ui.SVG {
	return ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">` + shapes + `</svg>`))
}

var (
	backIcon         = lineIcon(`<path d="M19 12H5m6-6-6 6 6 6"/>`)
	moreIcon         = lineIcon(`<g fill="currentColor" stroke="none"><circle cx="5" cy="12" r="1.6"/><circle cx="12" cy="12" r="1.6"/><circle cx="19" cy="12" r="1.6"/></g>`)
	accountIcon      = lineIcon(`<circle cx="12" cy="8" r="3.4"/><path d="M5 20v-1a7 7 0 0 1 14 0v1"/>`)
	globeIcon        = lineIcon(`<circle cx="12" cy="12" r="9"/><ellipse cx="12" cy="12" rx="4" ry="9"/><path d="M3 12h18"/>`)
	addDirectoryIcon = lineIcon(`<path d="M3 7V5.5A1.5 1.5 0 0 1 4.5 4H9l2.2 2.5h8.3A1.5 1.5 0 0 1 21 8v10.5a1.5 1.5 0 0 1-1.5 1.5h-15A1.5 1.5 0 0 1 3 18.5z"/><path d="M12 10v7m-3.5-3.5h7"/>`)
	pinIcon          = lineIcon(`<path d="M8 3h8l-1 7 3 4v1H6v-1l3-4zM12 15v6"/>`)
	minimizeIcon     = lineIcon(`<path d="M5 12h14"/>`)
	closeIcon        = lineIcon(`<path d="m6 6 12 12M18 6 6 18"/>`)
	settingsIcon     = lineIcon(`<path d="m10 3-.6 2.2-1.6.9-2.2-.5-2 3.4 1.6 1.6v2.8l-1.6 1.6 2 3.4 2.2-.5 1.6.9L10 21h4l.6-2.2 1.6-.9 2.2.5 2-3.4-1.6-1.6v-2.8L20.4 9l-2-3.4-2.2.5-1.6-.9L14 3z"/><circle cx="12" cy="12" r="3"/>`)
)

func withIcon(c *ui.Context, button ui.Element, icon *ui.SVG, label, hint string) ui.Element {
	return button.Size(28, 24).Padding(0).AlignItems(ui.Center).Justify(ui.Center).Label(label).Tooltip(hint).Children(func() {
		ui.Icon(c, icon).Size(16, 16)
	})
}

func windowAction(c *ui.Context, button ui.Element, icon *ui.SVG, label, hint string, active bool) ui.Element {
	t := c.Theme()
	button = withIcon(c, button, icon, label, hint).Size(30, 30).Radius(6).TextColor(t.TextMuted)
	if active || button.Pressed() {
		button.Background(t.SurfacePressed).TextColor(t.Text)
	} else if button.Hovered() {
		button.Background(t.SurfaceHover).TextColor(t.Text)
	}
	return button
}
