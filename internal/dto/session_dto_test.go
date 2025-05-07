package dto

import (
	"encoding/json"
	"exaroton-wa-bot/internal/constants"
	"exaroton-wa-bot/internal/constants/errs"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupTestContext() (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set up the session store
	store := sessions.NewCookieStore([]byte("test-secret"))
	e.Use(session.MiddlewareWithConfig(session.Config{
		Store: store,
	}))

	// Initialize the middleware for this request
	h := session.MiddlewareWithConfig(session.Config{
		Store: store,
	})(func(c echo.Context) error {
		return nil
	})
	h(c)

	return c, rec
}

func TestWebSession_GetUser(t *testing.T) {
	tests := []struct {
		name          string
		setupSession  func(echo.Context)
		expectedError error
		expectedUser  *UserClaims
	}{
		{
			name: "success_user_retrieval",
			setupSession: func(c echo.Context) {
				sess, _ := session.Get(constants.AuthCookieName, c)
				user := &UserClaims{
					ID:       1,
					Username: "testuser",
				}
				userJson, _ := json.Marshal(user)
				sess.Values[constants.AuthCookieKey] = userJson
				sess.Save(c.Request(), c.Response().Writer)
			},
			expectedError: nil,
			expectedUser: &UserClaims{
				ID:       1,
				Username: "testuser",
			},
		},
		{
			name: "no_session",
			setupSession: func(c echo.Context) {
				// No setup needed
			},
			expectedError: errs.ErrUserNotLoggedIn,
			expectedUser:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := setupTestContext()

			if tt.setupSession != nil {
				tt.setupSession(c)
			}

			ws := NewWebSession()

			user, err := ws.GetUser(c)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}
		})
	}
}

func TestWebSession_SetUser(t *testing.T) {
	tests := []struct {
		name          string
		user          *UserClaims
		expDuration   time.Duration
		expectedError error
	}{
		{
			name: "success",
			user: &UserClaims{
				ID:       1,
				Username: "testuser",
			},
			expDuration:   24 * time.Hour,
			expectedError: nil,
		},
		{
			name:          "nil_user",
			user:          nil,
			expDuration:   24 * time.Hour,
			expectedError: nil, // json.Marshal will handle nil gracefully
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := setupTestContext()

			ws := NewWebSession()

			err := ws.SetUser(c, tt.user, tt.expDuration)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
