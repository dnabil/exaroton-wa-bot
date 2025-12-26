package handler

import (
	"errors"
	"exaroton-wa-bot/internal/constants"
	"exaroton-wa-bot/internal/constants/errs"
	"exaroton-wa-bot/internal/dto"
	"exaroton-wa-bot/pages"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
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

func (w *Web) SettingsExarotonUpdate() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(dto.SettingsExarotonReq)
		if err := w.shouldBind(c, req); err != nil {
			if errors.As(err, &validation.Errors{}) {
				return w.SettingsExarotonPage(&dto.SettingsExarotonPageData{
					HttpCode:   http.StatusBadRequest,
					Validation: err.(validation.Errors),
				})(c)
			}

			return err
		}

		acc, err := w.svc.ServerSettingsService.ValidateExarotonAPIKey(c.Request().Context(), req.APIKey)
		if err != nil {
			if errors.Is(err, errs.ErrUnauthorized) {
				return w.SettingsExarotonPage(&dto.SettingsExarotonPageData{
					HttpCode: http.StatusBadRequest,
					APIKey:   req.APIKey,
					Validation: map[string]error{
						"api_key": errors.New(constants.MsgErrInvalidAPIKey),
					},
				})(c)
			}
			return err
		}

		if err := w.svc.ServerSettingsService.UpdateExarotonAPIKey(c.Request().Context(), req.APIKey); err != nil {
			return err
		}

		return w.SettingsExarotonPage(&dto.SettingsExarotonPageData{
			APIKey:      req.APIKey,
			AccountInfo: acc,
		})(c)
	}
}

func (w *Web) SettingsExarotonValidateApiKey() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(dto.SettingsExarotonReq)
		if err := w.shouldBind(c, req); err != nil {
			if errors.As(err, &validation.Errors{}) {
				return w.SettingsExarotonPage(&dto.SettingsExarotonPageData{
					HttpCode:   http.StatusBadRequest,
					Validation: err.(validation.Errors),
				})(c)
			}

			return err
		}

		acc, err := w.svc.ServerSettingsService.ValidateExarotonAPIKey(c.Request().Context(), req.APIKey)
		if err != nil {
			return err
		}

		return w.SettingsExarotonPage(&dto.SettingsExarotonPageData{
			APIKey:      req.APIKey,
			AccountInfo: acc,
		})(c)
	}
}
