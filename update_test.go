package main

import "testing"

func TestStandaloneDownloadsDoNotRelaunchAnOutdatedFilename(t *testing.T) {
	for _, tc := range []struct {
		name, goos string
		want       bool
	}{
		{"gus.exe", "windows", true},
		{"GUS.EXE", "windows", true},
		{"gus-1.4.0-windows-amd64.exe", "windows", false},
		{"gus-1.4.0-windows-arm64.exe", "windows", false},
		{"gus", "linux", true},
		{"gus-1.4.0-linux-amd64", "linux", false},
		{"gus", "darwin", true},
		{"renamed-app", "darwin", false},
	} {
		if got := updateExecutableName(tc.name, tc.goos); got != tc.want {
			t.Errorf("%s/%s: update allowed = %v, want %v", tc.goos, tc.name, got, tc.want)
		}
	}
}
