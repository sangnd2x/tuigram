package styles

import "github.com/charmbracelet/lipgloss"

var (
	ActiveBorder   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62"))
	InactiveBorder = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240"))

	SelectedChat = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	NormalChat   = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	UnreadBadge  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9"))

	OutgoingMsg = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	IncomingMsg = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	SenderName  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("213"))
	Timestamp   = lipgloss.NewStyle().Faint(true)
	MediaLabel  = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("178"))

	StatusBar = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)
