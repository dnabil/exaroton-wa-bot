package seeder

import (
	"errors"
	"exaroton-wa-bot/internal/database/entity"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var constantSeeders = []seederFunc{
	constantSeedUser,
}

func RunConstantSeeder(db *gorm.DB) error {
	tx := db.Begin()
	defer tx.Rollback()

	for _, seeder := range constantSeeders {
		if err := seeder(tx); err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}

// only seed if user with id 0 does not exist.
func constantSeedUser(tx *gorm.DB) error {
	user := &entity.User{}
	if err := tx.Where("id = ?", 1).First(user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	if user.ID != 0 {
		return nil
	}

	user = &entity.User{
		ID:       1,
		Username: "admin",
		Password: "admin",
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return tx.Create(user).Error
}
