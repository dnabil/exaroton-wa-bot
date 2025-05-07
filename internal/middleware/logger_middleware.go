package middleware

import (
	"exaroton-wa-bot/internal/config"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (m *Middleware) Logger() echo.MiddlewareFunc {
	logger := slog.Default()

	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogMethod: true,
		LogURI:    true,
		LogError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []slog.Attr{
				slog.Int(config.KeyLogStatus, v.Status),
				slog.String(config.KeyLogMethod, v.Method),
				slog.String(config.KeyLogURI, v.URI),
			}

			if v.Error != nil {
				config.ErrLog(c.Request().Context(), v.Error, nil, attrs...)
			} else {
				logger.LogAttrs(c.Request().Context(), slog.LevelInfo, "REQUEST", attrs...)
			}

			return nil
		},
	})
}
