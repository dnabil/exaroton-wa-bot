package dto

import (
	"encoding/gob"
	"encoding/json"
	"errors"
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
	"github.com/stretchr/testify/require"
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
	_ = h(c)

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
				require.NoError(t, sess.Save(c.Request(), c.Response().Writer))
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

func TestWeSession_GetSetFlash(t *testing.T) {
	gob.Register(WebFlashMessage{})

	tests := []struct {
		name         string
		setFlashMsgs WebFlashMessage
	}{
		{
			name: "success",
			setFlashMsgs: WebFlashMessage{
				"success": "Test message",
			},
		},
		{
			name: "success_multiple",
			setFlashMsgs: WebFlashMessage{
				"info":    "First message",
				"warning": "Second message",
			},
		},
		{
			name:         "success_empty",
			setFlashMsgs: WebFlashMessage{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := setupTestContext()

			ws := NewWebSession()

			for key, msg := range tt.setFlashMsgs {
				err := ws.SetFlash(c, key, msg)
				require.NoError(t, err)
			}

			res, err := ws.GetFlash(c)
			require.NoError(t, err)

			// assert
			assert.Len(t, res, len(tt.setFlashMsgs))
			for i, msg := range tt.setFlashMsgs {
				assert.Equal(t, msg, res[i])
			}

			// check if the flash messages are deleted after getting them
			res, err = ws.GetFlash(c)
			require.NoError(t, err)
			assert.Len(t, res, 0)
		})
	}
}

func TestWeSession_GetSetValidationError(t *testing.T) {
	gob.Register(WebValidationErrors{})

	tests := []struct {
		name       string
		setValErrs []map[string]error
	}{
		{
			name: "success",
			setValErrs: []map[string]error{
				{
					"email":    errors.New("invalid email address"),
					"password": errors.New("Password must be at least 8 characters"),
				},
			},
		},
		{
			name: "success_override",
			setValErrs: []map[string]error{
				{
					"email":    errors.New("invalid email address"),
					"password": errors.New("Password must be at least 8 characters"),
				},
			},
		},
		{
			name: "success_override",
			setValErrs: []map[string]error{
				{
					"email":    errors.New("Invalid email address"),
					"password": errors.New("Password must be at least 8 characters"),
				},
				{
					"email":    errors.New("OVERRIDEN VALUE"),
					"password": errors.New("OVERRIDEN VALUE"),
				},
			},
		},
		{
			name:       "success_empty",
			setValErrs: []map[string]error{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := setupTestContext()

			ws := NewWebSession()

			for _, valErr := range tt.setValErrs {
				err := ws.SetValidationError(c, valErr)
				require.NoError(t, err)
			}

			res, err := ws.GetValidationError(c)
			require.NoError(t, err)

			// get combined errors from testcase
			combinedErrs := make(WebValidationErrors)
			for _, valErr := range tt.setValErrs {
				for k, v := range valErr {
					combinedErrs[k] = v.Error()
				}
			}

			assert.Equal(t, combinedErrs, res)

			// check if the flash messages are deleted after getting them
			res, err = ws.GetValidationError(c)
			require.NoError(t, err)
			assert.Len(t, res, 0)
		})
	}
}

type oldInpputStructTest struct {
	Email    string
	Password string
}

func (o oldInpputStructTest) ToMap() map[string]string {
	return map[string]string{
		"email":    o.Email,
		"password": o.Password,
	}
}

func TestWeSession_GetSetOldInput(t *testing.T) {
	gob.Register(WebOldInput{})

	tests := []struct {
		name        string
		setOldInput []oldInpputStructTest
	}{
		{
			name: "success",
			setOldInput: []oldInpputStructTest{
				{
					Email:    "asd@",
					Password: "Password",
				},
			},
		},
		{
			name: "success_override",
			setOldInput: []oldInpputStructTest{
				{
					Email:    "asd@",
					Password: "Password",
				},
				{
					Email:    "OVERRIDEN VALUE",
					Password: "OVERRIDEN VALUE",
				},
			},
		},
		{
			name:        "success_empty",
			setOldInput: []oldInpputStructTest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := setupTestContext()

			ws := NewWebSession()

			for _, oldInput := range tt.setOldInput {
				err := ws.SetOldInput(c, oldInput)
				require.NoError(t, err)
			}

			res, err := ws.GetOldInput(c)
			require.NoError(t, err)

			// get combined oldinput from testcase
			combinedErrs := make(WebOldInput)
			for _, valErr := range tt.setOldInput {
				valErrMap := valErr.ToMap()
				for k, v := range valErrMap {
					combinedErrs[k] = v
				}
			}

			assert.Equal(t, combinedErrs, res)

			// check if the flash messages are deleted after getting them
			res, err = ws.GetOldInput(c)
			require.NoError(t, err)
			assert.Len(t, res, 0)
		})
	}
}
