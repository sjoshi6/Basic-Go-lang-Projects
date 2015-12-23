package api

import (
	"go-lbapp/generics"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Central router for all API requests
var router = mux.NewRouter().StrictSlash(true)
var routes = Routes{
	Route{
		"signup",
		"POST",
		"/v1/signup",
		CreateAccount,
	},
	Route{
		"login",
		"POST",
		"/v1/login",
		ConfirmCredentials,
	},
	Route{
		"create_event",
		"POST",
		"/v1/create_event",
		CreateEvent,
	},
}

// StartServer : Start the API Server by calling this function
func StartServer(controller chan generics.SyncMsg) {

	for _, route := range routes {

		var handler http.Handler
		handler = route.HandlerFunc
		handler = logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	// This route is essential to view the monitoring stats for the app.
	router.Handle("/debug/vars", http.DefaultServeMux)

	// Starting the api server
	log.Fatal(http.ListenAndServe(":8000", router))

	// to exit the main function
	controller <- generics.SyncMsg{}
}

// Used to log route access times
func logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
