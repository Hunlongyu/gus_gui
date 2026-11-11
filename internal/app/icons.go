package app

import "github.com/egoist/mygo/ui"

func lineIcon(shapes string) *ui.SVG {
	return ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">` + shapes + `</svg>`))
}

var (
	folderIcon       = lineIcon(`<path d="M3 7V5.5A1.5 1.5 0 0 1 4.5 4H9l2.2 2.5h8.3A1.5 1.5 0 0 1 21 8v10.5a1.5 1.5 0 0 1-1.5 1.5h-15A1.5 1.5 0 0 1 3 18.5z"/>`)
	layersIcon       = lineIcon(`<path d="m12 3 9 5-9 5-9-5zM3 12l9 5 9-5M3 16l9 5 9-5"/>`)
	alertIcon        = lineIcon(`<path d="m12 3 10 18H2zM12 9v5m0 3v.1"/>`)
	githubIcon       = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M12 .5C5.37.5 0 5.87 0 12.5c0 5.3 3.44 9.8 8.21 11.39.6.11.82-.26.82-.58v-2.23c-3.34.73-4.04-1.42-4.04-1.42-.55-1.39-1.33-1.76-1.33-1.76-1.09-.75.08-.73.08-.73 1.21.08 1.84 1.24 1.84 1.24 1.07 1.83 2.81 1.3 3.49.99.11-.78.42-1.3.76-1.6-2.67-.3-5.47-1.34-5.47-5.93 0-1.31.47-2.38 1.24-3.22-.12-.3-.54-1.53.12-3.18 0 0 1.01-.32 3.3 1.23A11.5 11.5 0 0 1 12 6.3c1.02 0 2.05.14 3.01.4 2.29-1.55 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.77.84 1.24 1.91 1.24 3.22 0 4.6-2.81 5.62-5.49 5.92.43.38.82 1.1.82 2.22v3.3c0 .32.22.7.83.58A12 12 0 0 0 24 12.5c0-6.63-5.37-12-12-12z"/></svg>`))
	gitlabIcon       = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path fill="#e24329" d="m12 22-4.4-13h8.8z"/><path fill="#fc6d26" d="M12 22 2 9h5.6zM12 22 22 9h-5.6z"/><path fill="#fca326" d="M12 22 1 14l1-5zM12 22l11-8-1-5z"/><path fill="#e24329" d="m2 9 2-7 3.6 7zM22 9l-2-7-3.6 7z"/></svg>`))
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
