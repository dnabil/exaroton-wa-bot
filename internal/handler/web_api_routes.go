package handler

func (web *Web) LoadAPIRoutes() {
	apiGroup := web.Router.Group("/api")

	// middlewares
	apiGroup.Use(web.middleware.Session())
	apiGroup.Use(web.middleware.Logger())
	apiGroup.Use(web.middleware.Recover())

	// other middlewares
	authMdw := web.middleware.Auth()
	waAuthMdw := web.middleware.WhatsappLoggedIn(false)

	// settings
	settingsGroup := apiGroup.Group("/settings", authMdw, waAuthMdw)
	{
		// server settings ()
		serverGroup := settingsGroup.Group("/server")
		{
			serverGroup.POST("/exaroton/save", web.APISettingsExarotonUpdate())
			serverGroup.POST("/exaroton/validate", web.APISettingsExarotonValidateApiKey())
		}

		// whatsapp settings
		whatsappGroup := settingsGroup.Group("/whatsapp")
		{
			whatsappGroup.POST("/logout", web.APIWhatsappLogout())
			whatsappGroup.GET("/groups", web.APIGetWhatsappGroups())
			whatsappGroup.POST("/groups/whitelist", web.APIWhatsappGroupWhitelist())
		}
	}

}
