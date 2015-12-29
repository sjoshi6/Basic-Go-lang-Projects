package main

import (
	"fmt"
	"go-lbapp/api"
	"go-lbapp/generics"
	"log"
	"os"
	"time"
)

func main() {

	// Accept port number from command line arguments and start server
	cmdLineArgs := os.Args[1:]

	if len(cmdLineArgs) < 1 {
		log.Fatal("Usage : ./go-lbapp <portnumber>")
	}

	port := ":" + string(cmdLineArgs[0])

	controller := make(chan generics.SyncMsg)
	fmt.Println("Go API Server - Logs", time.Now())
	fmt.Printf("API Server started at - %s\n", port)

	go api.StartServer(port, controller)
	<-controller
}
