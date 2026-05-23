package messages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tg "github.com/sangnguyen/tui-telegram/telegram"
	"github.com/sangnguyen/tui-telegram/ui/input"
	appkeys "github.com/sangnguyen/tui-telegram/ui/keys"
	"github.com/sangnguyen/tui-telegram/ui/styles"
)

type Model struct {
	viewport   viewport.Model
	input      input.Model
	messages   []tg.MsgItem
	chatTitle  string
	ChatID     int64
	inputFocus bool // true = typing; false = scrolling viewport
	focused    bool // pane is focused at all
	width      int
	height     int
}

func New(width, height int) Model {
	vp := viewport.New(width-4, height-7)
	inp := input.New(width - 4)
	return Model{
		viewport: vp,
		input:    inp,
		width:    width,
		height:   height,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch kmsg := msg.(type) {
	case tea.KeyMsg:
		if !m.focused {
			break
		}
		if m.inputFocus {
			if key.Matches(kmsg, appkeys.Default.Back) {
				m.inputFocus = false
				m.input.SetFocused(false)
				break
			}
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			if key.Matches(kmsg, appkeys.Default.FocusInput) {
				m.inputFocus = true
				m.input.SetFocused(true)
				break
			}
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	default:
		if m.inputFocus {
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	border := styles.InactiveBorder
	if m.focused {
		border = styles.ActiveBorder
	}

	header := styles.SenderName.Render(m.chatTitle)
	if m.chatTitle == "" {
		header = styles.Timestamp.Render("Select a chat")
	}

	hint := ""
	if m.focused && !m.inputFocus {
		hint = styles.Timestamp.Render(" [i] to type")
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		header+hint,
		m.viewport.View(),
		m.input.View(),
	)
	return border.Width(m.width-2).Height(m.height-2).Render(content)
}

func (m *Model) SetMessages(msgs []tg.MsgItem) {
	m.messages = msgs
	m.refreshViewport()
}

func (m *Model) AppendMessage(msg tg.MsgItem) {
	m.messages = append(m.messages, msg)
	m.refreshViewport()
}

func (m *Model) SetChat(id int64, title string) {
	m.ChatID = id
	m.chatTitle = title
	m.messages = nil
	m.refreshViewport()
}

func (m *Model) SetFocused(v bool) {
	m.focused = v
	if !v {
		m.inputFocus = false
		m.input.SetFocused(false)
	}
}

func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.viewport.Width = w - 4
	m.viewport.Height = h - 7
	m.input.SetWidth(w - 4)
}

func (m *Model) refreshViewport() {
	var lines []string
	for _, msg := range m.messages {
		ts := styles.Timestamp.Render(msg.Date.Format("15:04"))
		sender := styles.SenderName.Render(msg.SenderName)
		body := msg.Text
		if msg.Media != "" {
			if body != "" {
				body += " " + styles.MediaLabel.Render(msg.Media)
			} else {
				body = styles.MediaLabel.Render(msg.Media)
			}
		}
		line := fmt.Sprintf("%s %s: %s", ts, sender, body)
		if msg.Outgoing {
			styled := styles.OutgoingMsg.Render(line)
			line = lipgloss.PlaceHorizontal(m.viewport.Width, lipgloss.Right, styled)
		} else {
			line = styles.IncomingMsg.Render(line)
		}
		lines = append(lines, line)
	}
	m.viewport.SetContent(strings.Join(lines, "\n"))
	m.viewport.GotoBottom()
}
