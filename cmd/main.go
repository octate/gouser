package main

import (
	"fmt"
	"os"
)

func main() {
	switch os.Getenv("MODE") {

	case "server":
		serverRun()

	default:
		fmt.Println("Unknown mode. Exiting.")
	}
}
