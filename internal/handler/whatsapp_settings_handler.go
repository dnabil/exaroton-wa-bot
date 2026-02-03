package handler

import (
	"exaroton-wa-bot/internal/constants/messages"
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

		return c.Render(http.StatusOK, pages.SettingsWhatsapp, dto.SettingsWhatsappPageData{
			PhoneNumber: number,
		})
	}
}

func (w *Web) APIWhatsappLogout() echo.HandlerFunc {
	return func(c echo.Context) error {
		return w.svc.AuthService.WhatsappLogout(c.Request().Context())
	}
}

func (w *Web) APIWhatsappGroupWhitelist() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(dto.WhitelistWhatsappGroupReq)

		err := w.shouldBind(c, req)
		if err != nil {
			return err
		}

		if err = w.svc.WhatsappService.WhitelistGroup(c.Request().Context(), req); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, &dto.APIResponse{
			Success: true,
			Message: messages.GroupWhitelistSuccess,
		})
	}
}

func (w *Web) APIGetWhatsappGroups() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(dto.GetWhatsappGroupReq)

		err := w.shouldBind(c, req)
		if err != nil {
			return err
		}

		res, err := w.svc.WhatsappService.GetGroups(c.Request().Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, &dto.APIResponse{
			Success: true,
			Data:    res,
		})
	}
}
