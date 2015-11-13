package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	router := createRouter()
	fmt.Println("Go Server - Logs", time.Now())
	log.Fatal(http.ListenAndServe(":8080", router))
}
