package main

import (
	"fmt"
	"os"

	"github.com/soveran/redisurl"
)

const (
	redisURL = "redis://152.46.16.250:6379"
)

func main() {

	args := os.Args
	if len(args) != 2 {
		fmt.Println("Usage: main.go <master/slave>")
		os.Exit(1)
	}

	// Capturing which mode to launch go server in
	masterslvToggle := args[1]
	fmt.Printf("Mode Selected %s \n", masterslvToggle)

	// Use this connection only for setup activities of the node. No more communication should happen through this
	managerConn, err := redisurl.ConnectToURL("redis://localhost:6379")
	if err != nil {

		fmt.Println(err)
		os.Exit(1)

	}

	// Before function exits close the connection
	defer managerConn.Close()

	if masterslvToggle == "slave" {

		ipaddr := GetIPAddress()
		RegisterSlave(managerConn, ipaddr)
		Slave()

	} else if masterslvToggle == "master" {
		Master()
	} else {
		fmt.Println("Incorrect command line argument. Either use master or slave")
		os.Exit(1)
	}

}
