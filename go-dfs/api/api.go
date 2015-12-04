package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	redisURL           = "redis://152.46.16.250:6379"
	masterMessageQueue = "master_message"
)

// Add new routes here
var router = mux.NewRouter().StrictSlash(true)
var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		HomePage,
	},
	Route{
		"AllPaths",
		"GET",
		"/all",
		GetAll,
	},
}

// StartServer : Begin the API Server
func StartServer(exit *bool) {

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
	log.Fatal(http.ListenAndServe(":8080", router))

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
