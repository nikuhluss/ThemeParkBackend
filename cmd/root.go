package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	testing = false
	dokku   = false
)

func init() {
	viper.AutomaticEnv()
	rootCmd.PersistentFlags().BoolVarP(&testing, "testing", "t", false, "set testing flag")
	rootCmd.PersistentFlags().BoolVar(&dokku, "dokku", false, "set dokku flag")
}

var rootCmd = &cobra.Command{
	Use:   "backend",
	Short: "Backend server and data generator for database schema",
	Run: func(cmd *cobra.Command, args []string) {
		serverCmd.Run(cmd, args)
	},
}

// Execute executes the root cmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
