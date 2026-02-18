package wahandler

import (
	"errors"
	"exaroton-wa-bot/internal/config/warouter"
	"exaroton-wa-bot/internal/constants/errs"
	"exaroton-wa-bot/internal/dto"
	"exaroton-wa-bot/internal/helper"
	"log/slog"
)

func errHandler(c *warouter.Context, err error) {
	var resp dto.WhatsappMessage

	switch {
	case errors.Is(err, errs.ErrServerNotFound),
		errors.Is(err, errs.ErrCommandNotFound),
		errors.Is(err, errs.ErrServerIsAlreadyStopping),
		errors.Is(err, errs.ErrForbidden):
		resp.Conversation = helper.Ptr(err.Error())
	}

	if (resp != dto.WhatsappMessage{}) {
		_, _ = c.SendMessage(c, c.Chat, &resp)
		return
	}

	// debug
	slog.Warn("Unhandled error in wa handler", "error", err.Error())
}
