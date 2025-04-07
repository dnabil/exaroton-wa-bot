package middleware

import (
	"exaroton-wa-bot/internal/config"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Loggger(cfg *config.Cfg) echo.MiddlewareFunc {
	logger := slog.Default()

	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []slog.Attr{
				slog.String(config.KeyLogURI, v.URI),
				slog.Int(config.KeyLogStatus, v.Status),
			}

			if v.Error != nil {
				config.ErrLog(c.Request().Context(), cfg.Args, v.Error, nil, attrs...)
			} else {
				logger.LogAttrs(c.Request().Context(), slog.LevelInfo, "REQUEST", attrs...)
			}

			return nil
		},
	})
}
