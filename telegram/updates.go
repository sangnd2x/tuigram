package telegram

import (
	"context"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gotd/td/tg"
)

type Sender interface {
	Send(tea.Msg)
}

type UpdateBridge struct {
	mu        sync.Mutex
	cacheMu   sync.RWMutex
	program   Sender
	userCache map[int64]string
}

func NewUpdateBridge() *UpdateBridge {
	return &UpdateBridge{userCache: make(map[int64]string)}
}

func (b *UpdateBridge) SetProgram(p Sender) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.program = p
}

func (b *UpdateBridge) send(msg tea.Msg) {
	b.mu.Lock()
	p := b.program
	b.mu.Unlock()
	if p != nil {
		p.Send(msg)
	}
}

func (b *UpdateBridge) PopulateUsers(users map[int64]string) {
	b.cacheMu.Lock()
	defer b.cacheMu.Unlock()
	for id, name := range users {
		b.userCache[id] = name
	}
}

func (b *UpdateBridge) LookupUser(id int64) string {
	b.cacheMu.RLock()
	defer b.cacheMu.RUnlock()
	return b.userCache[id]
}

func (b *UpdateBridge) HandleNewMessage(_ context.Context, msgClass tg.MessageClass) error {
	msg, ok := msgClass.(*tg.Message)
	if !ok {
		return nil
	}
	chatID := extractChatID(msg.PeerID)

	senderName := "Unknown"
	if msg.Out {
		senderName = "You"
	} else if from, ok := msg.GetFromID(); ok {
		if p, ok := from.(*tg.PeerUser); ok {
			if name := b.LookupUser(p.UserID); name != "" {
				senderName = name
			}
		}
	}

	item := MsgItem{
		ID:         msg.ID,
		SenderName: senderName,
		Text:       msg.Message,
		Media:      parseMedia(msg.Media),
		Date:       time.Unix(int64(msg.Date), 0),
		Outgoing:   msg.Out,
	}
	b.send(NewMessageMsg{ChatID: chatID, Item: item})
	return nil
}

func extractChatID(peer tg.PeerClass) int64 {
	switch p := peer.(type) {
	case *tg.PeerUser:
		return p.UserID
	case *tg.PeerChat:
		return p.ChatID
	case *tg.PeerChannel:
		return p.ChannelID
	}
	return 0
}

// Tea message types used by the update bridge and app
type NewMessageMsg struct {
	ChatID int64
	Item   MsgItem
}
