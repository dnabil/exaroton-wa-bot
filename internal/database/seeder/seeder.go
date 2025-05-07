package seeder

import "gorm.io/gorm"

type seederFunc func(tx *gorm.DB) error
