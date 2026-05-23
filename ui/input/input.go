package input

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	appkeys "github.com/sangnguyen/tui-telegram/ui/keys"
)

type SubmitMsg struct {
	Text string
}

type Model struct {
	ti      textinput.Model
	focused bool
	width   int
}

func New(width int) Model {
	ti := textinput.New()
	ti.Placeholder = "Type a message…"
	ti.Width = width - 4
	return Model{ti: ti, width: width}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}
	if kmsg, ok := msg.(tea.KeyMsg); ok && key.Matches(kmsg, appkeys.Default.Enter) {
		text := m.ti.Value()
		if text == "" {
			return m, nil
		}
		m.ti.SetValue("")
		return m, func() tea.Msg { return SubmitMsg{Text: text} }
	}
	var cmd tea.Cmd
	m.ti, cmd = m.ti.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.ti.View()
}

func (m Model) Value() string {
	return m.ti.Value()
}

func (m *Model) SetFocused(v bool) {
	m.focused = v
	if v {
		m.ti.Focus()
	} else {
		m.ti.Blur()
	}
}

func (m *Model) SetWidth(w int) {
	m.width = w
	m.ti.Width = w - 4
}

func (m *Model) Reset() {
	m.ti.SetValue("")
}
