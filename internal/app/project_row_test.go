package app

import (
	"math"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"gus/internal/gitconfig"
	"gus/internal/identity"
	"gus/internal/state"
)

func TestProjectColumnsUseExtraWindowWidth(t *testing.T) {
	for _, lang := range []string{"zh-CN", "en"} {
		p := &application{preferences: state.Preferences{Language: lang}, projects: []gitconfig.Project{
			{Path: "bookmarks-manager", Current: identity.Resolved{Name: "Hunlongyu", Email: "hunlongyu@gmail.com"}},
			{Path: "music-dl-frontend", Current: identity.Resolved{Name: "ddmoyu", Email: "daydaymoyu@gmail.com"}},
		}}
		var repoText, nameText, emailText float32
		v := ui.NewTester(func(c *ui.Context) {
			p.view(c)
			repoText, _ = c.MeasureText(0, ui.Span{Text: "music-dl-frontend", Weight: 700})
			nameText, _ = c.MeasureText(0, ui.Span{Text: "Hunlongyu"})
			emailText, _ = c.MeasureText(0, ui.Span{Text: "daydaymoyu@gmail.com", Size: 13})
		}, minPanelWidth, panelHeight)
		find := func(label string) ui.Rect {
			t.Helper()
			b, ok := v.Find(label)
			if !ok {
				t.Fatalf("missing %q", label)
			}
			return b
		}
		narrow := find("music-dl-frontend")
		for _, width := range []int{700, 1000, minPanelWidth, 700} {
			v.SetSize(width, panelHeight)
			repo, name, email := find("music-dl-frontend"), find("Hunlongyu"), find("daydaymoyu@gmail.com")
			if width >= 700 && (repo.W <= narrow.W || repo.W < repoText || name.W < nameText || email.W < emailText) {
				t.Fatalf("wide window still truncates text: width=%d repo=%v/%v name=%v/%v email=%v/%v", width, repo.W, repoText, name.W, nameText, email.W, emailText)
			}
			for _, pair := range []struct {
				label string
				cell  ui.Rect
			}{{"仓库", repo}, {"名字", name}, {"邮箱", email}} {
				header := find(p.tr(pair.label))
				if math.Abs(float64(header.X-pair.cell.X)) > 1 || math.Abs(float64(header.W-pair.cell.W)) > 1 {
					t.Fatalf("header/body column mismatch: %s %+v %+v", pair.label, header, pair.cell)
				}
			}
			platforms := find(p.tr("托管平台：") + p.pushPlatforms(p.projects[0]))
			source := find(p.tr("身份来源") + " · " + p.tr("未设置"))
			if email.W <= 0 || email.X+email.W > platforms.X || platforms.X+platforms.W > source.X || source.X+source.W > float32(width) {
				t.Fatalf("columns overlap at width %d", width)
			}
		}
	}
}

func TestPlatformStackExpandsOnHover(t *testing.T) {
	for _, lang := range []string{"zh-CN", "en"} {
		r := gitconfig.Project{Path: "example", Remotes: []gitconfig.Remote{{PushURLs: []string{
			"git@github.com:team/project", "git@gitee.com:team/project", "git@gitlab.com:team/project",
			"git@codeberg.org:team/project", "https://git.company.example/team/project", "git@github.com:another/project",
		}}}}
		p := &application{preferences: state.Preferences{Language: lang}, projects: []gitconfig.Project{r}}
		v := ui.NewTester(p.view, minPanelWidth, minPanelHeight)
		if !v.HasText("+3") || slices.Contains(v.Texts(), "GitHub") {
			t.Fatal("stack should stay compact", v.Texts())
		}
		label := p.tr("托管平台：") + p.pushPlatforms(r)
		bounds, ok := v.Find(label)
		if !ok || bounds.W > 56 || bounds.X+bounds.W > minPanelWidth {
			t.Fatalf("stack outside column: %+v", bounds)
		}
		v.Move(bounds.X+bounds.W/2, bounds.Y+bounds.H/2)
		v.Advance(time.Second)
		v.Frame()
		for _, name := range []string{"GitHub", "Gitee", "GitLab", "Codeberg", "git.company.example"} {
			box, ok := v.Find(name)
			if !ok || box.X < 0 || box.Y < 0 || box.X+box.W > minPanelWidth || box.Y+box.H > minPanelHeight {
				t.Errorf("full platform missing or clipped: %s %+v", name, box)
			}
		}
		count := 0
		for _, text := range v.Texts() {
			if text == "GitHub" {
				count++
			}
		}
		if count != 1 {
			t.Fatal("platform not deduplicated", v.Texts())
		}
		v.Move(2, 2)
		v.Frame()
		if slices.Contains(v.Texts(), "GitHub") {
			t.Fatal("platform tooltip stayed open after leaving")
		}
	}
}

func TestSourceSymbolsExplainTheirScopes(t *testing.T) {
	for _, tc := range []struct{ nameScope, emailScope, source string }{
		{"global", "global", "继承"}, {"local", "local", "专用"},
		{"worktree", "worktree", "专用"}, {"global", "local", "混合"}, {"", "", "未设置"},
	} {
		p := &application{projects: []gitconfig.Project{{Path: "example", Current: identity.Resolved{Name: "Example", Email: "example@example.org", NameScope: tc.nameScope, EmailScope: tc.emailScope}}}}
		v := ui.NewTester(p.view, panelWidth, panelHeight)
		label := p.tr("身份来源") + " · " + p.tr(tc.source)
		bounds, ok := v.Find(label)
		if !ok || bounds.W != 24 || bounds.H != 24 || tc.source != "未设置" && slices.Contains(v.Texts(), p.tr(tc.source)) {
			t.Fatalf("source should be an icon: %s %+v", label, bounds)
		}
		v.Move(bounds.X+12, bounds.Y+12)
		v.Advance(time.Second)
		v.Frame()
		for _, want := range []string{label, p.tr("姓名：") + p.displayValue(tc.nameScope), p.tr("邮箱：") + p.displayValue(tc.emailScope)} {
			if !strings.Contains(strings.Join(v.Texts(), "\n"), want) {
				t.Fatalf("source tooltip missing %q: %q", want, v.Texts())
			}
		}
	}
}
