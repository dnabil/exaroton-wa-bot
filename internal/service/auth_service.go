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
	"go.mau.fi/whatsmeow/types"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Login(ctx context.Context, req *dto.UserLoginReq) (*dto.UserClaims, time.Duration, error)
	WhatsappLogin(ctx context.Context) (<-chan whatsmeow.QRChannelItem, error)
	WhatsappLogout(ctx context.Context) error
	WhatsappIsLoggedIn(ctx context.Context) bool
	GetWhatsappPhoneNumber() string
	GetWhatsappGroups(ctx context.Context) ([]*dto.WhatsappGroupInfo, error)

	// just filters based on whitelisted groups in db
	FilterWhatsappWhitelistedGroups(ctx context.Context, allGroups []*dto.WhatsappGroupInfo) ([]*dto.WhatsappGroupInfo, error)

	GetWhatsappWhitelistedGroupJIDs(ctx context.Context) ([]*dto.WhatsappWhitelistedGroup, error)
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

func (s *AuthService) GetWhatsappPhoneNumber() string {
	return s.waRepo.GetPhoneNumber()
}

func (s *AuthService) GetWhatsappGroups(ctx context.Context) ([]*dto.WhatsappGroupInfo, error) {
	groups, err := s.waRepo.GetGroups(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*dto.WhatsappGroupInfo, len(groups))
	for i, group := range groups {
		res[i] = dto.NewWhatsappGroupInfo(group)
	}

	return res, nil
}

func (s *AuthService) FilterWhatsappWhitelistedGroups(ctx context.Context, allGroups []*dto.WhatsappGroupInfo) ([]*dto.WhatsappGroupInfo, error) {
	tx := s.tx.Begin(ctx)
	defer func() {
		if rbErr := s.tx.Rollback(tx); rbErr != nil {
			slog.ErrorContext(ctx, rbErr.Error())
		}
	}()

	whitelistedJIDS, err := s.waRepo.GetWhitelistedGroupJIDs(ctx, tx)
	if err != nil {
		return nil, err
	}

	// create jidsMap based on whitelistedJIDs
	jidsMap := make(map[types.JID]bool)
	for _, jid := range whitelistedJIDS {
		jidsMap[types.NewJID(jid.JID, jid.ServerJID)] = true
	}

	// filter groups based on jids
	filteredGroups := make([]*dto.WhatsappGroupInfo, 0, len(allGroups))
	for _, group := range allGroups {
		if jidsMap[group.JID] {
			filteredGroups = append(filteredGroups, group)
		}
	}

	// TODO: sync the not existing jids (delete or update? not decided yet)

	return filteredGroups, nil
}

func (s *AuthService) WhatsappLogout(ctx context.Context) error {
	return s.waRepo.Logout(ctx)
}

func (s *AuthService) GetWhatsappWhitelistedGroupJIDs(ctx context.Context) ([]*dto.WhatsappWhitelistedGroup, error) {
	tx := s.tx.Begin(ctx)
	defer func() {
		if rbErr := s.tx.Rollback(tx); rbErr != nil {
			slog.ErrorContext(ctx, rbErr.Error())
		}
	}()

	entities, err := s.waRepo.GetWhitelistedGroupJIDs(ctx, tx)
	if err != nil {
		return nil, err
	}

	res := make([]*dto.WhatsappWhitelistedGroup, len(entities))
	for i, entity := range entities {
		res[i] = dto.NewWhatsappWhitelistedGroup(entity)
	}

	return res, nil
}
