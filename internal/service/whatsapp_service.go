package service

import (
	"context"
	"exaroton-wa-bot/internal/dto"
	"exaroton-wa-bot/internal/repository"
	"log/slog"

	"go.mau.fi/whatsmeow/types"
)

type IWhatsappService interface {
	WhitelistGroup(ctx context.Context, req *dto.WhitelistWhatsappGroupReq) error
	GetGroups(ctx context.Context, req *dto.GetWhatsappGroupReq) ([]*dto.WhatsappGroupInfo, error)
}

type WhatsappService struct {
	*svcTmpl
	waRepo repository.IWhatsappRepo
}

func NewWhatsappService(svcTmpl *svcTmpl, waRepo repository.IWhatsappRepo) IWhatsappService {
	return &WhatsappService{
		svcTmpl: svcTmpl,
		waRepo:  waRepo,
	}
}

func (s *WhatsappService) WhitelistGroup(ctx context.Context, req *dto.WhitelistWhatsappGroupReq) error {
	tx := s.tx.Begin(ctx)
	defer func() {
		if rbErr := s.tx.Rollback(tx); rbErr != nil {
			slog.ErrorContext(ctx, rbErr.Error())
		}
	}()

	if err := s.waRepo.WhitelistGroup(ctx, tx, req); err != nil {
		return err
	}

	return s.tx.Commit(tx)
}

func (s *WhatsappService) GetGroups(ctx context.Context, req *dto.GetWhatsappGroupReq) ([]*dto.WhatsappGroupInfo, error) {
	allGroups, err := s.waRepo.GetGroups(ctx)
	if err != nil {
		return nil, err
	}

	if req.Whitelist == nil {
		dtoGroups := make([]*dto.WhatsappGroupInfo, len(allGroups))
		for i, g := range allGroups {
			dtoGroups[i] = dto.NewWhatsappGroupInfo(g)
		}

		return dtoGroups, nil
	}

	// whitelist filter
	tx := s.tx.Begin(ctx)
	defer func() {
		if rbErr := s.tx.Rollback(tx); rbErr != nil {
			slog.ErrorContext(ctx, rbErr.Error())
		}
	}()

	jids, err := s.waRepo.GetWhitelistedJIDs(ctx, tx)
	if err != nil {
		return nil, err
	}

	whitelistMap := make(map[types.JID]bool)
	for _, j := range jids {
		whitelistMap[types.NewJID(j.JID, j.ServerJID)] = true
	}

	isWhitelist := *req.Whitelist
	filteredGroups := make([]*dto.WhatsappGroupInfo, 0)
	for _, g := range allGroups {
		_, inWhitelist := whitelistMap[g.JID]
		if (isWhitelist && inWhitelist) || (!isWhitelist && !inWhitelist) {
			filteredGroups = append(filteredGroups, dto.NewWhatsappGroupInfo(g))
		}
	}

	// might need to make the logic above cleaner later
	// but rn it'll do

	return filteredGroups, nil
}
