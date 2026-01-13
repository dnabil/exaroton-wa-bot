package handler

import (
	"exaroton-wa-bot/internal/dto"
	"net/http"

	"github.com/labstack/echo/v4"
)

// binds and validates
func (w *Web) shouldBind(c echo.Context, req any) error {
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// setting old input if request is mappable
	if reqMappable, ok := req.(dto.Mappable); ok {
		if err = w.session.SetOldInput(c, reqMappable); err != nil {
			return err
		}
	}

	return c.Validate(req)
}
