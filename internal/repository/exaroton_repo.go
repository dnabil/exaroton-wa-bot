package repository

import (
	"context"
	"errors"
	"exaroton-wa-bot/internal/constants/errs"
	"exaroton-wa-bot/internal/dto"
	"exaroton-wa-bot/internal/helper"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"pkg.icikowski.pl/exaroton"
)

type IExarotonRepo interface {
	ValidateApiKey(ctx context.Context, apiKey string) (*dto.ExarotonAccountInfo, error)
	ListServers(ctx context.Context, apiKey string) ([]*dto.ExarotonServerInfo, error)
	StartServer(ctx context.Context, apiKey string, serverID string) (err error)
	StopServer(ctx context.Context, apiKey string, serverID string) (err error)
	GetServerInfo(ctx context.Context, apiKey string, serverID string) (*dto.ExarotonServerInfo, error)
}

func newExarotonRepo() IExarotonRepo {
	return &ExarotonRepo{}
}

type ExarotonRepo struct{}

func (r *ExarotonRepo) ValidateApiKey(ctx context.Context, apiKey string) (*dto.ExarotonAccountInfo, error) {
	client, err := exaroton.NewClient(apiKey)
	if err != nil {
		return nil, err
	}

	acc, raw, err := client.GetAccount(ctx)
	if err := handleExarotonError(err, helper.Deref(raw).Error); err != nil {
		if errors.Is(err, errs.ErrForbidden) {
			return nil, errs.ErrGSInvalidAPIKey
		}

		return nil, fmt.Errorf("exaroton repo ValidateApiKey error: %w", err)
	}

	return &dto.ExarotonAccountInfo{
		Name:     acc.Name,
		Email:    acc.Email,
		Verified: acc.Verified,
		Credits:  acc.Credits,
	}, nil
}

func (r *ExarotonRepo) ListServers(ctx context.Context, apiKey string) ([]*dto.ExarotonServerInfo, error) {
	client, err := exaroton.NewClient(apiKey)
	if err != nil {
		return nil, err
	}

	serversResponse, raw, err := client.GetServers(ctx)
	if err := handleExarotonError(err, helper.Deref(raw).Error); err != nil {
		return nil, fmt.Errorf("exaroton repo ListServers error: %w", err)
	}

	servers := make([]*dto.ExarotonServerInfo, 0, len(serversResponse))
	for _, srv := range serversResponse {
		servers = append(servers, &dto.ExarotonServerInfo{
			ID:      srv.ID,
			Name:    srv.Name,
			Address: srv.Address,
			Motd:    srv.Motd,
			Status:  dto.ServerStatus(srv.Status),
			Host:    srv.Host,
			Port:    srv.Port,
			Shared:  srv.Shared,
			Players: dto.ExarotonServerPlayers{
				Max:   srv.Players.Max,
				Count: srv.Players.Count,
				List:  srv.Players.List,
			},
			Software: &dto.ExarotonServerSoftware{
				ID:      srv.Software.ID,
				Name:    srv.Software.Name,
				Version: srv.Software.Version,
			},
		})
	}

	return servers, nil
}

func (r *ExarotonRepo) StartServer(ctx context.Context, apiKey string, serverID string) (err error) {
	client, err := exaroton.NewClient(apiKey)
	if err != nil {
		return err
	}

	serverAPI := client.Server(serverID)
	raw, err := serverAPI.Start(ctx)
	if err := handleExarotonError(err, helper.Deref(raw).Error); err != nil {
		return fmt.Errorf("exaroton repo StartServer error: %w", err)
	}

	return nil
}

func (r *ExarotonRepo) StopServer(ctx context.Context, apiKey string, serverID string) (err error) {
	client, err := exaroton.NewClient(apiKey)
	if err != nil {
		return err
	}

	serverAPI := client.Server(serverID)
	raw, err := serverAPI.Stop(ctx)
	if err := handleExarotonError(err, helper.Deref(raw).Error); err != nil {
		if errors.Is(err, errs.ErrAlreadyReported) {
			return errs.ErrServerIsAlreadyStopping
		}
		return fmt.Errorf("exaroton repo StopServer error: %w", err)
	}

	return nil
}

func (r *ExarotonRepo) GetServerInfo(ctx context.Context, apiKey string, serverID string) (*dto.ExarotonServerInfo, error) {
	client, err := exaroton.NewClient(apiKey)
	if err != nil {
		return nil, err
	}

	serverAPI := client.Server(serverID)
	result, raw, err := serverAPI.GetServer(ctx)
	if err := handleExarotonError(err, helper.Deref(raw).Error); err != nil {
		return nil, fmt.Errorf("exaroton repo GetServerInfo error: %w", err)
	}

	return dto.NewExarotonServerInfo(result), nil
}

// =================================================================
// Helpers
// =================================================================

func handleExarotonError(err error, msg *string) error {
	if err == nil {
		return nil
	}

	httpCode, ok := dissectExarotonError(err)
	if !ok {
		return fmt.Errorf("not an exaroton error: %w", err)
	}

	errMsg := "(exaroton doesn't give error message)"
	if msg != nil {
		errMsg = *msg
	}

	err = errors.New(errMsg)

	switch httpCode {
	case http.StatusUnauthorized:
		err = errs.ErrUnauthorized
	case http.StatusForbidden:
		err = errs.ErrForbidden
	// unique case, somehow 208 is considered as an error
	case http.StatusAlreadyReported:
		err = errs.ErrAlreadyReported
	}

	return fmt.Errorf("error: %w, http code: %d, error message: %s", err, httpCode, errMsg)
}

// dissectExarotonError takes an error and tries to extract the status code from it.
// If extraction fails, it returns 0 and false. Otherwise, it returns the status code
// and true.
func dissectExarotonError(err error) (int, bool) {
	matches := regexp.MustCompile(`API error: \[(\d+)]`).FindStringSubmatch(err.Error())
	if len(matches) != 2 {
		return 0, false
	}

	code, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, false
	}

	return code, true
}
