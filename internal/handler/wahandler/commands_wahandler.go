package wahandler

import (
	"exaroton-wa-bot/internal/config/warouter"
	"exaroton-wa-bot/internal/constants/errs"
	"exaroton-wa-bot/internal/constants/messages"
	"exaroton-wa-bot/internal/dto"
	"exaroton-wa-bot/internal/helper"
	"exaroton-wa-bot/internal/service/command"
)

func (h *WaHandler) StartServer() warouter.HandlerFunc {
	return func(c *warouter.Context) error {
		startCmd, ok := h.cmdRegis.Get(command.StartServerCmdName)
		if !ok {
			return errs.ErrCommandNotFound
		}

		_, err := c.SendMessage(c, c.Chat, &dto.WhatsappMessage{
			Conversation: helper.Ptr(messages.ServerIsStarting),
		})
		if err != nil {
			return err
		}

		// iz blocking, will poll then return last server status.
		res := startCmd.Execute(c, c.Args)
		if res.Error != nil {
			return res.Error
		}

		_, err = c.SendMessage(c, c.Chat, &dto.WhatsappMessage{
			Conversation: &res.Text,
		})

		return err
	}
}

func (h *WaHandler) StopServer() warouter.HandlerFunc {
	return func(c *warouter.Context) error {
		stopCmd, ok := h.cmdRegis.Get(command.StopServerCmdName)
		if !ok {
			return errs.ErrCommandNotFound
		}

		res := stopCmd.Execute(c, c.Args)
		if res.Error != nil {
			return res.Error
		}

		_, err := c.SendMessage(c, c.Chat, &dto.WhatsappMessage{
			Conversation: &res.Text,
		})

		return err
	}
}

func (h *WaHandler) HelpCommand() warouter.HandlerFunc {
	return func(c *warouter.Context) error {
		helpCmd, ok := h.cmdRegis.Get(command.HelpCmdName)
		if !ok {
			return errs.ErrCommandNotFound
		}

		res := helpCmd.Execute(c, c.Args)
		if res.Error != nil {
			return res.Error
		}

		_, err := c.SendMessage(c, c.Chat, &dto.WhatsappMessage{
			Conversation: &res.Text,
		})

		return err
	}
}

func (h *WaHandler) ListServers() warouter.HandlerFunc {
	return func(c *warouter.Context) error {
		listCmd, ok := h.cmdRegis.Get(command.ListServerCmdName)
		if !ok {
			return errs.ErrCommandNotFound
		}

		res := listCmd.Execute(c, c.Args)
		if res.Error != nil {
			return res.Error
		}

		_, err := c.SendMessage(c, c.Chat, &dto.WhatsappMessage{
			Conversation: &res.Text,
		})

		return err
	}
}

func (h *WaHandler) ServerInfo() warouter.HandlerFunc {
	return func(c *warouter.Context) error {
		infoCmd, ok := h.cmdRegis.Get(command.InfoCmdName)
		if !ok {
			return errs.ErrCommandNotFound
		}

		res := infoCmd.Execute(c, c.Args)
		if res.Error != nil {
			return res.Error
		}

		_, err := c.SendMessage(c, c.Chat, &dto.WhatsappMessage{
			Conversation: &res.Text,
		})

		return err
	}
}

func (h *WaHandler) ListPlayers() warouter.HandlerFunc {
	return func(c *warouter.Context) error {
		return nil
	}
}
