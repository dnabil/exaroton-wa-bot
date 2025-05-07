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
