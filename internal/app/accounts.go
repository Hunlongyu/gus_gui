package app

import (
	"net/mail"
	"slices"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"gus/internal/identity"
)

func (p *application) identityMenu(path string, current identity.Resolved, available bool) func(*ui.Menu) {
	return func(m *ui.Menu) {
		disabled := p.busy || !available
		local := identity.RepositoryScope(current.NameScope) || identity.RepositoryScope(current.EmailScope)
		if path != "" {
			if m.Item(p.tr("项目详情")).Chosen() {
				p.openProject(path)
			}
			m.Separator()
			if m.Item(p.tr("继承全局")).Radio(!local).Disabled(disabled || !local).Chosen() {
				p.applyIdentity(path, nil)
			}
			m.Separator()
		}
		for _, a := range p.accounts {
			if m.Item(a.Label + " · " + a.Name).Subtitle(a.Email).Radio(current.Matches(a) && (path == "" || local)).Disabled(disabled).Chosen() {
				p.applyIdentity(path, &a)
			}
		}
		if available && current.Name != "" && current.Email != "" && !slices.ContainsFunc(p.accounts, current.Matches) {
			m.Separator()
			if m.Item(p.tr("保存当前身份为账号…")).Chosen() {
				p.beginAccount(-1)
				p.draft.Name, p.draft.Email = current.Name, current.Email
			}
		} else if len(p.accounts) == 0 && m.Item(p.tr("新增账号…")).Chosen() {
			p.beginAccount(-1)
		}
		if path != "" {
			m.Separator()
			if m.Item(p.tr("移出目录列表")).Disabled(p.busy).Chosen() {
				p.updateUI(func() {
					if !p.busy {
						p.removeDirectory(path)
					}
				})
			}
		}
	}
}
func (p *application) beginAccount(index int) {
	p.page, p.editing, p.formError = "account", index, ""
	p.draft = identity.Account{}
	if index >= 0 {
		p.draft = p.accounts[index]
	}
}

func (p *application) uses(id int) int {
	count := 0
	for _, r := range p.projects {
		if r.Current.Matches(p.accounts[id]) && (identity.RepositoryScope(r.Current.NameScope) || identity.RepositoryScope(r.Current.EmailScope)) {
			count++
		}
	}
	return count
}

func (p *application) removeAccount(id int) {
	if p.global.Matches(p.accounts[id]) || p.uses(id) > 0 {
		return
	}
	list := slices.Delete(slices.Clone(p.accounts), id, id+1)
	if err := p.persist(list, p.projects); err != nil {
		p.setNotice(err.Error())
	}
}
func (p *application) accountList(c *ui.Context) {
	t := c.Theme()
	ui.Row(c).Padding(2, 0, 8).Gap(8).AlignItems(ui.Center).Children(func() {
		windowAction(c, ui.ButtonBase(c), backIcon, p.tr("返回项目"), p.tr("返回项目"), false).OnClick(func() { p.page = "" })
		ui.Text(c, p.tr("账号管理")).Grow(1).Bold().FontSize(18)
		ui.PrimaryButton(c, p.tr("＋ 新增账号")).Padding(5, 10).OnClick(func() { p.beginAccount(-1) })
	})
	remove := -1
	ui.Scroll(c.Key("accounts")).Grow(1).Children(func() {
		for i, a := range p.accounts {
			count := p.uses(i)
			usage := p.tr("未使用")
			if count > 0 {
				usage = strconv.Itoa(count) + p.tr(" 个项目")
			}
			if p.global.Matches(a) {
				usage = p.tr("全局默认")
			}
			menu := func(m *ui.Menu) {
				if m.Item(p.tr("编辑账号")).Chosen() {
					p.beginAccount(i)
				}
				if m.Item(p.tr("删除账号")).Disabled(p.global.Matches(a) || count > 0).Chosen() {
					remove = i
				}
			}
			row := ui.Row(c.Key(a.Label)).Height(73).Shrink(0).Gap(10).Padding(0, 5).AlignItems(ui.Center).
				Border(1, t.Border).BorderWidth(0, 0, 1, 0)
			row.Focusable().ContextMenu(menu).Children(func() {
				ink := ui.RGB(92, 80, 153)
				if p.global.Matches(a) {
					ink = ui.RGB(43, 105, 80)
				}
				if t.Dark {
					ink = ink.Mix(t.Text, 0.60)
				}
				ui.Box(c).Size(35, 35).Radius(9).Background(t.Surface.Mix(ink, 0.12)).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
					ui.Icon(c, accountIcon).Size(19, 19).TextColor(ink)
				})
				ui.Column(c).Grow(1).MinWidth(0).Gap(4).Children(func() {
					ui.Row(c).Gap(7).AlignItems(ui.Center).Children(func() {
						ui.Text(c, a.Label).Bold().FontSize(15).MaxWidth(118).MaxLines(1).Ellipsis("…").Tooltip(a.Label)
						ui.Text(c, p.displayValue(a.Name)).Grow(1).MinWidth(0).MaxLines(1).Ellipsis("…").FontSize(13).TextColor(t.TextMuted)
					})
					ui.Text(c, p.displayValue(a.Email)).FontSize(13).TextColor(t.TextMuted).MaxLines(1).Ellipsis("…")
				})
				tag := ui.Text(c, usage).Shrink(0).FontSize(12).Radius(5).Padding(2, 7).TextColor(t.TextMuted).Background(t.Background)
				if p.global.Matches(a) {
					tag.TextColor(ink).Background(t.Surface.Mix(ink, 0.08))
				} else {
					ui.Button(c, p.tr("设为全局")).Shrink(0).FontSize(12).Padding(4, 8).
						Label(a.Label + " · " + p.tr("设为全局")).
						Tooltip(p.tr("设为 Git 全局提交身份，未单独配置身份的项目将使用此账号。")).
						Disabled(p.busy).OnClick(func() { p.applyIdentity("", &a) })
				}
				windowAction(c, ui.ButtonBase(c), moreIcon, a.Label+p.tr("账号操作"), p.tr("账号操作"), false).Menu(menu)
			})
		}
	})
	if remove >= 0 {
		p.removeAccount(remove)
	}
}

func (p *application) accountForm(c *ui.Context) {
	title := p.tr("新增账号")
	if p.editing >= 0 {
		title = p.tr("编辑账号")
	}
	t := c.Theme()
	ui.Row(c).Padding(2, 0, 8).Gap(8).AlignItems(ui.Center).Children(func() {
		windowAction(c, ui.ButtonBase(c), backIcon, p.tr("返回账号管理"), p.tr("返回账号管理"), false).OnClick(func() { p.page = "accounts" })
		ui.Text(c, title).Bold().FontSize(18)
	})
	ui.Scroll(c.Key("account-fields")).Grow(1).Padding(4, 5).Gap(16).Children(func() {
		accountField(c, p.tr("账号标签"), p.tr("例如：工作"), &p.draft.Label)
		accountField(c, p.tr("名字"), "Lin Chen", &p.draft.Name)
		accountField(c, p.tr("邮箱"), "lin@example.org", &p.draft.Email)
		if p.formError != "" {
			ui.Text(c, p.formError).TextColor(t.Danger).FontSize(13)
		}
	})
	ui.Row(c).Justify(ui.End).Gap(8).Padding(10, 0, 0).Border(1, t.Border).BorderWidth(1, 0, 0, 0).Children(func() {
		ui.Button(c, p.tr("取消")).Padding(5, 14).OnClick(func() { p.page = "accounts" })
		ui.PrimaryButton(c, p.tr("保存账号")).Padding(5, 14).OnClick(p.saveAccount)
	})
}

func accountField(c *ui.Context, label, example string, value *string) {
	ui.Field(c.Key(label), label, func() {
		ui.TextInput(c, value).Label(label).Placeholder(example).Height(36).Radius(6).Padding(6, 10)
	}).Gap(5)
}

func (p *application) saveAccount() {
	a := identity.Account{Label: strings.TrimSpace(p.draft.Label), Name: strings.TrimSpace(p.draft.Name), Email: strings.TrimSpace(p.draft.Email)}
	if a.Label == "" || a.Name == "" {
		p.formError = p.tr("请填写账号标签和名字。")
		return
	}
	address, err := mail.ParseAddress(a.Email)
	if err != nil || address.Address != a.Email {
		p.formError = p.tr("请输入有效的邮箱地址。")
		return
	}
	for i, existing := range p.accounts {
		if i != p.editing && existing.Label == a.Label {
			p.formError = p.tr("该账号标签已存在。")
			return
		}
	}
	list := slices.Clone(p.accounts)
	if p.editing < 0 {
		list = append(list, a)
	} else {
		list[p.editing] = a
	}
	if err := p.persist(list, p.projects); err != nil {
		p.formError = err.Error()
		return
	}
	p.page = "accounts"
}
