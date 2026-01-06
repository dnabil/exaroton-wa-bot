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
	"github.com/stretchr/testify/require"
	"go.mau.fi/whatsmeow"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupTestAuthService(t *testing.T) (
	IAuthService,
	*mockRepo.MockSqlTx,
	*mockRepo.MockIUserRepo,
	*mockRepo.MockIWhatsappRepo,
	*config.Cfg,
) {
	mockUserRepo := mockRepo.NewMockIUserRepo(t)
	mockWaRepo := mockRepo.NewMockIWhatsappRepo(t)
	mockSqlTx := mockRepo.NewMockSqlTx(t)

	cfg := &config.Cfg{
		Koanf: koanf.New("."),
	}
	require.NoError(t, cfg.Set(config.KeyAuthDuration, "24h"))

	svcTmpl := &svcTmpl{
		cfg: cfg,
		tx:  mockSqlTx,
	}

	authSvc := NewAuthService(svcTmpl, mockWaRepo, mockUserRepo)

	return authSvc, mockSqlTx, mockUserRepo, mockWaRepo, cfg
}

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name          string
		req           *dto.UserLoginReq
		mockSetup     func(*mockRepo.MockSqlTx, *mockRepo.MockIUserRepo)
		expectedUser  *dto.UserClaims
		expectedError error
	}{
		{
			name: "success_login",
			req: &dto.UserLoginReq{
				Username: "testuser",
				Password: "password123",
			},
			mockSetup: func(mockSqlTx *mockRepo.MockSqlTx, mockUserRepo *mockRepo.MockIUserRepo) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				// transaction
				mockSqlTx.EXPECT().Begin(mock.Anything).Return(new(gorm.DB))
				mockSqlTx.EXPECT().Rollback(mock.Anything).Return(nil)

				mockUserRepo.EXPECT().
					GetUserByUsername(mock.Anything, mock.Anything, "testuser").
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
			mockSetup: func(mockSqlTx *mockRepo.MockSqlTx, mockUserRepo *mockRepo.MockIUserRepo) {
				// transaction
				mockSqlTx.EXPECT().Begin(mock.Anything).Return(new(gorm.DB))
				mockSqlTx.EXPECT().Rollback(mock.Anything).Return(nil)

				mockUserRepo.EXPECT().
					GetUserByUsername(mock.Anything, mock.Anything, "nonexistent").
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
			mockSetup: func(mockSqlTx *mockRepo.MockSqlTx, mockUserRepo *mockRepo.MockIUserRepo) {
				// transaction
				mockSqlTx.EXPECT().Begin(mock.Anything).Return(new(gorm.DB))
				mockSqlTx.EXPECT().Rollback(mock.Anything).Return(nil)

				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				mockUserRepo.EXPECT().
					GetUserByUsername(mock.Anything, mock.Anything, "testuser").
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
			mockSetup: func(mockSqlTx *mockRepo.MockSqlTx, mockUserRepo *mockRepo.MockIUserRepo) {
				// transaction
				mockSqlTx.EXPECT().Begin(mock.Anything).Return(new(gorm.DB))
				mockSqlTx.EXPECT().Rollback(mock.Anything).Return(nil)

				mockUserRepo.EXPECT().
					GetUserByUsername(mock.Anything, mock.Anything, "testuser").
					Return(nil, assert.AnError)
			},
			expectedUser:  nil,
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authSvc, mockSqlTx, mockUserRepo, _, _ := setupTestAuthService(t)

			if tt.mockSetup != nil {
				tt.mockSetup(mockSqlTx, mockUserRepo)
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
	mockSqlTx := mockRepo.NewMockSqlTx(t)

	// Create config without auth duration
	cfg := &config.Cfg{
		Koanf: koanf.New("."),
	}

	// Create service template
	svcTmpl := &svcTmpl{
		cfg: cfg,
		tx:  mockSqlTx,
	}

	// Create auth service
	authSvc := NewAuthService(svcTmpl, mockWaRepo, mockUserRepo)

	// Setup mock for successful login
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 1)

	// transaction mock
	mockSqlTx.EXPECT().Begin(mock.Anything).Return(new(gorm.DB))
	mockSqlTx.EXPECT().Rollback(mock.Anything).Return(nil)

	mockUserRepo.EXPECT().
		GetUserByUsername(mock.Anything, mock.Anything, "testuser").
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
			authSvc, _, _, mockWaRepo, _ := setupTestAuthService(t)

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
			authSvc, _, _, mockWaRepo, _ := setupTestAuthService(t)

			if tt.mockSetup != nil {
				tt.mockSetup(mockWaRepo)
			}

			isLoggedIn := authSvc.WhatsappIsLoggedIn(context.Background())
			assert.Equal(t, tt.expectedValue, isLoggedIn)
		})
	}
}
