package main

import (
	"fmt"
	"os"
	"time"

	"github.com/soveran/redisurl"
)

const (
	redisURL           = "redis://152.46.16.250:6379"
	masterMessageQueue = "master_message"
)

func main() {

	// Common Variable to manage all processes
	exit := false

	args := os.Args
	if len(args) != 2 {
		fmt.Println("Usage: main.go <master/slave>")
		os.Exit(1)
	}

	// Capturing which mode to launch go server in
	masterslvToggle := args[1]
	fmt.Printf("Mode Selected %s \n", masterslvToggle)

	// Use this connection only for setup activities of the node. No more communication should happen through this
	managerConn, err := redisurl.ConnectToURL(redisURL)
	if err != nil {

		fmt.Println(err)
		os.Exit(1)

	}

	// Before function exits close the connection
	defer managerConn.Close()

	// go Channel for commands common for master and slave
	commandChan := make(chan string)
	go CommandLineInput(commandChan, &exit)
	go CmdHandler(commandChan, &exit)

	// Get Ip Address and key / value for this connection
	ipaddr := GetIPAddress()
	key := "online." + ipaddr
	val := ipaddr + ":8000"

	if masterslvToggle == "slave" {

		fmt.Printf("New Client Started at %s \n", ipaddr)

		// Register Slave to Redis DB
		go RegisterSlave(managerConn, key, val)

		// Start the main slave process
		Slave(ipaddr, &exit)

		// Send Heartbeats
		go SendHeartBeat(managerConn, key, val, &exit)

	} else if masterslvToggle == "master" {

		newSlaveChan := make(chan string)
		fmt.Printf("Master Started at %s \n", ipaddr)
		go ReceiveMessages(newSlaveChan, ipaddr)
		go HandleNewSlaves(newSlaveChan)
		Master()

	} else {

		fmt.Println("Incorrect command line argument. Either use master or slave")
		os.Exit(1)
	}

	for !exit {

		time.Sleep(1 * time.Second)
	}

	// Remove the user before slave function exits
	managerConn.Do("SREM", "online_slaves", ipaddr)
	managerConn.Do("DEL", key)

}
