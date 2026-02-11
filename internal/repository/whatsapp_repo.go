package repository

import (
	"context"
	"exaroton-wa-bot/internal/database/entity"
	"exaroton-wa-bot/internal/dto"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"gorm.io/gorm"
)

type IWhatsappRepo interface {
	RegisterEventHandler(f func(any)) uint32
	UnregisterEventHandler(handlerID uint32) bool
	Disconnect()
	Login(ctx context.Context) (<-chan whatsmeow.QRChannelItem, error)
	Logout(ctx context.Context) error
	IsLoggedIn() bool
	GetPhoneNumber() string // self
	GetSelfLID() *dto.WhatsappJID
	GetGroups(ctx context.Context) ([]*types.GroupInfo, error)
	GetWhitelistedGroupJIDs(ctx context.Context, tx *gorm.DB) ([]*entity.WhatsappWhitelistedGroup, error)
	WhitelistGroup(ctx context.Context, tx *gorm.DB, req *dto.WhitelistWhatsappGroupReq) error
	UnwhitelistGroup(ctx context.Context, tx *gorm.DB, req *dto.UnwhitelistWhatsappGroupReq) error

	// IsSyncComplete returns true if the sync is complete and false otherwise.
	IsSyncComplete(ctx context.Context) bool
}

type whatsappRepo struct {
	waClient *waClient // represents a single whatsapp device/account.
}

func newWhatsappRepo(waClient *waClient) IWhatsappRepo {
	return &whatsappRepo{
		waClient: waClient,
	}
}

func (r *whatsappRepo) RegisterEventHandler(f func(any)) uint32 {
	return r.waClient.client.RegisterEventHandler(f)
}

func (r *whatsappRepo) UnregisterEventHandler(handlerID uint32) bool {
	return r.waClient.client.UnregisterEventHandler(handlerID)
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

func (r *whatsappRepo) GetPhoneNumber() string {
	return r.waClient.GetPhoneNumber()
}

func (r *whatsappRepo) GetSelfLID() *dto.WhatsappJID {
	lid, err := r.waClient.GetSelfLID()
	if err != nil {
		return nil
	}

	return &dto.WhatsappJID{
		User:       lid.User,
		RawAgent:   lid.RawAgent,
		Device:     lid.Device,
		Integrator: lid.Integrator,
		Server:     lid.Server,
	}
}

func (r *whatsappRepo) GetGroups(ctx context.Context) ([]*types.GroupInfo, error) {
	return r.waClient.GetGroups(ctx)
}

func (r *whatsappRepo) GetWhitelistedGroupJIDs(ctx context.Context, tx *gorm.DB) ([]*entity.WhatsappWhitelistedGroup, error) {
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

func (r *whatsappRepo) IsSyncComplete(ctx context.Context) bool {
	return r.waClient.IsSyncComplete(ctx)
}
