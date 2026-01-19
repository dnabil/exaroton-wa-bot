package handler

import (
	"exaroton-wa-bot/internal/dto"
	"exaroton-wa-bot/pages"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (w *Web) SettingsWhatsappPage() echo.HandlerFunc {
	return func(c echo.Context) error {
		number, err := w.svc.AuthService.GetWhatsappPhoneNumber(c.Request().Context())
		if err != nil {
			return err
		}

		// TODO: list group tapi berdasarkan whitelist

		return c.Render(http.StatusOK, pages.SettingsWhatsapp, dto.SettingsWhatsappPageData{
			PhoneNumber: number,
		})
	}
}
