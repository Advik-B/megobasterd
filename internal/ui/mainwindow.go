package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Advik-B/megobasterd/internal/config"
)

// MainWindow represents the main application window
type MainWindow struct {
	app      fyne.App
	window   fyne.Window
	config   *config.Config
	
	// Tabs
	downloadTab *fyne.Container
	uploadTab   *fyne.Container
	
	// Download list
	downloadTable *widget.Table
	
	// Upload list
	uploadTable *widget.Table
}

// NewMainWindow creates a new main window
func NewMainWindow(app fyne.App, cfg *config.Config) *MainWindow {
	w := app.NewWindow("MegaBasterd - Go Edition")
	
	mw := &MainWindow{
		app:    app,
		window: w,
		config: cfg,
	}
	
	mw.setupUI()
	return mw
}

// setupUI initializes the user interface
func (mw *MainWindow) setupUI() {
	// Create tabs
	mw.downloadTab = mw.createDownloadsTab()
	mw.uploadTab = mw.createUploadsTab()
	
	tabs := container.NewAppTabs(
		container.NewTabItem("Downloads", mw.downloadTab),
		container.NewTabItem("Uploads", mw.uploadTab),
	)
	
	// Create menu
	mw.createMenu()
	
	// Set content
	mw.window.SetContent(tabs)
	mw.window.Resize(fyne.NewSize(1024, 768))
	
	// Setup system tray (if supported)
	mw.setupSystemTray()
}

// createDownloadsTab creates the downloads tab
func (mw *MainWindow) createDownloadsTab() *fyne.Container {
	// Sample data for demonstration
	downloadData := [][]string{
		{"example1.zip", "150 MB", "0%", "0 MB/s", "Queued"},
		{"example2.mp4", "500 MB", "0%", "0 MB/s", "Queued"},
	}
	
	// Create download table
	mw.downloadTable = widget.NewTable(
		func() (int, int) { return len(downloadData), 5 },
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			if id.Row < len(downloadData) && id.Col < len(downloadData[id.Row]) {
				label.SetText(downloadData[id.Row][id.Col])
			}
		},
	)
	
	// Set column widths
	mw.downloadTable.SetColumnWidth(0, 200) // Name
	mw.downloadTable.SetColumnWidth(1, 100) // Size
	mw.downloadTable.SetColumnWidth(2, 80)  // Progress
	mw.downloadTable.SetColumnWidth(3, 100) // Speed
	mw.downloadTable.SetColumnWidth(4, 100) // Status
	
	// Create toolbar
	addBtn := widget.NewButton("Add Download", mw.showAddDownloadDialog)
	pauseBtn := widget.NewButton("Pause", func() {
		// TODO: Implement pause
	})
	resumeBtn := widget.NewButton("Resume", func() {
		// TODO: Implement resume
	})
	removeBtn := widget.NewButton("Remove", func() {
		// TODO: Implement remove
	})
	
	toolbar := container.NewHBox(
		addBtn,
		pauseBtn,
		resumeBtn,
		removeBtn,
	)
	
	return container.NewBorder(toolbar, nil, nil, nil, mw.downloadTable)
}

// createUploadsTab creates the uploads tab
func (mw *MainWindow) createUploadsTab() *fyne.Container {
	// Sample data for demonstration
	uploadData := [][]string{
		{"file1.zip", "100 MB", "0%", "0 MB/s", "Queued"},
	}
	
	// Create upload table
	mw.uploadTable = widget.NewTable(
		func() (int, int) { return len(uploadData), 5 },
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			if id.Row < len(uploadData) && id.Col < len(uploadData[id.Row]) {
				label.SetText(uploadData[id.Row][id.Col])
			}
		},
	)
	
	// Set column widths
	mw.uploadTable.SetColumnWidth(0, 200) // Name
	mw.uploadTable.SetColumnWidth(1, 100) // Size
	mw.uploadTable.SetColumnWidth(2, 80)  // Progress
	mw.uploadTable.SetColumnWidth(3, 100) // Speed
	mw.uploadTable.SetColumnWidth(4, 100) // Status
	
	// Create toolbar
	addBtn := widget.NewButton("Add Upload", mw.showAddUploadDialog)
	pauseBtn := widget.NewButton("Pause", func() {
		// TODO: Implement pause
	})
	removeBtn := widget.NewButton("Remove", func() {
		// TODO: Implement remove
	})
	
	toolbar := container.NewHBox(
		addBtn,
		pauseBtn,
		removeBtn,
	)
	
	return container.NewBorder(toolbar, nil, nil, nil, mw.uploadTable)
}

// createMenu creates the application menu
func (mw *MainWindow) createMenu() {
	// File menu
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Settings", mw.showSettingsDialog),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", func() {
			mw.app.Quit()
		}),
	)
	
	// Help menu
	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About", mw.showAboutDialog),
	)
	
	mainMenu := fyne.NewMainMenu(fileMenu, helpMenu)
	mw.window.SetMainMenu(mainMenu)
}

// setupSystemTray sets up system tray icon and menu
func (mw *MainWindow) setupSystemTray() {
	// System tray is handled by Fyne's desktop extension
	// This would be implemented when running on desktop platforms
}

// Show displays the main window
func (mw *MainWindow) Show() {
	mw.window.Show()
}

// showAddDownloadDialog shows dialog to add a new download
func (mw *MainWindow) showAddDownloadDialog() {
	// TODO: Implement add download dialog
	// For now, just show a simple entry dialog
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Enter MEGA URL...")
	
	content := container.NewVBox(
		widget.NewLabel("Add Download"),
		urlEntry,
	)
	
	dialog := widget.NewModalPopUp(content, mw.window.Canvas())
	dialog.Show()
}

// showAddUploadDialog shows dialog to add a new upload
func (mw *MainWindow) showAddUploadDialog() {
	// TODO: Implement add upload dialog
	content := container.NewVBox(
		widget.NewLabel("Add Upload"),
		widget.NewLabel("File selection dialog would appear here"),
	)
	
	dialog := widget.NewModalPopUp(content, mw.window.Canvas())
	dialog.Show()
}

// showSettingsDialog shows the settings dialog
func (mw *MainWindow) showSettingsDialog() {
	// TODO: Implement full settings dialog
	content := container.NewVBox(
		widget.NewLabel("Settings"),
		widget.NewLabel("Settings dialog to be implemented"),
	)
	
	dialog := widget.NewModalPopUp(content, mw.window.Canvas())
	dialog.Show()
}

// showAboutDialog shows the about dialog
func (mw *MainWindow) showAboutDialog() {
	content := container.NewVBox(
		widget.NewLabel("MegaBasterd - Go Edition"),
		widget.NewLabel("Version 1.0.0"),
		widget.NewLabel(""),
		widget.NewLabel("A MEGA downloader/uploader written in Go"),
		widget.NewLabel(""),
		widget.NewLabel("Original by tonikelope"),
		widget.NewLabel("Go port by Advik-B"),
	)
	
	dialog := widget.NewModalPopUp(content, mw.window.Canvas())
	dialog.Show()
}
