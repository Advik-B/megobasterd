# Build System Documentation

This project uses Python scripts instead of Make for cross-platform build automation.

## Why Python Instead of Make?

1. **Cross-platform compatibility** - Works identically on Windows, macOS, and Linux
2. **Better error handling** - Python provides clearer error messages
3. **No Make dependency** - Only requires Python 3, which is commonly available
4. **More readable** - Python syntax is clearer than Makefile syntax
5. **Extensible** - Easy to add complex build logic when needed

## Quick Reference

### Main Build Script: `build.py`

```bash
# Development
python3 build.py dev                  # Run dev server with hot reload

# Building
python3 build.py build                # Standard build
python3 build.py build --upx          # Optimized build with UPX compression
python3 build.py build --clean        # Clean before building

# Testing
python3 build.py test                 # Run all tests
python3 build.py test --coverage      # Run tests with coverage report

# Maintenance
python3 build.py clean                # Remove build artifacts
python3 build.py deps                 # Install all dependencies
python3 build.py lint                 # Run linter (requires golangci-lint)

# Utilities
python3 build.py doctor               # Check Wails setup
python3 build.py install-wails        # Install Wails CLI
python3 build.py generate             # Generate Wails bindings
```

### Quick Run Script: `run.py`

For convenience, `run.py` is a wrapper around `build.py`:

```bash
./run.py dev        # Same as: python3 build.py dev
./run.py build      # Same as: python3 build.py build
./run.py test       # Same as: python3 build.py test
```

## Available Commands

### `dev`
Starts the Wails development server with hot reload for both backend and frontend.

**What it does:**
- Compiles Go backend
- Starts Vite dev server for frontend
- Watches for file changes
- Auto-reloads on changes

**Example:**
```bash
python3 build.py dev
```

### `build`
Builds the application for the current platform.

**Options:**
- `--clean` - Clean build artifacts before building
- `--upx` - Compress the binary with UPX

**Examples:**
```bash
python3 build.py build                # Standard build
python3 build.py build --upx          # Optimized build
python3 build.py build --clean --upx  # Clean + optimized build
```

**Output:** Binary in `build/bin/megobasterd` (or `.exe` on Windows)

### `test`
Runs all Go tests in the `internal/` and `pkg/` directories.

**Options:**
- `--coverage` - Generate HTML coverage report

**Examples:**
```bash
python3 build.py test                 # Run tests
python3 build.py test --coverage      # Run with coverage
```

**Note:** Tests exclude `cmd/megobasterd` as it requires frontend to be built.

### `clean`
Removes all build artifacts and generated files.

**Removes:**
- `build/` - Compiled binaries
- `frontend/dist/` - Frontend build output
- `frontend/node_modules/` - Frontend dependencies
- `coverage.out` - Go coverage data
- `coverage.html` - Coverage report

**Example:**
```bash
python3 build.py clean
```

### `deps`
Installs all project dependencies.

**What it does:**
1. Downloads Go modules (`go mod download`)
2. Tidies Go modules (`go mod tidy`)
3. Installs frontend dependencies (`npm install`)

**Example:**
```bash
python3 build.py deps
```

### `lint`
Runs golangci-lint on the codebase.

**Requirements:** `golangci-lint` must be installed

**Example:**
```bash
python3 build.py lint
```

### `doctor`
Checks the Wails development environment setup.

**What it checks:**
- Go installation and version
- Node.js installation and version
- Platform-specific dependencies (gcc, webkit, etc.)
- Wails configuration

**Example:**
```bash
python3 build.py doctor
```

### `install-wails`
Installs the Wails CLI tool.

**Example:**
```bash
python3 build.py install-wails
```

### `generate`
Generates Wails bindings (JavaScript/TypeScript bindings for Go functions).

**Example:**
```bash
python3 build.py generate
```

## Typical Workflows

### First Time Setup
```bash
# 1. Install Wails CLI
python3 build.py install-wails

# 2. Check environment
python3 build.py doctor

# 3. Install dependencies
python3 build.py deps

# 4. Run tests to verify
python3 build.py test

# 5. Start development
python3 build.py dev
```

### Daily Development
```bash
# Start dev server
python3 build.py dev

# In another terminal, run tests when making changes
python3 build.py test
```

### Before Committing
```bash
# Run tests
python3 build.py test

# Run linter (if installed)
python3 build.py lint

# Verify build works
python3 build.py build
```

### Creating a Release
```bash
# Clean everything
python3 build.py clean

# Install fresh dependencies
python3 build.py deps

# Run full test suite with coverage
python3 build.py test --coverage

# Build optimized binary
python3 build.py build --clean --upx
```

## Customization

The build script is designed to be easily customizable. To add new commands:

1. Add a new function in `build.py`
2. Add the command to the argument parser
3. Add the command to the main() function's command dispatcher

Example:
```python
def my_command():
    """My custom command"""
    print_header("Running My Command")
    # Your logic here
    return True

# In main():
parser.add_argument(
    "command",
    choices=["dev", "build", ..., "my-command"],
    help="Command to run"
)

# In command dispatcher:
elif args.command == "my-command":
    success = my_command()
```

## Troubleshooting

### "python3: command not found"
- **Windows:** Use `python` instead of `python3`
- **Linux/Mac:** Install Python 3: `sudo apt install python3` or `brew install python3`

### "wails: command not found"
Run: `python3 build.py install-wails`

### Build fails on Linux
Install platform dependencies:
```bash
sudo apt install gcc libgtk-3-dev libwebkit2gtk-4.0-dev
```

### Frontend dependencies fail
```bash
cd frontend
npm install
```

### Tests fail
Make sure you've run `python3 build.py deps` first to install dependencies.

## Migration from Makefile

If you're familiar with the old Makefile, here's the mapping:

| Makefile | Python Script |
|----------|---------------|
| `make dev` | `python3 build.py dev` |
| `make build` | `python3 build.py build` |
| `make build-prod` | `python3 build.py build --upx` |
| `make test` | `python3 build.py test` |
| `make test-coverage` | `python3 build.py test --coverage` |
| `make clean` | `python3 build.py clean` |
| `make deps` | `python3 build.py deps` |
| `make lint` | `python3 build.py lint` |
| `make install-wails` | `python3 build.py install-wails` |
| `make generate` | `python3 build.py generate` |
| `make doctor` | `python3 build.py doctor` |

## Advanced Usage

### Running with Different Python Versions
```bash
# Use specific Python version
python3.11 build.py dev

# Use Python from virtual environment
source venv/bin/activate
python build.py dev
```

### Combining with Other Tools
```bash
# Run in background (Linux/Mac)
python3 build.py dev &

# Time a build
time python3 build.py build

# Redirect output
python3 build.py test > test_results.txt 2>&1
```

### Environment Variables
The build script respects standard Go and Node.js environment variables:

```bash
# Use custom GOPATH
GOPATH=/custom/path python3 build.py build

# Use specific Node version (with nvm)
nvm use 18
python3 build.py deps
```

## Platform-Specific Notes

### Windows
- Use `python` instead of `python3`
- Use backslashes in paths when needed
- May need to run as Administrator for some operations

### macOS
- Ensure Xcode Command Line Tools are installed
- May need to allow unsigned apps in Security settings

### Linux
- Install platform dependencies first (see Troubleshooting)
- Make sure you have gcc and gtk development libraries

## Getting Help

```bash
# Show all available commands
python3 build.py --help

# Show detailed examples
python3 build.py --help | less
```

For issues, see the main README-WAILS.md or open an issue on GitHub.
