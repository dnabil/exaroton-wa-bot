package handler

import (
	"errors"
	"exaroton-wa-bot/internal/constants"
	"exaroton-wa-bot/internal/constants/errs"
	"exaroton-wa-bot/internal/dto"
	"exaroton-wa-bot/pages"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"go.mau.fi/whatsmeow"
)

func (w *Web) UserLoginPage(data *dto.LoginPageData) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, pages.Login, data)
	}
}

func (w *Web) UserLogin() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(dto.UserLoginReq)
		if err := w.shouldBind(c, req); err != nil {
			if errors.As(err, &validation.Errors{}) {
				return w.UserLoginPage(&dto.LoginPageData{
					Validation: (err.(validation.Errors)),
				})(c)
			}

			return err
		}

		userClaims, expDuration, err := w.svc.AuthService.Login(c.Request().Context(), req)
		if err != nil {
			if errors.Is(err, errs.ErrLoginFailed) {
				return w.UserLoginPage(&dto.LoginPageData{
					Validation: map[string]error{
						"password": err,
					},
				})(c)
			}

			return err
		}

		if err := w.session.SetUser(c, userClaims, expDuration); err != nil {
			return err
		}

		return c.Redirect(http.StatusSeeOther, homepageRoute.Path)
	}
}

func (w *Web) WhatsappLoginPage() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, pages.WhatsappLogin, &dto.WhatsappLoginPageData{
			WSPath: waLoginQRRoute.Path,
		})
	}
}

func (w *Web) WhatsappQRLogin() echo.HandlerFunc {
	return func(c echo.Context) error {
		ws, err := wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()

		qrChan, err := w.svc.AuthService.WhatsappLogin(c.Request().Context())
		if err != nil {
			ws.WriteJSON(dto.WhatsappQRWSRes{
				Event: dto.WhatsappQREventError,
				Error: err.Error(),
			})
			return err
		}

		// if already logged in, send success event
		if qrChan == nil {
			return ws.WriteJSON(dto.WhatsappQRWSRes{
				Event:   dto.WhatsappQREventSuccess,
				Message: constants.MsgWALoginSuccess,
			})
		}

		for qr := range qrChan {
			var err error
			switch qr.Event {

			// keep sending qr codes
			case whatsmeow.QRChannelEventCode: // "code"
				err = ws.WriteJSON(dto.WhatsappQRWSRes{
					Event: qr.Event,
					Code:  qr.Code,
				})
				if err != nil {
					return err
				}
				continue

			// pairing error
			case whatsmeow.QRChannelEventError:
				err = ws.WriteJSON(dto.WhatsappQRWSRes{
					Event: dto.WhatsappQREventError,
					Error: errs.ErrWAQRError.Error(),
				})

			// login success
			case whatsmeow.QRChannelSuccess.Event:
				err = ws.WriteJSON(dto.WhatsappQRWSRes{
					Event:   dto.WhatsappQREventSuccess,
					Message: constants.MsgWALoginSuccess,
				})

			// pairing timeout
			case whatsmeow.QRChannelTimeout.Event:
				err = ws.WriteJSON(dto.WhatsappQRWSRes{
					Event: dto.WhatsappQREventTimeout,
					Error: errs.ErrWAQRTimeout.Error(),
				})

			// unexpected event, pairing already happened
			case whatsmeow.QRChannelErrUnexpectedEvent.Event:
				err = ws.WriteJSON(dto.WhatsappQRWSRes{
					Event: dto.WhatsappQREventIsPairing,
					Error: errs.ErrWAQRIsPairing.Error(),
				})

			// client outdated
			case whatsmeow.QRChannelClientOutdated.Event:
				err = ws.WriteJSON(dto.WhatsappQRWSRes{
					Event: dto.WhatsappQREventClientOutdated,
					Error: errs.ErrWAQRClientOutdated.Error(),
				})

			// multidevice not enabled
			case whatsmeow.QRChannelScannedWithoutMultidevice.Event:
				err = ws.WriteJSON(dto.WhatsappQRWSRes{
					Event: dto.WhatsappQREventMultideviceNotEnabled,
					Error: errs.ErrWAQREnableMultidevice.Error(),
				})
			}

			if err != nil {
				return err
			}

			return nil // quit (success/errors)

		} // end of qrChan loop

		return nil
	}
}
