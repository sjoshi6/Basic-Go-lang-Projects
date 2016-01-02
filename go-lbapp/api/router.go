package api

import (
	"go-lbapp/controllers"
	"go-lbapp/generics"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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
