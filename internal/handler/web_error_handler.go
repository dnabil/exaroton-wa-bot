package handler

import (
	"errors"
	"exaroton-wa-bot/internal/constants/errs"
	"exaroton-wa-bot/internal/dto"
	"log/slog"
	"strings"

	// exaroton-wa-bot/internal/constants/errs
	"exaroton-wa-bot/pages"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

// error handler implementation
func errorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if err == nil || c.Response().Committed {
			return
		}

		// ensure any errors in this function is logged
		var err2 error
		defer func() {
			if err2 != nil {
				slog.ErrorContext(c.Request().Context(), "error in error handler", "error", err2)
			}
		}()

		// web socket errors should be handled by the handler itself
		// so no need to handle them here
		if strings.HasPrefix(c.Request().Header.Get("Upgrade"), "websocket") {
			slog.Error("web socket error", "error", err)
			return
		}

		if strings.HasPrefix(c.Path(), "/api") {
			err2 = apiErrorHandler(err, c)
			return
		}

		err2 = webErrorHandler(err, c)
	}
}

func webErrorHandler(err error, c echo.Context) error {
	httpErr, ok := err.(*echo.HTTPError)
	if !ok {
		httpErr = &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			// Internal: err,
		}
	}

	// custom error check
	switch {
	// already logged in
	case errors.Is(err, errs.ErrUserAlreadyLoggedIn) || errors.Is(err, errs.ErrWAAlreadyLoggedIn):
		return c.Redirect(http.StatusSeeOther, homepageRoute.Path)

	// is not logged in
	case errors.Is(err, errs.ErrUserNotLoggedIn):
		return c.Redirect(http.StatusSeeOther, loginPageRoute.Path)
	case errors.Is(err, errs.ErrWANotLoggedIn):
		return c.Redirect(http.StatusSeeOther, waLoginPageRoute.Path)

	case errors.Is(err, errs.ErrLoginFailed):
		return c.Redirect(http.StatusSeeOther, loginPageRoute.Path)
	}

	// end of custom error check

	return c.Render(httpErr.Code, pages.Error, echo.Map{
		"Code":        httpErr.Code,
		"Message":     httpErr.Message,
		"Description": "Something went wrong :(",
	})
}
func apiErrorHandler(err error, c echo.Context) error {
	httpErr, ok := err.(*echo.HTTPError)
	if !ok {
		httpErr = &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			// Internal: err,
		}
	}

	// optional
	var data any = nil

	// custom error check
	switch {
	case errors.As(err, &validation.Errors{}):
		httpErr.Code = http.StatusBadRequest
		httpErr.Message = http.StatusText(http.StatusBadRequest)
		data = err.(validation.Errors)

	case errors.Is(err, errs.ErrUnauthorized):
		httpErr.Code, httpErr.Message = http.StatusUnauthorized, errs.ErrUnauthorized.Error()
	case errors.Is(err, errs.ErrForbidden):
		httpErr.Code, httpErr.Message = http.StatusForbidden, errs.ErrForbidden.Error()

	// already logged in
	case errors.Is(err, errs.ErrUserAlreadyLoggedIn):
		httpErr.Code, httpErr.Message = http.StatusForbidden, errs.ErrUserAlreadyLoggedIn.Error()
	case errors.Is(err, errs.ErrWAAlreadyLoggedIn):
		httpErr.Code, httpErr.Message = http.StatusForbidden, errs.ErrWAAlreadyLoggedIn.Error()

	// is not logged in
	case errors.Is(err, errs.ErrUserNotLoggedIn):
		httpErr.Code, httpErr.Message = http.StatusUnauthorized, errs.ErrUserNotLoggedIn.Error()
	case errors.Is(err, errs.ErrWANotLoggedIn):
		httpErr.Code, httpErr.Message = http.StatusUnauthorized, errs.ErrWANotLoggedIn.Error()

	case errors.Is(err, errs.ErrLoginFailed):
		httpErr.Code, httpErr.Message = http.StatusUnauthorized, errs.ErrLoginFailed.Error()

	// game server related errors
	case errors.Is(err, errs.ErrGSInvalidAPIKey):
		httpErr.Code, httpErr.Message = http.StatusUnauthorized, errs.ErrGSInvalidAPIKey.Error()
	}

	// end of custom error check

	return c.JSON(httpErr.Code, dto.APIResponse{
		Message: httpErr.Message.(string),
		Data:    data,
	})
}
