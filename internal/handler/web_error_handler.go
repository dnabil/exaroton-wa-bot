package handler

import (
	"errors"
	"exaroton-wa-bot/internal/constants/errs"
	"log/slog"
	"strings"

	// exaroton-wa-bot/internal/constants/errs
	"exaroton-wa-bot/pages"
	"net/http"

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

		httpErr, ok := err.(*echo.HTTPError)
		if !ok {
			httpErr = &echo.HTTPError{
				Code:     http.StatusInternalServerError,
				Message:  http.StatusText(http.StatusInternalServerError),
				Internal: err,
			}
		}

		// custom error check
		switch {
		// already logged in
		case errors.Is(err, errs.ErrUserAlreadyLoggedIn) || errors.Is(err, errs.ErrWAAlreadyLoggedIn):
			err2 = c.Redirect(http.StatusSeeOther, homepageRoute.Path)
			return

		// is not logged in
		case errors.Is(err, errs.ErrUserNotLoggedIn):
			err2 = c.Redirect(http.StatusSeeOther, loginPageRoute.Path)
			return

		// is not logged in (whatsapp)
		case errors.Is(err, errs.ErrWANotLoggedIn):
			err2 = c.Redirect(http.StatusSeeOther, waLoginPageRoute.Path)
			return

		case errors.Is(err, errs.ErrLoginFailed):
			err2 = c.Redirect(http.StatusSeeOther, loginPageRoute.Path)
			return
		}

		// end of custom error check

		_ = c.Render(httpErr.Code, pages.Error, echo.Map{
			"Code":        httpErr.Code,
			"Message":     httpErr.Message,
			"Description": "Something went wrong :(",
		})
	}
}
