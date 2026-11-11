package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPanelSizeSurvivesRestartAndOtherSettingsChanges(t *testing.T) {
	t.Setenv("APPDATA", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("HOME", t.TempDir())
	p := &application{}
	p.restore()
	p.preferences = preferences{Language: "en", Appearance: "dark"}
	account := identity{Label: "test", Name: "Tester", Email: "tester@example.org"}
	projects := []project{{path: filepath.Join(t.TempDir(), "project")}}
	if err := p.persist([]identity{account}, projects); err != nil {
		t.Fatal(err)
	}
	if err := p.rememberPanelSize(780, 560); err != nil {
		t.Fatal(err)
	}
	reopened := &application{}
	reopened.restore()
	if reopened.storageError != "" || reopened.windowSize != (panelSize{780, 560}) ||
		len(reopened.accounts) != 1 || reopened.accounts[0] != account || len(reopened.projects) != 1 ||
		reopened.projects[0].path != projects[0].path || reopened.preferences != p.preferences {
		t.Fatalf("size or existing data lost: %+v", reopened)
	}
	// Account/project writes must preserve the size, and transient invalid
	// dimensions (such as minimizing) must not replace the last normal size.
	if err := reopened.persist(nil, nil); err != nil {
		t.Fatal(err)
	}
	if err := reopened.rememberPanelSize(0, 0); err != nil {
		t.Fatal(err)
	}
	third := &application{}
	third.restore()
	if third.windowSize != reopened.windowSize || len(third.accounts) != 0 || len(third.projects) != 0 {
		t.Fatal("subsequent writes lost the saved window size")
	}
	// A failed save must leave the old size eligible for a retry.
	before := reopened.windowSize
	reopened.stateFile = filepath.Join(t.TempDir(), "blocked", "state.json")
	if err := os.WriteFile(filepath.Dir(reopened.stateFile), []byte("file"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := reopened.rememberPanelSize(900, 600); err == nil || reopened.windowSize != before {
		t.Fatal("failed save was treated as persisted")
	}
}
