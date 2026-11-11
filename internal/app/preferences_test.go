package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"gus/internal/state"
)

func TestLanguageAndThemeSettings(t *testing.T) {
	history, err := os.ReadFile("../../CHANGELOG.md")
	if err != nil {
		t.Fatal(err)
	}
	p := &application{releaseHistory: string(history), stateFile: filepath.Join(t.TempDir(), "state.json"), page: "settings"}
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
	data, err := state.Load(p.stateFile)
	if err != nil || data.Preferences != (state.Preferences{Language: "en", Appearance: "dark"}) {
		t.Fatal(data, err)
	}
	click("Dark")
	click("Light")
	if px := v.Image().RGBAAt(2, 2); px.R < 200 {
		t.Fatal("light appearance", px)
	}
	click("Release notes")
	if notes := p.notes(p.releaseHistory); strings.Contains(notes, "Windows 托盘") || !strings.Contains(notes, "Simplified Chinese") {
		t.Fatal(notes)
	}
	click("Back to settings")
	click("Back to projects")
	if !v.HasText("Projects") {
		t.Fatal(v.Texts())
	}
}
