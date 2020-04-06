package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	testing = false
)

var rootCmd = &cobra.Command{
	Use:   "backend",
	Short: "Backend server and data generator for database schema",
	Run: func(cmd *cobra.Command, args []string) {
		serverCmd.Run(cmd, args)
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&testing, "testing", "t", false, "Set testing flag")
}

// Execute executes the root cmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
