package main

import (
	"fmt"
	"go-keystore/database/mysql"
	"go-keystore/rpcserv"
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

		rpcserv.StartRPCServer()

	case "client":
		testRPCClient("localhost")

	default:
		log.Fatal("Incorrect choice of command line parameter")
	}

}

// testRPCClient : Test if RPC connections work
func testRPCClient(hostname string) {
	r, _ := rpcserv.NewClient(hostname)
	pair, _ := r.Get("aa")
	fmt.Printf("Key : %s has Value : %s \n", pair.Key, pair.Value)
}
