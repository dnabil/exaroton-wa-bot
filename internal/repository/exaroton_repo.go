package repository

import (
	"context"
	"errors"
	"exaroton-wa-bot/internal/constants/errs"
	"exaroton-wa-bot/internal/dto"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"pkg.icikowski.pl/exaroton"
)

type IExarotonRepo interface {
	ValidateApiKey(ctx context.Context, apiKey string) (*dto.ExarotonAccountInfo, error)
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
	if err := handleExarotonError(err, raw.Error); err != nil {
		return nil, fmt.Errorf("exaroton repo ValidateApiKey error: %w", err)
	}

	return &dto.ExarotonAccountInfo{
		Name:     acc.Name,
		Email:    acc.Email,
		Verified: acc.Verified,
		Credits:  acc.Credits,
	}, nil
}

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
