package service

import (
	"context"
	"exaroton-wa-bot/internal/constants"
	"exaroton-wa-bot/internal/constants/errs"
	"exaroton-wa-bot/internal/database/entity"
	"exaroton-wa-bot/internal/dto"
	"exaroton-wa-bot/internal/helper"
	"exaroton-wa-bot/internal/repository"
	"log/slog"
	"time"
)

type IServerSettingsService interface {
	GetExarotonAPIKey(ctx context.Context) (string, error)
	UpdateExarotonAPIKey(ctx context.Context, apiKey string) error

	ValidateExarotonAPIKey(ctx context.Context, apiKey string) (*dto.ExarotonAccountInfo, error)
	ListExarotonServer(ctx context.Context) ([]*dto.ExarotonServerInfo, error)
	StartExarotonServer(ctx context.Context, serverIdx uint, opts ...StartExarotonServerOption) *dto.StartExarotonServerRes
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

type (
	startExarotonServerConfig struct {
		poll         bool
		timeout      time.Duration
		interval     time.Duration
		useOwnCredit bool
	}

	StartExarotonServerOption func(*startExarotonServerConfig)
)

func WithPolling(timeout, interval time.Duration) StartExarotonServerOption {
	// default values
	interval = helper.If((interval <= 5*time.Second), (5 * time.Second), interval)
	timeout = helper.If((timeout <= 10*time.Second), (10 * time.Second), timeout)

	// enforce interval <= timeout
	if interval > timeout {
		interval = timeout
	}

	return func(c *startExarotonServerConfig) {
		c.poll = true
		c.timeout = timeout
		c.interval = interval
	}
}

func WithOwnCredit() StartExarotonServerOption {
	return func(c *startExarotonServerConfig) {
		c.useOwnCredit = true
	}
}

func (s *ServerSettingsService) StartExarotonServer(ctx context.Context, serverIdx uint, opts ...StartExarotonServerOption) (res *dto.StartExarotonServerRes) {
	tx, res := s.tx.Begin(ctx), new(dto.StartExarotonServerRes)
	defer func() {
		if rbErr := s.tx.Rollback(tx); rbErr != nil {
			slog.ErrorContext(ctx, rbErr.Error())
		}
	}()

	// opt
	cfg := new(startExarotonServerConfig)
	for _, opt := range opts {
		opt(cfg)
	}

	settings, err := s.serverSettingsRepo.Get(ctx, tx, constants.ExarotonAPIKey)
	if err != nil {
		res.Err = err
		return
	}
	if settings == nil {
		res.Err = errs.ErrGSEmptyAPIKey
		return
	}

	apiKey := settings.Value

	servers, err := s.exarotonRepo.ListServers(ctx, apiKey)
	if err != nil {
		res.Err = err
		return
	}

	if serverIdx >= uint(len(servers)) {
		res.Err = errs.ErrServerNotFound
		return
	}

	// start the server
	startServerReq := dto.StartExarotonServerReq{UseOwnCredit: cfg.useOwnCredit}
	err = s.exarotonRepo.StartServer(ctx, apiKey, servers[serverIdx].ID, startServerReq)
	if err != nil {
		res.Err = err
		return
	}

	statusCh := make(chan dto.ServerStatus, 1)

	// polling to check status
	if cfg.poll {
		go func() {
			status := dto.ServerStatusStarting
			defer close(statusCh)

			ticker := time.NewTicker(cfg.interval)
			defer ticker.Stop()

			pollCtx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
			defer cancel()

			for {
				select {
				case <-pollCtx.Done():
					statusCh <- status
					return

				case <-ticker.C:
					srv, err := s.exarotonRepo.GetServerInfo(pollCtx, apiKey, servers[serverIdx].ID)
					if err != nil {
						slog.ErrorContext(ctx, "polling error", "err", err)
						return
					}

					statusCh <- srv.Status

					if srv.Status == dto.ServerStatusOnline ||
						srv.Status == dto.ServerStatusCrashed {
						return
					}
				}
			}
		}()
	}

	return &dto.StartExarotonServerRes{
		Status: statusCh,
		Err:    nil,
	}
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
