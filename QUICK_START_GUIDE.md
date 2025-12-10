# Quick Start Guide - MegaBasterd Golang Port

This is a quick reference guide for getting started with the MegaBasterd Golang port. For comprehensive details, see:
- **GOLANG_PORTING_PLAN.md** - Complete 16-week porting strategy
- **UI_FRAMEWORKS_COMPARISON.md** - Detailed UI framework analysis

---

## TL;DR - Executive Summary

**What**: Port MegaBasterd (Java/Swing desktop app) to Golang  
**Why**: Better performance, easier distribution, modern codebase  
**Recommended UI**: Fyne (cross-platform, pure Go, easy to learn)  
**Alternative UI**: Wails (modern web-based UI if team knows React/Vue)  
**Timeline**: 16 weeks  
**Team Size**: 2-3 developers recommended  

---

## Current State Analysis

- **Language**: Java 8+ (currently ~30,000 lines of code)
- **Build Tool**: Maven
- **UI Framework**: Java Swing (106 UI component imports)
- **Files**: 61 Java source files
- **Key Features**:
  - MEGA file download/upload with encryption
  - Multi-threaded transfers (6 workers default)
  - Proxy support with smart rotation
  - Video streaming server
  - SQLite database for persistence
  - System tray integration
  - Multi-language support
  - Clipboard monitoring

---

## Recommended Tech Stack for Go Port

### Core
```go
Language:     Go 1.21+
Build:        Go modules (go.mod) + Makefile
Package Mgr:  go get / go mod
```

### Key Libraries
```go
// Cryptography
golang.org/x/crypto

// Database
github.com/mattn/go-sqlite3

// HTTP Client
github.com/go-resty/resty/v2

// JSON
github.com/json-iterator/go

// Configuration
github.com/spf13/viper

// Logging
go.uber.org/zap

// Concurrency
golang.org/x/sync
```

### UI Framework (Choose ONE)

**Option 1: Fyne** ⭐ RECOMMENDED
```bash
go get fyne.io/fyne/v2
```
- Pure Go
- Cross-platform (Windows, macOS, Linux)
- Material Design UI
- Easy learning curve
- System tray support

**Option 2: Wails** (Alternative for modern web UI)
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```
- Use React/Vue/Svelte for frontend
- Go backend
- Modern, beautiful UIs
- Slightly more complex

---

## Project Structure

```
megobasterd-go/
├── cmd/megobasterd/main.go          # Entry point
├── internal/                         # Private app code
│   ├── api/                         # MEGA API client
│   ├── crypto/                      # Encryption utilities
│   ├── downloader/                  # Download management
│   ├── uploader/                    # Upload management
│   ├── database/                    # SQLite layer
│   ├── streaming/                   # Video streaming
│   ├── proxy/                       # Proxy management
│   └── ui/                          # UI layer
├── pkg/                             # Public packages
│   ├── models/                      # Data structures
│   └── utils/                       # Utilities
├── assets/                          # Images, icons
├── translations/                    # i18n files
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## Development Phases (16 Weeks)

### Phase 1-2: Foundation (Weeks 1-2)
- [ ] Set up Go project structure
- [ ] Choose UI framework (Fyne recommended)
- [ ] Set up build system (Makefile)
- [ ] Configure CI/CD

### Phase 3-5: Core Backend (Weeks 3-5)
- [ ] Port cryptography (AES, RSA, PBKDF2)
- [ ] Implement MEGA API client
- [ ] Write unit tests

### Phase 6-8: Transfers (Weeks 6-8)
- [ ] Download manager with multi-threading
- [ ] Upload manager with encryption
- [ ] Progress tracking and resume

### Phase 9-12: GUI (Weeks 9-12)
- [ ] Main window with tabs
- [ ] Download/upload lists
- [ ] All dialogs (settings, about, etc.)
- [ ] System tray integration

### Phase 13: Config & i18n (Week 13)
- [ ] Configuration management
- [ ] Multi-language support

### Phase 14-15: Testing (Weeks 14-15)
- [ ] Comprehensive testing
- [ ] Performance benchmarks
- [ ] Bug fixes

### Phase 16: Release (Week 16)
- [ ] Documentation
- [ ] Cross-platform builds
- [ ] Distribution packages

---

## Quick Start - Proof of Concept (Fyne)

Create a simple PoC to validate the approach:

```go
// main.go
package main

import (
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

func main() {
    // Create app
    a := app.New()
    w := a.NewWindow("MegaBasterd Go - PoC")

    // Add download button
    addBtn := widget.NewButton("Add Download", func() {
        // TODO: Show add dialog
    })

    // Download table (simplified)
    data := [][]string{
        {"file1.zip", "100 MB", "50%", "1.2 MB/s"},
        {"file2.mp4", "500 MB", "75%", "2.5 MB/s"},
    }

    table := widget.NewTable(
        func() (int, int) { return len(data), 4 },
        func() fyne.CanvasObject {
            return widget.NewLabel("cell")
        },
        func(id widget.TableCellID, cell fyne.CanvasObject) {
            label := cell.(*widget.Label)
            label.SetText(data[id.Row][id.Col])
        },
    )

    // Layout
    content := container.NewBorder(
        addBtn,  // top
        nil,     // bottom
        nil,     // left
        nil,     // right
        table,   // center
    )

    w.SetContent(content)
    w.Resize(fyne.NewSize(800, 600))
    w.ShowAndRun()
}
```

Run it:
```bash
go mod init github.com/yourusername/megobasterd-go
go get fyne.io/fyne/v2
go run main.go
```

---

## Development Setup

### Prerequisites

**All Platforms:**
- Go 1.21 or later
- Git

**For Fyne (Linux):**
```bash
sudo apt-get install gcc libgl1-mesa-dev xorg-dev
```

**For Fyne (macOS):**
```bash
xcode-select --install
```

**For Fyne (Windows):**
- Install MinGW-w64 (for CGO)
- Or use TDM-GCC

### Installation

```bash
# Clone repository
git clone https://github.com/Advik-B/megobasterd-go
cd megobasterd-go

# Initialize Go module
go mod init github.com/Advik-B/megobasterd-go

# Install dependencies
go mod tidy

# Build
go build -o bin/megobasterd cmd/megobasterd/main.go

# Run
./bin/megobasterd
```

---

## Java to Go Translation Guide

### Common Patterns

**Java Swing → Fyne Mapping:**

```
// Java
JFrame frame = new JFrame("Title");
JButton button = new JButton("Click");
button.addActionListener(e -> handleClick());

// Go (Fyne)
window := app.NewWindow("Title")
button := widget.NewButton("Click", handleClick)
```

**Concurrency:**

```java
// Java
ExecutorService executor = Executors.newCachedThreadPool();
executor.submit(() -> doWork());

// Go
go doWork()
// or with error handling
g, ctx := errgroup.WithContext(context.Background())
g.Go(func() error {
    return doWork()
})
err := g.Wait()
```

**Synchronized/Locks:**

```java
// Java
synchronized(this) {
    count++;
}

// Go
mu.Lock()
count++
mu.Unlock()
```

**File I/O:**

```java
// Java
Files.readAllBytes(Paths.get("file.txt"))

// Go
os.ReadFile("file.txt")
```

---

## Testing Strategy

### Unit Tests
```go
// internal/crypto/crypto_test.go
func TestAESEncryption(t *testing.T) {
    plaintext := []byte("test data")
    key := make([]byte, 32)
    
    encrypted, err := Encrypt(plaintext, key)
    assert.NoError(t, err)
    
    decrypted, err := Decrypt(encrypted, key)
    assert.NoError(t, err)
    assert.Equal(t, plaintext, decrypted)
}
```

### Running Tests
```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Build & Distribution

### Single Platform
```bash
# Current platform
go build -o megobasterd cmd/megobasterd/main.go
```

### Cross-Platform
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o megobasterd-linux

# Windows
GOOS=windows GOARCH=amd64 go build -o megobasterd.exe

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o megobasterd-macos-intel

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o megobasterd-macos-arm
```

### Fyne Packaging
```bash
# Creates platform-specific packages
fyne package -os linux -icon assets/icon.png
fyne package -os windows -icon assets/icon.png
fyne package -os darwin -icon assets/icon.png
```

---

## Performance Benchmarks

Before porting, establish baseline metrics from Java version:

- **Download Speed**: ___ MB/s (6 workers)
- **Upload Speed**: ___ MB/s (6 workers)
- **Memory Usage**: ___ MB (idle), ___ MB (active)
- **CPU Usage**: ___ % (downloading)
- **Startup Time**: ___ seconds
- **Binary Size**: ___ MB

Target: Match or exceed Java version performance.

---

## Resources

### Documentation
- [Go Documentation](https://go.dev/doc/)
- [Fyne Documentation](https://developer.fyne.io/)
- [Wails Documentation](https://wails.io/docs/)
- [Effective Go](https://go.dev/doc/effective_go)

### Learning Resources
- [Go by Example](https://gobyexample.com/)
- [Fyne Examples](https://github.com/fyne-io/examples)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)

### Community
- [r/golang](https://reddit.com/r/golang)
- [Gophers Slack](https://gophers.slack.com/)
- [Fyne Discord](https://discord.gg/fyne)

---

## Next Steps

1. **Review Documentation**
   - Read GOLANG_PORTING_PLAN.md thoroughly
   - Review UI_FRAMEWORKS_COMPARISON.md
   - Understand current Java codebase

2. **Team Setup**
   - Assign roles (backend, UI, testing)
   - Set up development environments
   - Create task tracking (GitHub Issues/Projects)

3. **Proof of Concept** (Week 1)
   - Build minimal UI with Fyne
   - Test MEGA API connection
   - Validate crypto compatibility
   - Test on all target platforms

4. **Go/No-Go Decision** (End of Week 1)
   - Review PoC results
   - Confirm UI framework choice
   - Commit to full migration or adjust plan

5. **Start Development** (Week 2+)
   - Follow phase plan in GOLANG_PORTING_PLAN.md
   - Regular standups and progress reviews
   - Continuous testing and integration

---

## Frequently Asked Questions

**Q: Why Go over staying with Java?**
A: Better performance, smaller binaries, easier distribution, modern tooling, strong concurrency support.

**Q: Will it be faster than Java?**
A: Likely yes for I/O-heavy operations. Go's goroutines are more efficient than Java threads for many concurrent operations.

**Q: How big will the binary be?**
A: ~10-30 MB (vs Java requiring JRE ~200MB+). With Fyne, expect 15-25 MB depending on platform.

**Q: Can we keep the same UI?**
A: Similar, but not identical. Fyne uses Material Design. Wails allows custom designs.

**Q: What about Java features we use?**
A: All Java features used in MegaBasterd have Go equivalents (crypto, networking, SQLite, etc.).

**Q: Is Go harder to learn than Java?**
A: No, Go is actually simpler. Smaller language, easier concurrency, no generics complexity (until recently).

**Q: Cross-platform build complexity?**
A: Easier than Java! Single command cross-compilation. No JRE installation needed.

---

## Support & Contact

For questions during the port:
- File issues on GitHub
- Review existing Java code for reference
- Check Go documentation
- Ask in Gophers Slack

Good luck with the port! 🚀
