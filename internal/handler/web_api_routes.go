package handler

import "github.com/labstack/echo/v4"

// shared api routes
var APIRoutes = new(struct {
	// TODO: remove all direct references to routes in DTOs (use the variables instead)

	WaLoginQRRoute     *echo.Route
	WaLoginNumberRoute *echo.Route
})

func (web *Web) LoadAPIRoutes() {
	apiGroup := web.Router.Group("/api")

	// middlewares
	apiGroup.Use(web.middleware.Session())
	apiGroup.Use(web.middleware.Logger())
	apiGroup.Use(web.middleware.Recover())

	// other middlewares
	authMdw := web.middleware.Auth()
	waAuthMdw := web.middleware.WhatsappLoggedIn(false)
	waGuestMdw := web.middleware.WhatsappLoggedIn(true)

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
			whatsappGroup.GET("/is-sync", web.APIWhatsappIsSync())
			whatsappGroup.POST("/logout", web.APIWhatsappLogout())
			whatsappGroup.GET("/groups", web.APIGetWhatsappGroups())
			whatsappGroup.POST("/groups/whitelist", web.APIWhatsappGroupWhitelist())
			whatsappGroup.DELETE("/groups/whitelist", web.APIWhatsappGroupUnwhitelist())
		}
	}

	// whatsapp login
	waLoginGroup := apiGroup.Group("/whatsapp/login", authMdw, waGuestMdw)
	{
		APIRoutes.WaLoginQRRoute = waLoginGroup.GET("/qr", web.APIWhatsappQRLogin())
		APIRoutes.WaLoginNumberRoute = waLoginGroup.GET("/number", nil)
	}
}
