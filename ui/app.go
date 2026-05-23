package ui

import (
	"context"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sangnguyen/tui-telegram/telegram"
	"github.com/sangnguyen/tui-telegram/ui/chatlist"
	"github.com/sangnguyen/tui-telegram/ui/folders"
	"github.com/sangnguyen/tui-telegram/ui/input"
	appkeys "github.com/sangnguyen/tui-telegram/ui/keys"
	"github.com/sangnguyen/tui-telegram/ui/messages"
	"github.com/sangnguyen/tui-telegram/ui/styles"
	tg "github.com/gotd/td/tg"
)

type focusedPane int

const (
	paneChatList focusedPane = iota
	paneMessages
)

type DialogsLoadedMsg struct {
	Dialogs []telegram.DialogItem
}

type FoldersLoadedMsg struct {
	Folders []telegram.FolderItem
}

type HistoryLoadedMsg struct {
	ChatID int64
	Items  []telegram.MsgItem
}

type MessageSentMsg struct {
	Item telegram.MsgItem
}

type ErrMsg struct {
	Err error
}

func (e ErrMsg) Error() string { return e.Err.Error() }

type AppModel struct {
	chatList     chatlist.Model
	msgPane      messages.Model
	help         help.Model
	folderBar    folders.Model
	focus        focusedPane
	tgClient     *telegram.ClientWrapper
	ctx          context.Context
	ready        bool
	width        int
	height       int
	statusMsg    string
	activePeer   tg.InputPeerClass
	allDialogs   []telegram.DialogItem
	activeFolder telegram.FolderItem
}

func New(tgClient *telegram.ClientWrapper, ctx context.Context) AppModel {
	return AppModel{
		tgClient:  tgClient,
		ctx:       ctx,
		focus:     paneChatList,
		help:      help.New(),
		folderBar: folders.New(0),
	}
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		leftW := msg.Width * 30 / 100
		rightW := msg.Width - leftW
		bodyH := msg.Height - 1 // minus status bar
		m.chatList = chatlist.New(leftW, bodyH-1) // minus folder bar row
		m.msgPane = messages.New(rightW, bodyH)
		m.chatList.SetFocused(true)
		m.folderBar = folders.New(leftW)
		m.help.Width = msg.Width
		m.ready = true
		return m, tea.Batch(
			loadDialogsCmd(m.tgClient, m.ctx),
			loadFoldersCmd(m.tgClient, m.ctx),
		)

	case tea.KeyMsg:
		if key.Matches(msg, appkeys.Default.ToggleHelp) {
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		}
		if key.Matches(msg, appkeys.Default.Quit) {
			return m, tea.Quit
		}
		if key.Matches(msg, appkeys.Default.SwitchPane) {
			if m.focus == paneChatList {
				m.focus = paneMessages
				m.chatList.SetFocused(false)
				m.msgPane.SetFocused(true)
			} else {
				m.focus = paneChatList
				m.chatList.SetFocused(true)
				m.msgPane.SetFocused(false)
			}
			return m, nil
		}
		var folderCmd tea.Cmd
		m.folderBar, folderCmd = m.folderBar.Update(msg)
		cmds = append(cmds, folderCmd)

		if m.focus == paneChatList {
			var cmd tea.Cmd
			m.chatList, cmd = m.chatList.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			var cmd tea.Cmd
			m.msgPane, cmd = m.msgPane.Update(msg)
			cmds = append(cmds, cmd)
		}

	case chatlist.SelectChatMsg:
		m.activePeer = msg.Item.InputPeer
		m.msgPane.SetChat(msg.Item.ID, msg.Item.Name)
		m.focus = paneMessages
		m.chatList.SetFocused(false)
		m.msgPane.SetFocused(true)
		cmds = append(cmds, loadHistoryCmd(m.tgClient, m.ctx, msg.Item.InputPeer, msg.Item.ID, msg.Item.Name))

	case DialogsLoadedMsg:
		m.allDialogs = msg.Dialogs
		cmds = append(cmds, m.chatList.SetItems(m.filteredDialogs()))

	case folders.SelectFolderMsg:
		m.activeFolder = msg.Folder
		cmds = append(cmds, m.chatList.SetItems(m.filteredDialogs()))

	case FoldersLoadedMsg:
		m.folderBar.SetFolders(msg.Folders)

	case HistoryLoadedMsg:
		if msg.ChatID == m.msgPane.ChatID {
			m.msgPane.SetMessages(msg.Items)
		}

	case telegram.NewMessageMsg:
		if msg.ChatID == m.msgPane.ChatID {
			m.msgPane.AppendMessage(msg.Item)
		}

	case input.SubmitMsg:
		if m.activePeer != nil {
			cmds = append(cmds, sendMessageCmd(m.tgClient, m.ctx, m.activePeer, msg.Text))
		}

	case MessageSentMsg:
		m.msgPane.AppendMessage(msg.Item)

	case ErrMsg:
		m.statusMsg = "Error: " + msg.Err.Error()

	default:
		// Propagate to both panes for things like cursor blink
		var cmd tea.Cmd
		m.chatList, cmd = m.chatList.Update(msg)
		cmds = append(cmds, cmd)
		m.msgPane, cmd = m.msgPane.Update(msg)
		cmds = append(cmds, cmd)
		var folderCmd tea.Cmd
		m.folderBar, folderCmd = m.folderBar.Update(msg)
		cmds = append(cmds, folderCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m AppModel) View() string {
	if !m.ready {
		return "Loading…"
	}
	left := lipgloss.JoinVertical(lipgloss.Left, m.folderBar.View(), m.chatList.View())
	right := m.msgPane.View()
	body := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	helpView := m.help.View(appkeys.Default)
	status := styles.StatusBar.Render(m.statusMsg + helpView)
	return lipgloss.JoinVertical(lipgloss.Left, body, status)
}

func (m AppModel) filteredDialogs() []telegram.DialogItem {
	if len(m.activeFolder.PeerIDs) == 0 {
		return m.allDialogs
	}
	var filtered []telegram.DialogItem
	for _, d := range m.allDialogs {
		if m.activeFolder.PeerIDs[d.ID] {
			filtered = append(filtered, d)
		}
	}
	return filtered
}

func loadDialogsCmd(w *telegram.ClientWrapper, ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		dialogs, err := w.LoadDialogs(ctx)
		if err != nil {
			return ErrMsg{Err: err}
		}
		return DialogsLoadedMsg{Dialogs: dialogs}
	}
}

func loadFoldersCmd(w *telegram.ClientWrapper, ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		fs, err := w.LoadFolders(ctx)
		if err != nil {
			return ErrMsg{Err: err}
		}
		return FoldersLoadedMsg{Folders: fs}
	}
}

func loadHistoryCmd(w *telegram.ClientWrapper, ctx context.Context, peer tg.InputPeerClass, chatID int64, chatName string) tea.Cmd {
	return func() tea.Msg {
		items, err := w.LoadHistory(ctx, peer, 50, chatName)
		if err != nil {
			return ErrMsg{Err: err}
		}
		return HistoryLoadedMsg{ChatID: chatID, Items: items}
	}
}

func sendMessageCmd(w *telegram.ClientWrapper, ctx context.Context, peer tg.InputPeerClass, text string) tea.Cmd {
	return func() tea.Msg {
		item, err := w.SendMessage(ctx, peer, text)
		if err != nil {
			return ErrMsg{Err: err}
		}
		return MessageSentMsg{Item: item}
	}
}
