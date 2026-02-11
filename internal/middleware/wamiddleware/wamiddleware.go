package wamiddleware

import (
	"exaroton-wa-bot/internal/config"
	"exaroton-wa-bot/internal/config/warouter"
	"exaroton-wa-bot/internal/constants/errs"
	"exaroton-wa-bot/internal/service"
)

type Middleware struct {
	cfg               *config.Cfg
	authSvc           service.IAuthService
	serverSettingsSvc service.IServerSettingsService
}

func NewMiddleware(
	cfg *config.Cfg,
	authSvc service.IAuthService,
	serverSettingsSvc service.IServerSettingsService,
) *Middleware {
	return &Middleware{
		cfg:               cfg,
		authSvc:           authSvc,
		serverSettingsSvc: serverSettingsSvc,
	}
}

// ===================================================
// Middlewares
// ===================================================

// - Valid Exaroton API Key middleware: ValidExarotonAPIKeyMiddleware

// WhitelistedWAGroup returns a middleware that checks if the user (group) is in a whitelisted.
func (m *Middleware) WhitelistedWAGroup() warouter.MiddlewareFunc {
	return func(next warouter.HandlerFunc) warouter.HandlerFunc {
		return func(c *warouter.Context) error {
			whitelisted := false
			whitelistedGroups, err := m.authSvc.GetWhatsappWhitelistedGroupJIDs(c.Context)
			if err != nil {
				return err
			}

			for _, g := range whitelistedGroups {
				if g.UserJID == c.ChatUserJID && g.ServerJID == c.CharServerJID {
					whitelisted = true
					break
				}
			}

			if whitelisted {
				return next(c)
			}

			return errs.ErrWAGroupNotWhitelisted
		}
	}
}

// // must start tagging the bot's number and commands with a "/" prefix
// func (m *Middleware) CommandPrefix() warouter.MiddlewareFunc {
// 	return func(next warouter.HandlerFunc) warouter.HandlerFunc {
// 		return func(c *warouter.Context, args []string) error {
// 			if len(args) < 2 {
// 				return errs.ErrInvalidCommandPrefix
// 			}

// 			phoneNumber, err := m.authSvc.GetWhatsappPhoneNumber(c)
// 			if err != nil {
// 				return err
// 			}

// 			if !strings.Contains(args[0], phoneNumber) {
// 				return errs.ErrInvalidCommandPrefix
// 			}

// 			if !strings.HasPrefix(args[1], "/") {
// 				return errs.ErrInvalidCommandPrefix
// 			}

// 			return next(c, args)
// 		}
// 	}
// }

func (m *Middleware) ValidExarotonAPIKey() warouter.MiddlewareFunc {
	return func(next warouter.HandlerFunc) warouter.HandlerFunc {
		return func(c *warouter.Context) error {
			apiKey, err := m.serverSettingsSvc.GetExarotonAPIKey(c)
			if err != nil {
				return err
			}

			if apiKey == "" {
				return errs.ErrGSEmptyAPIKey
			}

			accInfo, err := m.serverSettingsSvc.ValidateExarotonAPIKey(c, apiKey)
			if err != nil || accInfo == nil {
				return err
			}

			return next(c)
		}
	}
}
