package loadbalancer

import (
	"fmt"
	"go-keystore/database/redis"
	"log"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

var conn redis.Conn
var RedisURL string

// StartLoadBalancer : Used for round robin circulation of requests for go-keystore
func StartLoadBalancer() {

	// Get IP Address to establish redis conn
	ipAddr := GetIPAddress()
	RedisURL = fmt.Sprintf("redis://%s:6379", ipAddr)

	redisconn, err := db.RedisConn(RedisURL)
	conn = redisconn

	log.Printf("Connecting to redis at : %s \n", RedisURL)
	if err != nil {
		log.Println("Could not connect to redis DB")
		log.Fatal(err)
	}

	StartHTTPServer()

}

// StartHTTPServer : Start a MUX router to serve next request
func StartHTTPServer() {

	var router = mux.NewRouter().StrictSlash(true)
	router = AddRoutes(router)

	log.Println("Router ready to start. Initializing...")
	// Starting the api server
	log.Fatal(http.ListenAndServe(":8000", router))

}
