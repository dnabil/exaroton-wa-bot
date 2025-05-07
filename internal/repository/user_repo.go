package repository

import (
	"context"
	"errors"
	"exaroton-wa-bot/internal/database/entity"

	"gorm.io/gorm"
)

type IUserRepo interface {
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
}

type UserRepo struct {
	*repoTmpl
}

func newUserRepo(tmpl *repoTmpl) *UserRepo {
	return &UserRepo{repoTmpl: tmpl}
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	user := &entity.User{}

	if err := r.db.Where("username = ?", username).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
