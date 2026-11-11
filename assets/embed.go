// Package assets embeds the platform badges used by the desktop UI.
package assets

import "embed"

// Platforms keeps badges available offline without querying remote hosts.
//
//go:embed platforms/*.svg platforms/*.png
var Platforms embed.FS
