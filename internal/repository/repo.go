package repository

import (
	"context"
	"exaroton-wa-bot/internal/config"
	"log/slog"

	"gorm.io/gorm"
)

type Repo struct {
	WhatsappRepo       IWhatsappRepo
	UserRepo           IUserRepo
	ServerSettingsRepo IServerSettingsRepo
	ExarotonRepo       IExarotonRepo
}

func New(db *gorm.DB, waDB *config.WhatsappDB) (*Repo, error) {
	whatsappRepo, err := newWhatsappRepo(waDB)
	if err != nil {
		slog.Error("failed to create whatsapp repo", "error", err)
		return nil, err
	}

	return &Repo{
		WhatsappRepo:       whatsappRepo,
		UserRepo:           newUserRepo(),
		ServerSettingsRepo: newServerSettingsRepo(),
		ExarotonRepo:       newExarotonRepo(),
	}, nil
}

// before using repo functions, should be started using this.
type SqlTx interface {
	Begin(ctx context.Context) *gorm.DB
	Commit(tx *gorm.DB) error
	Rollback(tx *gorm.DB) error
}

func NewSqlTx(db *gorm.DB) SqlTx {
	return &sqlTx{db: db}
}

type sqlTx struct {
	db *gorm.DB
}

func (s *sqlTx) Begin(ctx context.Context) *gorm.DB {
	return s.db.Begin()
}

func (s *sqlTx) Commit(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (s *sqlTx) Rollback(tx *gorm.DB) error {
	return tx.Rollback().Error
}
