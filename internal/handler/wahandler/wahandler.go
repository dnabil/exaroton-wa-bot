package wahandler

import (
	"exaroton-wa-bot/internal/config"
	"exaroton-wa-bot/internal/config/warouter"
	"exaroton-wa-bot/internal/middleware/wamiddleware"
	"exaroton-wa-bot/internal/service"
)

type WaHandler struct {
	router            *warouter.Router
	cfg               *config.Cfg
	mdw               *wamiddleware.Middleware
	authSvc           service.IAuthService
	serverSettingsSvc service.IServerSettingsService

	// event handler codes
	HandlerCodeCommandWA uint32
}

func NewWAHandler(cfg *config.Cfg, wa warouter.WhatsappService, authSvc service.IAuthService, serverSettingsSvc service.IServerSettingsService) *WaHandler {
	router := warouter.NewRouter(cfg, wa)
	router.ErrorHandlerFunc = errHandler

	h := &WaHandler{
		router:            router,
		cfg:               cfg,
		mdw:               wamiddleware.NewMiddleware(cfg, authSvc, serverSettingsSvc),
		authSvc:           authSvc,
		serverSettingsSvc: serverSettingsSvc,
	}

	h.LoadCommandRoutes()

	return h
}

func (h *WaHandler) Run() {
	h.router.Run()
}

func (h *WaHandler) Stop() {
	h.router.Stop()
}
