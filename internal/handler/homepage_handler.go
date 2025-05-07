package handler

import (
	"exaroton-wa-bot/pages"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (w *Web) HomePage() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, pages.Index, nil)
	}
}
