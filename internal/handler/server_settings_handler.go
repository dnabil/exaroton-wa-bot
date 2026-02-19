package handler

import (
	"exaroton-wa-bot/internal/constants/messages"
	"exaroton-wa-bot/internal/dto"
	"exaroton-wa-bot/pages"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (w *Web) SettingsExarotonPage(data *dto.SettingsExarotonPageData) echo.HandlerFunc {
	return func(c echo.Context) error {
		// init data
		if data == nil {
			data = new(dto.SettingsExarotonPageData)
		}

		// init statusCode
		if data.HttpCode == 0 {
			data.HttpCode = http.StatusOK
		}

		apiKey, err := w.svc.ServerSettingsService.GetExarotonAPIKey(c.Request().Context())
		if err != nil {
			return err
		}

		data.APIKey = apiKey

		return c.Render(data.HttpCode, pages.SettingsExaroton, data)
	}
}

func (w *Web) APISettingsExarotonUpdate() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(dto.SettingsExarotonReq)
		if err := w.shouldBind(c, req); err != nil {
			return err
		}

		_, err := w.svc.ServerSettingsService.ValidateExarotonAPIKey(c.Request().Context(), req.APIKey)
		if err != nil {
			return err
		}

		if err := w.svc.ServerSettingsService.UpdateExarotonAPIKey(c.Request().Context(), req.APIKey); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, dto.APIResponse{
			Success: true,
			Message: messages.ResourceCreated,
		})
	}
}

func (w *Web) APISettingsExarotonValidateApiKey() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(dto.SettingsExarotonReq)
		if err := w.shouldBind(c, req); err != nil {
			return err
		}

		acc, err := w.svc.ServerSettingsService.ValidateExarotonAPIKey(c.Request().Context(), req.APIKey)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, dto.APIResponse{
			Success: true,
			Message: messages.ValidKey,
			Data:    acc,
		})
	}
}
