package main

import (
	"time"

	"github.com/egoist/mygo"
)

type panelSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (p *application) trackPanelSize() {
	if s := p.windowSize; s.Width > 0 && s.Height > 0 && !p.win.IsMaximized() && !p.win.IsFullScreen() {
		b := p.win.Bounds()
		area := mygo.Screen.DisplayNearestPoint(mygo.Point{X: b.X + b.Width/2, Y: b.Y + b.Height/2}).WorkArea
		p.win.SetSize(min(max(s.Width, minPanelWidth), max(area.Width, minPanelWidth)),
			min(max(s.Height, minPanelHeight), max(area.Height, minPanelHeight)))
	}
	// MyGo saves StateKey geometry on quit. Also persist after resizing so a
	// tray-only session or interrupted process does not lose the user's size.
	p.win.OnResize(func() {
		if p.resizeTimer != nil {
			p.resizeTimer.Stop()
		}
		p.resizeTimer = time.AfterFunc(350*time.Millisecond, func() { p.updateUI(p.savePanelSize) })
	})
	p.win.OnHide(p.savePanelSize)
}

func (p *application) savePanelSize() {
	if p.resizeTimer != nil {
		p.resizeTimer.Stop()
		p.resizeTimer = nil
	}
	if p.win == nil || p.win.IsMinimized() || p.win.IsMaximized() || p.win.IsFullScreen() {
		return
	}
	w, h := p.win.Size()
	if err := p.rememberPanelSize(w, h); err != nil {
		p.setNotice(err.Error())
	}
}

func (p *application) rememberPanelSize(width, height int) error {
	next := panelSize{width, height}
	if width < minPanelWidth || height < minPanelHeight || next == p.windowSize {
		return nil
	}
	previous := p.windowSize
	p.windowSize = next
	if err := p.persist(p.accounts, p.projects); err != nil {
		p.windowSize = previous
		return err
	}
	return nil
}
