package app

import (
	"github.com/egoist/mygo/ui"

	"gus/assets"
)

func loadPlatformSVG(name string) *ui.SVG {
	data, err := assets.Platforms.ReadFile("platforms/" + name + ".svg")
	if err != nil {
		panic(err)
	}
	return ui.MustParseSVG(data)
}

func loadPlatformImage(name string) ui.ImageSource {
	data, err := assets.Platforms.ReadFile("platforms/" + name)
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
