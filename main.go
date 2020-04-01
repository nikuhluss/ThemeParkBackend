package main

import (
	"fmt"
	"os"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/server"
)

func main() {
	err := server.Start(":5000")
	if err != nil {
		fmt.Printf("error stsarting server: %s\n", err)
		os.Exit(1)
	}
}
