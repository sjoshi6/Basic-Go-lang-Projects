package main

import (
	"fmt"
	"go-lbapp/controllers"
	"go-lbapp/generics"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
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

	go StartServer(port, controller)
	<-controller
}

// StartServer : Start the API Server by calling this function
func StartServer(port string, controller chan generics.SyncMsg) {

	// Creating a new mux router
	var router = mux.NewRouter().StrictSlash(true)

	// Add APP routes from various controllers
	router = controllers.SetBaseRoutes(router)
	router = controllers.SetResourceRoutes(router)

	// This route is essential to view the monitoring stats for the app.
	router.Handle("/debug/vars", http.DefaultServeMux)

	log.Fatal(http.ListenAndServe(port, router))

	// to exit the main function
	controller <- generics.SyncMsg{}
}
