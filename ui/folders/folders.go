package folders

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tg "github.com/sangnguyen/tui-telegram/telegram"
)

type SelectFolderMsg struct {
	Folder tg.FolderItem
}

var (
	activeTab   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Underline(true)
	inactiveTab = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	tabBar      = lipgloss.NewStyle().Padding(0, 1)
)

type Model struct {
	folders []tg.FolderItem
	active  int
	width   int
}

func New(width int) Model {
	return Model{
		folders: []tg.FolderItem{{ID: -1, Title: "All Chats"}},
		width:   width,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if kmsg, ok := msg.(tea.KeyMsg); ok {
		switch kmsg.String() {
		case "]":
			if m.active < len(m.folders)-1 {
				m.active++
				return m, func() tea.Msg { return SelectFolderMsg{Folder: m.folders[m.active]} }
			}
		case "[":
			if m.active > 0 {
				m.active--
				return m, func() tea.Msg { return SelectFolderMsg{Folder: m.folders[m.active]} }
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	var tabs []string
	for i, f := range m.folders {
		label := f.Title
		if i == m.active {
			tabs = append(tabs, activeTab.Render(label))
		} else {
			tabs = append(tabs, inactiveTab.Render(label))
		}
	}
	row := strings.Join(tabs, inactiveTab.Render(" │ "))
	return tabBar.Width(m.width - 2).Render(row)
}

func (m *Model) SetFolders(fs []tg.FolderItem) {
	m.folders = fs
	if m.active >= len(m.folders) {
		m.active = 0
	}
}

func (m Model) Active() tg.FolderItem {
	if len(m.folders) == 0 {
		return tg.FolderItem{ID: -1, Title: "All Chats"}
	}
	return m.folders[m.active]
}
