package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Advik-B/megobasterd/internal/config"
	"github.com/Advik-B/megobasterd/internal/ui"
)

const (
	// VERSION represents the current version of MegaBasterd
	VERSION = "1.0.0-go"
)

func main() {
	fmt.Printf("MegaBasterd Go Edition v%s\n", VERSION)
	fmt.Println("Starting application...")

	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Warning: Could not load config: %v. Using defaults.", err)
		cfg = config.GetDefault()
	}

	// Create Fyne application
	a := app.NewWithID("com.megobasterd.go")
	
	// Create main window
	mainWindow := ui.NewMainWindow(a, cfg)
	
	// Show and run
	mainWindow.Show()
	a.Run()
}
