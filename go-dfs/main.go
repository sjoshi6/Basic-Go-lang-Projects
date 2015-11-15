package main

import (
	"fmt"
	"os"
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

	if masterslvToggle == "slave" {
		Slave()

	} else if masterslvToggle == "master" {
		Master()
	} else {
		fmt.Println("Incorrect command line argument. Either use master or slave")
		os.Exit(1)
	}

}
