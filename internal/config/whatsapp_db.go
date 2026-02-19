package config

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type WhatsappDB struct {
	// container is just another name for whatsapp database.
	Container    *sqlstore.Container
	ClientLogger waLog.Logger
}

func InitWhatsappDB(cfg *Cfg) (*WhatsappDB, error) {
	var db *sql.DB
	var err error

	switch cfg.String(keyWADBDialect) {
	case sqliteDialect:
		db, err = initWhatsappSQLite(cfg)
	default:
		return nil, fmt.Errorf("unsupported dialect for wa-db: %s", cfg.String(keyWADBDialect))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to init wa-db: %w", err)
	}

	slog.Debug("Whatsapp database initialized")

	dbLogger, clientLogger := initWhatsappLogger(cfg)

	waContainer := sqlstore.NewWithDB(db, cfg.String(keyWADBDialect), dbLogger)

	// upgrade database schema if needed (won't do anything if not needed)
	if err = waContainer.Upgrade(context.TODO()); err != nil {
		return nil, fmt.Errorf("failed to upgrade whatsapp db: %w", err)
	}

	return &WhatsappDB{
		Container:    waContainer,
		ClientLogger: clientLogger,
	}, nil
}

func initWhatsappSQLite(cfg *Cfg) (*sql.DB, error) {
	dbPath := cfg.String(keyWASQLiteDBPath)
	if dbPath == "" {
		dbPath = "../db/wa-data.db"
		slog.Warn("SQlite db path for whatsapp is not set. Using default sqlite db path: " + (dbPath))
	}

	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	slog.Debug("Initializing WhatsApp Database", "path", absPath)

	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create db folder: %w", err)
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", absPath))
	if err != nil {
		return nil, err
	}

	return db, nil
}
