package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// binds and validates
func (w *Web) shouldBind(c echo.Context, req any) error {
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.Validate(req)
}
