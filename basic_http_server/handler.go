package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"

	"github.com/gorilla/mux"
)

// Index : Used for handling the index route
func Index(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))

}

// TodoIndex : Used for handling the todo access
func TodoIndex(w http.ResponseWriter, r *http.Request) {

	todos := Todos{
		Todo{Name: "Write Presentation"},
		Todo{Name: "AutoSPark complete"},
	}

	json.NewEncoder(w).Encode(todos)
}

// TodoShow : handle individual todo
func TodoShow(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	todoid := vars["todoid"]
	fmt.Fprintf(w, "Todo show: %q", todoid)

}

//HelloHandle Function to handle hello world
func HelloHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World !!!")
}
