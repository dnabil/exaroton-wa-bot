package service

import (
	"context"
	"exaroton-wa-bot/internal/config"
	"exaroton-wa-bot/internal/constants/errs"
	"exaroton-wa-bot/internal/dto"
	"exaroton-wa-bot/internal/repository"
	"log/slog"
	"time"

	"go.mau.fi/whatsmeow"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Login(ctx context.Context, req *dto.UserLoginReq) (*dto.UserClaims, time.Duration, error)
	WhatsappLogin(ctx context.Context) (<-chan whatsmeow.QRChannelItem, error)
	WhatsappIsLoggedIn(ctx context.Context) bool
}

type AuthService struct {
	*svcTmpl
	waRepo   repository.IWhatsappRepo
	userRepo repository.IUserRepo
}

func NewAuthService(svcTmpl *svcTmpl, waRepo repository.IWhatsappRepo, userRepo repository.IUserRepo) IAuthService {
	return &AuthService{
		svcTmpl:  svcTmpl,
		waRepo:   waRepo,
		userRepo: userRepo,
	}
}

func (s *AuthService) Login(ctx context.Context, req *dto.UserLoginReq) (*dto.UserClaims, time.Duration, error) {
	tx := s.tx.Begin(ctx)
	defer func() {
		if rbErr := s.tx.Rollback(tx); rbErr != nil {
			slog.ErrorContext(ctx, rbErr.Error())
		}
	}()

	user, err := s.userRepo.GetUserByUsername(ctx, tx, req.Username)
	if err != nil {
		return nil, 0, err
	}

	if user == nil {
		return nil, 0, errs.ErrLoginFailed
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, 0, errs.ErrLoginFailed
	}

	expDuration := s.cfg.Duration(config.KeyAuthDuration)
	if expDuration == 0 {
		expDuration = 12 * time.Hour
		slog.WarnContext(ctx, errs.ErrWarnJwtDurationNotSet.Error())
	}

	return dto.NewUserClaims(user), expDuration, nil
}

func (s *AuthService) WhatsappLogin(ctx context.Context) (<-chan whatsmeow.QRChannelItem, error) {
	return s.waRepo.Login(ctx)
}

func (s *AuthService) WhatsappIsLoggedIn(ctx context.Context) bool {
	return s.waRepo.IsLoggedIn()
}
