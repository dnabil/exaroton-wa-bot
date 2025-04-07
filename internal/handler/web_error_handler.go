package handler

import (
	"errors"
	"exaroton-wa-bot/pages"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

func webErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
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

		switch {
		case errors.As(err, &validation.Errors{}):
			httpErr.Code = http.StatusBadRequest
			httpErr.Message = "validation fail"
		}

		// end of custom error check

		c.Render(httpErr.Code, pages.Error, echo.Map{
			"Code":        httpErr.Code,
			"Message":     httpErr.Message,
			"Description": "Something went wrong :(",
		})
	}
}
