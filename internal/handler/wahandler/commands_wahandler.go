package wahandler

import (
	"exaroton-wa-bot/internal/config/warouter"
)

func (h *WaHandler) StartServer() warouter.HandlerFunc {
	return func(c *warouter.Context) error {
		return nil
	}
}
func (h *WaHandler) StopServer() warouter.HandlerFunc {
	return func(c *warouter.Context) error {
		return nil
	}
}

func (h *WaHandler) HelpCommand() warouter.HandlerFunc {
	return func(c *warouter.Context) error {
		return nil
	}
}

func (h *WaHandler) ListServers() warouter.HandlerFunc {
	return func(c *warouter.Context) error {
		return nil
	}
}

func (h *WaHandler) ServerInfo() warouter.HandlerFunc {
	return func(c *warouter.Context) error {
		return nil
	}
}

func (h *WaHandler) ListPlayers() warouter.HandlerFunc {
	return func(c *warouter.Context) error {
		return nil
	}
}
