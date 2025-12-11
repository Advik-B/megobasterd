# MegaBasterd - Go Edition (Wails + Svelte)

A modern, cross-platform MEGA downloader built with Go and Wails, featuring a beautiful Svelte frontend.

## Features

✅ **Working Download Functionality**
- Real MEGA API integration
- Multi-threaded downloads
- Progress tracking with speed calculation
- Pause/Resume capability
- Real-time UI updates

✅ **Modern UI**
- Beautiful gradient design
- Real-time progress bars
- Download status badges
- Responsive layout

✅ **Cross-Platform**
- Windows, macOS, Linux support
- Native performance with Go backend
- Web technologies for UI (Svelte)

## Technology Stack

- **Backend:** Go 1.21+
- **Frontend:** Svelte + Vite
- **Framework:** Wails v2
- **API Client:** go-resty/resty
- **Configuration:** Viper
- **Crypto:** golang.org/x/crypto

## Quick Start

### Prerequisites

1. **Go 1.21+**
   ```bash
   go version
   ```

2. **Node.js 16+**
   ```bash
   node --version
   npm --version
   ```

3. **Wails CLI**
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

4. **Platform-specific requirements:**
   - **Linux:** `gcc`, `gtk3`, `webkit2gtk`
     ```bash
     sudo apt install gcc libgtk-3-dev libwebkit2gtk-4.0-dev
     ```
   - **macOS:** Xcode command line tools
   - **Windows:** MinGW-w64 or MSVC

### Installation

```bash
# Clone the repository
git clone https://github.com/Advik-B/megobasterd.git
cd megobasterd

# Install dependencies
make deps

# Run in development mode
make dev
```

### Building

```bash
# Build for current platform
make build

# The executable will be in build/bin/
```

## Usage

1. **Launch the application**
   ```bash
   make dev  # Development mode with hot reload
   # or
   ./build/bin/megobasterd  # Production build
   ```

2. **Add a download**
   - Paste a MEGA URL (e.g., `https://mega.nz/file/...#...`)
   - Click "Add Download"
   - Watch the download progress in real-time!

3. **Manage downloads**
   - **Pause:** Click the "Pause" button while downloading
   - **Remove:** Click the "Remove" button to delete from list

## Project Structure

```
megobasterd/
├── cmd/megobasterd/          # Main application entry point
│   └── main.go
├── internal/
│   ├── app/                  # Wails app backend
│   │   └── app.go           # Download logic & API integration
│   ├── api/                  # MEGA API client
│   ├── crypto/               # Encryption utilities
│   ├── config/               # Configuration management
│   └── downloader/           # Download manager
├── frontend/
│   ├── src/
│   │   ├── App.svelte       # Main Svelte component
│   │   └── main.js          # Entry point
│   ├── public/
│   │   └── index.html
│   ├── package.json
│   ├── vite.config.js
│   └── svelte.config.js
├── wails.json                # Wails configuration
├── go.mod
├── Makefile
└── README-WAILS.md          # This file
```

## Development

### Running in Dev Mode

```bash
make dev
```

This starts:
- Go backend with hot reload
- Vite dev server for frontend
- Auto-refresh on code changes

### Building for Production

```bash
# Optimized build with UPX compression
make build-prod

# Or standard build
make build
```

### Running Tests

```bash
# Run all tests
make test

# With coverage
make test-coverage
```

## Download Functionality

The application implements real MEGA download functionality:

1. **URL Parsing:** Extracts file ID and encryption key from MEGA URLs
2. **API Communication:** Retrieves download URL from MEGA API
3. **HTTP Download:** Streams file content with progress tracking
4. **Real-time Updates:** UI updates every 500ms with:
   - Downloaded bytes
   - Total file size
   - Download speed
   - Progress percentage
   - Status (queued, downloading, completed, failed, paused)

### Code Example

```go
// Backend: internal/app/app.go
func (a *App) AddDownload(url string) (*Download, error) {
    // Parse MEGA URL
    fileID, key, err := a.ParseMegaURL(url)
    
    // Get download URL from MEGA API
    downloadURL, err := a.client.GetDownloadURL(ctx, fileID, key)
    
    // Start download in background with progress tracking
    go a.startDownload(download)
    
    return download, nil
}
```

```svelte
<!-- Frontend: frontend/src/App.svelte -->
<script>
  import { AddDownload, GetDownloads } from '../wailsjs/go/app/App';
  import { EventsOn } from '../wailsjs/runtime/runtime';
  
  // Listen for real-time updates
  EventsOn('download-update', (download) => {
    // Update UI with latest progress
  });
</script>
```

## Configuration

Default download path: `~/Downloads`

To customize, edit `~/.megobasterd/megobasterd.yaml`:

```yaml
download_path: "/path/to/downloads"
max_workers: 6
language: "en"
```

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make dev` | Run with hot reload |
| `make build` | Build for current platform |
| `make build-prod` | Optimized build with compression |
| `make test` | Run all tests |
| `make test-coverage` | Run tests with coverage report |
| `make deps` | Install all dependencies |
| `make clean` | Remove build artifacts |
| `make doctor` | Check Wails setup |

## Troubleshooting

### "wails: command not found"
```bash
make install-wails
# or
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Build fails on Linux
Install required dependencies:
```bash
sudo apt install gcc libgtk-3-dev libwebkit2gtk-4.0-dev
```

### Frontend not loading
```bash
cd frontend
npm install
npm run build
```

## Contributing

This is a port of the original MegaBasterd. Contributions welcome!

## License

GPL v3 - Same as original MegaBasterd

## Credits

- **Original MegaBasterd:** tonikelope
- **Go Port:** Advik-B
- **UI Framework:** Wails (Go + Svelte)

## Disclaimer

You are not authorized to use MegaBasterd in any way that violates [MEGA's terms of use](https://mega.io/terms).
