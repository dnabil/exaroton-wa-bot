package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func MigrateDB(ctx context.Context, dialect string, db *sql.DB) error {
	fsys, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to  migrations filesystem: %w", err)
	}

	pvd, err := goose.NewProvider(goose.Dialect(dialect), db, fsys)
	if err != nil {
		return fmt.Errorf("failed to create migrations provider: %w", err)
	}

	migrations, err := pvd.Up(ctx)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	for _, migration := range migrations {
		slog.Debug("Migration applied", "data", migration.String())
	}

	return nil
}
