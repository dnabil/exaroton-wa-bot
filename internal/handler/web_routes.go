package handler

import (
	"exaroton-wa-bot/internal/config"
	"fmt"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

// shared web routes
var WebRoutes = new(struct {
	// TODO: remove all direct references to routes in DTOs (use the variables instead)

	HomepageRoute             *echo.Route
	LoginPageRoute            *echo.Route
	LoginRoute                *echo.Route
	WaLoginPageRoute          *echo.Route
	WaLoginPageQRRoute        *echo.Route
	WaLoginPageNumberRoute    *echo.Route
	WaLoginQRRoute            *echo.Route
	SettingsExarotonPageRoute *echo.Route
})

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
	webGroup.Use(web.middleware.FlashMessage())
	webGroup.Use(web.middleware.FlashOldInput())
	webGroup.Use(web.middleware.FlashValidationError())

	// other middlewares
	authMdw := web.middleware.Auth()
	guestMdw := web.middleware.Guest()
	waAuthMdw := web.middleware.WhatsappLoggedIn(false)
	waGuestMdw := web.middleware.WhatsappLoggedIn(true)

	// homepage route
	WebRoutes.HomepageRoute = webGroup.GET("/", web.HomePage(), authMdw, waAuthMdw)

	// user routes
	userGroup := webGroup.Group("/user")
	{
		WebRoutes.LoginPageRoute = userGroup.GET("/login", web.UserLoginPage(nil), guestMdw)
		WebRoutes.LoginRoute = userGroup.POST("/login", web.UserLogin(), guestMdw)
	}

	// whatsapp login routes
	waGroup := webGroup.Group("/whatsapp/login", authMdw, waGuestMdw)
	{
		WebRoutes.WaLoginPageRoute = waGroup.GET("/", web.WhatsappLoginPage())
		WebRoutes.WaLoginPageQRRoute = waGroup.GET("/qr", web.WhatsappLoginQRPage())
		WebRoutes.WaLoginPageNumberRoute = waGroup.GET("/number", nil) // TODO: implement
	}

	// settings
	settingsGroup := webGroup.Group("/settings", authMdw, waAuthMdw)
	{
		serverGroup := settingsGroup.Group("/server")
		{
			WebRoutes.SettingsExarotonPageRoute = serverGroup.GET("/exaroton", web.SettingsExarotonPage(nil))
		}

		// whatsapp settings
		whatsappGroup := settingsGroup.Group("/whatsapp")
		{
			whatsappGroup.GET("", web.SettingsWhatsappPage())
		}
	}
}
