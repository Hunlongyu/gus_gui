package app

import (
	"bytes"
	"image"
	"image/png"
	"runtime"
	"time"

	"github.com/egoist/mygo"
	"golang.org/x/image/draw"
)

func trayIcon(appIcon []byte) []byte {
	src, err := png.Decode(bytes.NewReader(appIcon))
	if err != nil {
		panic(err)
	}
	// Downsample the original before handing it to the notification area.
	// 64 pixels also covers high-DPI taskbars without a large native bitmap.
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	draw.CatmullRom.Scale(img, img.Bounds(), src, src.Bounds(), draw.Src, nil)
	var out bytes.Buffer
	if err := png.Encode(&out, img); err != nil {
		panic(err)
	}
	return out.Bytes()
}

// Linux desktops may have no StatusNotifier host even when AppIndicator is
// installed. Keep their window reachable from the taskbar at all times.
func (p *application) hideToTray() bool {
	return runtime.GOOS != "linux" && p.tray != nil
}

func (p *application) minimize() {
	p.savePanelSize()
	if p.hideToTray() {
		p.win.Hide()
	} else {
		p.win.Minimize()
	}
}

func (p *application) show() {
	if mygo.Updater.Enabled() && time.Since(p.lastUpdateCheck) >= 6*time.Hour {
		p.checkUpdates()
	}
	p.loginEnabled = mygo.App.OpenAtLogin()
	p.refresh()
	var a mygo.Rectangle
	if p.tray != nil {
		a = p.tray.Bounds()
	}
	if p.win.IsMinimized() {
		p.win.Restore()
	}
	if a.Width > 0 && a.Height > 0 {
		area := mygo.Screen.DisplayNearestPoint(mygo.Point{X: a.X + a.Width/2, Y: a.Y + a.Height/2}).WorkArea
		b := p.win.Bounds()
		x, y := a.X+(a.Width-b.Width)/2, a.Y+a.Height+6
		if y+b.Height > area.Y+area.Height {
			y = a.Y - b.Height - 6
		}
		p.win.SetPosition(max(area.X, min(x, area.X+area.Width-b.Width)), max(area.Y, min(y, area.Y+area.Height-b.Height)))
	} else {
		p.win.Center()
	}
	p.win.Show()
}
