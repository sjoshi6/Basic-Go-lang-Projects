package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/soveran/redisurl"
)

const (
	redisURL           = "redis://152.46.16.250:6379"
	masterMessageQueue = "master_message"
)

var wg sync.WaitGroup

func main() {

	// Common Variable to manage all processes
	exit := false

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage: main.go <master/slave/client/api-server>")
		os.Exit(1)
	}

	// Capturing which mode to launch go server in
	mode := args[1]
	fmt.Printf("Mode Selected %s \n", mode)

	// Manage name of user
	username := "test"
	if len(args) == 3 {
		username = args[2]
	}

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

	// Get Ip Address and key / value for this connection
	ipaddr := GetIPAddress()
	key := "online." + ipaddr
	val := ipaddr + ":8000"

	switch mode {
	case "slave":

		// Start command line driver
		go CommandLineInput(commandChan, &exit)
		go CmdHandler(commandChan, &exit)

		fmt.Printf("New Client Started at %s \n", ipaddr)

		// Register Slave to Redis DB
		go RegisterSlave(managerConn, key, val)

		// Start File Server
		go StartFileServer()

		// Start the main slave process
		Slave(ipaddr, &exit)

		// Send Heartbeats
		go SendHeartBeat(managerConn, key, val, &exit)

	case "master":

		go CommandLineInput(commandChan, &exit)
		go CmdHandler(commandChan, &exit)

		newSlaveChan := make(chan string)
		fmt.Printf("Master Started at %s \n", ipaddr)
		go ReceiveMessages(newSlaveChan, ipaddr)
		go HandleNewSlaves(newSlaveChan)
		go GetFileIPServer()

		for !exit {
			time.Sleep(1 * time.Second)
		}

	case "client":

		go FileSystemCommandHandler(&exit, username)

	case "api-server":

		go CommandLineInput(commandChan, &exit)
		go CmdHandler(commandChan, &exit)
		go StartServer(&exit)

	default:

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
