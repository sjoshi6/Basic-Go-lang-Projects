package loadbalancer

import (
	"encoding/json"
	"go-keystore/database/redis"
	"go-keystore/model"
	"log"
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

var routes = model.Routes{
	model.Route{
		Name:        "Next",
		Method:      "GET",
		Pattern:     "/v1/nextNode",
		HandlerFunc: roundRobin,
	},
	model.Route{
		Name:        "Register",
		Method:      "POST",
		Pattern:     "/v1/register",
		HandlerFunc: register,
	},
}

// AddRoutes : Add routes to the router
func AddRoutes(router *mux.Router) *mux.Router {

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

	return router
}

/*
   Used to register a new storage node on loadbalancer
   Should be called by the node managing data partitioning only not by the new
   storage node
*/
func register(w http.ResponseWriter, r *http.Request) {

	redisconn, err := db.RedisConn(RedisURL)
	defer redisconn.Close()

	if err != nil {
		log.Println(err)
		ThrowInternalErrAndExit(w)
		return
	}

	log.Println(r.Body)

	var comm BalancerCommunication

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&comm)

	if err != nil {
		log.Println("Unable to decode the input json")
		log.Println(err)

		ThrowInternalErrAndExit(w)
		return
	}

	log.Printf("Retrived value successfully %s \n", comm.IPAddr)

	// Insert the new ip in redis
	_, err = redisconn.Do("RPUSH", "serverList", comm.IPAddr)
	if err != nil {
		log.Println(err)

		ThrowInternalErrAndExit(w)
		return
	}

	log.Printf("Registration Successful for %s \n", comm.IPAddr)

	RespondSuccessAndExit(w, "Registered the node successfully")
}

// Used to return back a new storage node IP on every request
func roundRobin(w http.ResponseWriter, r *http.Request) {

	elem, err := redis.String(conn.Do("LPOP", "serverList"))
	if err != nil {
		ThrowInternalErrAndExit(w)
		return
	}

	log.Printf("Retrived value successfully %s \n", elem)
	_, err = conn.Do("RPUSH", "serverList", elem)

	reply := BalancerCommunication{IPAddr: elem}
	jsonReply, err := json.Marshal(reply)
	if err != nil {
		ThrowInternalErrAndExit(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonReply)
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
