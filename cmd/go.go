package cmd

import (
	"fmt"

	"github.com/Abdu-Rauf/speedext/speedtest"
	"github.com/spf13/cobra"
)

var goCmd = &cobra.Command{
	Use:   "go",
	Short: "Measure using Go",
	Run: func(cmd *cobra.Command, args []string) {
		download := speedtest.RunDownload()
		upload := speedtest.RunUpload()
		fmt.Printf("Download: %.2f Mbps\n", download)
		fmt.Printf("Upload:   %.2f Mbps\n", upload)

	},
}

func init() {
	rootCmd.AddCommand(goCmd)
}
