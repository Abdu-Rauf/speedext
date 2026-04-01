Write-Host "→ Creating virtual environment..."
python -m venv .venv

Write-Host "→ Installing Playwright into venv..."
.venv\Scripts\pip install playwright
.venv\Scripts\playwright install

Write-Host "→ Installing Go binary..."
go install .

Write-Host "Done. cd into the project directory and run: speedext --help"