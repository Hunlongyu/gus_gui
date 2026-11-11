package app

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"gus/internal/gitconfig"
)

func TestLiveProjectSearch(t *testing.T) {
	p := &application{projects: []gitconfig.Project{{Path: "C:/work/alpha服务"}, {Path: "C:/alpha/bravo"}, {Path: "C:/work/ALPHAbet"}}}
	v := ui.NewTester(p.view, panelWidth, panelHeight)
	if err := v.Click("搜索项目名"); err != nil {
		t.Fatal(err)
	}
	v.Type("ALPHA")
	if !v.HasText("2 / 3") || v.HasText("bravo") || !v.HasText("ALPHAbet") {
		t.Fatalf("partial name filter: %q", v.Texts())
	}
	v.Type("服务")
	if !v.HasText("1 / 3") || !v.HasText("alpha服务") {
		t.Fatalf("Chinese query: %q", v.Texts())
	}
	v.Type("X")
	if !v.HasText("没有匹配的项目") {
		t.Fatal(v.Texts())
	}
	v.Key(0, ui.KeyEscape)
	if p.projectQuery != "" || !v.HasText("bravo") {
		t.Fatal("clear", v.Texts())
	}
}
