package middleware

import (
	"exaroton-wa-bot/internal/config"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func Recover(cfg *config.Cfg) echo.MiddlewareFunc {
	mdwConfig := echoMiddleware.DefaultRecoverConfig
	mdwConfig.LogErrorFunc = func(c echo.Context, err error, stack []byte) error {
		config.ErrLog(c.Request().Context(), cfg.Args, err, stack)
		return err // let custom http error handler properly response to client.
	}

	return echoMiddleware.RecoverWithConfig(mdwConfig)
}
