package config

import (
	"context"
	"exaroton-wa-bot/internal/database"
	"exaroton-wa-bot/internal/database/seeder"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// supported dialects
const (
	sqliteDialect = "sqlite3"
)

func InitDB(ctx context.Context, cfg *Cfg) (*gorm.DB, error) {
	slog.Debug("Initializing database connection")

	dialect := cfg.String(keyDBDialect)

	var db *gorm.DB
	var err error

	switch dialect {
	case sqliteDialect:
		db, err = initSQLiteDB(cfg)
	default:
		return nil, fmt.Errorf("unsupported DB dialect: %s", cfg.String(keyDBDialect))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	slog.Debug("Database connection established")

	// migrate database
	sqlDB, err := db.DB()
	if err := database.MigrateDB(ctx, dialect, sqlDB); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// seed constants
	if err := seeder.RunConstantSeeder(db); err != nil {
		return nil, fmt.Errorf("failed to seed constants to database: %w", err)
	}

	return db, err
}

func initSQLiteDB(cfg *Cfg) (*gorm.DB, error) {
	dbPath := cfg.String(keySQLiteDBPath)
	if dbPath == "" {
		dbPath = "../db/data.db"
		slog.Warn("SQLite db path is not set. Using default sqlite db path: " + (dbPath))
	}

	slog.Debug("SQLite path", "path", dbPath)

	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create db folder: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(absPath), getGormConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	return db, nil
}

// gorm config
func getGormConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
}
