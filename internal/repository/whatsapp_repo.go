package repository

import (
	"context"
	"exaroton-wa-bot/internal/config"
	"log/slog"

	"go.mau.fi/whatsmeow"
)

type IWhatsappRepo interface {
	Disconnect()
	Login(ctx context.Context) (<-chan whatsmeow.QRChannelItem, error)
	IsLoggedIn() bool
}

type whatsappRepo struct {
	waClient *waClient // represents a single whatsapp device/account.
}

func newWhatsappRepo(waDB *config.WhatsappDB) (IWhatsappRepo, error) {
	waClient, err := NewWAClient(waDB)
	if err != nil {
		slog.Error("failed to create whatsapp client in repo", "error", err)
		return nil, err
	}

	return &whatsappRepo{
		waClient: waClient,
	}, nil
}

func (r *whatsappRepo) Disconnect() {
	r.waClient.Disconnect()
}

func (r *whatsappRepo) Login(ctx context.Context) (<-chan whatsmeow.QRChannelItem, error) {
	return r.waClient.Login(ctx)
}

func (r *whatsappRepo) IsLoggedIn() bool {
	return r.waClient.IsLoggedIn()
}
