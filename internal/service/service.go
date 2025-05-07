package service

import (
	"exaroton-wa-bot/internal/config"
	"exaroton-wa-bot/internal/repository"

	"gorm.io/gorm"
)

type Service struct {
	AuthService IAuthService
}

func New(cfg *config.Cfg, db *gorm.DB, repo *repository.Repo) *Service {
	svcTmpl := newSvcTmpl(cfg, db)

	// register services here...
	return &Service{
		AuthService: NewAuthService(svcTmpl, repo.WhatsappRepo, repo.UserRepo),
	}
}

func newSvcTmpl(cfg *config.Cfg, db *gorm.DB) *svcTmpl {
	return &svcTmpl{
		cfg: cfg,
		db:  db,
	}
}

// svcTmpl is a template for services.
// It contains the common fields for all services.
// It is used to avoid code duplication.
type svcTmpl struct {
	cfg *config.Cfg

	// use for transactions only.
	// for testing, mock the underlying sql.DB
	db *gorm.DB
}
