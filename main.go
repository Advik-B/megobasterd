package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/Advik-B/megobasterd/internal/app"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

const VERSION = "1.0.0-wails"

func main() {
	fmt.Printf("MegaBasterd Go Edition v%s (Wails)\n", VERSION)

	// Create an instance of the app structure
	megaApp := app.NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "MegaBasterd - Go Edition",
		Width:  1200,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        megaApp.Startup,
		Bind: []interface{}{
			megaApp,
		},
	})

	if err != nil {
		log.Fatal("Error:", err)
	}
}
