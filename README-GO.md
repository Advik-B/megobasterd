# MegaBasterd - Go Edition

MegaBasterd ported to Go with Fyne UI framework.

## Project Status

**Phase 1-4 Implementation (In Progress)**

Currently implemented:
- ✅ Phase 1: Project structure and foundation
- ✅ Phase 2: UI framework selection (Fyne chosen)
- ✅ Phase 3: Core modules (crypto, API client) - Basic implementation
- ✅ Phase 4: GUI implementation - Basic UI structure

Still in development:
- 🔨 Full MEGA API integration
- 🔨 Complete download/upload functionality
- 🔨 Database layer
- 🔨 Streaming server
- 🔨 Proxy management

See [docs/GOLANG_PORTING_PLAN.md](docs/GOLANG_PORTING_PLAN.md) for the complete implementation plan.

## Quick Start

### Prerequisites

- Go 1.21 or later
- For Linux: `gcc`, `libgl1-mesa-dev`, `xorg-dev`
- For macOS: Xcode command line tools
- For Windows: MinGW-w64 (for CGO)

### Installation

```bash
# Clone the repository
git clone https://github.com/Advik-B/megobasterd.git
cd megobasterd

# Install dependencies
make deps

# Build the application
make build

# Run the application
make run
```

### Development

```bash
# Run tests
make test

# Run with coverage
make test-coverage

# Run linter (requires golangci-lint)
make lint

# Development build with race detector
make dev
```

## Project Structure

```
megobasterd/
├── cmd/megobasterd/          # Main application entry point
├── internal/                 # Private application code
│   ├── api/                 # MEGA API client
│   ├── crypto/              # Encryption/decryption utilities
│   ├── downloader/          # Download management
│   ├── uploader/            # Upload management (TODO)
│   ├── database/            # SQLite database layer (TODO)
│   ├── streaming/           # Video streaming server (TODO)
│   ├── proxy/               # Proxy management (TODO)
│   ├── ui/                  # UI layer (Fyne)
│   └── config/              # Configuration management
├── pkg/                     # Public/reusable packages
│   ├── models/              # Data models (TODO)
│   └── utils/               # Utility functions (TODO)
├── assets/                  # Static resources (icons, images)
├── translations/            # i18n files (TODO)
├── scripts/                 # Build and deployment scripts
├── docs/                    # Documentation
├── go.mod                   # Go module definition
├── go.sum                   # Go dependencies
├── Makefile                 # Build automation
└── README.md               # This file
```

## Documentation

Complete documentation is available in the `docs/` folder:

- **[PORTING_DOCS_INDEX.md](docs/PORTING_DOCS_INDEX.md)** - Documentation index and navigation
- **[QUICK_START_GUIDE.md](docs/QUICK_START_GUIDE.md)** - Quick reference guide
- **[GOLANG_PORTING_PLAN.md](docs/GOLANG_PORTING_PLAN.md)** - Complete porting plan
- **[UI_FRAMEWORKS_COMPARISON.md](docs/UI_FRAMEWORKS_COMPARISON.md)** - UI framework analysis
- **[JAVA_TO_GO_REFERENCE.md](docs/JAVA_TO_GO_REFERENCE.md)** - Java to Go translation guide

## Technology Stack

- **Language:** Go 1.21+
- **UI Framework:** Fyne v2.4.3
- **HTTP Client:** go-resty/resty
- **Database:** mattn/go-sqlite3
- **Logging:** uber.org/zap
- **Configuration:** spf13/viper
- **Cryptography:** golang.org/x/crypto

## Features (Planned)

- ✅ Cross-platform GUI (Windows, macOS, Linux)
- ✅ Material Design UI
- 🔨 Multi-threaded downloads with resume capability
- 🔨 Multi-threaded uploads with encryption
- 🔨 MEGA API integration
- 🔨 Proxy support with smart rotation
- 🔨 Video streaming server
- 🔨 Multi-language support
- 🔨 System tray integration

## Building for Distribution

```bash
# Build for all platforms
make build-all

# Package with Fyne (creates platform-specific packages)
make package-fyne
```

## Contributing

This is a port of the original MegaBasterd by tonikelope. See the planning documents in `docs/` for implementation guidelines.

## License

GPL v3 - Same as original MegaBasterd

## Credits

- **Original MegaBasterd:** tonikelope
- **Go Port:** Advik-B and contributors

## Important Notice

You are not authorized to use MegaBasterd in any way that violates [MEGA's terms of use](https://mega.io/terms).
