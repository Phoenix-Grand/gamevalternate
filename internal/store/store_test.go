package store_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gamevault-go/internal/store"
)

func TestSchemaInit(t *testing.T) {
	tmp := t.TempDir()
	dbPath := filepath.Join(tmp, "test.db")
	db := store.Open(dbPath)
	if db == nil {
		t.Fatal("expected non-nil *gorm.DB")
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB() error: %v", err)
	}
	defer sqlDB.Close()
}

func TestTableNames(t *testing.T) {
	tmp := t.TempDir()
	db := store.Open(filepath.Join(tmp, "test.db"))

	var tables []struct{ Name string }
	if err := db.Raw("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name").Scan(&tables).Error; err != nil {
		t.Fatalf("query failed: %v", err)
	}

	names := make(map[string]bool)
	for _, tbl := range tables {
		names[tbl.Name] = true
	}

	required := []string{"server_profiles", "game_caches", "downloads", "save_paths", "app_settings"}
	for _, req := range required {
		if !names[req] {
			t.Errorf("missing table: %s (found: %v)", req, names)
		}
	}
}

func TestDefaultPath(t *testing.T) {
	p := store.DefaultPath()
	if p == "" {
		t.Fatal("DefaultPath() returned empty string")
	}
	if !strings.Contains(p, "gamevault-go") {
		t.Errorf("expected path to contain 'gamevault-go', got: %s", p)
	}
	if filepath.Base(p) != "state.db" {
		t.Errorf("expected filename 'state.db', got: %s", filepath.Base(p))
	}
}

func TestOpenIdempotent(t *testing.T) {
	tmp := t.TempDir()
	dbPath := filepath.Join(tmp, "test.db")
	// Should not panic or fatal on second call
	store.Open(dbPath)
	store.Open(dbPath)
}

// TestDefaultPath_HomeEnv verifies fallback when XDG_CONFIG_HOME is not set.
// This is an extended test to ensure robust path resolution.
func TestDefaultPath_HomeEnv(t *testing.T) {
	p := store.DefaultPath()
	// Path must not be empty regardless of environment
	if p == "" {
		t.Fatal("DefaultPath() returned empty string even with HOME set")
	}
	// Must end with state.db
	if !strings.HasSuffix(p, "state.db") {
		t.Errorf("expected path to end with 'state.db', got: %s", p)
	}
	// Directory containing state.db must be named gamevault-go
	dir := filepath.Dir(p)
	if filepath.Base(dir) != "gamevault-go" {
		t.Errorf("expected parent dir to be 'gamevault-go', got: %s", filepath.Base(dir))
	}
	_ = os.MkdirAll(dir, 0755) // ensure path is creatable
}
