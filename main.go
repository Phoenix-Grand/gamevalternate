package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"gamevault-go/internal/app"
	"gamevault-go/internal/store"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	db := store.Open(store.DefaultPath())
	application := app.NewApp(db)

	err := wails.Run(&options.App{
		Title:     "GameVault",
		Width:     1200,
		Height:    800,
		MinWidth:  900,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  application.OnStartup,
		OnShutdown: application.OnShutdown,
		Bind: []interface{}{
			application,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
