package main

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// Queue from UI callbacks, including completed background jobs. Drain once
// during a frame so rendering cannot continually restart a toast's timer.
func (p *application) notifyError(message string) {
	if message == "" {
		return
	}
	p.pendingToasts = append(p.pendingToasts, ui.Toast{
		ID: "error:" + message, Type: "error", Title: p.tr("操作未完成"),
		Description: message, Timeout: 8 * time.Second,
	})
}

func (p *application) setNotice(message string) {
	p.notice = message
	p.notifyError(message)
}

func (p *application) showNotifications(c *ui.Context) {
	for _, notification := range p.pendingToasts {
		c.AddToast(notification)
	}
	p.pendingToasts = nil
	t := c.Theme()
	ui.ToastViewportBase(c, func(viewport ui.Element, toasts []ui.Toast) {
		viewport.Padding(14).Gap(8)
		for _, notification := range toasts {
			toast := ui.ToastBase(c, notification)
			toast.Root.Row().Gap(10).Padding(12).Radius(8).Background(t.Surface).
				Border(1, t.Border).Shadow(0, 3, 14, 0, ui.RGBA(0, 0, 0, 0.18)).Children(func() {
				if notification.Type == "error" {
					ui.Icon(c, alertIcon).Size(16, 16).TextColor(t.Danger)
				}
				ui.Column(c).Grow(1).MinWidth(0).Gap(4).Children(func() {
					ui.Text(c, notification.Title).Bold().FontSize(13).TextColor(t.Text).MaxLines(1).Ellipsis("…")
					if notification.Description != "" {
						ui.Text(c, notification.Description).FontSize(12).TextColor(t.TextMuted).
							MaxLines(3).Ellipsis("…").Tooltip(notification.Description)
					}
					if notification.Action != "" {
						toast.ActionButton().Label(notification.Action).Children(func() { ui.Text(c, notification.Action).FontSize(12).TextColor(platformInk(c)) })
					}
				})
				withIcon(c, toast.CloseButton(), closeIcon, p.tr("关闭通知"), p.tr("关闭通知")).Size(20, 20)
			})
		}
	})
}
