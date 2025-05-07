package middleware

import (
	"exaroton-wa-bot/internal/config"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (m *Middleware) Session() echo.MiddlewareFunc {
	return session.MiddlewareWithConfig(session.Config{
		Store: sessions.NewCookieStore([]byte(m.cfg.MustString(config.KeySessionSecret))),
	})
}
