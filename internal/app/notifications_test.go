package app

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"gus/internal/gitconfig"
	"gus/internal/state"
)

func TestNotificationsFloatDismissAndExpire(t *testing.T) {
	for _, lang := range []string{"zh-CN", "en"} {
		p := &application{preferences: state.Preferences{Language: lang}, projects: []gitconfig.Project{{Path: "example"}}}
		v := ui.NewTester(p.view, minPanelWidth, minPanelHeight)
		before, _ := v.Find("example")
		message := "Cannot open workspace: not a git repository. " + strings.Repeat("Long directory name / ", 25)
		p.setNotice(message)
		v.Frame()
		after, _ := v.Find("example")
		if before != after {
			t.Fatalf("notification moved project rows: %+v -> %+v", before, after)
		}
		close, ok := v.Find(p.tr("关闭通知"))
		if !ok || close.Y < minPanelHeight/2 || close.X+close.W > minPanelWidth || close.Y+close.H > minPanelHeight {
			t.Fatalf("toast must float inside the bottom of the window: %+v", close)
		}
		if err := v.Click(p.tr("关闭通知")); err != nil {
			t.Fatal(err)
		}
		v.Frame()
		if v.HasText(p.tr("操作未完成")) {
			t.Fatal("dismissed notification remains visible")
		}
		// The same error from a later operation must display again.
		p.setNotice(message)
		v.Frame()
		if !v.HasText(p.tr("操作未完成")) {
			t.Fatal("repeated operation did not notify")
		}
		v.Move(2, 2)
		v.Advance(9 * time.Second)
		v.Frame()
		if v.HasText(p.tr("操作未完成")) || p.notice == "" {
			t.Fatal("toast should expire without clearing stored operation state")
		}
		v.Frame()
		if v.HasText(p.tr("操作未完成")) {
			t.Fatal("rendering restarted an expired notification")
		}
	}
}
