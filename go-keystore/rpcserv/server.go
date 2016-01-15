package rpcserv

import (
	"go-keystore/config"
	"log"
	"net"
	"net/rpc"
)

// StartRPCServer : Start RPC server for prodcedure calls
func StartRPCServer() {

	// Get a new obj from factory & register it as an rpcObject
	rpcObj := NewRPC()
	rpc.Register(rpcObj)

	// Start a tcp connection to allow execution of rpc
	port := ":" + string(settings.ServerPort)
	conn, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starts accepting rpc connections")

	// Accept Remote prodcedure calls over this tcp conn
	rpc.Accept(conn)
}
