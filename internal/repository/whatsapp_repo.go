package repository

import (
	"context"
	"exaroton-wa-bot/internal/config"
	"exaroton-wa-bot/internal/database/entity"
	"exaroton-wa-bot/internal/dto"
	"log/slog"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"gorm.io/gorm"
)

type IWhatsappRepo interface {
	Disconnect()
	Login(ctx context.Context) (<-chan whatsmeow.QRChannelItem, error)
	Logout(ctx context.Context) error
	IsLoggedIn() bool
	GetPhoneNumber(ctx context.Context) (string, error)
	GetGroups(ctx context.Context) ([]*types.GroupInfo, error)
	GetWhitelistedJIDs(ctx context.Context, tx *gorm.DB) ([]*entity.WhatsappWhitelistedGroup, error)
	WhitelistGroup(ctx context.Context, tx *gorm.DB, req *dto.WhitelistWhatsappGroupReq) error
	UnwhitelistGroup(ctx context.Context, tx *gorm.DB, req *dto.UnwhitelistWhatsappGroupReq) error
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

func (r *whatsappRepo) Logout(ctx context.Context) error {
	return r.waClient.Logout(ctx)
}

func (r *whatsappRepo) IsLoggedIn() bool {
	return r.waClient.IsLoggedIn()
}

func (r *whatsappRepo) GetPhoneNumber(ctx context.Context) (string, error) {
	return r.waClient.GetPhoneNumber(ctx)
}

func (r *whatsappRepo) GetGroups(ctx context.Context) ([]*types.GroupInfo, error) {
	return r.waClient.GetGroups(ctx)
}

func (r *whatsappRepo) GetWhitelistedJIDs(ctx context.Context, tx *gorm.DB) ([]*entity.WhatsappWhitelistedGroup, error) {
	whitelistedGroups := make([]*entity.WhatsappWhitelistedGroup, 0)
	if err := tx.Find(&whitelistedGroups).Error; err != nil {
		return nil, err
	}

	return whitelistedGroups, nil
}

func (r *whatsappRepo) WhitelistGroup(ctx context.Context, tx *gorm.DB, req *dto.WhitelistWhatsappGroupReq) error {
	return tx.Create(entity.WhatsappWhitelistedGroup{
		JID:       req.User,
		ServerJID: req.Server,
	}).Error
}

func (r *whatsappRepo) UnwhitelistGroup(ctx context.Context, tx *gorm.DB, req *dto.UnwhitelistWhatsappGroupReq) error {
	return tx.Where(entity.WhatsappWhitelistedGroup{
		JID:       req.User,
		ServerJID: req.Server,
	}).Delete(&entity.WhatsappWhitelistedGroup{}).Error
}
