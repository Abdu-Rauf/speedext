package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func getPythonPath(pwd string) string {
	// windows
	winPath := filepath.Join(pwd, ".venv", "Scripts", "python.exe")
	if _, err := os.Stat(winPath); err == nil {
		return winPath
	}
	// unix
	return filepath.Join(pwd, ".venv", "bin", "python")
}

func runScraper(spath string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}
	// Get path of the virtual env
	pythonExe := getPythonPath(pwd)
	scriptPath := filepath.Join(pwd, spath)

	execCommand := exec.Command(pythonExe, scriptPath)
	execCommand.Dir = pwd

	op, err := execCommand.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("scraper failed: %s", string(exitErr.Stderr))
		}
		return "", err
	}
	return string(op), nil
}

var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Measure speed using Python scrapers",
	Run: func(cmd *cobra.Command, args []string) {
		ookla, err := runScraper("pyscrapers/ookla_ext.py")
		if err != nil {
			fmt.Println("ookla error:", err)
		} else {
			fmt.Printf("Ookla Download: %s", ookla)
		}

		fast, err := runScraper("pyscrapers/fast_ext.py")
		if err != nil {
			fmt.Println("fast error:", err)
		} else {
			fmt.Printf("Fast Download: %s", fast)
		}
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)
}
