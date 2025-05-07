package repository

import (
	"exaroton-wa-bot/internal/config"
	"log/slog"

	"gorm.io/gorm"
)

type Repo struct {
	WhatsappRepo IWhatsappRepo
	UserRepo     IUserRepo
}

func New(db *gorm.DB, waDB *config.WhatsappDB) (*Repo, error) {
	tmpl := &repoTmpl{db: db}

	whatsappRepo, err := newWhatsappRepo(waDB)
	if err != nil {
		slog.Error("failed to create whatsapp repo", "error", err)
		return nil, err
	}

	return &Repo{
		WhatsappRepo: whatsappRepo,
		UserRepo:     newUserRepo(tmpl),
	}, nil
}

// repoTmpl is a template for common fields for most repositories.
type repoTmpl struct {
	db *gorm.DB
}
