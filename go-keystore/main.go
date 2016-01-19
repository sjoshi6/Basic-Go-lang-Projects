package main

import (
	"fmt"
	"go-keystore/client"
	"go-keystore/database/postgres"
	"go-keystore/loadbalancer"
	"go-keystore/model"
	"go-keystore/server"
	"log"
	"os"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
}

func main() {

	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("Usage : ./go-keystore <server/client/dual>")
	}

	switch args[0] {

	case "server":

		err := db.CreateTableIfNotExists()
		if err != nil {
			log.Fatal(err)
		}

		server.StartRPCServer()

	case "client":
		testRPCClient("localhost")

	case "loadbalancer":

		log.Println("Starting Load Balancer...")
		loadbalancer.StartLoadBalancer()

	default:
		log.Fatal("Incorrect choice of command line parameter")
	}

}

// testRPCClient : Test if RPC connections work
func testRPCClient(hostname string) {

	r, _ := client.NewClient(hostname)

	// Test for Get
	pair, _ := r.Get("aa")
	fmt.Printf("Key : %s has Value : %s \n", pair.Key, pair.Value)

	// Test for Put
	keypair := &model.KeyPair{
		Key:   "qq",
		Value: "{\"qq\":\"saurabh\"}",
	}

	success, _ := r.Put(keypair)
	fmt.Println(success)

	// Test for Delete
	fmt.Println("Beginning test for delete.. Creating mock")
	keypair = &model.KeyPair{
		Key:   "test",
		Value: "{\"test\":\"saurabh\"}",
	}

	success, _ = r.Put(keypair)
	fmt.Println(success)

	fmt.Println("Insert Mock Complete")
	fmt.Println("Beginning delete op...")

	success, _ = r.Delete("test")
	fmt.Println(success)

}
