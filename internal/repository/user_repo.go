package repository

import (
	"context"
	"errors"
	"exaroton-wa-bot/internal/database/entity"

	"gorm.io/gorm"
)

type IUserRepo interface {
	GetUserByUsername(ctx context.Context, tx *gorm.DB, username string) (*entity.User, error)
}

type UserRepo struct{}

func newUserRepo() IUserRepo {
	return &UserRepo{}
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, tx *gorm.DB, username string) (*entity.User, error) {
	user := &entity.User{}

	if err := tx.Where("username = ?", username).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
