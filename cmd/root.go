package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "speedext",
	Short: "speed measurement tool",
}

func Execute() {
	rootCmd.Execute()
}
