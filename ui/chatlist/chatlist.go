package chatlist

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	tg "github.com/sangnguyen/tui-telegram/telegram"
	appkeys "github.com/sangnguyen/tui-telegram/ui/keys"
	"github.com/sangnguyen/tui-telegram/ui/styles"
)

type SelectChatMsg struct {
	Item tg.DialogItem
}

type Item struct {
	Dialog tg.DialogItem
}

func (i Item) Title() string {
	title := i.Dialog.Name
	if i.Dialog.UnreadCount > 0 {
		title = styles.UnreadBadge.Render(fmt.Sprintf("(%d) ", i.Dialog.UnreadCount)) + title
	}
	return title
}

func (i Item) Description() string {
	if i.Dialog.LastMessage == "" {
		return ""
	}
	msg := i.Dialog.LastMessage
	if len(msg) > 50 {
		msg = msg[:47] + "..."
	}
	return msg
}

func (i Item) FilterValue() string { return i.Dialog.Name }

type Model struct {
	list    list.Model
	focused bool
	width   int
	height  int
}

func New(width, height int) Model {
	delegate := list.NewDefaultDelegate()
	l := list.New([]list.Item{}, delegate, width-2, height-2)
	l.Title = "Chats"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	return Model{list: l, width: width, height: height}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if kmsg, ok := msg.(tea.KeyMsg); ok && m.focused {
		if key.Matches(kmsg, appkeys.Default.Enter) {
			if sel, ok := m.list.SelectedItem().(Item); ok {
				return m, func() tea.Msg { return SelectChatMsg{Item: sel.Dialog} }
			}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	border := styles.InactiveBorder
	if m.focused {
		border = styles.ActiveBorder
	}
	return border.Width(m.width - 2).Height(m.height - 2).Render(m.list.View())
}

func (m *Model) SetItems(dialogs []tg.DialogItem) tea.Cmd {
	items := make([]list.Item, len(dialogs))
	for i, d := range dialogs {
		items[i] = Item{Dialog: d}
	}
	return m.list.SetItems(items)
}

func (m *Model) SetFocused(v bool) {
	m.focused = v
}

func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.list.SetSize(w-4, h-4)
}

func (m Model) SelectedDialog() (tg.DialogItem, bool) {
	if sel, ok := m.list.SelectedItem().(Item); ok {
		return sel.Dialog, true
	}
	return tg.DialogItem{}, false
}

