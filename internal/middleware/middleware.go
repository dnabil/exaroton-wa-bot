package middleware

import (
	"exaroton-wa-bot/internal/config"
	"exaroton-wa-bot/internal/dto"
	"exaroton-wa-bot/internal/service"
)

type Middleware struct {
	cfg     *config.Cfg
	authSvc service.IAuthService
	session dto.WebSession
}

func NewMiddleware(
	cfg *config.Cfg,
	authSvc service.IAuthService,
	session dto.WebSession,
) *Middleware {
	return &Middleware{
		cfg:     cfg,
		authSvc: authSvc,
		session: session,
	}
}
