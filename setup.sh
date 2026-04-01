#!/bin/bash
set -e

echo "→ Checking prerequisites..."
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed"
    exit 1
fi

if ! command -v python3 &> /dev/null; then
    echo "Error: Python3 is not installed"
    exit 1
fi

echo "→ Creating virtual environment..."
python3 -m venv .venv

echo "→ Installing Playwright into venv..."
.venv/bin/pip install playwright
.venv/bin/playwright install

echo "→ Installing Go binary..."
go install .

echo ""
echo "Done. cd into the project directory and run: speedext --help"