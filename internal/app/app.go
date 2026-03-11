package app

import (
	"context"

	"gorm.io/gorm"
)

// App is the main application struct. Its methods are bound to the Wails
// frontend as TypeScript functions via wails.Run Bind field.
type App struct {
	ctx context.Context
	db  *gorm.DB
}

// NewApp creates a new App with the given database connection.
func NewApp(db *gorm.DB) *App {
	return &App{db: db}
}

// OnStartup is called by Wails after the window is created.
func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx
}

// OnShutdown is called by Wails when the window is closed.
// Phase 3+ flushes queues here.
func (a *App) OnShutdown(ctx context.Context) {
}
