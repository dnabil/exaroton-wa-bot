package service

import (
	"exaroton-wa-bot/internal/config"
	"exaroton-wa-bot/internal/repository"

	"gorm.io/gorm"
)

type Service struct {
	AuthService           IAuthService
	ServerSettingsService IServerSettingsService
	WhatsappService       IWhatsappService
}

func New(cfg *config.Cfg, db *gorm.DB, repo *repository.Repo) *Service {
	svcTmpl := newSvcTmpl(cfg, db)

	// register services here...
	return &Service{
		AuthService:           NewAuthService(svcTmpl, repo.WhatsappRepo, repo.UserRepo),
		ServerSettingsService: NewServerSettingsService(svcTmpl, repo.ServerSettingsRepo, repo.ExarotonRepo),
		WhatsappService:       NewWhatsappService(svcTmpl, repo.WhatsappRepo),
	}
}

func newSvcTmpl(cfg *config.Cfg, db *gorm.DB) *svcTmpl {
	sqlTx := repository.NewSqlTx(db)

	return &svcTmpl{
		cfg: cfg,
		tx:  sqlTx,
	}
}

// svcTmpl is a template for services.
// It contains the common fields for all services.
// It is used to avoid code duplication.
type svcTmpl struct {
	cfg *config.Cfg

	// maindb
	//
	// every start on service func should use this to begin transaction when using
	// repository layer for maindb.
	tx repository.SqlTx
}
