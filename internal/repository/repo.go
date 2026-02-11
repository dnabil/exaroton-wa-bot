package repository

import (
	"context"

	"gorm.io/gorm"
)

type Repo struct {
	WhatsappRepo       IWhatsappRepo
	UserRepo           IUserRepo
	ServerSettingsRepo IServerSettingsRepo
	ExarotonRepo       IExarotonRepo
}

func New(db *gorm.DB, waClient *waClient) (*Repo, error) {
	return &Repo{
		WhatsappRepo:       newWhatsappRepo(waClient),
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
