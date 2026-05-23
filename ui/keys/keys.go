package keys

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	SwitchPane key.Binding
	Up         key.Binding
	Down       key.Binding
	Enter      key.Binding
	FocusInput key.Binding
	Back       key.Binding
	Search     key.Binding
	ToggleHelp key.Binding
	Quit       key.Binding
}

var Default = KeyMap{
	SwitchPane: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch pane")),
	Up:         key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("up/k", "up")),
	Down:       key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("down/j", "down")),
	Enter:      key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select/send")),
	FocusInput: key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "type message")),
	Back:       key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Search:     key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
	ToggleHelp: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	Quit:       key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.SwitchPane, k.Enter, k.FocusInput, k.Back, k.Search, k.ToggleHelp, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter},
		{k.SwitchPane, k.FocusInput, k.Back},
		{k.Search, k.ToggleHelp, k.Quit},
	}
}
