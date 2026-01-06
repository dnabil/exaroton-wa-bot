package dto

import (
	"encoding/json"
	"exaroton-wa-bot/internal/constants"
	"exaroton-wa-bot/internal/constants/errs"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type WebSession interface {
	GetUser(c echo.Context) (*UserClaims, error)
	SetUser(c echo.Context, user *UserClaims, expDuration time.Duration) error

	// Get and clear flash message
	GetFlash(c echo.Context) ([]WebFlashMessage, error)
	//	Set flash message
	SetFlash(c echo.Context, flashMessage WebFlashMessage) error
	// Get and clear validation error
	GetValidationError(c echo.Context) (WebValidationErrors, error)
	// Set validation error
	SetValidationError(c echo.Context, valErr WebValidationErrors) error
	// Get and clear old input
	GetOldInput(c echo.Context) (WebOldInput, error)
	// Set old input
	SetOldInput(c echo.Context, oldInput WebOldInput) error
}

type webSession struct{}

func NewWebSession() WebSession {
	return &webSession{}
}

func (s *webSession) GetUser(c echo.Context) (*UserClaims, error) {
	sess, err := session.Get(constants.AuthCookieName, c)
	if err != nil {
		return nil, errs.ErrUserNotLoggedIn
	}

	val, ok := sess.Values[constants.AuthCookieKey].([]byte)
	if !ok {
		return nil, errs.ErrUserNotLoggedIn
	}

	userClaims := new(UserClaims)
	err = json.Unmarshal(val, userClaims)
	if err != nil {
		return nil, errs.ErrUserNotLoggedIn
	}

	return userClaims, nil
}

func (s *webSession) SetUser(c echo.Context, user *UserClaims, expDuration time.Duration) error {
	sess, err := session.Get(constants.AuthCookieName, c)
	if err != nil {
		return err
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}

	sess.Values[constants.AuthCookieKey] = userJson
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int(expDuration.Seconds()),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}

	err = sess.Save(c.Request(), c.Response().Writer)
	if err != nil {
		return fmt.Errorf("failed to save user session: %w", err)
	}

	return nil
}

// ==========================================
// Session messages

// ===============================
// Flash message types
type WebFlashMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type WebValidationErrors map[string]string

type WebOldInput map[string]string

// ===============================

var (
	sessionBaseName     = "session"
	sessionFlashName    = "_flash"
	sessionValErrName   = "_val_err"
	sessionOldInputName = "_old_input"
)

func (s *webSession) GetFlash(c echo.Context) ([]WebFlashMessage, error) {
	sess, err := session.Get(sessionBaseName, c)
	if err != nil {
		return nil, err
	}

	flashes := sess.Flashes(sessionFlashName) // alr delete after getting
	res := make([]WebFlashMessage, len(flashes))

	for i, f := range flashes {
		if flashMessage, ok := f.(WebFlashMessage); ok {
			res[i] = flashMessage
		}
	}

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *webSession) SetFlash(c echo.Context, flashMessage WebFlashMessage) error {
	sess, err := session.Get(sessionBaseName, c)
	if err != nil {
		return err
	}

	sess.AddFlash(flashMessage, sessionFlashName)

	return sess.Save(c.Request(), c.Response())
}

func (s *webSession) GetValidationError(c echo.Context) (WebValidationErrors, error) {
	sess, err := session.Get(sessionBaseName, c)
	if err != nil {
		return nil, err
	}

	flashes := sess.Flashes(sessionValErrName) // alr delete after getting
	res := make(WebValidationErrors)

	for _, f := range flashes {
		if valErr, ok := f.(WebValidationErrors); ok {
			for k, v := range valErr {
				res[k] = v
			}
		}
	}

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *webSession) SetValidationError(c echo.Context, valErr WebValidationErrors) error {
	sess, err := session.Get(sessionBaseName, c)
	if err != nil {
		return err
	}

	sess.AddFlash(valErr, sessionValErrName)

	return sess.Save(c.Request(), c.Response())
}

func (s *webSession) GetOldInput(c echo.Context) (WebOldInput, error) {
	sess, err := session.Get(sessionBaseName, c)
	if err != nil {
		return nil, err
	}

	flashes := sess.Flashes(sessionOldInputName) // alr delete after getting
	res := make(WebOldInput)

	for _, f := range flashes {
		if oldInput, ok := f.(WebOldInput); ok {
			for k, v := range oldInput {
				res[k] = v
			}
		}
	}

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *webSession) SetOldInput(c echo.Context, oldInput WebOldInput) error {
	sess, err := session.Get(sessionBaseName, c)
	if err != nil {
		return err
	}

	sess.AddFlash(oldInput, sessionOldInputName)

	return sess.Save(c.Request(), c.Response())
}

// ==========================================
