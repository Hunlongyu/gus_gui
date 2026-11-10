package main

import (
	"github.com/egoist/mygo/ui"
	"path/filepath"
	"strings"
	"testing"
)

func TestLanguageAndThemeSettings(t *testing.T) {
	p := &application{stateFile: filepath.Join(t.TempDir(), "state.json"), page: "settings"}
	v := ui.NewTester(p.view, panelWidth, panelHeight)
	click := func(s string) {
		t.Helper()
		if err := v.Click(s); err != nil {
			t.Fatal(err)
		}
		v.Frame()
	}
	click("简体中文")
	click("English")
	if !v.HasText("Settings") {
		t.Fatal(v.Texts())
	}
	click("System")
	click("Dark")
	if px := v.Image().RGBAAt(2, 2); px.R > 70 {
		t.Fatal("dark appearance", px)
	}
	data, err := loadData(p.stateFile)
	if err != nil || data.Preferences != (preferences{"en", "dark"}) {
		t.Fatal(data, err)
	}
	click("Dark")
	click("Light")
	if px := v.Image().RGBAAt(2, 2); px.R < 200 {
		t.Fatal("light appearance", px)
	}
	click("Release notes")
	if notes := p.notes(releaseHistory); strings.Contains(notes, "Windows 托盘") || !strings.Contains(notes, "Simplified Chinese") {
		t.Fatal(notes)
	}
	click("Back to settings")
	click("Back to projects")
	if !v.HasText("Projects") {
		t.Fatal(v.Texts())
	}
}
