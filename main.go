package main

import (
	_ "embed"

	"gus/internal/app"
)

// Use the same image for packaging, the window, and the tray.
//
//go:embed resources/icon.png
var appIcon []byte

//go:embed CHANGELOG.md
var releaseHistory string

func main() {
	app.Run(appIcon, releaseHistory)
}
