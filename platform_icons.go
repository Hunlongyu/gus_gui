package main

import (
	"embed"

	"github.com/egoist/mygo/ui"
)

// Local assets keep platform badges available offline and avoid sending
// repository hosts to a third-party favicon service.
//
//go:embed assets/platforms/*.svg assets/platforms/*.png
var platformIconFiles embed.FS

func loadPlatformSVG(name string) *ui.SVG {
	data, err := platformIconFiles.ReadFile("assets/platforms/" + name + ".svg")
	if err != nil {
		panic(err)
	}
	return ui.MustParseSVG(data)
}

func loadPlatformImage(name string) ui.ImageSource {
	data, err := platformIconFiles.ReadFile("assets/platforms/" + name)
	if err != nil {
		panic(err)
	}
	image, err := ui.DecodeImage(data)
	if err != nil {
		panic(err)
	}
	return image
}

var (
	codebergIcon    = loadPlatformSVG("codeberg")
	giteeIcon       = loadPlatformSVG("gitee")
	bitbucketIcon   = loadPlatformSVG("bitbucket")
	gitcodeIcon     = loadPlatformSVG("gitcode")
	azureDevopsIcon = loadPlatformSVG("azuredevops")
	gitgudImage     = loadPlatformImage("gitgud.png")
)
