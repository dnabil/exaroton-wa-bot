package database

import (
	"testing"
)

func TestMigrationsFS(t *testing.T) {
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		t.Fatalf("Failed to read migrations directory: %v", err)
	}

	if len(entries) == 0 {
		t.Fatal("No migration files found")
	}

	t.Logf("Found %d migration files", len(entries))
}
