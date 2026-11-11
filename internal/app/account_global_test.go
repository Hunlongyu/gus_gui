package app

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"gus/internal/gitconfig"
	"gus/internal/identity"
	"gus/internal/state"
)

func TestSetGlobalIdentityFromAccounts(t *testing.T) {
	dir := isolatedGit(t)
	inherited, local := filepath.Join(dir, "inherited"), filepath.Join(dir, "local")
	for _, root := range []string{inherited, local} {
		runGit(t, dir, "init", "-q", root)
	}
	personal := identity.Account{Label: "个人", Name: "Personal User", Email: "personal@example.org"}
	work := identity.Account{Label: "工作", Name: "Work User", Email: "work@example.org"}
	apply(t, local, &work)
	queue := make(chan func(), 4)
	p := &application{
		page: "accounts", accounts: []identity.Account{personal, work},
		projects: []gitconfig.Project{gitconfig.InspectProject(inherited), gitconfig.InspectProject(local)},
		updateUI: func(f func()) { queue <- f },
	}
	v := ui.NewTester(p.view, minPanelWidth, panelHeight)
	settle := func() {
		t.Helper()
		deadline := time.After(20 * time.Second)
		for p.busy {
			select {
			case f := <-queue:
				f()
				v.Frame()
			case <-deadline:
				t.Fatal("global identity update stalled")
			}
		}
	}
	click := func(a identity.Account) {
		t.Helper()
		if err := v.Click(a.Label + " · 设为全局"); err != nil {
			t.Fatal(err)
		}
		v.Frame()
		if !p.busy || !p.writing {
			t.Fatal("account action did not start the Git write")
		}
		settle()
	}
	for _, account := range []identity.Account{personal, work} {
		click(account)
		actual, err := gitconfig.ReadIdentity("")
		if err != nil || !actual.Matches(account) || !p.global.Matches(account) || p.notice != "" {
			t.Fatalf("global identity was not applied: %+v %v %s", actual, err, p.notice)
		}
		if !p.projects[0].Current.Matches(account) || p.projects[0].Current.Source() != "继承" ||
			!p.projects[1].Current.Matches(work) || p.projects[1].Current.Source() != "专用" {
			t.Fatal("inherited/local identity scopes were not preserved", p.projects)
		}
		if !v.HasText("全局默认") || v.HasText(account.Label+" · 设为全局") || p.page != "accounts" {
			t.Fatal("account status did not reflect the selected global identity", v.Texts())
		}
	}
	// A failed write keeps the current global account and reports a toast.
	lock := os.Getenv("GIT_CONFIG_GLOBAL") + ".lock"
	if err := os.WriteFile(lock, []byte("locked"), 0600); err != nil {
		t.Fatal(err)
	}
	click(personal)
	if p.notice == "" || !p.global.Matches(work) || !v.HasText("操作未完成") {
		t.Fatal("failed switch was not reported or changed the selected account")
	}
}

func TestGlobalAccountActionFitsMinimumWindow(t *testing.T) {
	for _, language := range []string{"zh-CN", "en"} {
		account := identity.Account{Label: "工作", Name: "Hunlongyu", Email: "hunlongyu@example.org"}
		p := &application{page: "accounts", preferences: state.Preferences{Language: language}, accounts: []identity.Account{account}, busy: true}
		v := ui.NewTester(p.view, minPanelWidth, minPanelHeight)
		label := account.Label + " · " + p.tr("设为全局")
		button, ok := v.Find(label)
		menu, menuOK := v.Find(account.Label + p.tr("账号操作"))
		email, emailOK := v.Find(account.Email)
		if !ok || !menuOK || !emailOK || email.W <= 0 || email.X+email.W > button.X ||
			button.W <= 0 || button.X+button.W > menu.X || menu.X+menu.W > minPanelWidth {
			t.Fatalf("account controls overlap in %s: email=%+v button=%+v menu=%+v", language, email, button, menu)
		}
		// Refreshes/writes disable the action, preventing a second Git mutation.
		_ = v.Click(label)
		v.Frame()
		if p.writing {
			t.Fatal("busy account action started another write")
		}
	}
}
