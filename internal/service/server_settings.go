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

	ValidateExarotonAPIKey(ctx context.Context, apiKey string) (*dto.ExarotonAccountInfo, error)
	ListExarotonServer(ctx context.Context) ([]*dto.ExarotonServerInfo, error)
	StartExarotonServer(ctx context.Context, serverIdx uint) error
	StopExarotonServer(ctx context.Context, serverIdx uint) error
	GetExarotonServerInfo(ctx context.Context, serverIdx uint) (*dto.ExarotonServerInfo, error)
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

func (s *ServerSettingsService) ListExarotonServer(ctx context.Context) ([]*dto.ExarotonServerInfo, error) {
	tx := s.tx.Begin(ctx)
	defer func() {
		if rbErr := s.tx.Rollback(tx); rbErr != nil {
			slog.ErrorContext(ctx, rbErr.Error())
		}
	}()

	settings, err := s.serverSettingsRepo.Get(ctx, tx, constants.ExarotonAPIKey)
	if err != nil {
		return nil, err
	}

	if settings == nil {
		return nil, errs.ErrGSEmptyAPIKey
	}

	apiKey := settings.Value

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

	if serverIdx >= uint(len(servers)) {
		return errs.ErrServerNotFound
	}

	return s.exarotonRepo.StartServer(ctx, apiKey, servers[serverIdx].ID)
}

func (s *ServerSettingsService) StopExarotonServer(ctx context.Context, serverIdx uint) error {
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

	if serverIdx >= uint(len(servers)) {
		return errs.ErrServerNotFound
	}

	return s.exarotonRepo.StopServer(ctx, apiKey, servers[serverIdx].ID)
}

func (s *ServerSettingsService) GetExarotonServerInfo(ctx context.Context, serverIdx uint) (*dto.ExarotonServerInfo, error) {
	tx := s.tx.Begin(ctx)
	defer func() {
		if rbErr := s.tx.Rollback(tx); rbErr != nil {
			slog.ErrorContext(ctx, rbErr.Error())
		}
	}()

	settings, err := s.serverSettingsRepo.Get(ctx, tx, constants.ExarotonAPIKey)
	if err != nil {
		return nil, err
	}

	if settings == nil {
		return nil, errs.ErrGSEmptyAPIKey
	}

	apiKey := settings.Value

	servers, err := s.exarotonRepo.ListServers(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	if serverIdx >= uint(len(servers)) {
		return nil, errs.ErrServerNotFound
	}

	return s.exarotonRepo.GetServerInfo(ctx, apiKey, servers[serverIdx].ID)
}
