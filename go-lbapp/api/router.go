package api

import (
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
		SignUp,
	},
	Route{
		"login",
		"GET",
		"/v1/login",
		Login,
	},
}

// StartServer : Start the API Server by calling this function
func StartServer() {

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

	// Starting the api server
	log.Fatal(http.ListenAndServe(":8000", router))

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
