package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/server"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		err := server.Start(":5000")
		if err != nil {
			fmt.Printf("error starting server: %s\n", err)
			os.Exit(1)
		}
	},
}
