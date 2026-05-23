package telegram

import (
	"context"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	appauth "github.com/sangnguyen/tui-telegram/auth"
	"github.com/sangnguyen/tui-telegram/config"
)

type ClientWrapper struct {
	Raw         *tg.Client
	Sender      *message.Sender
	tgClient    *telegram.Client
	Bridge      *UpdateBridge
}

func New(cfg *config.Config, sessionPath string, bridge *UpdateBridge) (*ClientWrapper, error) {
	storage := &session.FileStorage{Path: sessionPath}
	dispatcher := tg.NewUpdateDispatcher()
	dispatcher.OnNewMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
		return bridge.HandleNewMessage(ctx, update.Message)
	})
	dispatcher.OnNewChannelMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewChannelMessage) error {
		return bridge.HandleNewMessage(ctx, update.Message)
	})

	client := telegram.NewClient(cfg.APIID, cfg.APIHash, telegram.Options{
		SessionStorage: storage,
		UpdateHandler:  dispatcher,
	})

	return &ClientWrapper{
		tgClient: client,
		Bridge:   bridge,
	}, nil
}

func (w *ClientWrapper) Run(ctx context.Context, authenticator *appauth.Authenticator, onReady func(context.Context) error) error {
	return w.tgClient.Run(ctx, func(ctx context.Context) error {
		if err := w.tgClient.Auth().IfNecessary(ctx, auth.NewFlow(
			authenticator,
			auth.SendCodeOptions{},
		)); err != nil {
			return err
		}
		w.Raw = w.tgClient.API()
		w.Sender = message.NewSender(w.Raw)
		return onReady(ctx)
	})
}
