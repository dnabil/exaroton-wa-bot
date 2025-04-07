package handler

import (
	"exaroton-wa-bot/internal/config"
	"exaroton-wa-bot/internal/middleware"
	"exaroton-wa-bot/pages"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func (w *Web) routes() {
	staticDir, err := filepath.Abs(w.cfg.MustString(config.KeyPublicDir))
	if err != nil {
		panic(fmt.Sprintf("failed to get public dir (key: %s)", config.KeyPublicDir))
	}

	// static files
	w.Router.Static("/public", staticDir)
	w.Router.File("/favicon.ico", filepath.Join(staticDir, "favicon.ico"))

	// global middlewares
	w.Router.Use(middleware.Loggger(w.cfg))
	w.Router.Use(middleware.Recover(w.cfg))

	// routes
	w.Router.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, pages.Index, nil)
	})
}
