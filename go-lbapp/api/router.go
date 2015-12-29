package api

import (
	"go-lbapp/generics"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Central router for all API requests

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
	Route{
		"search_events",
		"GET",
		"/v1/search_events",
		SearchEventsByRange,
	},
}

// StartServer : Start the API Server by calling this function
func StartServer(port string, controller chan generics.SyncMsg) {

	// Starting the api server
	router := GetRouter()
	log.Fatal(http.ListenAndServe(port, router))

	// to exit the main function
	controller <- generics.SyncMsg{}
}

// GetRouter : Get an object of gorilla mux router
func GetRouter() *mux.Router {

	var router = mux.NewRouter().StrictSlash(true)
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

	return router
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
