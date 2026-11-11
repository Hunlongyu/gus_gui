package main

import (
	"math"
	"path/filepath"

	"github.com/egoist/mygo/ui"
)

type projectColumns struct {
	repo, name float32
}

func (p *application) projectColumns(c *ui.Context) projectColumns {
	// Reserve the panel padding/border, row padding, four gaps and icon columns.
	width, _ := c.Size()
	available := max(0, width-34-20-4*8-56-42)
	minimum := [3]float32{95, 64, 90}
	desired := minimum
	measure := func(column int, text string, size float32, weight int) {
		w, _ := c.MeasureText(0, ui.Span{Text: text, Size: size, Weight: weight})
		desired[column] = max(desired[column], float32(math.Ceil(float64(w)))+4)
	}
	for i, label := range []string{"仓库", "名字", "邮箱"} {
		measure(i, p.tr(label), 12, 400)
	}
	// Use all projects so filtering does not make the columns jump.
	for _, r := range p.projects {
		measure(0, filepath.Base(r.path), c.Theme().FontSize, 700)
		measure(1, p.displayValue(r.current.Name), c.Theme().FontSize, 400)
		measure(2, p.displayValue(r.current.Email), 13, 400)
	}
	base := minimum[0] + minimum[1] + minimum[2]
	wanted := desired[0] + desired[1] + desired[2]
	if available < wanted {
		// Share limited space in proportion to each column's unmet text width.
		ratio := max(0, (available-base)/max(1, wanted-base))
		for i := range desired {
			desired[i] = minimum[i] + (desired[i]-minimum[i])*ratio
		}
	} else {
		// Once text fits, give spare room to repositories and email, not short names.
		desired[0] += (available - wanted) * 0.55
	}
	return projectColumns{repo: desired[0], name: desired[1]}
}
