package app

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"gus/internal/gitconfig"
	"gus/internal/state"
)

func TestPanelAtMinimumSize(t *testing.T) {
	for _, lang := range []string{"zh-CN", "en"} {
		p := &application{preferences: state.Preferences{Language: lang}, projects: []gitconfig.Project{{Path: "example"}}}
		v := ui.NewTester(p.view, panelWidth, panelHeight)
		v.SetSize(minPanelWidth, minPanelHeight)
		inside := func(label string) {
			t.Helper()
			r, ok := v.Find(label)
			if !ok || r.W <= 0 || r.H <= 0 || r.X < 0 || r.Y < 0 || r.X+r.W > minPanelWidth || r.Y+r.H > minPanelHeight {
				t.Errorf("%s: %q outside minimum panel: %+v", lang, label, r)
			}
		}
		for _, label := range []string{"搜索项目名", "账号管理", "设置", "添加目录", "身份来源", "关闭"} {
			inside(p.tr(label))
		}
		p.beginAccount(-1)
		p.formError = p.tr("请输入有效的邮箱地址。")
		v.Frame()
		inside(p.tr("保存账号"))
		inside(p.tr("取消"))
	}
}
