#!/usr/bin/env python3
"""
Quick run script for MegaBasterd Go Edition
Shortcut to run common commands
"""

import sys
import subprocess
from pathlib import Path

# Change to script directory
script_dir = Path(__file__).parent
build_script = script_dir / "build.py"

# If no arguments, show help
if len(sys.argv) == 1:
    subprocess.run([sys.executable, str(build_script), "--help"])
    sys.exit(0)

# Pass all arguments to build.py
subprocess.run([sys.executable, str(build_script)] + sys.argv[1:])
