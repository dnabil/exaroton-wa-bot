package service

import (
	"context"
	"exaroton-wa-bot/internal/constants"
	"exaroton-wa-bot/internal/constants/errs"
	"exaroton-wa-bot/internal/database/entity"
	"exaroton-wa-bot/internal/dto"
	"exaroton-wa-bot/internal/repository"
	"log/slog"
)

type IServerSettingsService interface {
	GetExarotonAPIKey(ctx context.Context) (string, error)
	UpdateExarotonAPIKey(ctx context.Context, apiKey string) error

	// TODO: remove apiKey from the arg, orhcestrate in this layer
	ValidateExarotonAPIKey(ctx context.Context, apiKey string) (*dto.ExarotonAccountInfo, error)
	// TODO: remove apiKey from the arg, orhcestrate in this layer
	ListExarotonServer(ctx context.Context, apiKey string) ([]*dto.ExarotonServerInfo, error)
	StartExarotonServer(ctx context.Context, serverIdx uint) error
}

type ServerSettingsService struct {
	*svcTmpl
	serverSettingsRepo repository.IServerSettingsRepo
	exarotonRepo       repository.IExarotonRepo
}

func NewServerSettingsService(svcTmpl *svcTmpl, serverSettingsRepo repository.IServerSettingsRepo, exarotonRepo repository.IExarotonRepo) IServerSettingsService {
	return &ServerSettingsService{
		svcTmpl:            svcTmpl,
		serverSettingsRepo: serverSettingsRepo,
		exarotonRepo:       exarotonRepo,
	}
}

func (s *ServerSettingsService) GetExarotonAPIKey(ctx context.Context) (string, error) {
	tx := s.tx.Begin(ctx)
	defer func() {
		if rbErr := s.tx.Rollback(tx); rbErr != nil {
			slog.ErrorContext(ctx, rbErr.Error())
		}
	}()

	settings, err := s.serverSettingsRepo.Get(ctx, tx, constants.ExarotonAPIKey)
	if err != nil {
		return "", err
	}

	if settings == nil {
		return "", nil
	}

	return settings.Value, nil
}

func (s *ServerSettingsService) UpdateExarotonAPIKey(ctx context.Context, apiKey string) error {
	tx := s.tx.Begin(ctx)
	defer func() {
		if rbErr := s.tx.Rollback(tx); rbErr != nil {
			slog.ErrorContext(ctx, rbErr.Error())
		}
	}()

	err := s.serverSettingsRepo.Upsert(ctx, tx, &entity.ServerSettings{
		Key:   constants.ExarotonAPIKey,
		Value: apiKey,
	})
	if err != nil {
		return err
	}

	return s.tx.Commit(tx)
}

func (s *ServerSettingsService) ValidateExarotonAPIKey(ctx context.Context, apiKey string) (*dto.ExarotonAccountInfo, error) {
	return s.exarotonRepo.ValidateApiKey(ctx, apiKey)
}

func (s *ServerSettingsService) ListExarotonServer(ctx context.Context, apiKey string) ([]*dto.ExarotonServerInfo, error) {
	return s.exarotonRepo.ListServers(ctx, apiKey)
}

func (s *ServerSettingsService) StartExarotonServer(ctx context.Context, serverIdx uint) error {
	tx := s.tx.Begin(ctx)
	defer func() {
		if rbErr := s.tx.Rollback(tx); rbErr != nil {
			slog.ErrorContext(ctx, rbErr.Error())
		}
	}()

	settings, err := s.serverSettingsRepo.Get(ctx, tx, constants.ExarotonAPIKey)
	if err != nil {
		return err
	}

	if settings == nil {
		return errs.ErrGSEmptyAPIKey
	}

	apiKey := settings.Value

	servers, err := s.exarotonRepo.ListServers(ctx, apiKey)
	if err != nil {
		return err
	}

	return s.exarotonRepo.StartServer(ctx, apiKey, servers[serverIdx].ID)
}
