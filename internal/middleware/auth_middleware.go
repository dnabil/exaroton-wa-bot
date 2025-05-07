package middleware

import (
	"exaroton-wa-bot/internal/constants/errs"

	"github.com/labstack/echo/v4"
)

func (m *Middleware) Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, err := m.session.GetUser(c)
			if err != nil || user == nil {
				return errs.ErrUserNotLoggedIn
			}

			return next(c)
		}
	}
}

func (m *Middleware) Guest() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, _ := m.session.GetUser(c)
			if user != nil {
				return errs.ErrUserAlreadyLoggedIn
			}

			return next(c)
		}
	}
}

func (m *Middleware) WhatsappLoggedIn(isGuest bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isLoggedIn := m.authSvc.WhatsappIsLoggedIn(c.Request().Context())

			if (isGuest && !isLoggedIn) || (!isGuest && isLoggedIn) {
				return next(c)
			}

			if isGuest {
				return errs.ErrUserAlreadyLoggedIn
			}
			return errs.ErrWANotLoggedIn
		}
	}
}
