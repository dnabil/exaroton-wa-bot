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
