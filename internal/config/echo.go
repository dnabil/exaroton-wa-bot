package config

import "github.com/labstack/echo/v4"

func NewEcho() *echo.Echo {
	return echo.New()
}
