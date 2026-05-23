package telegram

import (
	"context"
	"time"

	"github.com/gotd/td/tg"
	"github.com/gotd/td/telegram/message/unpack"
)

type MsgItem struct {
	ID         int
	SenderName string
	Text       string
	Media      string
	Date       time.Time
	Outgoing   bool
}

func (w *ClientWrapper) LoadHistory(ctx context.Context, peer tg.InputPeerClass, limit int32, chatName string) ([]MsgItem, error) {
	result, err := w.Raw.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
		Peer:  peer,
		Limit: int(limit),
	})
	if err != nil {
		return nil, err
	}

	var rawMsgs []tg.MessageClass
	userMap := make(map[int64]*tg.User)

	switch r := result.(type) {
	case *tg.MessagesMessages:
		rawMsgs = r.Messages
		for _, u := range r.Users {
			if user, ok := u.(*tg.User); ok {
				userMap[user.ID] = user
			}
		}
	case *tg.MessagesMessagesSlice:
		rawMsgs = r.Messages
		for _, u := range r.Users {
			if user, ok := u.(*tg.User); ok {
				userMap[user.ID] = user
			}
		}
	case *tg.MessagesChannelMessages:
		rawMsgs = r.Messages
		for _, u := range r.Users {
			if user, ok := u.(*tg.User); ok {
				userMap[user.ID] = user
			}
		}
	}

	var items []MsgItem
	for _, m := range rawMsgs {
		msg, ok := m.(*tg.Message)
		if !ok {
			continue
		}

		senderName := chatName
		if senderName == "" {
			senderName = "Unknown"
		}
		if msg.Out {
			senderName = "You"
		} else if from, ok := msg.GetFromID(); ok {
			if p, ok := from.(*tg.PeerUser); ok {
				if u, ok := userMap[p.UserID]; ok {
					senderName = u.FirstName
					if u.LastName != "" {
						senderName += " " + u.LastName
					}
				}
			}
		}

		items = append(items, MsgItem{
			ID:         msg.ID,
			SenderName: senderName,
			Text:       msg.Message,
			Media:      parseMedia(msg.Media),
			Date:       time.Unix(int64(msg.Date), 0),
			Outgoing:   msg.Out,
		})
	}

	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}

	return items, nil
}

func parseMedia(media tg.MessageMediaClass) string {
	if media == nil {
		return ""
	}
	switch m := media.(type) {
	case *tg.MessageMediaPhoto:
		_ = m
		return "[Photo]"
	case *tg.MessageMediaDocument:
		if doc, ok := m.Document.(*tg.Document); ok {
			for _, attr := range doc.Attributes {
				switch a := attr.(type) {
				case *tg.DocumentAttributeFilename:
					return "[Document: " + a.FileName + "]"
				case *tg.DocumentAttributeVideo:
					_ = a
					return "[Video]"
				}
			}
		}
		return "[Document]"
	case *tg.MessageMediaUnsupported:
		return "[Unsupported media]"
	case *tg.MessageMediaGeo:
		return "[Location]"
	case *tg.MessageMediaContact:
		return "[Contact]"
	}
	return ""
}

func (w *ClientWrapper) SendMessage(ctx context.Context, peer tg.InputPeerClass, text string) (MsgItem, error) {
	msgID, err := unpack.MessageID(w.Sender.To(peer).Text(ctx, text))
	if err != nil {
		return MsgItem{}, err
	}
	return MsgItem{
		ID:         msgID,
		SenderName: "You",
		Text:       text,
		Date:       time.Now(),
		Outgoing:   true,
	}, nil
}
