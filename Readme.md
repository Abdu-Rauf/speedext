# Speedext

Speedext is a Go-based CLI speed testing tool that leverages the Cloudflare API to test speed. It provides the user with the option to run the native Go tester, run Python-based scrapers (e.g., Fast.com, Ookla), or display a comparison table.

## Features

- **Native Go Tester:** Directly tests download and upload speeds using Cloudflare's edge network (`__down` and `__up` endpoints) for high-performance, low-overhead results.
- **Scraper Mode:** Uses Python and Playwright to scrape speed test values from third-party sources such as Ookla and Fast.com.
- **CLI Interface:** Built with Cobra for clean, easy-to-use terminal commands.

## Prerequisites

Before running Speedext, ensure you have the following installed to support both the Go binary and the Python scrapers:
- **Go** (1.20 or newer)
- **Python** (3.10 or newer)

## Setup

Run the setup scripts to configure the Python virtual environment, install Playwright dependencies, and compile the Go binary.

**For Linux / macOS:**
```bash
chmod +x setup.sh
./setup.sh
```

**For Windows (PowerShell):**
```powershell
.\setup.ps1
```

## Usage

Because the scraper command relies on a local Python virtual environment (`.venv`), **you must run `speedext` from within the root of this project directory.**

Once installed via the setup script, ensure you are in the `speedext` folder and run:

```bash
# Run the fast Go-based Cloudflare speed test
speedext go

# Run the Python Playwright scrapers (Ookla, Fast.com)
speedext scrape

# View all available commands
speedext --help
```

## Changes / Improvements to make

- [ ] Allow flags to help users select which specific tester to scrape from (e.g., `--ookla` or `--fast`).
- [ ] Add upload speed metric parsing to the scrapers.
- [ ] Test embedding Python code/scripts into the Go binary to remove the local `.venv` dependency, allowing `speedext` to be run globally from any directory.