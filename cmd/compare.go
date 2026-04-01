package cmd

import (
	"fmt"

	"github.com/Abdu-Rauf/speedext/speedtest"
	"github.com/spf13/cobra"
)

var compCommand = &cobra.Command{
	Use:   "compare",
	Short: "A comparison of Speed testers",
	Run: func(cmd *cobra.Command, args []string) {
		goDown := speedtest.RunDownload()
		ookla, _ := runScraper("pyscrapers/ookla_ext.py")
		fast, _ := runScraper("pyscrapers/fast_ext.py")

		fmt.Printf("%-15s %-15s\n", "Tester", "Download Speed (Mbps)\n")
		fmt.Printf("%-15s %-15.2f\n", "Go Speedtest", goDown)
		fmt.Printf("%-15s %-15s\n", "Ookla", ookla)
		fmt.Printf("%-15s %-15s\n", "Fast", fast)
	},
}

func init() {
	rootCmd.AddCommand(compCommand)
}
