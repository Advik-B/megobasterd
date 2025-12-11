# MegaBasterd - Go Edition (Wails + Svelte)

A modern, cross-platform MEGA downloader built with Go and Wails, featuring a beautiful Svelte frontend.

## Features

✅ **Working Download Functionality**
- Real MEGA API integration
- HTTP streaming downloads
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
- Python-based build system (no Make dependency)

## Technology Stack

- **Backend:** Go 1.21+
- **Frontend:** Svelte 4 + Vite 5
- **Framework:** Wails v2.8.0
- **API Client:** go-resty/resty
- **Configuration:** Viper
- **Crypto:** golang.org/x/crypto
- **Build System:** Python 3 (cross-platform)

## Quick Start

### Prerequisites

1. **Python 3.6+** (for build scripts)
   ```bash
   python3 --version
   ```

2. **Go 1.21+**
   ```bash
   go version
   ```

3. **Node.js 16+**
   ```bash
   node --version
   npm --version
   ```

4. **Wails CLI**
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

5. **Platform-specific requirements:**
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
python3 build.py deps

# Run in development mode
python3 build.py dev
```

### Building

```bash
# Build for current platform
python3 build.py build

# Build with optimizations (UPX compression)
python3 build.py build --upx

# The executable will be in build/bin/
```

## Usage

1. **Launch the application**
   ```bash
   python3 build.py dev  # Development mode with hot reload
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
├── main.go                   # Main application entry point (Wails)
├── build.py                  # Cross-platform build script
├── run.py                    # Quick wrapper script
├── wails.json                # Wails configuration
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
├── build.py                  # Python build script (replaces Makefile)
├── run.py                    # Quick run script
└── README-WAILS.md           # This file
```

## Development

### Running in Dev Mode

```bash
python3 build.py dev
```

This starts:
- Go backend with hot reload
- Vite dev server for frontend
- Auto-refresh on code changes

### Building for Production

```bash
# Optimized build with UPX compression
python3 build.py build --upx

# Or standard build
python3 build.py build
```

### Running Tests

```bash
# Run all tests
python3 build.py test

# With coverage
python3 build.py test --coverage
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

## Build Script Commands

| Command | Description |
|---------|-------------|
| `python3 build.py dev` | Run with hot reload |
| `python3 build.py build` | Build for current platform |
| `python3 build.py build --upx` | Optimized build with UPX compression |
| `python3 build.py test` | Run all tests |
| `python3 build.py test --coverage` | Run tests with coverage report |
| `python3 build.py deps` | Install all dependencies |
| `python3 build.py clean` | Remove build artifacts |
| `python3 build.py doctor` | Check Wails setup |
| `python3 build.py install-wails` | Install Wails CLI |
| `python3 build.py generate` | Generate Wails bindings |

### Quick Run Script

For convenience, you can also use `run.py`:

```bash
./run.py dev          # Same as: python3 build.py dev
./run.py build        # Same as: python3 build.py build
./run.py test         # Same as: python3 build.py test
```

## Troubleshooting

### "no Go files in D:\Github\megobasterd" error

This error occurs when the project structure is incorrect. The fix has been applied:
- `main.go` is now at the project root (not in `cmd/megobasterd/`)
- `frontend/dist/` directory is created automatically
- Run `python3 build.py deps` to ensure everything is set up correctly

### "wails: command not found"
```bash
python3 build.py install-wails
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

### "pattern all:frontend/dist: no matching files found"
This happens if the frontend hasn't been built yet. Solutions:
1. Run `python3 build.py deps` to install frontend dependencies
2. The `frontend/dist/.gitkeep` ensures the directory exists
3. Wails will build the frontend automatically during `wails dev` or `wails build`

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
