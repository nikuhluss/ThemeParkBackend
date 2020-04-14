package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/server"
)

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().IntP("port", "p", 5000, "the port to use when starting the HTTP server")
	viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {

		port := viper.GetInt("port")
		bindAddress := fmt.Sprintf(":%d", port)

		err := server.Start(bindAddress, testing, dokku)
		if err != nil {
			fmt.Printf("error starting server: %s\n", err)
			os.Exit(1)
		}
	},
}
