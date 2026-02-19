package repository

import (
	"context"
	"errors"
	"exaroton-wa-bot/internal/database/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IServerSettingsRepo interface {
	Get(ctx context.Context, tx *gorm.DB, key string) (*entity.ServerSettings, error)
	Upsert(ctx context.Context, tx *gorm.DB, settings *entity.ServerSettings) error
}

type ServerSettingsRepo struct{}

func newServerSettingsRepo() IServerSettingsRepo {
	return &ServerSettingsRepo{}
}

func (r *ServerSettingsRepo) Get(ctx context.Context, tx *gorm.DB, key string) (*entity.ServerSettings, error) {
	settings := &entity.ServerSettings{}

	if err := tx.Where(&entity.ServerSettings{Key: key}).First(settings).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return settings, nil
}

func (r *ServerSettingsRepo) Upsert(ctx context.Context, tx *gorm.DB, settings *entity.ServerSettings) error {
	if settings == nil {
		return errors.New("upsert: settings cannot be nil")
	}

	entitySettings := &entity.ServerSettings{
		Key:   settings.Key,
		Value: settings.Value,
	}

	err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&entitySettings).Error
	if err != nil {
		return err
	}

	return nil
}
