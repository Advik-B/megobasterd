#!/usr/bin/env python3
"""
Build script for MegaBasterd Go Edition
Replaces Makefile with cross-platform Python script
"""

import argparse
import os
import shutil
import subprocess
import sys
from pathlib import Path


class Colors:
    """ANSI color codes for terminal output"""
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKCYAN = '\033[96m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'


def print_header(msg):
    """Print a colored header message"""
    print(f"\n{Colors.HEADER}{Colors.BOLD}=== {msg} ==={Colors.ENDC}\n")


def print_success(msg):
    """Print a success message"""
    print(f"{Colors.OKGREEN}✓ {msg}{Colors.ENDC}")


def print_error(msg):
    """Print an error message"""
    print(f"{Colors.FAIL}✗ {msg}{Colors.ENDC}", file=sys.stderr)


def print_info(msg):
    """Print an info message"""
    print(f"{Colors.OKCYAN}ℹ {msg}{Colors.ENDC}")


def run_command(cmd, cwd=None, check=True):
    """Run a shell command and return the result"""
    print_info(f"Running: {' '.join(cmd)}")
    try:
        result = subprocess.run(
            cmd,
            cwd=cwd,
            check=check,
            capture_output=False,
            text=True
        )
        return result.returncode == 0
    except subprocess.CalledProcessError as e:
        print_error(f"Command failed with exit code {e.returncode}")
        return False
    except FileNotFoundError:
        print_error(f"Command not found: {cmd[0]}")
        return False


def dev():
    """Run the application in development mode with hot reload"""
    print_header("Starting Development Server")
    return run_command(["wails", "dev"])


def build(clean=False, upx=False):
    """Build the application"""
    print_header("Building Application")
    
    cmd = ["wails", "build"]
    if clean:
        cmd.append("-clean")
    if upx:
        cmd.append("-upx")
    
    return run_command(cmd)


def test(coverage=False):
    """Run tests"""
    print_header("Running Tests")
    
    # Exclude cmd/megobasterd from tests as it requires frontend build
    
    if coverage:
        # Run tests with coverage
        cmd = ["go", "test", "-v", "-coverprofile=coverage.out", "./internal/...", "./pkg/..."]
        if not run_command(cmd):
            return False
        print_success("Tests passed with coverage")
        
        # Generate HTML coverage report
        print_info("Generating coverage report...")
        if run_command(["go", "tool", "cover", "-html=coverage.out", "-o", "coverage.html"]):
            print_success("Coverage report generated: coverage.html")
            return True
        return False
    else:
        # Run tests without coverage
        cmd = ["go", "test", "-v", "./internal/...", "./pkg/..."]
        if run_command(cmd):
            print_success("All tests passed")
            return True
        return False


def clean():
    """Clean build artifacts"""
    print_header("Cleaning Build Artifacts")
    
    paths_to_remove = [
        "build/",
        "frontend/dist/",
        "frontend/node_modules/",
        "coverage.out",
        "coverage.html"
    ]
    
    for path in paths_to_remove:
        full_path = Path(path)
        if full_path.exists():
            if full_path.is_dir():
                print_info(f"Removing directory: {path}")
                shutil.rmtree(full_path)
            else:
                print_info(f"Removing file: {path}")
                full_path.unlink()
            print_success(f"Removed: {path}")
        else:
            print_info(f"Skipping (not found): {path}")
    
    print_success("Clean complete")
    return True


def deps():
    """Install dependencies"""
    print_header("Installing Dependencies")
    
    # Go dependencies
    print_info("Installing Go dependencies...")
    if not run_command(["go", "mod", "download"]):
        return False
    if not run_command(["go", "mod", "tidy"]):
        return False
    print_success("Go dependencies installed")
    
    # Frontend dependencies
    frontend_dir = Path("frontend")
    if frontend_dir.exists():
        print_info("Installing frontend dependencies...")
        if run_command(["npm", "install"], cwd=frontend_dir):
            print_success("Frontend dependencies installed")
            return True
        return False
    else:
        print_warning("Frontend directory not found, skipping npm install")
        return True


def lint():
    """Run linter"""
    print_header("Running Linter")
    
    if run_command(["golangci-lint", "run", "./..."], check=False):
        print_success("Linting complete")
        return True
    else:
        print_error("Linting failed (or golangci-lint not installed)")
        return False


def install_wails():
    """Install Wails CLI"""
    print_header("Installing Wails")
    
    if run_command(["go", "install", "github.com/wailsapp/wails/v2/cmd/wails@latest"]):
        print_success("Wails installed successfully")
        return True
    return False


def generate():
    """Generate Wails bindings"""
    print_header("Generating Wails Bindings")
    
    if run_command(["wails", "generate", "module"]):
        print_success("Bindings generated")
        return True
    return False


def doctor():
    """Check Wails setup"""
    print_header("Checking Wails Setup")
    
    return run_command(["wails", "doctor"])


def main():
    """Main entry point"""
    parser = argparse.ArgumentParser(
        description="Build script for MegaBasterd Go Edition",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  python build.py dev              # Run development server
  python build.py build            # Build application
  python build.py build --upx      # Build with UPX compression
  python build.py test             # Run tests
  python build.py test --coverage  # Run tests with coverage
  python build.py clean            # Clean build artifacts
  python build.py deps             # Install dependencies
  python build.py lint             # Run linter
  python build.py doctor           # Check Wails setup
        """
    )
    
    parser.add_argument(
        "command",
        choices=["dev", "build", "test", "clean", "deps", "lint", "install-wails", "generate", "doctor"],
        help="Command to run"
    )
    
    parser.add_argument(
        "--coverage",
        action="store_true",
        help="Generate coverage report (for test command)"
    )
    
    parser.add_argument(
        "--clean",
        action="store_true",
        help="Clean before building (for build command)"
    )
    
    parser.add_argument(
        "--upx",
        action="store_true",
        help="Use UPX compression (for build command)"
    )
    
    args = parser.parse_args()
    
    # Change to script directory
    script_dir = Path(__file__).parent
    os.chdir(script_dir)
    
    # Execute command
    success = False
    
    if args.command == "dev":
        success = dev()
    elif args.command == "build":
        success = build(clean=args.clean, upx=args.upx)
    elif args.command == "test":
        success = test(coverage=args.coverage)
    elif args.command == "clean":
        success = clean()
    elif args.command == "deps":
        success = deps()
    elif args.command == "lint":
        success = lint()
    elif args.command == "install-wails":
        success = install_wails()
    elif args.command == "generate":
        success = generate()
    elif args.command == "doctor":
        success = doctor()
    
    # Exit with appropriate code
    if success:
        print(f"\n{Colors.OKGREEN}{Colors.BOLD}✓ Success!{Colors.ENDC}\n")
        sys.exit(0)
    else:
        print(f"\n{Colors.FAIL}{Colors.BOLD}✗ Failed!{Colors.ENDC}\n")
        sys.exit(1)


if __name__ == "__main__":
    main()
