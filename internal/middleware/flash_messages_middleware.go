package middleware

import (
	"exaroton-wa-bot/internal/constants"

	"github.com/labstack/echo/v4"
)

// FlashValidationError returns middleware for setting validation_error msg to context from session.
func (m *Middleware) FlashValidationError() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get validation errors from session
			valErr, err := m.session.GetValidationError(c)
			if err != nil {
				return err
			}

			c.Set(constants.FlashValErrCtxKey, valErr)

			return next(c)
		}
	}
}

// FlashOldInput returns middleware for setting old input to context from session.
func (m *Middleware) FlashOldInput() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get old input from session
			oldInput, err := m.session.GetOldInput(c)
			if err != nil {
				return err
			}

			c.Set(constants.FlashOldInputCtxKey, oldInput)

			return next(c)
		}
	}
}

// FlashOldInput returns middleware for setting old input to context from session.
func (m *Middleware) FlashMessage() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get old input from session
			flashMsg, err := m.session.GetFlash(c)
			if err != nil {
				return err
			}

			c.Set(constants.FlashMsgCtxKey, flashMsg)

			return next(c)
		}
	}
}
