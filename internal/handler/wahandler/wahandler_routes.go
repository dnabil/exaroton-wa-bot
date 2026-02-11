package wahandler

func (h *WaHandler) LoadCommandRoutes() {
	router := h.router
	mdw := h.mdw

	// middlewares
	router.Use(mdw.ValidExarotonAPIKey())
	router.Use(mdw.WhitelistedWAGroup())

	router.Register("/help", h.HelpCommand())    // shows the manual page/guide thru WhatsApp chat for commands available
	router.Register("/servers", h.ListServers()) // shows available server ids
	router.Register("/start", h.StartServer())   //  [server-id] starts the server specified by its id
	router.Register("/stop", h.StopServer())     // [server-id] stops the server specified by its id
	router.Register("/info", h.ServerInfo())     // [server-id] shows the current server info
	router.Register("/players", h.ListPlayers()) // [server-id] shows the players that are currently online on a server
}
