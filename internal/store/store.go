// Package store provides SQLite database initialization and schema management.
package store

import (
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// Open opens (or creates) the SQLite database at dbPath, runs AutoMigrate to
// ensure all schema tables exist, and returns the *gorm.DB handle.
// Calls log.Fatalf on any unrecoverable error.
func Open(dbPath string) *gorm.DB {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatalf("failed to create data dir: %v", err)
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	if err := db.AutoMigrate(
		&ServerProfile{},
		&GameCache{},
		&Download{},
		&SavePath{},
		&AppSetting{},
	); err != nil {
		log.Fatalf("failed to migrate schema: %v", err)
	}
	return db
}
