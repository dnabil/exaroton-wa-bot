package wahandler

import (
	"errors"
	"exaroton-wa-bot/internal/config/warouter"
	"exaroton-wa-bot/internal/constants/errs"
	"exaroton-wa-bot/internal/dto"
	"exaroton-wa-bot/internal/helper"
)

func errHandler(c *warouter.Context, err error) {
	var resp dto.WhatsappMessage

	switch {
	case errors.Is(err, errs.ErrServerNotFound):
		resp.Conversation = helper.Ptr(err.Error())
	}

	if (resp != dto.WhatsappMessage{}) {
		c.SendMessage(c, c.Chat, &resp)
	}
}
