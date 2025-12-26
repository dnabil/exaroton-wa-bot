package handler

import (
	"exaroton-wa-bot/internal/config"
	"fmt"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

// shared web routes
var (
	homepageRoute *echo.Route

	loginPageRoute *echo.Route
	loginRoute     *echo.Route

	waLoginPageRoute *echo.Route
	waLoginQRRoute   *echo.Route
)

func (web *Web) LoadRoutes() {
	// MAIN GROUP
	webGroup := web.Router.Group("")

	staticDir, err := filepath.Abs(web.cfg.MustString(config.KeyPublicDir))
	if err != nil {
		panic(fmt.Sprintf("failed to get public dir (key: %s)", config.KeyPublicDir))
	}

	// static files
	webGroup.Static("/public", staticDir)
	webGroup.File("/favicon.ico", filepath.Join(staticDir, "favicon.ico"))

	// middlewares
	webGroup.Use(web.middleware.Session())
	webGroup.Use(web.middleware.Logger())
	webGroup.Use(web.middleware.Recover())

	// other middlewares
	authMdw := web.middleware.Auth()
	guestMdw := web.middleware.Guest()
	waAuthMdw := web.middleware.WhatsappLoggedIn(false)
	waGuestMdw := web.middleware.WhatsappLoggedIn(true)

	// homepage route
	homepageRoute = webGroup.GET("/", web.HomePage(), authMdw, waAuthMdw)

	// user routes
	userGroup := webGroup.Group("/user")
	{
		loginPageRoute = userGroup.GET("/login", web.UserLoginPage(nil), guestMdw)
		loginRoute = userGroup.POST("/login", web.UserLogin(), guestMdw)
	}

	// whatsapp login routes
	waGroup := webGroup.Group("/wa")
	{
		waLoginPageRoute = waGroup.GET("/login", web.WhatsappLoginPage(), authMdw, waGuestMdw)
		waLoginQRRoute = waGroup.GET("/qr", web.WhatsappQRLogin(), authMdw, waGuestMdw)
	}

	// settings
	settingsGroup := webGroup.Group("/settings", authMdw, waAuthMdw)
	{
		// server settings ()
		serverGroup := settingsGroup.Group("/server")
		{
			serverGroup.GET("/exaroton", web.SettingsExarotonPage(nil))
			serverGroup.POST("/exaroton", web.SettingsExarotonUpdate())
			serverGroup.POST("/exaroton/me", web.SettingsExarotonValidateApiKey())
		}
	}
}
