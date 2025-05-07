package service

import (
	"context"
	"exaroton-wa-bot/internal/config"
	"exaroton-wa-bot/internal/constants/errs"
	"exaroton-wa-bot/internal/database/entity"
	"exaroton-wa-bot/internal/dto"
	mockRepo "exaroton-wa-bot/internal/mocks/repository"
	"testing"
	"time"

	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mau.fi/whatsmeow"
	"golang.org/x/crypto/bcrypt"
)

func setupTestAuthService(t *testing.T) (IAuthService, *mockRepo.MockIUserRepo, *mockRepo.MockIWhatsappRepo, *config.Cfg) {
	mockUserRepo := mockRepo.NewMockIUserRepo(t)
	mockWaRepo := mockRepo.NewMockIWhatsappRepo(t)

	cfg := &config.Cfg{
		Koanf: koanf.New("."),
	}
	cfg.Set(config.KeyAuthDuration, "24h")

	svcTmpl := &svcTmpl{
		cfg: cfg,
	}

	authSvc := NewAuthService(svcTmpl, mockWaRepo, mockUserRepo)

	return authSvc, mockUserRepo, mockWaRepo, cfg
}

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name          string
		req           *dto.UserLoginReq
		mockSetup     func(*mockRepo.MockIUserRepo)
		expectedUser  *dto.UserClaims
		expectedError error
	}{
		{
			name: "success_login",
			req: &dto.UserLoginReq{
				Username: "testuser",
				Password: "password123",
			},
			mockSetup: func(mockUserRepo *mockRepo.MockIUserRepo) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				mockUserRepo.EXPECT().
					GetUserByUsername(mock.Anything, "testuser").
					Return(&entity.User{
						ID:       1,
						Username: "testuser",
						Password: string(hashedPassword),
					}, nil)
			},
			expectedUser: &dto.UserClaims{
				ID:       1,
				Username: "testuser",
			},
			expectedError: nil,
		},
		{
			name: "user_not_found",
			req: &dto.UserLoginReq{
				Username: "nonexistent",
				Password: "password123",
			},
			mockSetup: func(mockUserRepo *mockRepo.MockIUserRepo) {
				mockUserRepo.EXPECT().
					GetUserByUsername(mock.Anything, "nonexistent").
					Return(nil, nil)
			},
			expectedUser:  nil,
			expectedError: errs.ErrLoginFailed,
		},
		{
			name: "wrong_password",
			req: &dto.UserLoginReq{
				Username: "testuser",
				Password: "wrongpassword",
			},
			mockSetup: func(mockUserRepo *mockRepo.MockIUserRepo) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				mockUserRepo.EXPECT().
					GetUserByUsername(mock.Anything, "testuser").
					Return(&entity.User{
						ID:       1,
						Username: "testuser",
						Password: string(hashedPassword),
					}, nil)
			},
			expectedUser:  nil,
			expectedError: errs.ErrLoginFailed,
		},
		{
			name: "get_user_error",
			req: &dto.UserLoginReq{
				Username: "testuser",
				Password: "password123",
			},
			mockSetup: func(mockUserRepo *mockRepo.MockIUserRepo) {
				mockUserRepo.EXPECT().
					GetUserByUsername(mock.Anything, "testuser").
					Return(nil, assert.AnError)
			},
			expectedUser:  nil,
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authSvc, mockUserRepo, _, _ := setupTestAuthService(t)

			if tt.mockSetup != nil {
				tt.mockSetup(mockUserRepo)
			}

			user, expDuration, err := authSvc.Login(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
				assert.Equal(t, 24*time.Hour, expDuration)
			}
		})
	}
}

func TestAuthService_Login_DefaultDuration(t *testing.T) {
	// Create mocks
	mockUserRepo := mockRepo.NewMockIUserRepo(t)
	mockWaRepo := mockRepo.NewMockIWhatsappRepo(t)

	// Create config without auth duration
	cfg := &config.Cfg{
		Koanf: koanf.New("."),
	}

	// Create service template
	svcTmpl := &svcTmpl{
		cfg: cfg,
	}

	// Create auth service
	authSvc := NewAuthService(svcTmpl, mockWaRepo, mockUserRepo)

	// Setup mock for successful login
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUserRepo.EXPECT().
		GetUserByUsername(mock.Anything, "testuser").
		Return(&entity.User{
			ID:       1,
			Username: "testuser",
			Password: string(hashedPassword),
		}, nil)

	// Test login with default duration
	user, expDuration, err := authSvc.Login(context.Background(), &dto.UserLoginReq{
		Username: "testuser",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.Equal(t, &dto.UserClaims{
		ID:       1,
		Username: "testuser",
	}, user)
	assert.Equal(t, 12*time.Hour, expDuration)
}

func TestAuthService_WhatsappLogin(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*mockRepo.MockIWhatsappRepo)
		expectedError error
	}{
		{
			name: "success_login",
			mockSetup: func(mockWaRepo *mockRepo.MockIWhatsappRepo) {
				qrChan := make(chan whatsmeow.QRChannelItem)
				mockWaRepo.EXPECT().
					Login(mock.Anything).
					Return((<-chan whatsmeow.QRChannelItem)(qrChan), nil)
			},
			expectedError: nil,
		},
		{
			name: "already_logged_in",
			mockSetup: func(mockWaRepo *mockRepo.MockIWhatsappRepo) {
				mockWaRepo.EXPECT().
					Login(mock.Anything).
					Return(nil, errs.ErrWAAlreadyLoggedIn)
			},
			expectedError: errs.ErrWAAlreadyLoggedIn,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authSvc, _, mockWaRepo, _ := setupTestAuthService(t)

			if tt.mockSetup != nil {
				tt.mockSetup(mockWaRepo)
			}

			qrChan, err := authSvc.WhatsappLogin(context.Background())

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, qrChan)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, qrChan)
			}
		})
	}
}

func TestAuthService_WhatsappIsLoggedIn(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*mockRepo.MockIWhatsappRepo)
		expectedValue bool
	}{
		{
			name: "is_logged_in",
			mockSetup: func(mockWaRepo *mockRepo.MockIWhatsappRepo) {
				mockWaRepo.EXPECT().
					IsLoggedIn().
					Return(true)
			},
			expectedValue: true,
		},
		{
			name: "is_not_logged_in",
			mockSetup: func(mockWaRepo *mockRepo.MockIWhatsappRepo) {
				mockWaRepo.EXPECT().
					IsLoggedIn().
					Return(false)
			},
			expectedValue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authSvc, _, mockWaRepo, _ := setupTestAuthService(t)

			if tt.mockSetup != nil {
				tt.mockSetup(mockWaRepo)
			}

			isLoggedIn := authSvc.WhatsappIsLoggedIn(context.Background())
			assert.Equal(t, tt.expectedValue, isLoggedIn)
		})
	}
}
