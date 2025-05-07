package config

import "github.com/labstack/echo/v4"

func NewEcho() *echo.Echo {
	e := echo.New()
	e.HTTPErrorHandler = nil
	return e
}
