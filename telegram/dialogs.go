package telegram

import (
	"context"

	"github.com/gotd/td/tg"
)

type DialogItem struct {
	ID          int64
	InputPeer   tg.InputPeerClass
	Name        string
	LastMessage string
	UnreadCount int32
	Date        int32
}

type FolderItem struct {
	ID      int
	Title   string
	PeerIDs map[int64]bool // nil = show all chats
}

func (w *ClientWrapper) LoadDialogs(ctx context.Context) ([]DialogItem, error) {
	result, err := w.Raw.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
		Limit:      100,
		OffsetPeer: &tg.InputPeerEmpty{},
	})
	if err != nil {
		return nil, err
	}

	var dialogs []tg.DialogClass
	var msgs []tg.MessageClass
	userMap := make(map[int64]*tg.User)
	chatMap := make(map[int64]*tg.Chat)
	channelMap := make(map[int64]*tg.Channel)

	switch r := result.(type) {
	case *tg.MessagesDialogs:
		dialogs = r.Dialogs
		msgs = r.Messages
		for _, u := range r.Users {
			if user, ok := u.(*tg.User); ok {
				userMap[user.ID] = user
			}
		}
		for _, c := range r.Chats {
			switch ch := c.(type) {
			case *tg.Chat:
				chatMap[ch.ID] = ch
			case *tg.Channel:
				channelMap[ch.ID] = ch
			}
		}
	case *tg.MessagesDialogsSlice:
		dialogs = r.Dialogs
		msgs = r.Messages
		for _, u := range r.Users {
			if user, ok := u.(*tg.User); ok {
				userMap[user.ID] = user
			}
		}
		for _, c := range r.Chats {
			switch ch := c.(type) {
			case *tg.Chat:
				chatMap[ch.ID] = ch
			case *tg.Channel:
				channelMap[ch.ID] = ch
			}
		}
	case *tg.MessagesDialogsNotModified:
		return nil, nil
	}

	msgMap := make(map[int]string)
	for _, m := range msgs {
		if msg, ok := m.(*tg.Message); ok {
			msgMap[msg.ID] = msg.Message
		}
	}

	var items []DialogItem
	for _, d := range dialogs {
		dialog, ok := d.(*tg.Dialog)
		if !ok {
			continue
		}
		var name string
		var inputPeer tg.InputPeerClass
		var id int64

		switch peer := dialog.Peer.(type) {
		case *tg.PeerUser:
			id = peer.UserID
			if u, ok := userMap[peer.UserID]; ok {
				name = u.FirstName
				if u.LastName != "" {
					name += " " + u.LastName
				}
				inputPeer = &tg.InputPeerUser{UserID: u.ID, AccessHash: u.AccessHash}
			}
		case *tg.PeerChat:
			id = peer.ChatID
			if c, ok := chatMap[peer.ChatID]; ok {
				name = c.Title
				inputPeer = &tg.InputPeerChat{ChatID: c.ID}
			}
		case *tg.PeerChannel:
			id = peer.ChannelID
			if ch, ok := channelMap[peer.ChannelID]; ok {
				name = ch.Title
				inputPeer = &tg.InputPeerChannel{ChannelID: ch.ID, AccessHash: ch.AccessHash}
			}
		}

		if name == "" || inputPeer == nil {
			continue
		}

		items = append(items, DialogItem{
			ID:          id,
			InputPeer:   inputPeer,
			Name:        name,
			LastMessage: msgMap[dialog.TopMessage],
			UnreadCount: int32(dialog.UnreadCount),
		})
	}

	nameMap := make(map[int64]string, len(userMap))
	for id, u := range userMap {
		name := u.FirstName
		if u.LastName != "" {
			name += " " + u.LastName
		}
		nameMap[id] = name
	}
	w.Bridge.PopulateUsers(nameMap)

	return items, nil
}

func (w *ClientWrapper) LoadFolders(ctx context.Context) ([]FolderItem, error) {
	result, err := w.Raw.MessagesGetDialogFilters(ctx)
	if err != nil {
		return nil, err
	}
	folders := []FolderItem{{ID: -1, Title: "All Chats"}}
	for _, f := range result.Filters {
		switch v := f.(type) {
		case *tg.DialogFilter:
			peerIDs := make(map[int64]bool)
			for _, p := range v.PinnedPeers {
				if id := inputPeerID(p); id != 0 {
					peerIDs[id] = true
				}
			}
			for _, p := range v.IncludePeers {
				if id := inputPeerID(p); id != 0 {
					peerIDs[id] = true
				}
			}
			folders = append(folders, FolderItem{ID: v.ID, Title: v.Title.Text, PeerIDs: peerIDs})
		case *tg.DialogFilterChatlist:
			peerIDs := make(map[int64]bool)
			for _, p := range v.PinnedPeers {
				if id := inputPeerID(p); id != 0 {
					peerIDs[id] = true
				}
			}
			for _, p := range v.IncludePeers {
				if id := inputPeerID(p); id != 0 {
					peerIDs[id] = true
				}
			}
			folders = append(folders, FolderItem{ID: v.ID, Title: v.Title.Text, PeerIDs: peerIDs})
		}
	}
	return folders, nil
}

func inputPeerID(peer tg.InputPeerClass) int64 {
	switch p := peer.(type) {
	case *tg.InputPeerUser:
		return p.UserID
	case *tg.InputPeerChat:
		return p.ChatID
	case *tg.InputPeerChannel:
		return p.ChannelID
	}
	return 0
}
